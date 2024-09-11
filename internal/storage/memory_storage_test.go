package storage

import (
	"eth-parser/pkg/models"
	"reflect"
	"strings"
	"testing"
)

func TestMemoryStorage(t *testing.T) {

	t.Run("CurrentBlock", func(t *testing.T) {
		ms := NewMemoryStorage()
		if ms.GetCurrentBlock() != 0 {
			t.Errorf("Initial current block should be 0, got %d", ms.GetCurrentBlock())
		}
		ms.SetCurrentBlock(100)
		if ms.GetCurrentBlock() != 100 {
			t.Errorf("Current block should be 100, got %d", ms.GetCurrentBlock())
		}
	})

	t.Run("Subscribe", func(t *testing.T) {
		ms := NewMemoryStorage()

		testCases := []struct {
			address  string
			expected bool
		}{
			{"0xdAC17F958D2ee523a2206206994597C13D831ec7", true},
			{"0xdAC17F958D2ee523a2206206994597C13D831ec7", false}, // Already subscribed
			{"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", true},
		}

		for _, tc := range testCases {
			if result := ms.Subscribe(tc.address); result != tc.expected {
				t.Errorf("Subscribe(%s) = %v, want %v", tc.address, result, tc.expected)
			}
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		ms := NewMemoryStorage()

		address := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		ms.Subscribe(address)

		if !ms.Unsubscribe(address) {
			t.Error("Failed to unsubscribe existing address")
		}
		if ms.Unsubscribe(address) {
			t.Error("Unsubscribed non-existing address")
		}
	})

	t.Run("IsSubscribed", func(t *testing.T) {
		ms := NewMemoryStorage()
		address := "0xdAC17F958D2ee523a2206206994597C13D831ec7"

		if ms.IsSubscribed(address) {
			t.Error("IsSubscribed returned true for non-subscribed address")
		}

		ms.Subscribe(address)
		if !ms.IsSubscribed(address) {
			t.Error("IsSubscribed returned false for subscribed address")
		}

		ms.Unsubscribe(address)
		if ms.IsSubscribed(address) {
			t.Error("IsSubscribed returned true for unsubscribed address")
		}
	})

	t.Run("GetSubscribeList", func(t *testing.T) {
		ms := NewMemoryStorage()
		addresses := []string{
			"0xdAC17F958D2ee523a2206206994597C13D831ec7",
			"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		}

		for _, addr := range addresses {
			ms.Subscribe(addr)
		}

		list := ms.GetSubscribeList()
		if len(list) != len(addresses) {
			t.Errorf("GetSubscribeList() returned %d addresses, want %d", len(list), len(addresses))
		}

		for _, addr := range addresses {
			found := false
			for _, listAddr := range list {
				if strings.EqualFold(addr, listAddr) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Address %s not found in subscribe list", addr)
			}
		}

	})

	t.Run("AddTransaction", func(t *testing.T) {
		ms := NewMemoryStorage()
		from := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
		to := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		ms.Subscribe(from)
		ms.Subscribe(to)

		tx := models.Transaction{From: from, To: to, Value: "100"}
		ms.AddTransaction(tx)

		fromTxs := ms.GetTransactions(from)
		if len(fromTxs) != 1 {
			t.Errorf("Expected 1 transaction for 'from' address, got %d", len(fromTxs))
		}

		toTxs := ms.GetTransactions(to)
		if len(toTxs) != 1 {
			t.Errorf("Expected 1 transaction for 'to' address, got %d", len(toTxs))
		}

		if !reflect.DeepEqual(fromTxs[0], tx) || !reflect.DeepEqual(toTxs[0], tx) {
			t.Error("Stored transaction does not match the original")
		}
	})

	t.Run("GetTransactions", func(t *testing.T) {
		ms := NewMemoryStorage()
		address := "0xabc"
		ms.Subscribe(address)

		if txs := ms.GetTransactions(address); len(txs) != 0 {
			t.Error("GetTransactions returned non-empty list for address with no transactions")
		}

		tx := models.Transaction{From: address, To: "0xdef", Value: "100"}
		ms.AddTransaction(tx)

		txs := ms.GetTransactions(address)
		if len(txs) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(txs))
		}
		if !reflect.DeepEqual(txs[0], tx) {
			t.Error("Retrieved transaction does not match the original")
		}
	})

	t.Run("CaseSensitivity", func(t *testing.T) {
		ms := NewMemoryStorage()
		lowerAddress := "0xabc123"
		upperAddress := "0xABC123"

		ms.Subscribe(lowerAddress)
		if !ms.IsSubscribed(upperAddress) {
			t.Error("Subscription should be case-insensitive")
		}

		ms.AddTransaction(models.Transaction{From: lowerAddress, To: "0xdef", Value: "100"})
		if txs := ms.GetTransactions(upperAddress); len(txs) != 1 {
			t.Error("Transaction retrieval should be case-insensitive")
		}
	})
}
