package storage

import (
	"eth-parser/pkg/models"
	"strings"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()

	t.Run("Subscribe", func(t *testing.T) {
		if !ms.Subscribe("0xdAC17F958D2ee523a2206206994597C13D831ec7") {
			t.Error("Failed to subscribe new address")
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		if !ms.Unsubscribe("0xdAC17F958D2ee523a2206206994597C13D831ec7") {
			t.Error("Failed to unsubscribe existing address")
		}
		if ms.Unsubscribe("0xdAC17F958D2ee523a2206206994597C13D831ec7") {
			t.Error("Unsubscribed non-existing address")
		}
	})

	t.Run("IsSubscribed", func(t *testing.T) {
		ms.Subscribe("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		if !ms.IsSubscribed("0xdAC17F958D2ee523a2206206994597C13D831ec7") {
			t.Error("IsSubscribed returned false for subscribed address")
		}
		ms.Unsubscribe("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		if ms.IsSubscribed("0xdAC17F958D2ee523a2206206994597C13D831ec7") {
			t.Error("IsSubscribed returned true for non-subscribed address")
		}
	})

	t.Run("GetSubscribeList", func(t *testing.T) {
		ms.Subscribe("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		if ms.GetSubscribeList()[0] != strings.ToLower("0xdAC17F958D2ee523a2206206994597C13D831ec7") {
			t.Error("GetSubscribeList does not match")
		}
	})

	t.Run("AddTransaction", func(t *testing.T) {
		ms.Subscribe("0xabc")
		ms.AddTransaction(models.Transaction{From: "0xabc", To: "0xdef", Value: "100"})
		txs := ms.GetTransactions("0xabc")
		if len(txs) != 1 {
			t.Error("Expected 1 transaction, got", len(txs))
		}
	})
}
