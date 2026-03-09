package api

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port               int
	lastReceivedData   map[string]interface{}
	lastReceivedDataMu sync.RWMutex
}

func NewServer(port int) *Server {
	return &Server{
		port:             port,
		lastReceivedData: make(map[string]interface{}),
	}
}

func (s *Server) Start() error {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "rabbit-http-bridge-go",
		})
	})

	r.POST("/api/data1", s.handleData1Post)
	r.GET("/api/data1", s.handleData1Get)
	r.POST("/api/data2", s.handleData2Post)
	r.GET("/api/data2", s.handleData2Get)

	log.Printf("Starting API server on port %d...", s.port)
	return r.Run(":" + strconv.Itoa(s.port))
}

func (s *Server) handleData1Post(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("--- DATA RECEIVED AT HTTP1 ---")
	s.lastReceivedDataMu.Lock()
	s.lastReceivedData["data1"] = data
	s.lastReceivedDataMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "success", "received": data})
}

func (s *Server) handleData1Get(c *gin.Context) {
	s.lastReceivedDataMu.RLock()
	data, ok := s.lastReceivedData["data1"]
	s.lastReceivedDataMu.RUnlock()

	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "No data received yet for Queue 1. Push a message to RabbitMQ!"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (s *Server) handleData2Post(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("--- DATA RECEIVED AT HTTP2 ---")
	s.lastReceivedDataMu.Lock()
	s.lastReceivedData["data2"] = data
	s.lastReceivedDataMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "success", "received": data})
}

func (s *Server) handleData2Get(c *gin.Context) {
	s.lastReceivedDataMu.RLock()
	data, ok := s.lastReceivedData["data2"]
	s.lastReceivedDataMu.RUnlock()

	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "No data received yet for Queue 2. Push a message to RabbitMQ!"})
		return
	}
	c.JSON(http.StatusOK, data)
}
