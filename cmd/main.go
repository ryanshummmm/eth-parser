package main

import (
	"context"
	"eth-parser/internal/api"
	"eth-parser/internal/config"
	"eth-parser/internal/ethereum"
	"eth-parser/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Initialize config
	cfg := config.Load()

	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Initialize storage
	memoryStorage := storage.NewMemoryStorage()
	// Initialize parser
	parser := ethereum.NewEthParser(memoryStorage)
	// Initialize API handler
	handler := api.NewHandler(parser, logger)

	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/current-block", handler.GetCurrentBlockHandler)
	mux.HandleFunc("/subscribe-list", handler.GetSubscribeListHandler)
	mux.HandleFunc("/subscribe", handler.SubscribeHandler)
	mux.HandleFunc("/unsubscribe", handler.UnsubscribeHandler)
	mux.HandleFunc("/transactions", handler.GetTransactionsHandler)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: mux,
	}

	// Start the background task
	parser.Start()

	// Start the server
	go func() {
		logger.Printf("Server starting on %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	// Stop the parser
	parser.Stop()

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exiting")
}
