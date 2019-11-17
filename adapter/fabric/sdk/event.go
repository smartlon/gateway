/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
	"time"
)

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
	select {
	case ccEvent,ok := <-notifierCreateChannel:
		if !ok {
			return errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event")
		}
		fmt.Printf("received chaincode event CreateChannel:  %v\n", ccEvent)
	case ccEvent,ok := <-notifierDeliveryLogistics:
		if !ok {
			return errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event")
		}
		fmt.Printf("received chaincode event DeliveryLogistics:  %v\n", ccEvent)
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






