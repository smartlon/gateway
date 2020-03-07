package ports

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/pkg/errors"
	"github.com/smartlon/gateway/adapter/ports/fabric/sdk"
	"github.com/smartlon/gateway/log"
	"github.com/smartlon/gateway/types"
)

func init() {
	builder := func(config AdapterConfig) (AdapterService, error) {
		a := &FabAdaptor{config: &config}
		a.Start()
		a.Sync()
		//a.Subscribe(config.Listener)
		return a, nil
	}
	GetPortsIncetance().RegisterBuilder("fabric", builder)
}

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
	config      *AdapterConfig
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
func (a *FabAdaptor) Subscribe(listener EventsListener) {
	log.Infof("event subscribe: %s", GetAdapterKey(a))
	peerUrl := fmt.Sprintf("%s:%d",a.GetIP(),a.GetPort())
	action,err := sdk.NewChaincodeEventAction()
	if err != nil {
		log.Error(err)
	}
	action.Set(sdk.Config().ChannelID,ChaincodeID,[]sdk.Args{})

	ec, err := action.EventClient(peerUrl,event.WithBlockEvents())
	if err != nil {
		fmt.Println("failed to create client")
		panic(err)
	}

	_, notifierCreateChannel, err := ec.RegisterChaincodeEvent(ChaincodeID, `{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: CreateChannel")
		panic(err)

	}
	_, notifierDeliveryLogistics, err := ec.RegisterChaincodeEvent(ChaincodeID, `{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: DeliveryLogistics")
		panic(err)
	}
	var iotapayload types.IotaPayload
	go func () {
		for {
			event :=&types.Event{}
			select {
			case ccEvent,ok := <-notifierCreateChannel:
				if !ok {
					panic(errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event"))
				}
				err := json.Unmarshal(ccEvent.Payload,&iotapayload)
				if err != nil {
					fmt.Println(err.Error())
				}
				event.Func = "FABRICCREATE"
				fmt.Printf("received chaincode event CreateChannel:  %v\n", iotapayload)


			case ccEvent,ok := <-notifierDeliveryLogistics:
				if !ok {
					panic(errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event"))
				}
				err := json.Unmarshal(ccEvent.Payload,&iotapayload)
				if err != nil {
					fmt.Println(err.Error())
				}
				event.Func = "FABRICDELEVERY"
				fmt.Printf("received chaincode event DeliveryLogistics:  %v\n", iotapayload)

			//case <-time.After(time.Second * 3):
			//	fmt.Println("Exit while waiting for chaincode event")
			}
			event.IotaPayload=iotapayload
			event.From = a.GetChainName()
			event.To = "supply-iota"
			event.NodeAddress = GetAdapterKey(a)
			listener(event,a)
		}
	}()
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
	totalNumber = GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = Consensus2of3(totalNumber)
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
