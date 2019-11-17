package sdk

import (
	"fmt"
	"github.com/iotaledger/iota.go/account"
	. "github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/trinary"
	"log"
)

func GetInMemorySeed(seed string)(inMemorySeed trinary.Trytes){
	seedProv := account.NewInMemorySeedProvider(seed)
	inMemorySeed, err := seedProv.Seed()
	if err != nil {
		log.Fatal(err)
	}
	return inMemorySeed
}

func GetNewAddress(seed string,api *API)(trinary.Hashes){
	// GetNewAddress retrieves the first unspent from address through IRI
	addresses, err := api.GetNewAddress(seed, GetNewAddressOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nYour new address: ", addresses[0])
	return addresses
}