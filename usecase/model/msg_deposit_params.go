package model

import "github.com/crypto-com/chain-indexing/usecase/coin"

type MsgDepositParams struct {
	ProposalId string    `json:"proposalId"`
	Depositor  string    `json:"depositor"`
	Amount     coin.Coin `json:"amount"`
}
