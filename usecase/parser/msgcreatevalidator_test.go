package parser_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/infrastructure/tendermint"
	"github.com/crypto-com/chainindex/usecase/coin"
	command_usecase "github.com/crypto-com/chainindex/usecase/command"
	"github.com/crypto-com/chainindex/usecase/event"
	"github.com/crypto-com/chainindex/usecase/model"
	"github.com/crypto-com/chainindex/usecase/parser"
	usecase_parser_test "github.com/crypto-com/chainindex/usecase/parser/test"
)

var _ = Describe("ParseMsgCommands", func() {
	Describe("MsgBeginRedelegate", func() {

		It("should parse Msg commands when there is staking.MsgCreateValidator in the transaction", func() {
			txDecoder := parser.NewTxDecoder("basetrcro")
			block, _, _ := tendermint.ParseBlockResp(strings.NewReader(usecase_parser_test.TX_MSG_CREATE_VALIDATOR_BLOCK_RESP))
			blockResults, _ := tendermint.ParseBlockResultsResp(strings.NewReader(usecase_parser_test.TX_MSG_CREATE_VALIDATOR_BLOCK_RESULTS_RESP))

			cmds := parser.ParseMsgToCommands(
				txDecoder,
				block,
				blockResults,
			)
			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			for i, s := range cmds {
				fmt.Printf("%d this_cmd Name= %s  %+v\n", i, s.Name(), s)
			}

			fmt.Println("############################################################################")
			this_cmd := cmds[2]
			Expect(this_cmd.Name()).To(Equal("CreateMsgCreateValidator"))

			fmt.Printf("%d\n", len(cmds))
			Expect(cmds).To(HaveLen(5))
			commiossionrates := model.CommissionRates{
				Rate:          "0.100000000000000000",
				MaxRate:       "0.200000000000000000",
				MaxChangeRate: "0.010000000000000000",
			}

			Expect(this_cmd).To(Equal(command_usecase.NewCreateMsgCreateValidator(
				event.MsgCommonParams{
					BlockHeight: int64(503978),
					TxHash:      "E69985AC8168383A81B7952DBE03EB9B3400FF80AEC0F362369DD7F38B1C2FE9",
					TxSuccess:   true,
					MsgIndex:    0,
				},
				model.MsgCreateValidatorParams{
					CommissionRates:  commiossionrates,
					DelegatorAddress: "tcro1fmprm0sjy6lz9llv7rltn0v2azzwcwzvk2lsyn",
					ValidatorAddress: "tcrocncl1fmprm0sjy6lz9llv7rltn0v2azzwcwzvr4ufus",
					PubKey:           "tcrocnclconspub1zcjduepqa5rksn4ds9u6jmmg4n86d9wct7wmj23pyqe6p7e252lffzqsgcvqxm5lc2",
					Amount:           coin.MustNewCoinFromString("10"),
				},
			)))

			Expect(1).To(Equal(1))
		})
	})
})
