package event_test

import (
	event_entity "github.com/crypto-com/chainindex/entity/event"
	"github.com/crypto-com/chainindex/test/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/usecase/coin"
	event_usecase "github.com/crypto-com/chainindex/usecase/event"
	"github.com/crypto-com/chainindex/usecase/model"
)

var _ = Describe("Event", func() {
	registry := event_entity.NewRegistry()
	event_usecase.RegisterEvents(registry)

	Describe("En/DecodeTransactionFailed", func() {
		It("should able to encode and decode to the same Event", func() {
			anyTxHash := factory.RandomTxHash()
			anyHeight := int64(1000)
			anyParams := model.CreateTransactionParams{
				TxHash:    anyTxHash,
				Code:      0,
				Log:       "{\"events\":[]}",
				MsgCount:  1,
				Fee:       coin.MustNewCoinFromString("1000"),
				GasWanted: "200000",
				GasUsed:   "10000",
			}
			event := event_usecase.NewTransactionFailed(anyHeight, anyParams)

			encoded, err := event.ToJSON()
			Expect(err).To(BeNil())

			decodedEvent, err := registry.DecodeByType(
				event_usecase.TRANSACTION_FAILED, 1, []byte(encoded),
			)
			Expect(err).To(BeNil())
			Expect(decodedEvent).To(Equal(event))
			typedEvent, _ := decodedEvent.(*event_usecase.TransactionFailed)
			Expect(typedEvent.TxHash).To(Equal(anyTxHash))
		})
	})
})
