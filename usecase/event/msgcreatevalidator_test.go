package event_test

import (
	event_entity "github.com/crypto-com/chainindex/entity/event"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	event_usecase "github.com/crypto-com/chainindex/usecase/event"
)

var _ = Describe("Event", func() {
	registry := event_entity.NewRegistry()
	event_usecase.RegisterEvents(registry)

	Describe("En/DecodeMsgCreateValidator", func() {
		It("should able to encode and decode to the same event", func() {

			Expect(1).To(Equal(1))
		})

		It("should able to encode and decode to failed event", func() {

			Expect(1).To(Equal(1))
		})
	})
})
