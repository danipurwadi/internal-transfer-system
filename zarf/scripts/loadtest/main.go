package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	baseURL     = "http://localhost:8080"
	numAccounts = 300
	// initialBalance =
	numTransfers = 10000
	concurrency  = 100 // Increased concurrency
)

// RequestResult stores the outcome of a single HTTP request.

type RequestResult struct {
	Latency    time.Duration
	StatusCode int
	Successful bool
	Endpoint   string
	Error      error
}

func main() {
	fmt.Println("Starting load test...")
	startTime := time.Now()

	resultsChan := make(chan RequestResult, numAccounts+numTransfers)

	// Create accounts
	fmt.Println("Creating accounts...")
	accountIDs := createAccounts(resultsChan)

	// Perform transfers
	if len(accountIDs) > 1 {
		fmt.Println("Performing transfers...")
		performTransfers(accountIDs, resultsChan)
	}

	close(resultsChan)

	totalTime := time.Since(startTime)

	// Collect results
	results := []RequestResult{}
	for result := range resultsChan {
		results = append(results, result)
	}

	fmt.Println("\nLoad test finished.")
	printReport(results, totalTime)
}

func createAccounts(resultsChan chan<- RequestResult) []int {
	var wg sync.WaitGroup
	accountIDs := make(chan int, numAccounts)
	sem := make(chan struct{}, concurrency)

	for i := 1; i <= numAccounts; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(accountID int) {
			defer wg.Done()
			createAccount(accountID, resultsChan)
			accountIDs <- accountID
			<-sem
		}(i)
	}

	wg.Wait()
	close(accountIDs)

	var ids []int
	for id := range accountIDs {
		ids = append(ids, id)
	}
	return ids
}

func createAccount(accountID int, resultsChan chan<- RequestResult) {
	startTime := time.Now()
	url := fmt.Sprintf("%s/accounts", baseURL)
	payload := map[string]interface{}{
		"account_id":      accountID,
		"initial_balance": fmt.Sprintf("%.5f", (rand.Float64() * 10000)),
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	latency := time.Since(startTime)
	if err != nil {
		resultsChan <- RequestResult{Latency: latency, Successful: false, Endpoint: "/accounts", Error: err}
		return
	}
	defer resp.Body.Close()

	resultsChan <- RequestResult{
		Latency:    latency,
		StatusCode: resp.StatusCode,
		Successful: resp.StatusCode == http.StatusCreated,
		Endpoint:   "/accounts",
	}
}

func performTransfers(accountIDs []int, resultsChan chan<- RequestResult) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	for i := 0; i < numTransfers; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()

			fromAccount := accountIDs[i%len(accountIDs)]
			toAccount := accountIDs[(i+1)%len(accountIDs)]

			if fromAccount == toAccount {
				toAccount = accountIDs[(i+2)%len(accountIDs)]
			}

			transferFunds(fromAccount, toAccount, "10", resultsChan)
			<-sem
		}(i)
	}

	wg.Wait()
}

func transferFunds(from, to int, amount string, resultsChan chan<- RequestResult) {
	startTime := time.Now()
	url := fmt.Sprintf("%s/transactions", baseURL)
	payload := map[string]interface{}{
		"source_account_id":      from,
		"destination_account_id": to,
		"amount":                 amount,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	latency := time.Since(startTime)
	if err != nil {
		resultsChan <- RequestResult{Latency: latency, Successful: false, Endpoint: "/transactions", Error: err}
		return
	}
	defer resp.Body.Close()

	resultsChan <- RequestResult{
		Latency:    latency,
		StatusCode: resp.StatusCode,
		Successful: resp.StatusCode == http.StatusCreated,
		Endpoint:   "/transactions",
	}
}

func printReport(results []RequestResult, totalTime time.Duration) {
	endpointResults := make(map[string][]RequestResult)
	for _, r := range results {
		endpointResults[r.Endpoint] = append(endpointResults[r.Endpoint], r)
	}

	totalRequests := len(results)
	successfulRequests := 0
	for _, r := range results {
		if r.Successful {
			successfulRequests++
		}
	}

	fmt.Println("\n--- Overall Performance ---")
	fmt.Printf("Total Time Taken: %v\n", totalTime.Round(time.Millisecond))
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful Requests: %d\n", successfulRequests)
	fmt.Printf("Failed Requests: %d\n", totalRequests-successfulRequests)
	if totalRequests > 0 {
		fmt.Printf("Success Rate: %.2f%%\n", float64(successfulRequests)/float64(totalRequests)*100)
	}
	fmt.Printf("Requests Per Second (RPS): %.2f\n", float64(totalRequests)/totalTime.Seconds())

	for endpoint, res := range endpointResults {
		fmt.Printf("\n--- Endpoint: %s ---\n", endpoint)

		if len(res) == 0 {
			fmt.Println("No requests made to this endpoint.")
			continue
		}

		latencies := make([]time.Duration, len(res))
		var totalLatency time.Duration
		for i, r := range res {
			latencies[i] = r.Latency
			totalLatency += r.Latency
		}

		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

		avgLatency := totalLatency / time.Duration(len(res))
		minLatency := latencies[0]
		maxLatency := latencies[len(latencies)-1]
		p50Latency := percentile(latencies, 50)
		p90Latency := percentile(latencies, 90)
		p95Latency := percentile(latencies, 95)

		fmt.Printf("Total Requests: %d\n", len(res))
		fmt.Printf("Average Latency: %v\n", avgLatency.Round(time.Microsecond))
		fmt.Printf("Min Latency: %v\n", minLatency.Round(time.Microsecond))
		fmt.Printf("Max Latency: %v\n", maxLatency.Round(time.Microsecond))
		fmt.Printf("P50 Latency: %v\n", p50Latency.Round(time.Microsecond))
		fmt.Printf("P90 Latency: %v\n", p90Latency.Round(time.Microsecond))
		fmt.Printf("P95 Latency: %v\n", p95Latency.Round(time.Microsecond))
	}

	for _, r := range results {
		if !r.Successful {
			fmt.Println("\n--- Error Details ---")
			fmt.Printf("Endpoint: %s\n", r.Endpoint)
			fmt.Printf("Status Code: %d\n", r.StatusCode)
			fmt.Printf("Error: %v\n", r.Error)
		}
	}
}

func percentile(latencies []time.Duration, p int) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	index := (len(latencies) * p) / 100
	if index >= len(latencies) {
		index = len(latencies) - 1
	}
	return latencies[index]
}
