package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ساختار برای مدیریت درخواست‌های در حال انتظار
type PollingManager struct {
	mu      sync.Mutex
	clients map[string]chan string
}

func NewPollingManager() *PollingManager {
	return &PollingManager{
		clients: make(map[string]chan string),
	}
}

// ثبت یک کلاینت جدید برای دریافت پیام
func (pm *PollingManager) RegisterClient(clientID string) chan string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	ch := make(chan string, 1)
	pm.clients[clientID] = ch
	return ch
}

// ارسال پیام به کلاینت خاص
func (pm *PollingManager) SendMessage(clientID, message string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if ch, exists := pm.clients[clientID]; exists {
		ch <- message
		delete(pm.clients, clientID) // بعد از ارسال پیام، کلاینت را حذف می‌کنیم
	}
}

// حذف کلاینت‌ها در صورت تایم‌اوت
func (pm *PollingManager) RemoveClient(clientID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.clients, clientID)
}

func main() {
	r := gin.Default()
	pm := NewPollingManager()
	// Endpoint برای دریافت پیام‌های Long-Polling
	r.GET("/poll/:id", func(c *gin.Context) {
		clientID := c.Param("id")
		// ثبت کلاینت جدید
		ch := pm.RegisterClient(clientID)

		// استفاده از select برای منتظر ماندن تا دریافت پیام یا تایم‌اوت
		select {
		case msg := <-ch:
			c.JSON(http.StatusOK, gin.H{"message": msg})
			fmt.Println("68:", msg)
		case <-time.After(30 * time.Second): // بعد از ۳۰ ثانیه تایم‌اوت می‌شود
			pm.RemoveClient(clientID)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout"})
		}
	})

	// Endpoint برای ارسال پیام به کلاینت خاص
	r.POST("/send/:id", func(c *gin.Context) {
		clientID := c.Param("id")
		var json struct {
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		pm.SendMessage(clientID, json.Message)
		c.JSON(http.StatusOK, gin.H{"status": "message sent"})
	})

	// اجرای سرور روی پورت 8090
	r.Run(":8090")
}
