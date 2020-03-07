package consensus

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smartlon/gateway/config"
	"github.com/smartlon/gateway/log"
	"github.com/smartlon/gateway/queue"
	"github.com/smartlon/gateway/types"
)

// StartQcpConsume Start to consume tx msg
func StartQcpConsume(conf *config.Config) (err error) {
	qsconfigs := conf.Qscs
	if len(qsconfigs) < 2 {
		return errors.New("config error , at least two chain targets ")
	}
	var subjects string
	es := make(chan error, 1024) //TODO 1024参数按需修改
	defer close(es)
	for i, qsconfig := range qsconfigs {
		for j := i + 1; j < len(qsconfigs); j++ {
			qcpConsume(qsconfigs[j].Name, qsconfig.Name, es)
			qcpConsume(qsconfig.Name, qsconfigs[j].Name, es)
			subjects += fmt.Sprintf("[%s] [%s]", qsconfigs[j].Name+"2"+qsconfig.Name, qsconfig.Name+"2"+qsconfigs[j].Name)
		}
	}
	return
}

func qcpConsume(from, to string, e chan<- error) {
	var i int64
	listener := func(data []byte, consumer queue.Consumer) {
		i++
		tx := types.Event{}
		err := json.Unmarshal(data, &tx)
		if err != nil {
			return
		}
		log.Infof("[#%d] Consume subject [%s] nodeAddress '%s' event '%s'",
			i, consumer.Subject(), tx.NodeAddress,string(data))
	}
	subject := from + "2" + to
	consumer, err := queue.NewConsumer(subject)
	if err != nil {
		e <- err
	}
	if consumer == nil {
		e <- fmt.Errorf("New consumer error: get nil")
		return
	}
	consumer.Subscribe(listener)
	return
}