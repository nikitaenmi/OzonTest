package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nikitaenmi/OzonTest/internal/config"

	"github.com/joho/godotenv"
)

type PaymentRequest struct {
	Provider string  `json:"provider"`
	Amount   float64 `json:"amount"`
	Date     string  `json:"date"`
	Currency string  `json:"currency"`
}

type PaymentResponse struct {
	ID       int     `json:"id"`
	Provider string  `json:"provider"`
	Amount   float64 `json:"amount"`
	Date     string  `json:"date"`
	Currency string  `json:"currency"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type testResult struct {
	successCount int32
	errorCount   int32
	mu           sync.Mutex
	requestTypes map[string]int32
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("failed to load env", "error", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	result := testResult{
		requestTypes: make(map[string]int32),
	}
	sem := make(chan struct{}, cfg.LoadTest.Concurrency)
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < cfg.LoadTest.TotalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			req, expectedStatus, reqType := generateRequest(id)
			result.recordRequestType(reqType)

			if err := sendRequest(cfg.LoadTest.BaseURL, cfg.LoadTest.Timeout, id, req, expectedStatus, reqType); err != nil {
				atomic.AddInt32(&result.errorCount, 1)
				slog.Error("request failed", "request_id", id, "error", err)
			} else {
				atomic.AddInt32(&result.successCount, 1)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	result.printStats(duration, cfg.LoadTest.TotalRequests)
}

func (r *testResult) recordRequestType(reqType string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requestTypes[reqType]++
}

func (r *testResult) printStats(duration time.Duration, total int) {
	slog.Info("load test completed",
		"total_requests", total,
		"successful", r.successCount,
		"errors", r.errorCount,
		"duration", duration,
		"rps", float64(total)/duration.Seconds())

	slog.Info("request distribution", "distribution", r.requestTypes)
}

func generateRequest(id int) (PaymentRequest, int, string) {
	switch id % 5 {
	case 0:
		return PaymentRequest{
			Provider: "tbank",
			Amount:   float64(10 + id%100),
			Date:     "01/01/2024",
			Currency: "USD",
		}, http.StatusOK, "valid_tbank_usd"
	case 1:
		return PaymentRequest{
			Provider: "alpha",
			Amount:   float64(10 + id%100),
			Date:     "01/01/2024",
			Currency: "EUR",
		}, http.StatusOK, "valid_alpha_eur"
	case 2:
		return PaymentRequest{
			Provider: "tbank",
			Amount:   15001,
			Date:     "01/01/2024",
			Currency: "USD",
		}, http.StatusBadRequest, "amount_exceeded"
	case 3:
		return PaymentRequest{
			Provider: "alpha",
			Amount:   -50,
			Date:     "01/01/2024",
			Currency: "EUR",
		}, http.StatusBadRequest, "negative_amount"
	default:
		return PaymentRequest{
			Provider: "tbank",
			Amount:   100,
			Date:     "01/01/2024",
			Currency: "ABC",
		}, http.StatusBadRequest, "invalid_currency"
	}
}

func sendRequest(baseURL string, timeout time.Duration, id int, req PaymentRequest, expectedStatus int, reqType string) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("status %d: %s", resp.StatusCode, errorResp.Message)
		}
		return fmt.Errorf("unexpected status: %d, expected: %d", resp.StatusCode, expectedStatus)
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var paymentResp PaymentResponse
		if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
