package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Serajian/GO-long-polling-gin/longpolling"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	pm := longpolling.NewPollingManager()

	// Long-poll endpoint
	r.GET("/poll/:id", Getter(pm))

	// Send a message to a specific client ID
	r.POST("/send/:id", Sender(pm))

	// Run server on port 8090
	_ = r.Run(":8090")
}

func Getter(pm *longpolling.PollingManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Param("id")
		ch := pm.RegisterClient(clientID)

		fmt.Printf("1*** client %s is registered and ready for get msg.\n", clientID)

		select {
		case msg := <-ch:
			c.JSON(http.StatusOK, gin.H{"message": msg})
			fmt.Printf("3*** client %s get this msg: %s.\n", clientID, msg)
		case <-time.After(30 * time.Second):
			pm.RemoveClient(clientID)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout"})
		case <-c.Request.Context().Done(): // client closed connection
			pm.RemoveClient(clientID)
			// don’t write response because client is gone
		}
	}
}

func Sender(pm *longpolling.PollingManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Param("id")
		var json struct {
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		pm.SendMessage(clientID, json.Message)

		fmt.Printf("2*** send this msg: '%s' to client by this id: %s.\n", json.Message, clientID)

		c.JSON(http.StatusOK, gin.H{"status": "message sent"})
	}
}
