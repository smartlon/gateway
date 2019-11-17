package sdk

import (
	"fmt"
	"github.com/iotaledger/iota.go/account"
	"github.com/iotaledger/iota.go/account/builder"
	"github.com/iotaledger/iota.go/account/event"
	"github.com/iotaledger/iota.go/account/store/badger"
	"github.com/iotaledger/iota.go/account/timesrc"
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/pow"
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

func GetNewAddress(seed string,api *api.API)(trinary.Hashes){
	// GetNewAddress retrieves the first unspent from address through IRI
	addresses, err := api.GetNewAddress(seed, api.GetNewAddressOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nYour new address: ", addresses[0])
	return addresses
}

func NewAccount() (account.Account, error) {
	_, proofOfWorkFunc := pow.GetFastestProofOfWorkImpl()
	//endpoint := beego.AppConfig.String("endpoint")
	iotaAPI, err := api.ComposeAPI(api.HTTPClientSettings{
		URI:ENDPOINT,
		LocalProofOfWorkFunc: proofOfWorkFunc,
	})
	if err != nil {

	}
	store, err := badger.NewBadgerStore("/home/lgao/go/src/github.com/smartlon/gateway/adapter/iota/db")
	if err != nil {

	}
	defer store.Close()

	em := event.NewEventMachine()
	// create an accurate time source (in this case Google's NTP server).
	timesource := timesrc.NewNTPTimeSource("time.google.com")
	account, err := builder.NewBuilder().
		// the IOTA API to use
		WithAPI(iotaAPI).
		// the database onject to use
		WithStore(store).
		// the seed of the account
		WithSeed(SEED).
		// the minimum weight magnitude for the Devnet
		WithMWM(9).
		// Load the EventMachine
		WithEvents(em).
		// the time source to use during input selection
		WithTimeSource(timesource).
		// load the default plugins that enhance the functionality of the account
		WithDefaultPlugins().
		Build()
	if err != nil {
		return nil, err
	}
	return account,nil
}