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
	beego.Router("/fabric/requestlogistic", &controller.LogisticsController{},"post:RequestLogistic")
	beego.Router("/fabric/transitlogistics", &controller.LogisticsController{},"post:TransitLogistics")
	beego.Router("/fabric/deliverylogistics", &controller.LogisticsController{},"post:DeliveryLogistics")
	beego.Router("/fabric/querylogistics", &controller.LogisticsController{},"post:QueryLogistics")
	beego.Router("/fabric/registeruser", &controller.LogisticsController{},"post:RegisterUser")

}
