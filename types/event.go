package types

type IotaPayload struct {
	ContainerID string `json:"ContainerID"`
	Seed        string `json:"Seed"`
	MamState    string `json:"MamState"`
	Root        string `json:"Root"`
	Mode       	string `json:"Mode"`
	SideKey     string `json:"SideKey"`
}


type GatewayEventDataTx struct {
	From      string `json:"from"` //qsc name 或 qos
	To        string `json:"to"`   //qsc name 或 qos
	Func    string  `json:"func"`
	Address  string  `json:"address"`
	IotaPayload IotaPayload `json:"iotaPayload"` //TxQcp 做 sha256
}

// Event cache tx event tags and node info
type Event struct {
	NodeAddress        string               `json:"node"` //event 源地址
	GatewayEventDataTx `json:"eventDataTx"` //event 事件
}



