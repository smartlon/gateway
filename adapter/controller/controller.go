package controller

import (
	. "github.com/iotaledger/iota.go/api"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/mam/v1"
	. "github.com/smartlon/gateway/adapter/iota/sdk"
	"sync"
)

type LogisticsController struct {
	beego.Controller
}

type IotaClient struct {
	Api *API
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
	err = transmitter.SetMode(mam.ChannelModeRestricted,SIDEKEYPRIVATE)
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

func (lc *LogisticsController) CreateMAM(){

}
func (lc *LogisticsController) BlockMAM(){

}

type MAMTransmitReq struct {
	Message string `json:"message"`
}

func (lc *LogisticsController) MAMTransmit(){
	mamReqBytes := lc.Ctx.Input.RequestBody
	var mamReq MAMTransmitReq
	err := json.Unmarshal(mamReqBytes,&mamReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	message := mamReq.Message
	root, err := iotaClient.Transmitter.Transmit(message)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("root: "+ root)
	lc.Data["json"] = map[string]interface{}{"root": root,"sidekey": SIDEKEYPRIVATE, "msg": err}
	lc.ServeJSON()

}
type MAMReceiveReq struct {
	Root string `json:"root"`
	SideKey string `json:"sidekey"`
}

func (lc *LogisticsController) MAMReceive(){
	mamReqBytes := lc.Ctx.Input.RequestBody
	var mamReq MAMReceiveReq
	err := json.Unmarshal(mamReqBytes,&mamReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(mamReq)
	root := mamReq.Root
	sideKey := mamReq.SideKey
	err = iotaClient.Receiver.SetMode(mam.ChannelModeRestricted,sideKey)
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
	lc.Data["json"] = map[string]interface{}{"mammessage": messages, "msg": err}
	lc.ServeJSON()
}

func (lc *LogisticsController) GetNodeInfo(){
	nodeInfo, err := iotaClient.Api.GetNodeInfo()
	if err != nil {
		fmt.Println(err.Error())
	}
	nodeInfoBytes,err := json.Marshal(nodeInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	lc.Ctx.Output.Body(nodeInfoBytes)
}
