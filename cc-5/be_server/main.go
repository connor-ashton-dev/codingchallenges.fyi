package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	args := os.Args
	server := gin.Default()
	PORT := ":" + args[2]

	server.GET("/ping", func(c *gin.Context) {
		delay := time.Duration(1 * time.Second)
		time.Sleep(delay)
		c.JSON(http.StatusOK, gin.H{
			"message": "hi from BE at port" + PORT,
		})
	})

	server.GET("/checkhealth", func(c *gin.Context) {
		// return with status 200 but no message
		c.JSON(http.StatusOK, gin.H{})
	})

	err := server.Run(PORT)
	if err != nil {
		log.Fatal("Error with server:", err)
	}
}
