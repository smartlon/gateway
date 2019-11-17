package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	. "github.com/smartlon/gateway/adapter/iota/sdk"
)

type LogisticsController struct {
	beego.Controller
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
	root,err := MAMTransmit(message)
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
	messages := MAMReceive(sideKey,root)
	lc.Data["json"] = map[string]interface{}{"mammessage": messages, "msg": err}
	lc.ServeJSON()
}

func (lc *LogisticsController) GetNodeInfo(){
	nodeInfoBytes,err := NodeInfo()
	var code,message,ret string
	if err != nil {
		code = "201"
		message = "failed to get node info"
		ret = err.Error()
	}else {
		code = "200"
		message = "successed to get node info"
		ret = string(nodeInfoBytes)
	}

	lc.Data["json"] = map[string]interface{}{"code": code,"message": message, "result": ret}
	lc.ServeJSON()
}
