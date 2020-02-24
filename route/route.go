package route

import (
	"encoding/json"
	"errors"

	"github.com/smartlon/gateway/log"
	"github.com/smartlon/gateway/queue"
	"github.com/smartlon/gateway/types"
)

//type route struct{}

// Event2queue produce event to message queue (Nats)
func Event2queue(nats string, event *types.Event) (subject string, err error) {

	if event == nil {

		return "", errors.New("event is nil")
	}
	if event.From == "" || event.To == "" || event.NodeAddress == "" {
		return "", errors.New("event data is empty")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	subject = event.From + "2" + event.To

	producer, err := queue.NewProducer(subject)
	if err != nil {
		return "", err
	}
	producer.Produce(data)

	log.Infof("routed event from[%s]  to subject [%s] ", event.NodeAddress,  subject)

	return subject, nil
}
