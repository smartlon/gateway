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

type MAMClient struct {
	Api *mam.API
	Receiver *mam.Receiver
	Transmitter *mam.Transmitter
}

var mamClient MAMClient
var iotaApi *api.API
var once *sync.Once

func init(){
	account, iotaMAMApi,err := NewAccount()
	if err != nil {
		log.Fatal(err)
	}
	account.Start()
	defer account.Shutdown()
	receiver := mam.NewReceiver(iotaMAMApi)
	transmitter := mam.NewTransmitter(iotaMAMApi,GetInMemorySeed(SEED),MWM,consts.SecurityLevelMedium)
	if err != nil {
		fmt.Println(err.Error())
	}
	once = &sync.Once{}
	once.Do(func() {
		iotaApi,err = NewConnection()
		if err != nil {
			fmt.Println(err.Error())
		}
		mamClient = MAMClient{
			&iotaMAMApi,
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
	err :=  mamClient.Transmitter.SetMode(mam.ChannelModeRestricted,sideKey)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	root, err := mamClient.Transmitter.Transmit(message)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return root, nil
}

func  MAMReceive(sideKey, root string) ([]string, error ){
	err := mamClient.Receiver.SetMode(mam.ChannelModeRestricted,sideKey)
	if err != nil {
		fmt.Println(err.Error())
		return []string{},err
	}
	messages := make([]string,0)
	message := make([]string,0)
	for root != "" {
		root,message, err = mamClient.Receiver.Receive(root)
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
	nodeInfo, err := iotaApi.GetNodeInfo()
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
