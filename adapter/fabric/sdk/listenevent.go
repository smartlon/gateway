/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"

)

func listener(action *chaincodeInvokeAction) {

	ec, err := action.EventClient()
	if err != nil {
		fmt.Println("failed to create client")
	}

	registrationCreateChannel, notifierCreateChannel, err := ec.RegisterChaincodeEvent("log", `{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event")
	}
	registrationDeliveryLogistics, notifierDeliveryLogistics, err := ec.RegisterChaincodeEvent("log", `{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event")
	}
	defer unregister(ec,[]fab.Registration{registrationCreateChannel,registrationDeliveryLogistics})

	select {
	case ccEvent := <-notifierCreateChannel:
		fmt.Printf("received chaincode event %v\n", ccEvent)
	case ccEvent := <-notifierDeliveryLogistics:
		fmt.Printf("received chaincode event %v\n", ccEvent)
	case <-time.After(time.Second * 5):
		fmt.Println("timeout while waiting for chaincode event")
	}

	// Timeout is expected since there is no event producer

	// Output: timeout while waiting for chaincode event

}

func unregister(ec *event.Client,registrations []fab.Registration) {
	for _,registration := range registrations {
		ec.Unregister(registration)
	}
}





