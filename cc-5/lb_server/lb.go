package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

type LoadBalancer struct {
	Servers         []int
	CheckServers    []int
	CurrServer      int32 // change to int32 for atomic operations
	CurrCheckServer int
	mu              sync.Mutex
}

func NewLB() *LoadBalancer {
	lb := LoadBalancer{}
	lb.Servers = []int{3000, 3001, 3002, 3003}
	lb.CheckServers = lb.Servers
	lb.CurrServer = 0
	lb.CurrCheckServer = 0
	lb.mu = sync.Mutex{}
	return &lb
}

func (lb *LoadBalancer) SwitchCheckServer() {
	lb.CurrCheckServer = (lb.CurrCheckServer + 1) % len(lb.Servers)
}

func (lb *LoadBalancer) SwitchServer() {
	atomic.AddInt32(&lb.CurrServer, 1) // safely increment CurrServer
	if atomic.LoadInt32(&lb.CurrServer) >= int32(len(lb.Servers)) {
		atomic.StoreInt32(&lb.CurrServer, 0) // reset to 0 if it exceeds the slice length
	}
}

func checkHealth(server int) int {
	url := fmt.Sprintf("http://localhost:%d/checkhealth", server)
	resp, err := http.Get(url) // Simplified to http.Get for a GET request with no headers or body
	if err != nil {
		// Log the error and return a status indicating that the health check couldn't be performed
		fmt.Printf("Error performing health check on server %d: %s\n", server, err)
		return http.StatusServiceUnavailable // Indicates that the health check couldn't be completed
	}
	defer resp.Body.Close() // Ensure we close the body to avoid leaking resources

	// You may want to read the body if your health check endpoint returns useful information.
	// For now, we're just using the status code as an indication of health.

	return resp.StatusCode
}

// Perform health checks on all servers in parallel
func performHealthChecks(lb *LoadBalancer) {
	var wg sync.WaitGroup

	for i := range lb.CheckServers {
		wg.Add(1)
		go func(serverIndex int) {
			defer wg.Done()

			server := lb.CheckServers[serverIndex]
			status := checkHealth(server)
			if status != http.StatusOK {
				fmt.Println("Server is unhealthy:", server)
				removeServer(lb, server)
			} else {
				fmt.Println("Server is healthy:", server)
				addServer(lb, server)
			}
		}(i)
	}

	wg.Wait() // Wait for all health checks to complete
	fmt.Println("-----NEW CYCLE-----")
}

// Safely remove a server from the list
func removeServer(lb *LoadBalancer, server int) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, srv := range lb.Servers {
		if srv == server {
			lb.Servers = append(lb.Servers[:i], lb.Servers[i+1:]...)
			break
		}
	}
}

// Safely add a server to the list
func addServer(lb *LoadBalancer, server int) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	found := false
	for _, srv := range lb.Servers {
		if srv == server {
			found = true
			break
		}
	}
	if !found {
		lb.Servers = append(lb.Servers, server)
	}
}
