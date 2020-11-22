package event

import (
	"bytes"

	"github.com/crypto-com/chainindex/usecase/coin"
	"github.com/crypto-com/chainindex/usecase/model"

	entity_event "github.com/crypto-com/chainindex/entity/event"
	jsoniter "github.com/json-iterator/go"
	"github.com/luci/go-render/render"
)

const MSG_CREATE_VALIDATOR = "MsgCreateValidator"
const MSG_CREATE_VALIDATOR_CREATED = "MsgCreateValidatorCreated"
const MSG_CREATE_VALIDATOR_FAILED = "MsgCreateValidatorFailed"

type MsgCreateValidator struct {
	MsgBase

	CommissionRates  model.CommissionRates `json:"commission"`
	DelegatorAddress string                `json:"delegatorAddress"`
	ValidatorAddress string                `json:"validatorAddress"`
	PubKey           string                `json:"pubKey"`
	Amount           coin.Coin             `json:"amount"`
}

func NewMsgCreateValidator(msgCommonParams MsgCommonParams, params model.MsgCreateValidatorParams) *MsgCreateValidator {
	return &MsgCreateValidator{
		NewMsgBase(MsgBaseParams{
			MsgName:         MSG_DELEGATE,
			Version:         1,
			MsgCommonParams: msgCommonParams,
		}),

		params.CommissionRates,
		params.DelegatorAddress,
		params.ValidatorAddress,
		params.PubKey,
		params.Amount,
	}
}

// ToJSON encodes the event into JSON string payload
func (event *MsgCreateValidator) ToJSON() (string, error) {
	encoded, err := jsoniter.Marshal(event)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func (event *MsgCreateValidator) String() string {
	return render.Render(event)
}

func DecodeMsgCreateValidator(encoded []byte) (entity_event.Event, error) {
	jsonDecoder := jsoniter.NewDecoder(bytes.NewReader(encoded))
	jsonDecoder.DisallowUnknownFields()

	var event *MsgCreateValidator
	if err := jsonDecoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}