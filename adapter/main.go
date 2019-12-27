package main

import (
	"github.com/astaxie/beego"
	_ "github.com/smartlon/gateway/adapter/rooter"
)

type Action struct {
	From  string  `json:"From"`
	To  string  `json:"To"`
	Func  string  `json:"Func"`
}

func main() {
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.Listen.HTTPPort = 8081
	//go sdk.StartTxFeed()
	beego.Run()

	//var action *Action
	//var args *sdk.Args
	//if action == nil {
	//	action = &Action{
	//		"Fabric",
	//		"Iota",
	//		"Create",
	//	}
	//	args = &sdk.Args{
	//		"invoke",
	//		[]string{"RequestLogistic","product1","medical","seller1","xian","buyer1","beijing"},
	//	}
	//}
	//
	//actionBytes, _ := json.Marshal(action)
	//argsBytes, _ := json.Marshal(args)
	//fmt.Println(string(actionBytes))
	//fmt.Println(string(argsBytes))
}