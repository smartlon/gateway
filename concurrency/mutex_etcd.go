package concurrency

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/QOSGroup/cassini/log"
	v3 "github.com/coreos/etcd/clientv3"
	v3c "github.com/coreos/etcd/clientv3/concurrency"
)

// NewEtcdMutex new a mutex for a etcd implementation.
func NewEtcdMutex(chainID string, addrs []string) (*EtcdMutex, error) {
	cli, err := v3.New(v3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Errorf("New client error: %s", err)
		return nil, err
	}

	var sess *v3c.Session
	sess, err = v3c.NewSession(cli, v3c.WithTTL(5))
	if err != nil {
		cli.Close()
		log.Errorf("New session error: %s", err)
		return nil, err
	}

	m := &EtcdMutex{chainID: chainID,
		client:  cli,
		session: sess}

	return m, nil
}

// EtcdMutex implements a distributed lock based on etcd.
type EtcdMutex struct {
	chainID  string
	client   *v3.Client
	session  *v3c.Session
	mutex    *v3c.Mutex
	locked   bool
}

// Lock get lock
func (e *EtcdMutex)Lock(address string) (string, error) {
	e.mutex = v3c.NewMutex(e.session, e.chainID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()

	var err error
	err = e.mutex.Lock(ctx)
	if err != nil {
		log.Errorf("Lock error: %s", err)
		return "", err
	}
	iotaPayload,err := e.get(address)

	e.locked = true

	if err != nil {
		defer func() {
			err := e.Unlock(false,address)
			if err != nil {
				log.Error("Unlock error: ", err)
			}
		}()
		log.Error(err)
		return "", err
	}

	log.Debugf("Get lock success, %s: %s", e.chainID, address)
	return iotaPayload, nil
}

// Update update the lock
func (e *EtcdMutex) Update(address string,iota string) error {
	mux := v3c.NewMutex(e.session, e.chainID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()

	var err error
	err = mux.Lock(ctx)
	if err != nil {
		log.Errorf("Update lock address(%s) error: %s", address, err)
		return err
	}
	defer func() {
		mux.Unlock(ctx)
	}()
	err = e.put(address,iota)
	if err != nil {
		log.Errorf("Update address(%s), iota: %s, error: %s",
			address, iota, err)
		return err
	}
	log.Debugf("Upadte lock success, %s: %s", e.chainID, address)
	return nil
}

// Unlock unlock the lock
func (e *EtcdMutex) Unlock(success bool,address string) (err error) {
	if !e.locked {
		return nil
	}
	if success {
		err = e.delete(address)
		if err != nil {
			log.Errorf("Put key value error when unlock: ", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()
	err = e.mutex.Unlock(ctx)
	if err != nil {
		log.Errorf("Unlock error: ", err)
		return
	}
	e.locked = false
	log.Debugf("Unlock success, %s: %s", e.chainID, address)
	return
}


func (e *EtcdMutex) get(address string) (string,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()
	// var resp *v3.GetResponse
	key := generateKey(e.chainID,address)
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		log.Error("Get key value error: ", err)
		return "",err
	}
	for _, kv := range resp.Kvs {
		if strings.EqualFold(string(kv.Key), key ) {
			return string(kv.Value),nil
		}
	}
	return "",fmt.Errorf("iotaPayload is not in the etcd ,address %s",address)
}

func (e *EtcdMutex) put(address string, iotaPayload string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()
	// var resp *v3.PutResponse
	key := generateKey(e.chainID,address)
	_, err := e.client.Put(ctx, key, iotaPayload)
	if err != nil {
		log.Error("Put key value error: ", err)
	}
	return err
}

func (e *EtcdMutex) delete(address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置2s超时
	defer cancel()
	// var resp *v3.PutResponse
	key := generateKey(e.chainID,address)
	_, err := e.client.Delete(ctx, key)
	if err != nil {
		log.Error("delete key value error: ", err)
	}
	return err
}

// Close close the lock
func (e *EtcdMutex) Close() error {
	e.session.Close()
	return e.client.Close()
}

func generateKey(chainid string, address string) string {
	return fmt.Sprintf("%s/%s",chainid,address)
}
