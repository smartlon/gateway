package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/smartlon/gateway/adapter/ports/fabric/sdk"
	"github.com/smartlon/gateway/adapter/ports"
	"github.com/smartlon/gateway/log"
)

//func init() {
//	builder := func(config ports.AdapterConfig) (ports.AdapterService, error) {
//		a := &FabAdaptor{config: &config}
//		a.Start()
//		a.Sync()
//		//a.Subscribe(config.Listener)
//		return a, nil
//	}
//	ports.GetPortsIncetance().RegisterBuilder("fabric", builder)
//}

const (

	ChainType = "fabric"
	// ChainName name of hyperledger fabric
	ChainName = "supplychain"
	// ChaincodeID id of chaincode
	ChaincodeID = "supcc"
)

// ChainResult result of hypelrdger fabric chaincode
type ChainResult struct {
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	ErrString string `json:"error,omitempty"`
	// Result    interface{} `json:"result,omitempty"`
}

// FabAdaptor provides adapter for hyperledger fabric
type FabAdaptor struct {
	config      *ports.AdapterConfig
}

// Start fabric adapter service
func (a *FabAdaptor) Start() error {
	return nil
}

// Sync status for fabric adapter service
func (a *FabAdaptor) Sync() error {

	return nil
}

// Stop fabric adapter service
func (a *FabAdaptor) Stop() error {
	return nil
}

// Subscribe events from fabric chain
func (a *FabAdaptor) Subscribe(listener ports.EventsListener) {
	log.Infof("event subscribe: %s", ports.GetAdapterKey(a))
	peerUrl := fmt.Sprintf("%s:%d",a.GetIP(),a.GetPort())
	go sdk.Listener(ChaincodeID,peerUrl,listener,a)
}

// SubmitTx submit Tx to hyperledger fabric chain
func (a *FabAdaptor) SubmitTx(tx string) (string, error) {
	var argsArray []sdk.Args
	err := json.Unmarshal([]byte(tx),&argsArray)
	if err != nil {
		return "",err
	}
	return sdk.ChaincodeInvoke(ChaincodeID,argsArray)
}

func (a *FabAdaptor) ObtainTx(tx string) (string, error) {
	return "",nil
}


// Count Calculate the total and consensus number for chain
func (a *FabAdaptor) Count() (totalNumber int, consensusNumber int) {
	totalNumber = ports.GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = ports.Consensus2of3(totalNumber)
	return
}

// GetChainName returns chain name
func (a *FabAdaptor) GetChainName() string {
	return a.config.ChainName
}

// GetIP returns chain node ip
func (a *FabAdaptor) GetIP() string {
	return a.config.IP
}

// GetPort returns chain node port
func (a *FabAdaptor) GetPort() int {
	return a.config.Port
}
