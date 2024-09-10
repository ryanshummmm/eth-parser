package main

import (
	"eth-parser/internal/api"
	"eth-parser/internal/ethereum"
	"eth-parser/internal/storage"
	"log"
	"net/http"
)

func main() {

	memoryStorage := storage.NewMemoryStorage()
	parser := ethereum.NewEthParser(memoryStorage)
	handler := api.NewHandler(parser)

	http.HandleFunc("/current-block", handler.GetCurrentBlockHandler)
	http.HandleFunc("/subscribe-list", handler.GetSubscribeListHandler)
	http.HandleFunc("/subscribe", handler.SubscribeHandler)
	http.HandleFunc("/unsubscribe", handler.UnsubscribeHandler)
	http.HandleFunc("/transactions", handler.GetTransactionsHandler)

	// Start the background task
	parser.Start()
	defer parser.Stop()

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
