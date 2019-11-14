package sdk

import (
	. "github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/pow"
)



//NewConnection establishes a connection with the given endpoint

func NewConnection()(*API, error){
	_, proofOfWorkFunc := pow.GetFastestProofOfWorkImpl()
	//endpoint := beego.AppConfig.String("endpoint")
	api, err := ComposeAPI(HTTPClientSettings{
		URI:ENDPOINT,
		LocalProofOfWorkFunc: proofOfWorkFunc,
	})
	if err != nil {
		return nil, err
	}
	return api,nil


}

