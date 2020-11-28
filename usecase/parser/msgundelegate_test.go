package parser_test

import (
	"github.com/crypto-com/chain-indexing/internal/utctime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chain-indexing/entity/command"
	"github.com/crypto-com/chain-indexing/usecase/coin"
	command_usecase "github.com/crypto-com/chain-indexing/usecase/command"
	"github.com/crypto-com/chain-indexing/usecase/event"
	"github.com/crypto-com/chain-indexing/usecase/model"
	"github.com/crypto-com/chain-indexing/usecase/parser"
	usecase_parser_test "github.com/crypto-com/chain-indexing/usecase/parser/test"
)

var _ = Describe("ParseMsgCommands", func() {
	Describe("MsgDelegate", func() {

		It("should parse Msg commands when there is staking.MsgUndelegate in the transaction", func() {
			txDecoder := parser.NewTxDecoder("basetrcro")
			block, _ := mustParseBlockResp(usecase_parser_test.TX_MSG_UNDELEGATE_BLOCK_RESP)
			blockResults := mustParseBlockResultsResp(usecase_parser_test.TX_MSG_UNDELEGATE_BLOCK_RESULTS_RESP)

			cmds, err := parser.ParseBlockResultsTxsMsgToCommands(
				txDecoder,
				block,
				blockResults,
			)
			Expect(err).To(BeNil())
			Expect(cmds).To(HaveLen(1))
			Expect(cmds).To(Equal([]command.Command{command_usecase.NewCreateMsgUndelegate(
				event.MsgCommonParams{
					BlockHeight: int64(374371),
					TxHash:      "0F525EFC1DD9C319E9036C35CF1656E09480B308301BB3A46F850AE482A3875C",
					TxSuccess:   true,
					MsgIndex:    0,
				},
				model.MsgUndelegateParams{
					DelegatorAddress: "tcro1gs80n8fpc5mc3ywkgfy93l23tg0gdqj5w2ll64",
					ValidatorAddress: "tcrocncl1j7pej8kplem4wt50p4hfvndhuw5jprxxxtenvr",
					Amount:           coin.MustNewCoinFromString("1000000000"),
					UnbondCompleteAt: utctime.FromUnixNano(int64(1605152654000000000)),
				},
			)}))
		})
	})
})
