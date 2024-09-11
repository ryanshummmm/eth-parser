package api

import (
	"encoding/json"
	"eth-parser/common"
	"eth-parser/internal/ethereum"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	parser ethereum.Parser
	logger *log.Logger
}

func NewHandler(parser ethereum.Parser, logger *log.Logger) *Handler {
	return &Handler{parser: parser, logger: logger}
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf("Health check: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{"status": "OK"}
	w.Header().Set(common.HeaderContentTypeKey, common.ApplicationJsonContentType)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Health check: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Println("Health check: OK")
}

func (h *Handler) GetCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf("Get current block: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentBlock := h.parser.GetCurrentBlock()
	err := json.NewEncoder(w).Encode(map[string]int64{"currentBlock": currentBlock})
	if err != nil {
		h.logger.Printf("Get current block: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetSubscribeListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf("Get subscribe list: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	list := h.parser.GetSubscribeList()
	response := map[string][]string{"subscribedAddresses": list}
	w.Header().Set(common.HeaderContentTypeKey, common.ApplicationJsonContentType)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Get subscribe list: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Printf("Get subscribe list: Returned %d addresses", len(list))

}

func (h *Handler) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Printf("Subscribe: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Subscribe: Error decoding request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	address := strings.ToLower(req.Address)
	success := h.parser.Subscribe(address)
	w.Header().Set(common.HeaderContentTypeKey, common.ApplicationJsonContentType)

	if err := json.NewEncoder(w).Encode(map[string]bool{"subscribed": success}); err != nil {
		h.logger.Printf("Subscribe: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Printf("Subscribe: Address %s, Success: %v", address, success)
}

func (h *Handler) UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Printf("Unsubscribe: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Unsubscribe: Error decoding request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set(common.HeaderContentTypeKey, common.ApplicationJsonContentType)
	success := h.parser.Unsubscribe(req.Address)
	if err := json.NewEncoder(w).Encode(map[string]bool{"unsubscribed": success}); err != nil {
		h.logger.Printf("Unsubscribe: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Printf("Unsubscribe: Address %s, Success: %v", req.Address, success)
}

func (h *Handler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf("Get transactions: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		h.logger.Println("Get transactions: No address provided")
		http.Error(w, "No address provided", http.StatusBadRequest)
		return
	}

	w.Header().Set(common.HeaderContentTypeKey, common.ApplicationJsonContentType)
	transactions := h.parser.GetTransactions(address)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		h.logger.Printf("Get transactions: Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Printf("Get transactions: Returned %d transactions for address %s", len(transactions), address)
}
