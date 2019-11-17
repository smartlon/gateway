package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/mam/v1"
	"log"
	"sync"
)

type IotaClient struct {
	Api *mam.API
	Receiver *mam.Receiver
	Transmitter *mam.Transmitter
}

var iotaClient IotaClient
var once *sync.Once

func init(){
	account, iotaApi,err := NewAccount()
	if err != nil {
		log.Fatal(err)
	}
	account.Start()
	defer account.Shutdown()
	receiver := mam.NewReceiver(iotaApi)
	transmitter := mam.NewTransmitter(iotaApi,GetInMemorySeed(SEED),MWM,consts.SecurityLevelMedium)
	if err != nil {
		fmt.Println(err.Error())
	}
	once = &sync.Once{}
	once.Do(func() {
		iotaClient = IotaClient{
			&iotaApi,
			receiver,
			transmitter,
		}
	})

}

func  CreateMAM(messageBytes []byte, sideKey string)(string, error){
	message := string(messageBytes)
	return mamSend(message,sideKey)
}
func  BlockMAM(messageBytes []byte, sideKey string)(string, error){
	message := string(messageBytes)
	return mamSend(message,sideKey)
}

func MAMTransmit(message, sideKey string) (string, error){
	return mamSend(message,sideKey)
}

func mamSend( message, sideKey string) (string, error ){
	err :=  iotaClient.Transmitter.SetMode(mam.ChannelModeRestricted,sideKey)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	root, err := iotaClient.Transmitter.Transmit(message)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return root, nil
}

func  MAMReceive(sideKey, root string) ([]string, error ){
	err := iotaClient.Receiver.SetMode(mam.ChannelModeRestricted,sideKey)
	if err != nil {
		fmt.Println(err.Error())
		return []string{},err
	}
	messages := make([]string,0)
	message := make([]string,0)
	for root != "" {
		root,message, err = iotaClient.Receiver.Receive(root)
		if err != nil {
			fmt.Println(err.Error())
			return []string{},err
		}
		for _,str := range message {
			messages = append(messages,str)
		}
	}
	return messages,nil
}


func NodeInfo() ([]byte, error) {
	nodeInfo, err := iotaClient.Api.GetNodeInfo()
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	nodeInfoBytes,err := json.Marshal(nodeInfo)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	return nodeInfoBytes,nil
}
