package storage

import (
	"eth-parser/pkg/models"
	"strings"
	"sync"
)

type MemoryStorage struct {
	currentBlock        int64
	subscribedAddresses map[string]bool
	transactions        map[string][]models.Transaction
	mu                  sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		subscribedAddresses: make(map[string]bool),
		transactions:        make(map[string][]models.Transaction),
	}
}

func (ms *MemoryStorage) GetCurrentBlock() int64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.currentBlock
}

func (ms *MemoryStorage) SetCurrentBlock(block int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.currentBlock = block
}

func (ms *MemoryStorage) GetSubscribeList() []string {
	keys := make([]string, 0, len(ms.subscribedAddresses))
	for key := range ms.subscribedAddresses {
		keys = append(keys, key)
	}
	return keys
}

func (ms *MemoryStorage) Subscribe(address string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	address = strings.ToLower(address)
	if !ms.subscribedAddresses[address] {
		ms.subscribedAddresses[address] = true
	}
	return true
}

func (ms *MemoryStorage) Unsubscribe(address string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	address = strings.ToLower(address)

	if ms.subscribedAddresses[address] {
		delete(ms.subscribedAddresses, address)
		return true
	}
	return false
}

func (ms *MemoryStorage) IsSubscribed(address string) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	address = strings.ToLower(address)
	return ms.subscribedAddresses[address]
}

func (ms *MemoryStorage) GetTransactions(address string) []models.Transaction {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	address = strings.ToLower(address)
	return ms.transactions[address]
}

func (ms *MemoryStorage) AddTransaction(tx models.Transaction) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.subscribedAddresses[tx.From] {
		ms.transactions[tx.From] = append(ms.transactions[tx.From], tx)
	}
	if ms.subscribedAddresses[tx.To] {
		ms.transactions[tx.To] = append(ms.transactions[tx.To], tx)
	}
}
