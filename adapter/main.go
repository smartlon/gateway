package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
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
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		//AllowOrigins:      []string{"https://192.168.0.102"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

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