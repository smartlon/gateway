package consensus

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smartlon/gateway/adapter/ports"
	"github.com/smartlon/gateway/config"
	"github.com/smartlon/gateway/log"
	"github.com/smartlon/gateway/queue"
	"github.com/smartlon/gateway/types"
	iotaSDK "github.com/smartlon/gateway/adapter/ports/iota/sdk"
	fabricSDK "github.com/smartlon/supplynetwork/fabric/sdk"
	"strconv"
	"time"
)

// StartQcpConsume Start to consume tx msg
func StartQcpConsume(conf *config.Config) (err error) {
	qsconfigs := conf.Qscs
	if len(qsconfigs) < 2 {
		return errors.New("config error , at least two chain targets ")
	}
	var subjects string
	es := make(chan error, 1024) //TODO 1024参数按需修改
	defer close(es)
	for i, qsconfig := range qsconfigs {
		for j := i + 1; j < len(qsconfigs); j++ {
			qcpConsume(qsconfigs[j].Name, qsconfig.Name, es)
			qcpConsume(qsconfig.Name, qsconfigs[j].Name, es)
			subjects += fmt.Sprintf("[%s] [%s]", qsconfigs[j].Name+"2"+qsconfig.Name, qsconfig.Name+"2"+qsconfigs[j].Name)
		}
	}
	return
}

func qcpConsume(from, to string, e chan<- error) {
	var i int64
	listener := func(data []byte, consumer queue.Consumer) {
		i++
		tx := types.Event{}
		err := json.Unmarshal(data, &tx)
		if err != nil {
			return
		}
		log.Infof("[#%d] Consume subject [%s] nodeAddress '%s' event '%s'",
			i, consumer.Subject(), tx.NodeAddress,string(data))
		ferry(from,to,tx)
	}
	subject := from + "2" + to
	consumer, err := queue.NewConsumer(subject)
	if err != nil {
		e <- err
	}
	if consumer == nil {
		e <- fmt.Errorf("New consumer error: get nil")
		return
	}
	consumer.Subscribe(listener)
	return
}

func ferry(from,to string, tx types.Event) {
	fromAd := getAdapter(from)
	toAd := getAdapter(to)
	switch tx.Func {
	case "FABRICCREATE":
		timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
		iotdata := &iotaSDK.IoTData{
			tx.IotaPayload.ContainerID,
			"",
			"",
			timestamp,
			"start",
		}
		iotdatabytes,err := json.Marshal(iotdata)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		transmitTxInfo := iotaSDK.TransmitTxInfo{
			Message:string(iotdatabytes),
			Mamstate:"",
			Seed:tx.IotaPayload.Seed,
			Mode:tx.IotaPayload.Mode,
			SideKey:tx.IotaPayload.SideKey,
			TransactionTag:tx.Func,
			Endpoint:"",
		}
		transmitTxInfoBytes,err := json.Marshal(transmitTxInfo)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		transmitReturnBytes,err :=toAd.SubmitTx(string(transmitTxInfoBytes))
		if err != nil {
			log.Error(err)
			panic(err)
		}
		var transmitReturn iotaSDK.TransmitResult
		err = json.Unmarshal([]byte(transmitReturnBytes),&transmitReturn)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		root :=transmitReturn.Root
		mamstat := transmitReturn.Mamstate
		var argsArray []fabricSDK.Args
		argsArray = append(argsArray, fabricSDK.Args{"InTransitLogistics",[]string{tx.IotaPayload.ContainerID,root,mamstat}})
		argsBytes,err := json.Marshal(argsArray)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		fabReturn, err := fromAd.SubmitTx(string(argsBytes))
		if err !=nil {
			log.Error("fabric  failed to callback for InTransitLogistics :  %v\n", err.Error())
		}
		log.Info(fabReturn)
	case "FABRICDELEVERY":
		timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
		iotdata := &iotaSDK.IoTData{
			tx.IotaPayload.ContainerID,
			"",
			"",
			timestamp,
			"start",
		}
		iotdatabytes,err := json.Marshal(iotdata)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		transmitTxInfo := iotaSDK.TransmitTxInfo{
			Message:string(iotdatabytes),
			Mamstate:"",
			Seed:tx.IotaPayload.Seed,
			Mode:tx.IotaPayload.Mode,
			SideKey:tx.IotaPayload.SideKey,
			TransactionTag:tx.Func,
			Endpoint:"",
		}
		transmitTxInfoBytes,err := json.Marshal(transmitTxInfo)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		_,err =toAd.SubmitTx(string(transmitTxInfoBytes))
		if err != nil {
			log.Error(err)
			panic(err)
		}
		receiveInfo := iotaSDK.RecieveTxInfo{
			Root:tx.IotaPayload.Root,
			Mode:tx.IotaPayload.Mode,
			SideKey:tx.IotaPayload.SideKey,
			Endpoint:"",
		}
		receiveInfobytes,err := json.Marshal(receiveInfo)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		temperature,err :=toAd.ObtainTx(string(receiveInfobytes))
		if err != nil {
			log.Error(err)
			panic(err)
		}
		var argsArray []fabricSDK.Args
		fmt.Println("SignLogistics temperature: ",temperature)
		argsArray = append(argsArray, fabricSDK.Args{"SignLogistics",[]string{tx.IotaPayload.ContainerID,temperature}})
		argsBytes,err := json.Marshal(argsArray)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		fabReturn, err := fromAd.SubmitTx(string(argsBytes))
		if err !=nil {
			log.Error("fabric  failed to callback for InTransitLogistics :  %v\n", err.Error())
		}
		log.Info(fabReturn)
	default:
		log.Error("could not find the relay func")
	}
}

func getAdapter(chainName string) ports.Adapter {
	ads,err :=ports.GetAdapters(chainName)
	if err != nil {
		log.Error(err)
	}
	for _, v := range ads {
		return v
	}
	return nil
}