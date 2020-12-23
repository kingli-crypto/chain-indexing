package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	account_view "github.com/crypto-com/chain-indexing/appinterface/projection/account/view"

	"github.com/crypto-com/chain-indexing/appinterface/projection/rdbprojectionbase"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	event_entity "github.com/crypto-com/chain-indexing/entity/event"
	applogger "github.com/crypto-com/chain-indexing/internal/logger"
	event_usecase "github.com/crypto-com/chain-indexing/usecase/event"
)

func ConvertToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

type Account struct {
	*rdbprojectionbase.Base

	rdbConn rdb.Conn
	logger  applogger.Logger
}

func NewAccount(logger applogger.Logger, rdbConn rdb.Conn) *Account {
	return &Account{
		rdbprojectionbase.NewRDbBase(rdbConn.ToHandle(), "Account"),

		rdbConn,
		logger,
	}
}

func (_ *Account) GetEventsToListen() []string {
	return []string{event_usecase.ACCOUNT_TRANSFERRED}
}

func (projection *Account) OnInit() error {
	return nil
}

func (projection *Account) HandleEvents(height int64, events []event_entity.Event) error {
	rdbTx, err := projection.rdbConn.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = rdbTx.Rollback()
		}
	}()

	rdbTxHandle := rdbTx.ToHandle()

	accountsView := account_view.NewAccounts(rdbTxHandle)

	for _, event := range events {
		if accountCreatedEvent, ok := event.(*event_usecase.AccountTransferred); ok {
			if handleErr := projection.handleAccountCreatedEvent(accountsView, accountCreatedEvent); handleErr != nil {
				return fmt.Errorf("error handling AccountCreatedEvent: %v", handleErr)
			}
		}
	}

	if err = projection.UpdateLastHandledEventHeight(rdbTxHandle, height); err != nil {
		return fmt.Errorf("error updating last handled event height: %v", err)
	}

	if err = rdbTx.Commit(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}
	committed = true

	return nil
}

func (projection *Account) handleAccountCreatedEvent(accountsView *account_view.Accounts, event *event_usecase.AccountTransferred) error {

	recipienterr := projection.writeAccountInfo(accountsView, event.Recipient)
	if recipienterr != nil {
		return recipienterr
	}

	sendererr := projection.writeAccountInfo(accountsView, event.Sender)
	if sendererr != nil {
		return sendererr
	}

	return nil
}

func GetAccountInfo(address string) (accounttype string, accountaddress string, pubkey string, accountnumber string, sequencenumber string, err error) {

	resp, err := http.Get(fmt.Sprintf("https://testnet-croeseid-1.crypto.com:1317/cosmos/auth/v1beta1/accounts/%s", address))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	outputbytes, _ := ioutil.ReadAll(resp.Body)

	var myjson map[string]interface{}

	if err := json.Unmarshal(outputbytes, &myjson); err != nil {
		panic(err)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, outputbytes, "", "\t")

	fmt.Println(string(prettyJSON.Bytes()))

	var myaccounttype string
	var myaddress string
	var mypubkeycontainer map[string]interface{}
	var mypubkey string
	var myaccountnumber string
	var mysequence string

	myaccount := myjson["account"].(map[string]interface{})
	myaccounttypemeta := myaccount["@type"].(string)
	myaccounttypename, nameok := myaccount["name"].(string)
	if !nameok {
		myaccounttypename = ""
	}
	myaccounttype = fmt.Sprintf("%s %s", myaccounttypemeta, myaccounttypename)
	mybaseaccount, mybaseaccountok := myaccount["base_account"].(map[string]interface{})

	if !mybaseaccountok {
		// normal account
		myaddress = myaccount["address"].(string)
		mypubkeycontainer = myaccount["pub_key"].(map[string]interface{})
		mypubkey = mypubkeycontainer["key"].(string)
		myaccountnumber = myaccount["account_number"].(string)
		mysequence = myaccount["sequence"].(string)

	} else {
		// base account ok
		myaddress = mybaseaccount["address"].(string)
		mypubkeycontainer, _ = mybaseaccount["pub_key"].(map[string]interface{})
		myaccountnumber = mybaseaccount["account_number"].(string)
		mysequence = mybaseaccount["sequence"].(string)

	}
	//fmt.Println(myaddress, mypubkey, myaccountnumber, mysequence)

	return myaccounttype, myaddress, mypubkey, myaccountnumber, mysequence, nil

}

func GetAccountBalance(address string, denom string) (retbalance string, retdenom string, err error) {

	resp, err := http.Get(fmt.Sprintf("https://testnet-croeseid-1.crypto.com:1317/cosmos/bank/v1beta1/balances/%s/%s", address, denom))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	outputbytes, _ := ioutil.ReadAll(resp.Body)

	var myjson map[string]interface{}

	if err := json.Unmarshal(outputbytes, &myjson); err != nil {
		panic(err)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, outputbytes, "", "\t")
	fmt.Println(string(prettyJSON.Bytes()))
	mybalance := myjson["balance"].(map[string]interface{})
	myamount := mybalance["amount"].(string)
	mydenom := mybalance["denom"].(string)

	return myamount, mydenom, nil
}

func (projection *Account) writeAccountInfo(accountsView *account_view.Accounts, whichaddress string) error {

	atype, aaddress, pubkey, aaccountnumber, asequenenumber, _ := GetAccountInfo(whichaddress)
	abalance, adenom, _ := GetAccountBalance(whichaddress, "basecro")

	if err := accountsView.Upsert(&account_view.Account{
		AccountType:    atype,
		AccountAddress: aaddress,
		Pubkey:         pubkey,
		AccountNumber:  ConvertToInt64(aaccountnumber),
		SequenceNumber: ConvertToInt64(asequenenumber),
		AccountBalance: ConvertToInt64(abalance),
		AccountDenom:   adenom,
	}); err != nil {
		return err
	}

	return nil
}
