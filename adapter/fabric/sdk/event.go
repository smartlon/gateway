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
)

type chaincodeListenAction struct {
	Action
	done       chan bool
}

func newChaincodeListenAction() (*chaincodeListenAction, error) {
	action := &chaincodeListenAction{done: make(chan bool)}
	err := action.Initialize()
	return action, err
}

func (a *chaincodeListenAction) Listener(chaincode string) error {

	ec, err := a.EventClient(event.WithBlockEvents())
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
	defer a.Unregister(ec,[]fab.Registration{registrationCreateChannel,registrationDeliveryLogistics})
	exitLisener :=a.WaitForExit()
	for {
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
		case _, _ = <-exitLisener:
			fmt.Println("Exit while waiting for chaincode event")
			return nil
		}
	}

}

func (a *chaincodeListenAction) Unregister(ec *event.Client,registrations []fab.Registration) {
	for _,registration := range registrations {
		ec.Unregister(registration)
	}
}


// WaitForEnter waits until the user presses Enter
func (a *chaincodeListenAction) WaitForExit() chan bool {
	a.done <- true
	return a.done
}




