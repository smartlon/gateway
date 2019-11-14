package controller

import (
	"../fabric/sdk"
	"encoding/json"
	"fmt"
	"../log"
)

const (
	CHAINCODEID = "log"
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
	var requestLogisticReq sdk.Args
	err := json.Unmarshal(requestLogisticReqBytes,&requestLogisticReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	var argsArray []sdk.Args
	argsArray = append(argsArray, requestLogisticReq)
	var ret string
	ret, err = sdk.ChaincodeInvoke(CHAINCODEID, argsArray)
	if err != nil {
		log.Error(err)
	}
	lc.Data["json"] = map[string]interface{}{"code": 200,"message": err, "result": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) TransitLogistics(){

}
func (lc *LogisticsController) DeliveryLogistics(){

}
func (lc *LogisticsController) QueryLogistics(){

}
