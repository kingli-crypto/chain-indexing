package handlers

import (
	"github.com/crypto-com/chainindex/appinterface/projection/view"
	"github.com/crypto-com/chainindex/appinterface/rdb"
	"github.com/crypto-com/chainindex/infrastructure/httpapi"
	applogger "github.com/crypto-com/chainindex/internal/logger"
	"github.com/valyala/fasthttp"
)

type InfoHandler struct {
	logger applogger.Logger

	statusView *view.Status
}

func NewInfoHandler(logger applogger.Logger, rdbHandle *rdb.Handle) *InfoHandler {
	return &InfoHandler{
		logger.WithFields(applogger.LogFields{
			"module": "InfoHandler",
		}),

		view.NewStatus(rdbHandle),
	}
}

func (handler *InfoHandler) GetLatestHeight(ctx *fasthttp.RequestCtx) {
	latestHeight, err := handler.statusView.FindBy("LatestHeight")

	if err != nil {
		handler.logger.Errorf("error fetching latest height: %v", err)
		httpapi.InternalServerError(ctx)
		return
	}

	status := GetLatestHeightStatus{
		LatestHeight: latestHeight,
	}

	httpapi.Success(ctx, status)
}

type GetLatestHeightStatus struct {
	LatestHeight string `json:"latestHeight"`
}
