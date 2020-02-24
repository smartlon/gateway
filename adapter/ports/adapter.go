package ports

import (
	"github.com/smartlon/gateway/types"
)


// Adapter Chain adapter interface for consensus engine ( consensus.ConsEngine )
// and ferry ( consensus.Ferry )
type Adapter interface {
	SubmitTx(tx string) (string, error)
	ObtainTx(tx string) (string, error)
	// Count Calculate the total and consensus number for chain
	Count() (totalNumber int, consensusNumber int)
	GetChainName() string
	GetIP() string
	GetPort() int
}

// EventsListener Listen Tx events from target chain
type EventsListener func(event *types.Event, adapter Adapter)

// AdapterController Chain adapter controller interface for adapter pool manager ( adapter.Ports )
type AdapterController interface {
	Start() error
	Sync() error
	Stop() error
	Subscribe(listener EventsListener)
}

/*
AdapterService Inner cache type ( AdapterController and Adapter )

Suitable for a variety of different block chain
*/
type AdapterService interface {
	AdapterController
	Adapter
}

// AdapterConfig is parameters for build an AdapterService
type AdapterConfig struct {
	ChainName string
	ChainType string
	IP        string
	Port      int
	Query     string
	Listener  EventsListener
}


