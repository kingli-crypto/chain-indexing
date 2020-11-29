package transasaction

import (
	"fmt"

	view2 "github.com/crypto-com/chain-indexing/appinterface/projection/transasaction/view"

	projection_entity "github.com/crypto-com/chain-indexing/entity/projection"

	"github.com/crypto-com/chain-indexing/appinterface/projection/rdbbase"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	event_entity "github.com/crypto-com/chain-indexing/entity/event"
	applogger "github.com/crypto-com/chain-indexing/internal/logger"
	"github.com/crypto-com/chain-indexing/internal/utctime"
	event_usecase "github.com/crypto-com/chain-indexing/usecase/event"
)

var _ projection_entity.Projection = &Transaction{}

type Transaction struct {
	*rdbbase.RDbBase

	rdbConn rdb.Conn
	logger  applogger.Logger
}

func NewTransaction(logger applogger.Logger, rdbConn rdb.Conn) *Transaction {
	return &Transaction{
		rdbbase.NewRDbBase(rdbConn.ToHandle(), "Transaction"),

		rdbConn,
		logger,
	}
}

func (_ *Transaction) GetEventsToListen() []string {
	return append([]string{
		event_usecase.BLOCK_CREATED,
		event_usecase.TRANSACTION_CREATED,
		event_usecase.TRANSACTION_FAILED,
	}, event_usecase.MSG_EVENTS...)
}

func (projection *Transaction) OnInit() error {
	return nil
}

func (projection *Transaction) HandleEvents(height int64, events []event_entity.Event) error {
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
	transactionsView := view2.NewTransactions(rdbTxHandle)

	var blockTime utctime.UTCTime
	var blockHash string
	txs := make([]view2.TransactionRow, 0)
	txMsgs := make(map[string][]event_usecase.MsgEvent)
	for _, event := range events {
		if blockCreatedEvent, ok := event.(*event_usecase.BlockCreated); ok {
			blockTime = blockCreatedEvent.Block.Time
			blockHash = blockCreatedEvent.Block.Hash
		} else if transactionCreatedEvent, ok := event.(*event_usecase.TransactionCreated); ok {
			txs = append(txs, view2.TransactionRow{
				BlockHeight:   height,
				BlockTime:     utctime.UTCTime{}, // placeholder
				Hash:          transactionCreatedEvent.TxHash,
				Success:       true,
				Code:          transactionCreatedEvent.Code,
				Log:           transactionCreatedEvent.Log,
				Fee:           transactionCreatedEvent.Fee,
				FeePayer:      transactionCreatedEvent.FeePayer,
				FeeGranter:    transactionCreatedEvent.FeeGranter,
				GasWanted:     transactionCreatedEvent.GasWanted,
				GasUsed:       transactionCreatedEvent.GasUsed,
				Memo:          transactionCreatedEvent.Memo,
				TimeoutHeight: transactionCreatedEvent.TimeoutHeight,
				Messages:      make([]view2.TransactionRowMessage, 0),
			})
		} else if transactionFailedEvent, ok := event.(*event_usecase.TransactionFailed); ok {
			txs = append(txs, view2.TransactionRow{
				BlockHeight:   height,
				BlockTime:     utctime.UTCTime{}, // placeholder
				Hash:          transactionFailedEvent.TxHash,
				Success:       false,
				Code:          transactionFailedEvent.Code,
				Log:           transactionFailedEvent.Log,
				Fee:           transactionFailedEvent.Fee,
				FeePayer:      transactionFailedEvent.FeePayer,
				FeeGranter:    transactionFailedEvent.FeeGranter,
				GasWanted:     transactionFailedEvent.GasWanted,
				GasUsed:       transactionFailedEvent.GasUsed,
				Memo:          transactionFailedEvent.Memo,
				TimeoutHeight: transactionFailedEvent.TimeoutHeight,
				Messages:      make([]view2.TransactionRowMessage, 0),
			})
		} else if msgEvent, ok := event.(event_usecase.MsgEvent); ok {
			if _, exist := txMsgs[msgEvent.TxHash()]; !exist {
				txMsgs[msgEvent.TxHash()] = make([]event_usecase.MsgEvent, 0)
			}
			txMsgs[msgEvent.TxHash()] = append(txMsgs[msgEvent.TxHash()], msgEvent)
		}
	}

	for _, tx := range txs {
		mutTx := tx
		mutTx.BlockTime = blockTime
		mutTx.BlockHash = blockHash

		for _, msg := range txMsgs[mutTx.Hash] {
			mutTx.Messages = append(mutTx.Messages, view2.TransactionRowMessage{
				Type:    msg.MsgType(),
				Content: msg,
			})
		}

		if insertErr := transactionsView.Insert(&mutTx); insertErr != nil {
			return fmt.Errorf("error inserting transaction into view: %v", insertErr)
		}
	}

	if err := projection.UpdateLastHandledEventHeight(rdbTxHandle, height); err != nil {
		return fmt.Errorf("error updating last handled event height: %v", err)
	}

	if err := rdbTx.Commit(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}
	committed = true
	return nil
}