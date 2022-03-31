package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hortelanobruno/foaas-api/domain/service/handler"
	"github.com/hortelanobruno/foaas-api/middleware"
	"github.com/hortelanobruno/foaas-api/ratelimiter"
)

type Server struct {
	messageHandler *handler.MessageHandler
	rateLimiter    ratelimiter.RateLimiter
}

func NewServer(messageHandler *handler.MessageHandler, rateLimiter ratelimiter.RateLimiter) *Server {
	return &Server{
		messageHandler: messageHandler,
		rateLimiter:    rateLimiter,
	}
}

func (s *Server) Start(port int) {
	engine := gin.Default()

	if s.rateLimiter != nil {
		engine.Use(middleware.RateLimiter(s.rateLimiter))
	}

	s.attachEndpoints(engine)
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}

func (s *Server) attachEndpoints(engine *gin.Engine) {
	engine.GET("/message", s.messageHandler.HandleGetMessage)
}
