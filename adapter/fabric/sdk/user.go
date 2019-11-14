package sdk

import (

	"github.com/pkg/errors"
	mspapi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

type userAction struct {
	Action
	numEnrolled uint32
	done       chan bool
}

func newUserAction() (*userAction, error) {
	action := &userAction{done: make(chan bool)}
	err := action.Initialize()
	//listener(action)
	return action, err
}

func (a *userAction)registerUser(username, pwd string) (mspapi.SigningIdentity, error) {
	user, err := a.newUser(a.OrgID(),username,pwd)
	if err != nil {
		return nil, errors.Errorf("error getting user: %s", err)
	}
	return user,nil
}