// Package concurrency provides support for distributed concurrency capabilities
//
// The goal is to simplify the development of concurrency capabilities.
package concurrency

import (
	"github.com/smartlon/gateway/types"
)

// Mutex distributed lock interface
type Mutex interface {
	// Lock get lock through chain id and sequence
	//
	// If it returned an error, indicates that the call failed.
	// Whether successful or not,
	// the current sequence saved in the distributed lock is returned.
	// Negative sequence(<0) are returned unless there are some unknown exceptions.
	Lock(address string) (string, error)

	// Update update the sequence saved in the distributed lock.
	Update(address string,iota string) error

	// Unlock after successfully acquiring the lock, the lock needs to be unlocked.
	//
	// If it returned an error, indicates that the call failed.
	Unlock(success bool,address string) error

	// Close close the lock
	Close() error
}

// NewMutex new mutex based on configuration.
func NewMutex(name, address string) (m Mutex, err error) {
	protocol, addrs := types.ParseAddrs(address)

	switch protocol {
	case "etcd":
		m, err = NewEtcdMutex(name, addrs)
	default:
		m = NewStandaloneMutex(name)
	}
	return
}
