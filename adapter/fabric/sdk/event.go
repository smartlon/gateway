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
	"time"
)


type mamchannel struct {
	Root    string     `json:",Root"`
	SideKey string     `json:"SideKey"`
}

// logisticstrans type
type logisticstrans struct {
	//product might be food,fish,phone,other itmes
	//Product id should be unique such as FISH123,Prawns456,ICECREAM789
	ProductID         string       `json:"ProductID"`
	ProductType       string       `json:"ProductType"`
	SellerID          string       `json:"SellerID"`
	SellerLocation    string       `json:"SellerLocation"`
	BuyerID           string       `json:"BuyerID"`
	BuyerLocation     string       `json:"BuyerLocation"`
	LogisticsID       string       `json:"LogisticsID"`
	LogisticsLocation string       `json:"LogisticsLocation"`
	JourneyStartTime  string       `json:",JourneyStartTime"`
	JourneyEndTime    string       `json:",JourneyEndTime"`
	Status            string       `json:"Status"`
	MAMChannel    mamchannel `json:"MAMChannel"`
}

func  listener(action *chaincodeInvokeAction,chaincode string) error {

	ec, err := action.EventClient(event.WithBlockEvents())
	if err != nil {
		fmt.Println("failed to create client")
		return err
	}

	registrationCreateChannel, notifierCreateChannel, err := ec.RegisterChaincodeEvent(chaincode, `{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: CreateChannel")
		return err
                
	}
	registrationDeliveryLogistics, notifierDeliveryLogistics, err := ec.RegisterChaincodeEvent(chaincode, `{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: DeliveryLogistics")
		return err
	}
	defer  unregister(ec,[]fab.Registration{registrationCreateChannel,registrationDeliveryLogistics})
	var logisticstran logisticstrans

	select {
	case ccEvent,ok := <-notifierCreateChannel:
		if !ok {
			return errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event")
		}
		err := json.Unmarshal(ccEvent.Payload,&logisticstran)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("received chaincode event CreateChannel:  %v\n", logisticstran)

	case ccEvent,ok := <-notifierDeliveryLogistics:
		if !ok {
			return errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event")
		}
		err := json.Unmarshal(ccEvent.Payload,&logisticstran)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("received chaincode event DeliveryLogistics:  %v\n", logisticstran)
	case <-time.After(time.Second * 10):
		fmt.Println("Exit while waiting for chaincode event")
		}
	return nil


}

func unregister(ec *event.Client,registrations []fab.Registration) {
	for _,registration := range registrations {
		ec.Unregister(registration)
	}
}






