package model

import (
	"github.com/crypto-com/chainindex/usecase/coin"
)

type MsgCreateValidatorParams struct {
	CommissionRates  CommissionRates `json:"commission"`
	DelegatorAddress string          `json:"delegatorAddress"`
	ValidatorAddress string          `json:"validatorAddress"`
	PubKey           string          `json:"pubKey"`
	Amount           coin.Coin       `json:"amount"`
}
