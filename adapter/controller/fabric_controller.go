package controller

import (
	"github.com/smartlon/gateway/adapter/fabric/sdk"
	"encoding/json"
	"fmt"
	"github.com/smartlon/gateway/adapter/log"
)

const (
	CHAINCODEID = "logistic"

)

type UserReq struct {
	UserName string `json:"UserName"`
	PassWord string `json:"PassWord"`
}

func (lc *LogisticsController) RegisterUser(){
	registerUserReqBytes := lc.Ctx.Input.RequestBody
	var registerUserReq UserReq
	err := json.Unmarshal(registerUserReqBytes,&registerUserReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	var ret string
	ret, err = sdk.RegisterUser(registerUserReq.UserName, registerUserReq.PassWord)
	message := "register user success"
	if err != nil {
		log.Error(err)
		message = "register user fail"
	}
	lc.Data["json"] = map[string]interface{}{"code": 200,"message": message, "result": ret}
	lc.ServeJSON()
}

func (lc *LogisticsController) RequestLogistic(){
	requestLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(requestLogisticReqBytes)
        
	lc.Data["json"] = map[string]interface{}{"code": code,"message": message, "result": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) TransitLogistics(){
	requestLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(requestLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"message": message, "result": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) DeliveryLogistics(){
	requestLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(requestLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"message": message, "result": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) QueryLogistics(){
	requestLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(requestLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"message": message, "result": ret}
	lc.ServeJSON()
}
func invokeController(requestLogisticReqBytes []byte)(code, message, ret string){
	var requestLogisticReq sdk.Args
	err := json.Unmarshal(requestLogisticReqBytes,&requestLogisticReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	var argsArray []sdk.Args
	argsArray = append(argsArray, requestLogisticReq)

	ret, err = sdk.ChaincodeInvoke(CHAINCODEID, argsArray)
	if err != nil {
		log.Error(err)
		message = err.Error()
		code = "201"
	}else {
		message = "request logistics success"
		code = "200"
	}
	return
}
