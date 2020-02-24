package concurrency

import (
	"fmt"
	"sync"
)

// NewStandaloneMutex new a mutex for a standalone implementation.
func NewStandaloneMutex(name string) *StandaloneMutex {
	iota := make(map[string]string)
	return &StandaloneMutex{
		chainID:  name,
		iotaPayload: iota}
}

// StandaloneMutex implements a standalone version for single process.
type StandaloneMutex struct {
	chainID  string
	iotaPayload map[string]string
	mux      sync.Mutex
}

// Lock get lock
func (s *StandaloneMutex) Lock(address string) (string, error) {
	s.mux.Lock()
	iotaPayload,ok := s.iotaPayload[address]
	if !ok {
		s.mux.Unlock()
		return "", fmt.Errorf("Wrong iotaPayload about address %s",address)
	}
	return iotaPayload, nil
}

// Unlock unlock the lock
func (s *StandaloneMutex) Unlock(success bool,address string) error {
	if success {
		delete(s.iotaPayload, address)
	}
	s.mux.Unlock()
	return nil
}

// Update update the sequence in lock
func (s *StandaloneMutex) Update(address string,iota string) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.iotaPayload[address] = iota
	return nil
}

// Close close the lock
func (s *StandaloneMutex) Close() error {
	return nil
}
