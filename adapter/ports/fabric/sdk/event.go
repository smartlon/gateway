/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
	"github.com/smartlon/gateway/adapter/ports"
	"github.com/smartlon/gateway/adapter/ports/fabric"
	"github.com/smartlon/gateway/log"
	"github.com/smartlon/gateway/types"
	"time"
)



func  Listener(chaincodeID,peerUrl string,listener ports.EventsListener,a *fabric.FabAdaptor) {
	if chaincodeID == "" {
		err := fmt.Errorf("must specify the chaincode ID")
		panic(err)
	}
	action, err := newChaincodeInvokeAction()
	defer action.Terminate()
	action.Set(Config().ChannelID,chaincodeID,[]Args{})
	if err != nil {
		log.Errorf("Error while initializing invokeAction: %v", err)
		panic(err)
	}
	ec, err := action.EventClient(peerUrl,event.WithBlockEvents())
	if err != nil {
		fmt.Println("failed to create client")
		panic(err)
	}

	registrationCreateChannel, notifierCreateChannel, err := ec.RegisterChaincodeEvent(chaincodeID, `{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: CreateChannel")
		panic(err)
                
	}
	registrationDeliveryLogistics, notifierDeliveryLogistics, err := ec.RegisterChaincodeEvent(chaincodeID, `{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: DeliveryLogistics")
		panic(err)
	}
	defer  unregister(ec,[]fab.Registration{registrationCreateChannel,registrationDeliveryLogistics})
	var iotapayload types.IotaPayload
	for {
		event :=&types.Event{}
		select {
		case ccEvent,ok := <-notifierCreateChannel:
			if !ok {
				panic(errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event"))
			}
			err := json.Unmarshal(ccEvent.Payload,&iotapayload)
			if err != nil {
				fmt.Println(err.Error())
			}
			event.Func = "FABRICCREATE"
			fmt.Printf("received chaincode event CreateChannel:  %v\n", iotapayload)


			//timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
			//iotdata := &sdk.IoTData{
			//	iotapayload.ContainerID,
			//	"",
			//	"",
			//	timestamp,
			//	"start",
			//}
			//iotdatabytes,_ := json.Marshal(iotdata)
			//mamstat, root,address:= sdk.MAMTransmit(string(iotdatabytes),"",iotapayload.Seed,iotapayload.Mode,iotapayload.SideKey,"FABRIC")
			//
			//var argsArray []Args
			//argsArray = append(argsArray, Args{"InTransitLogistics",[]string{iotapayload.ContainerID,root,mamstat}})
			//_, err = action.invoke(Config().ChannelID, chaincode, argsArray)
			//if err !=nil {
			//	fmt.Printf("fabric  failed to callback for InTransitLogistics :  %v\n", err.Error())
			//}
		case ccEvent,ok := <-notifierDeliveryLogistics:
			if !ok {
				panic(errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event"))
			}
			err := json.Unmarshal(ccEvent.Payload,&iotapayload)
			if err != nil {
				fmt.Println(err.Error())
			}
			event.Func = "FABRICDELEVERY"
			fmt.Printf("received chaincode event DeliveryLogistics:  %v\n", iotapayload)
			//timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
			//iotdata := &sdk.IoTData{
			//	iotapayload.ContainerID,
			//	"",
			//	"",
			//	timestamp,
			//	"end",
			//}
			//iotdatabytes,_ := json.Marshal(iotdata)
			//sdk.MAMTransmit(string(iotdatabytes),iotapayload.MamState,iotapayload.Seed,iotapayload.Mode,iotapayload.SideKey,"FABRIC")
			//var argsArray []Args
			//iotdatas := sdk.MAMReceive(iotapayload.Root,iotapayload.Mode,iotapayload.SideKey)
			//var temperature string
			//for _,iotdata := range iotdatas {
			//	var iot sdk.IoTData
			//	err := json.Unmarshal([]byte(iotdata),&iot)
			//	if err != nil {
			//		fmt.Printf("fabric  failed to callback for SignLogistics :  %v\n", err.Error())
			//	}
			//	if len(iot.Temperature) != 0 {
			//		temperature = temperature + iot.Temperature + ","
			//	}
			//}
			//temperature = strings.TrimSuffix(temperature,",")
			//fmt.Println("SignLogistics temperature: ",temperature)
			//argsArray = append(argsArray, Args{"SignLogistics",[]string{iotapayload.ContainerID,temperature}})
			//_, err = action.invoke(Config().ChannelID, chaincode, argsArray)
			//if err !=nil {
			//	fmt.Printf("fabric  failed to callback for SignLogistics :  %v\n", err.Error())
			//}
		case <-time.After(time.Second * 3):
			fmt.Println("Exit while waiting for chaincode event")
		}
		event.IotaPayload=iotapayload
		event.From = "FABRIC"
		event.To = "IOTA"
		event.NodeAddress = ports.GetAdapterKey(a)

		listener(event,a)
	}
}

func unregister(ec *event.Client,registrations []fab.Registration) {
	for _,registration := range registrations {
		ec.Unregister(registration)
	}
}






