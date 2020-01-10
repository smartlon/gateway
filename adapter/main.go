package main

import (
	"fmt"
	"github.com/QOSGroup/cassini/adapter/ports"
	"github.com/smartlon/gateway/concurrency"
	"github.com/smartlon/gateway/config"
	"github.com/smartlon/gateway/consensus"
	"github.com/smartlon/gateway/log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main()  {

	log.Info("Starting gateway...")

	var w sync.WaitGroup
	errChannel := make(chan error, 1)
	startLog(errChannel)
	startEtcd(&w)
	startAdapterPorts(errChannel, &w)
	w.Wait()
	startConsensus(&w)

	go func() {
		w.Wait()
		log.Info("gateway started ")
	}()
	return
}

func startLog(errChannel <-chan error) {
	go func() {
		for {
			select {
			case e, ok := <-errChannel:
				{
					if ok && e != nil {
						log.Error(e)
					}
				}
			}
		}
	}()
}

func startEtcd(w *sync.WaitGroup) {
	w.Add(1)
	go func() {
		etcd, err := concurrency.StartEmbedEtcd(config.GetConfig())
		if err != nil {
			panic(fmt.Errorf("Etcd server start error: %v", err))
		}
		w.Done()
		if etcd == nil {
			return
		}
		defer etcd.Close()
		select {
		case <-etcd.Server.ReadyNotify():
			log.Info("Etcd server is ready!")
		case <-time.After(60 * time.Second):
			etcd.Server.Stop() // trigger a shutdown
			log.Info("Etcd server took too long to start!")
		}
		err = <-etcd.Err()
		log.Error("Etcd running error: ", err)
	}()
}

func startAdapterPorts(errChannel chan<- error, w *sync.WaitGroup) {
	log.Info("Starting adapter ports...")
	w.Add(1)
	go func() {
		conf := config.GetConfig()
		for _, qsc := range conf.Qscs {
			for _, nodeAddr := range strings.Split(qsc.Nodes, ",") {
				if err := registerAdapter(
					nodeAddr, qsc, errChannel); err != nil {
					errChannel <- err
				}
			}
		}
		w.Done()
	}()
}

func registerAdapter(nodeAddr string, qsc *config.QscConfig,
	errChannel chan<- error) (err error) {
	defer func() {
		if o := recover(); o != nil {
			if err, ok := o.(error); ok {
				errChannel <- fmt.Errorf(
					"Register adapter error: %v", err)
			}
		}
	}()
	addrs := strings.Split(nodeAddr, ":")
	if len(addrs) != 2 {
		err = fmt.Errorf(
			"Chain(%s) node address(%s) parse error: %s",
			qsc.Name, nodeAddr,
			"invalid node address format")
		return
	}
	var port int
	port, err = strconv.Atoi(addrs[1])
	if err != nil {
		err = fmt.Errorf(
			"Chain(%s) node address(%s) parse error: %v",
			qsc.Name, nodeAddr, err)
		return
	}
	conf := &ports.AdapterConfig{
		ChainName: qsc.Name,
		ChainType: qsc.Type,
		IP:        addrs[0],
		Port:      port}
	if err = ports.RegisterAdapter(conf); err != nil {
		err = fmt.Errorf(
			"Register adapter error: %v", err)
	}
	return
}

func startConsensus(w *sync.WaitGroup) {
	log.Info("Starting qcp consumer...")
	//启动nats 消费
	w.Add(1)
	go func() {
		err := consensus.StartQcpConsume(config.GetConfig())
		if err != nil {
			log.Errorf("Start qcp consume error: %s", err)
			log.Flush()
			os.Exit(1)
		}
		w.Done()
	}()
}