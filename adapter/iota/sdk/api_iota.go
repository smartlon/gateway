package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/mam/v1"
	"sync"
	"github.com/iotaledger/iota.go/api"
)

type IotaClient struct {
	Api *api.API
	Receiver *mam.Receiver
	Transmitter *mam.Transmitter
}

var iotaClient IotaClient
var once *sync.Once

func init(){
	api,err := NewConnection()
	receiver := mam.NewReceiver(api)
	transmitter := mam.NewTransmitter(api,GetInMemorySeed(SEED),MWM,consts.SecurityLevelMedium)
	if err != nil {
		fmt.Println(err.Error())
	}
	once = &sync.Once{}
	once.Do(func() {
		iotaClient = IotaClient{
			api,
			receiver,
			transmitter,
		}
	})

}

func  CreateMAM(message []byte){
	//root, err := iotaClient.Transmitter.Transmit(message)
}
func  BlockMAM(message []byte){

}

func MAMTransmit(message string) (string, error){
	root, err := iotaClient.Transmitter.Transmit(message)
	err =  iotaClient.Transmitter.SetMode(mam.ChannelModeRestricted,SIDEKEYPRIVATE)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return root, nil
}

func  MAMReceive(sideKey, root string) []string {
	err := iotaClient.Receiver.SetMode(mam.ChannelModeRestricted,sideKey)
	if err != nil {
		fmt.Println(err.Error())
	}
	messages := make([]string,0)
	message := make([]string,0)
	for root != "" {
		root,message, err = iotaClient.Receiver.Receive(root)
		if err != nil {
			fmt.Println(err.Error())
		}
		for _,str := range message {
			messages = append(messages,str)
		}
	}
	return messages
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