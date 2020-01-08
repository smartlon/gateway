package controller

import (
	"crypto/rand"
	"github.com/smartlon/gateway/adapter/fabric/sdk"
	"encoding/json"
	"fmt"
	"github.com/smartlon/gateway/adapter/log"
	"math/big"
	"strconv"
	"time"
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
        
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) TransitLogistics(){
	transitLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(transitLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) DeliveryLogistics(){
	deliveryLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(deliveryLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) QueryLogistics(){
	queryLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}

func (lc *LogisticsController) RecordContainer(){
	recordContainerReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(recordContainerReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) QueryContainer(){
	queryContainerReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryContainerReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
type QueryResponse struct {
	Key    string `json:"Key"`
	Record    string `json:"Record"`
}

func (lc *LogisticsController) QueryAllContainers(){
	queryAllContainersReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllContainersReqBytes)
	var qr []QueryResponse
	err := json.Unmarshal([]byte(ret),&qr)
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(qr)
	var resp []string
	for _,v := range qr {
		resp = append(resp,v.Record)
	}
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": resp}
	lc.ServeJSON()
}
func (lc *LogisticsController) QueryAllLogistics(){
	queryAllLogisticsReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllLogisticsReqBytes)
	var qr []QueryResponse
	err := json.Unmarshal([]byte(ret),&qr)
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(qr)
	var resp []string
	for _,v := range qr {
		resp = append(resp,v.Record)
	}
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": resp}
	lc.ServeJSON()
}


func invokeController(invokeReqBytes []byte)(code int, message, ret string){
	var invokeReq sdk.Args
	err := json.Unmarshal(invokeReqBytes,&invokeReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	if invokeReq.Func == "RecordContainer" {
		timestamp := timeStamp()
		seed := generateRandomSeedString(81)
		sidekey := generateRandomSeedString(81)
		argscomposite := []string{timestamp,seed,sidekey}
		invokeReq.Args = append(invokeReq.Args,argscomposite...)
	}
	if invokeReq.Func == "TransitLogistics"  {
		timestamp := timeStamp()
		sidekey := generateRandomSeedString(81)
		argscomposite := []string{sidekey,timestamp}
		invokeReq.Args = append(invokeReq.Args,argscomposite...)
	}
	if  invokeReq.Func == "DeliveryLogistics" {
		timestamp := timeStamp()
		invokeReq.Args = append(invokeReq.Args,timestamp)
	}
	var argsArray []sdk.Args
	argsArray = append(argsArray, invokeReq)

	ret, err = sdk.ChaincodeInvoke(CHAINCODEID, argsArray)
	if err != nil {
		log.Error(err)
		message = err.Error()
		code = 201
	}else {
		message = "invoke " +invokeReq.Func+ " success"
		code = 200
	}
	return
}


func timeStamp() string {
	return strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
}

func generateRandomSeedString(length int) string {
	seed := ""
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ9"

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(27))
		if err != nil {
			fmt.Println(err)
		}
		seed += string(alphabet[n.Int64()])
	}
	return seed
}
