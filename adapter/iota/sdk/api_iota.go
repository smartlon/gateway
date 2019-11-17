package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/mam/v1"
	"github.com/smartlon/gateway/adapter/controller"
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
	_, iotaMAMApi,err := NewAccount()
	if err != nil {
		log.Fatal(err)
	}
	//account.Start()
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
func  BlockMAM(messageBytes []byte, root,sideKey string)(string,string, error){
	var err error
	err = mamClient.Receiver.SetMode(mam.ChannelModeRestricted,sideKey)
	if err != nil {
		fmt.Println(err.Error())
		return "","",err
	}
	var temperature string
	var iotData controller.IoTData
	var messages []string
	root,messages, err = mamClient.Receiver.Receive(root)
	for root != "" {
		root,messages, err = mamClient.Receiver.Receive(root)
		if err != nil {
			fmt.Println(err.Error())
			return "","",err
		}
		iotDataBytes := convertData(messages)
		json.Unmarshal(iotDataBytes,&iotData)
		temperature += iotData.Temperature + ","
	}
	temperature = temperature[:len(temperature)-1]
	message := string(messageBytes)
	root,err =mamSend(message,sideKey)
	return root,temperature,err
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

func  MAMReceive(sideKey, root string) (string, error ){
	var err error
	err = mamClient.Receiver.SetMode(mam.ChannelModeRestricted,sideKey)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}
	var iotDataAllBytes []byte

	var messages []string
	for root != "" {
		root,messages, err = mamClient.Receiver.Receive(root)
		if err != nil {
			fmt.Println(err.Error())
			return "",err
		}
		iotDataBytes := convertData(messages)
		iotDataAllBytes =append(iotDataAllBytes,iotDataBytes...)

	}
	return string(iotDataAllBytes),nil
}

func convertData(messages []string) []byte{
	var temp string
	for _,message := range messages {
		temp += message
	}
	return []byte(temp)
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
