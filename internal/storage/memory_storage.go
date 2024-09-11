package storage

import (
	"eth-parser/pkg/models"
	"strings"
	"sync"
)

type MemoryStorage struct {
	currentBlock        int64
	subscribedAddresses sync.Map
	transactions        sync.Map
	mu                  sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
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
	var addresses []string

	ms.subscribedAddresses.Range(func(key, value interface{}) bool {
		addresses = append(addresses, key.(string))
		return true
	})
	return addresses
}

func (ms *MemoryStorage) Subscribe(address string) bool {
	address = strings.ToLower(address)
	_, loaded := ms.subscribedAddresses.LoadOrStore(address, true)

	return !loaded
}

func (ms *MemoryStorage) Unsubscribe(address string) bool {
	address = strings.ToLower(address)
	_, loaded := ms.subscribedAddresses.LoadAndDelete(address)

	return loaded
}

func (ms *MemoryStorage) IsSubscribed(address string) bool {
	address = strings.ToLower(address)
	_, ok := ms.subscribedAddresses.Load(address)

	return ok
}

func (ms *MemoryStorage) GetTransactions(address string) []models.Transaction {
	address = strings.ToLower(address)

	if txs, ok := ms.transactions.Load(address); ok {
		return txs.([]models.Transaction)
	}
	return nil
}

func (ms *MemoryStorage) AddTransaction(tx models.Transaction) {
	ms.addTransactionForAddress(strings.ToLower(tx.From), tx)
	ms.addTransactionForAddress(strings.ToLower(tx.To), tx)
}

func (ms *MemoryStorage) addTransactionForAddress(address string, tx models.Transaction) {
	if ms.IsSubscribed(address) {
		ms.mu.Lock()
		defer ms.mu.Unlock()

		var txs []models.Transaction
		if existingTxs, ok := ms.transactions.Load(address); ok {
			txs = existingTxs.([]models.Transaction)
		}
		txs = append(txs, tx)
		ms.transactions.Store(address, txs)
	}
}
