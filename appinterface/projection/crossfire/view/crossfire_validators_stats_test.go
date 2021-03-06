package view_test

import (
	"github.com/crypto-com/chain-indexing/appinterface/projection/crossfire/constants"
	"github.com/crypto-com/chain-indexing/appinterface/projection/crossfire/view"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chain-indexing/infrastructure/pg"
	. "github.com/crypto-com/chain-indexing/test"
)

var _ = Describe("Crossfire Validators Stats", func() {
	WithTestPgxConn(func(conn *pg.PgxConn, migrate *pg.Migrate) {
		BeforeEach(func() {
			_ = migrate.Reset()
			migrate.MustUp()
		})

		AfterEach(func() {
			_ = migrate.Reset()
		})

		Describe("Increment", func() {
			It("should run increment correctly", func() {
				var err error

				operatorAddress := "tcrocncl14m5a4kxt2e82uqqs5gtqza29dm5wqzyalddug5"
				key := constants.ValidatorCommitmentKey(operatorAddress, constants.PHASE_1N2_COMMIT_PREFIX)
				crossfireValidatorsStatsView := view.NewCrossfireValidatorsStats(conn.ToHandle())

				// increment
				err = crossfireValidatorsStatsView.IncrementOne(key)
				Expect(err).To(BeNil())
				// should be 1
				actual, err := crossfireValidatorsStatsView.FindBy(key)
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(int64(1)))

				// increment
				err = crossfireValidatorsStatsView.IncrementOne(key)
				Expect(err).To(BeNil())
				// should be 2
				actual, err = crossfireValidatorsStatsView.FindBy(key)
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(int64(2)))
			})
		})
	})
})
