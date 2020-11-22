package event

import (
	"bytes"

	"github.com/crypto-com/chainindex/usecase/coin"
	"github.com/crypto-com/chainindex/usecase/model"

	entity_event "github.com/crypto-com/chainindex/entity/event"
	jsoniter "github.com/json-iterator/go"
	"github.com/luci/go-render/render"
)

const MSG_EDIT_VALIDATOR = "MsgEditValidator"
const MSG_EDIT_VALIDATOR_CREATED = "MsgEditValidatorCreated"
const MSG_EDIT_VALIDATOR_FAILED = "MsgEditValidatorFailed"

type MsgEditValidator struct {
	MsgBase

	DelegatorAddress string    `json:"delegatorAddress"`
	ValidatorAddress string    `json:"validatorAddress"`
	Amount           coin.Coin `json:"amount"`
}

func NewMsgEditValidator(msgCommonParams MsgCommonParams, params model.MsgDelegateParams) *MsgEditValidator {
	return &MsgEditValidator{
		NewMsgBase(MsgBaseParams{
			MsgName:         MSG_DELEGATE,
			Version:         1,
			MsgCommonParams: msgCommonParams,
		}),

		params.DelegatorAddress,
		params.ValidatorAddress,
		params.Amount,
	}
}

// ToJSON encodes the event into JSON string payload
func (event *MsgEditValidator) ToJSON() (string, error) {
	encoded, err := jsoniter.Marshal(event)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func (event *MsgEditValidator) String() string {
	return render.Render(event)
}

func DecodeMsgEditValidator(encoded []byte) (entity_event.Event, error) {
	jsonDecoder := jsoniter.NewDecoder(bytes.NewReader(encoded))
	jsonDecoder.DisallowUnknownFields()

	var event *MsgEditValidator
	if err := jsonDecoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}
