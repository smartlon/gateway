package sdk

import (
	"fmt"
	"github.com/iotaledger/iota.go/account"
	"github.com/iotaledger/iota.go/account/event"
	"github.com/iotaledger/iota.go/account/event/listener"
)

// NewEventLoggerPlugin ...
func NewEventLoggerPlugin(em event.EventMachine) account.Plugin {
	return &logplugin{em: em, exit: make(chan struct{})}
}

type logplugin struct {
	em   event.EventMachine
	acc  account.Account
	exit chan struct{}
}


func (l *logplugin) Name() string {
	return "eventLogger"
}

func (l *logplugin) Start(acc account.Account) error {
	l.acc = acc
	l.log()
	return nil
}
func (l *logplugin) Shutdown() error {
	l.exit <- struct{}{}
	return nil
}

func (l *logplugin) log() {
	lis := listener.NewChannelEventListener(l.em).All()

	go func() {
		defer lis.Close()
	exit:
		for {
			select {
			case ev := <-lis.Promoted:
				fmt.Printf("Promoted %s with %s\n", ev.BundleHash[:10], ev.PromotionTailTxHash)
			case ev := <-lis.Reattached:
				fmt.Printf("Reattached %s with %s\n", ev.BundleHash[:10], ev.ReattachmentTailTxHash)
			case ev := <-lis.SentTransfer:
				tail := ev[0]
				fmt.Printf("Sent %s with tail %s\n", tail.Bundle[:10], tail.Hash)
			case ev := <-lis.TransferConfirmed:
				tail := ev[0]
				fmt.Printf("Transfer confirmed %s with tail %s\n", tail.Bundle[:10], tail.Hash)
			case ev := <-lis.ReceivingDeposit:
				tail := ev[0]
				fmt.Printf("Receiving deposit %s with tail %s\n", tail.Bundle[:10], tail.Hash)
			case ev := <-lis.ReceivedDeposit:
				tail := ev[0]
				fmt.Printf("Received deposit %s with tail %s\n", tail.Bundle[:10], tail.Hash)
			case ev := <-lis.ReceivedMessage:
				tail := ev[0]
				fmt.Printf("Received msg %s with tail %s\n", tail.Bundle[:10], tail.Hash)
			case balanceCheck := <-lis.ExecutingInputSelection:
				fmt.Printf("Doing input selection (balance check: %v) \n", balanceCheck)
			case <-lis.PreparingTransfers:
				fmt.Printf("Preparing transfers\n")
			case <-lis.GettingTransactionsToApprove:
				fmt.Printf("Getting transactions to approve\n")
			case <-lis.AttachingToTangle:
				fmt.Printf("Doing proof of work\n")
			case err := <-lis.InternalError:
				fmt.Printf("Received internal error: %s\n", err.Error())
			case <-l.exit:
				break exit
			}
		}
	}()
}