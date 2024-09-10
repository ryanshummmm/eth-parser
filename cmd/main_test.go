package main

import (
	"bytes"
	"encoding/json"
	"eth-parser/internal/api"
	"eth-parser/internal/ethereum"
	"eth-parser/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	memStorage := storage.NewMemoryStorage()
	parser := ethereum.NewEthParser(memStorage)
	handler := api.NewHandler(parser)

	parser.Start()
	defer parser.Stop()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/current-block":
			handler.GetCurrentBlockHandler(w, r)
		case "/subscribe":
			handler.SubscribeHandler(w, r)
		case "/unsubscribe":
			handler.UnsubscribeHandler(w, r)
		case "/subscribe-list":
			handler.GetSubscribeListHandler(w, r)
		case "/transactions":
			handler.GetTransactionsHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	t.Run("GetCurrentBlock", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/current-block")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}
	})

	t.Run("CheckSubscribe", func(t *testing.T) {

		checkSubscribeList := func(expectedAddresses []string) {
			resp, err := http.Get(ts.URL + "/subscribe-list")
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK for subscribe list; got %v", resp.Status)
			}

			var result map[string][]string
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatal(err)
			}

			if len(result["subscribedAddresses"]) != len(expectedAddresses) {
				t.Errorf("Expected %d subscribed addresses, got %d", len(expectedAddresses), len(result["subscribedAddresses"]))
			}

			for _, addr := range expectedAddresses {
				addr = strings.ToLower(addr)
				found := false
				for _, subAddr := range result["subscribedAddresses"] {
					if addr == subAddr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected address %s not found in subscribe list", addr)
				}
			}
		}

		// Check initial empty list
		checkSubscribeList([]string{})

		address := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		payload := map[string]string{"address": address}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(ts.URL+"/subscribe", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}

		// Check list after subscribe
		checkSubscribeList([]string{address})

		// Unsubscribe
		resp, err = http.Post(ts.URL+"/unsubscribe", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK for unsubscribe; got %v", resp.Status)
		}

		// Check list after unsubscribe
		checkSubscribeList([]string{})
	})

}
