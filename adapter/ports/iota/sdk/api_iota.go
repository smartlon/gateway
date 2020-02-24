package sdk

import (
	"github.com/iotaledger/iota.go/mam/v1"
	"github.com/iotaledger/iota.go/trinary"
	"log"
)

type IoTData struct {
	ContainerID        string `json:"ContainerID"`
	Temperature string `json:"Temperature"`
	Location    string `json:"Location"`
	Time        string `json:"Time"`
	Status        string `json:"Status"`
}

type TransmitTxInfo struct {
	Message        string `json:"Message"`
	Mamstate       string `json:"Mamstate"`
	Seed           string `json:"Seed"`
	Mode           string `json:"Mode"`
	SideKey        string `json:"SideKey"`
	TransactionTag string `json:"TransactionTag"`
	Endpoint           string `json:"Endpoint"`
}

type TransmitResult struct {
	Root        string `json:"Root"`
	Mamstate       string `json:"Mamstate"`
	Address           string `json:"Address"`
}

type RecieveTxInfo struct {
	Root        string `json:"Root"`
	Mode           string `json:"Mode"`
	SideKey        string `json:"SideKey"`
	Endpoint           string `json:"Endpoint"`
}

func MAMTransmit(message string,mamstate string, seed string, mode string, sideKey string, transactionTag string,endpoint string) (string, string,string){
	var t *mam.Transmitter
	if mamstate != "" {
		mamState := StringToMamState(mamstate)
		t = ReconstructTransmitter(seed, mamState,endpoint)
	}
	transmitter, root := Publish(message, t, seed, mode, sideKey,transactionTag,endpoint)
	channel := transmitter.Channel()
	rootTrits,err := trinary.TrytesToTrits(root)
	if err != nil {
		log.Fatal(err)
	}
	address, err := makeAddress(mode, rootTrits, transmitter.Channel().SideKey)
	if err != nil {
		log.Fatal(err)
	}
	return MamStateToString(channel), root,address
}

func  MAMReceive( root string, mode string, sideKey string,endpoint string) ([]string ){
	channelMessages := Fetch(root,mode,sideKey,endpoint )
	return channelMessages
}

