package iota

import (
	"encoding/json"
	"fmt"
	"github.com/pebbe/zmq4"
	"github.com/smartlon/gateway/adapter/ports"
	"github.com/smartlon/gateway/adapter/ports/iota/sdk"
	"github.com/smartlon/gateway/log"
	"github.com/smartlon/gateway/types"
	"strings"
)

func init() {
	builder := func(config ports.AdapterConfig) (ports.AdapterService, error) {
		a := &IOTAAdaptor{config: &config}
		a.Start()
		a.Sync()
		return a, nil
	}
	ports.GetPortsIncetance().RegisterBuilder("iota", builder)
}

type IOTAAdaptor struct {
	config      *ports.AdapterConfig
}


func (a *IOTAAdaptor) Start() error {
	return nil
}


func (a *IOTAAdaptor) Sync() error {
	return nil
}


func (a *IOTAAdaptor) Stop() error {
	return nil
}

// Subscribe events from fabric chain
func (a *IOTAAdaptor) Subscribe(listener ports.EventsListener) {
	log.Infof("event subscribe: %s", ports.GetAdapterKey(a))
	dealt := a.GetPort() - 14625
	zmqPort := dealt * 1000 + 5556
	endpoint := ports.GetAdapterKey(a)
	zmqAddress := ports.GenEndpoint("tcp",a.GetIP(),zmqPort)
	go func() {
		socket, err := zmq4.NewSocket(zmq4.SUB)
		must(err)
		socket.SetSubscribe("tx")
		err = socket.Connect(zmqAddress)
		must(err)

		fmt.Printf("started tx feed\n")
		for {
			msg, err := socket.Recv(0)
			must(err)

			tx := sdk.BuildTxFromZMQData(msg)
			//fmt.Printf("received tx: %s\n",tx)
			if tx == nil {
				fmt.Printf("tx: receive error! message format error\n")
				continue
			} else if tx.Type == "tx_trytes" {
				//fmt.Printf("tx: trytes received. Skip.\n")
				continue
			}
			if strings.Contains(tx.Tag, sdk.ChainTag)  {
				log.Infof("received tx: %s\n",tx)
				gatewayEventDataTx:= types.GatewayEventDataTx{
					a.GetChainName(),
					sdk.ChainTag,
					tx.Tag[len(sdk.ChainTag):],
					tx.Address,
					types.IotaPayload{},
				}
				event := &types.Event{
					endpoint,
					gatewayEventDataTx,
				}
				listener(event,a)
			}

		}
	}()
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (a *IOTAAdaptor) SubmitTx(tx string) (string,error) {
	if tx == "" {
		return "",fmt.Errorf("submit tx with nil param")
	}
	transmitTxInfo := sdk.TransmitTxInfo{}
	err := json.Unmarshal([]byte(tx),&transmitTxInfo)
	if err != nil {
		return "",err
	}
	transmitTxInfo.Endpoint = ports.GenEndpoint("http",a.GetIP(),a.GetPort())
	mamstate,root,address := sdk.MAMTransmit(transmitTxInfo.Message,transmitTxInfo.Mamstate,transmitTxInfo.Seed,transmitTxInfo.Mode,transmitTxInfo.SideKey,transmitTxInfo.TransactionTag,transmitTxInfo.Endpoint)
	transmitResult := sdk.TransmitResult{
		root,
		mamstate,
		address,
	}
	resultByte,err := json.Marshal(transmitResult)
	if err != nil {
		return "",err
	}
	return string(resultByte),nil
}



func (a *IOTAAdaptor) ObtainTx(tx string) (string, error) {
	if tx == "" {
		return "",fmt.Errorf("submit tx with nil param")
	}
	recieveTxInfo := sdk.RecieveTxInfo{}
	err := json.Unmarshal([]byte(tx),&recieveTxInfo)
	if err != nil {
		return "",err
	}
	recieveTxInfo.Endpoint = ports.GenEndpoint("http",a.GetIP(),a.GetPort())
	iotdatas := sdk.MAMReceive(recieveTxInfo.Root,recieveTxInfo.Mode,recieveTxInfo.SideKey,recieveTxInfo.Endpoint)
	var temperature string
	for _,iotdata := range iotdatas {
		var iot sdk.IoTData
		err := json.Unmarshal([]byte(iotdata),&iot)
		if err != nil {
			return "",err
		}
		if len(iot.Temperature) != 0 {
			temperature = temperature + iot.Temperature + ","
		}
	}
	temperature = strings.TrimSuffix(temperature,",")
	return temperature, nil
}


// Count Calculate the total and consensus number for chain
func (a *IOTAAdaptor) Count() (totalNumber int, consensusNumber int) {
	totalNumber = ports.GetPortsIncetance().Count(a.GetChainName())
	consensusNumber = ports.Consensus2of3(totalNumber)
	return
}

// GetChainName returns chain name
func (a *IOTAAdaptor) GetChainName() string {
	return a.config.ChainName
}

// GetIP returns chain node ip
func (a *IOTAAdaptor) GetIP() string {
	return a.config.IP
}

// GetPort returns chain node port
func (a *IOTAAdaptor) GetPort() int {
	return a.config.Port
}