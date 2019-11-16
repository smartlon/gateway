package rooter

import (
	"github.com/astaxie/beego"
	controller "github.com/smartlon/gateway/adapter/controller"
	)

func init() {
	beego.Router("/iota/nodeinfo", &controller.LogisticsController{},"get:GetNodeInfo")
	beego.Router("/iota/createmam", &controller.LogisticsController{},"post:CreateMAM")
	beego.Router("/iota/blockmam", &controller.LogisticsController{},"post:BlockMAM")
	beego.Router("/iota/mamtransmit", &controller.LogisticsController{},"post:MAMTransmit")
	beego.Router("/iota/mamreceive", &controller.LogisticsController{},"post:MAMReceive")
	beego.Router("/faric/requestlogistic", &controller.LogisticsController{},"post:RequestLogistic")
	beego.Router("/faric/transitlogistics", &controller.LogisticsController{},"post:TransitLogistics")
	beego.Router("/faric/deliverylogistics", &controller.LogisticsController{},"post:DeliveryLogistics")
	beego.Router("/faric/querylogistics", &controller.LogisticsController{},"get:QueryLogistics")
	beego.Router("/faric/registeruser", &controller.LogisticsController{},"post:RegisterUser")

}