package main

import (
	"fmt"
	"time"

	"github.com/crypto-com/chainindex/appinterface/projection/view"
	"github.com/crypto-com/chainindex/appinterface/rdb"
	"github.com/crypto-com/chainindex/infrastructure/tendermint"
	"github.com/crypto-com/chainindex/infrastructure/util"
	applogger "github.com/crypto-com/chainindex/internal/logger"
)

const INFO_DEFAULT_POLLING_INTERVAL = 5 * time.Second

type InfoManager struct {
	rdbConn         rdb.Conn
	client          *tendermint.HTTPClient
	logger          applogger.Logger
	pollingInterval time.Duration
	viewStatus      *view.Status
}

func NewInfoManager(
	logger applogger.Logger,
	rdbConn rdb.Conn,
	tendermintRPCUrl string,
) *InfoManager {
	tendermintClient := tendermint.NewHTTPClient(tendermintRPCUrl)

	viewStatus := view.NewStatus(rdbConn.ToHandle())
	return &InfoManager{
		rdbConn:    rdbConn,
		client:     tendermintClient,
		viewStatus: viewStatus,
	}

}

func (manager *InfoManager) Run() error {
	go func() {
		for true {
			status, _ := manager.client.Status()
			result := (*status)["result"]
			syncInfo := result.(map[string]interface{})["sync_info"]
			latestHeight := syncInfo.(map[string]interface{})["latest_block_height"].(string)
			util.WriteLog("a.txt", fmt.Sprintf("index service run %s\n", latestHeight))
			manager.viewStatus.Insert("LatestHeight", latestHeight)

			found, _ := manager.viewStatus.FindBy("LatestHeight")
			fmt.Printf("################################")
			fmt.Printf("found value %s", found)
			util.WriteLog("a.txt", fmt.Sprintf("found value2 %s\n", found))

			time.Sleep(time.Millisecond * 1000)
		}
	}()
	return nil
}
