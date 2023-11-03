package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"` // Make sure the field is exported by capitalizing the first letter
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100, // Set MaxIdleConns to a reasonably high number
		MaxIdleConnsPerHost: 10,  // Set this to the number of your servers or whatever suits your load
		IdleConnTimeout:     90 * time.Second,
	},
}

func handlePing(lb *LoadBalancer, c *gin.Context) {
	currentServerIndex := atomic.LoadInt32(&lb.CurrServer)
	atomic.CompareAndSwapInt32(&lb.CurrServer, currentServerIndex, (currentServerIndex+1)%int32(len(lb.Servers)))

	currentServer := lb.Servers[currentServerIndex]

	url := fmt.Sprintf("http://localhost:%d/ping", currentServer)
	resp, err := httpClient.Get(url) // Use the global client to make the request
	if err != nil {
		fmt.Println("Error with request to BE server:", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "BE server not responding"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Response code not OK:", resp.StatusCode)
		c.JSON(resp.StatusCode, gin.H{"error": "BE server returned non-OK status"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var myResponse Response
	err = json.Unmarshal(body, &myResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response from BE server": myResponse.Message,
	})
}
