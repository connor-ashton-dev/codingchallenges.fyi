package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	PORT := ":1337"

	lb := NewLB()
	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		handlePing(lb, c)
	})

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				performHealthChecks(lb)
			}
		}
	}()

	err := server.Run(PORT)

	if err != nil {
		log.Fatal("Error with server:", err)
	}
}
