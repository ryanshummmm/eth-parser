package storage

import "eth-parser/pkg/models"

type Storage interface {
	GetCurrentBlock() int
	SetCurrentBlock(number int)
	GetSubscribeList() []string
	Subscribe(address string) bool
	Unsubscribe(address string) bool
	IsSubscribed(address string) bool
	GetTransactions(address string) []models.Transaction
	AddTransaction(tx models.Transaction)
}
