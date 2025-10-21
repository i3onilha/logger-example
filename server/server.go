package server

import (
	"logger-example/controller"
	appconfig "logger-example/config"

	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"go.uber.org/fx"
)

type Server struct {
	router *gin.Engine
	port   string
}

func NewServer(cfg *appconfig.Config, userController *controller.UserController) *Server {
	r := gin.Default()
	r.Use(gintrace.Middleware(cfg.DataDog.ServiceName))

	// Setup routes
	r.GET("/user/:id", userController.GetUser)

	return &Server{
		router: r,
		port:   cfg.Server.Port,
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.port)
}

// FX Module for server
var Module = fx.Module("server",
	fx.Provide(NewServer),
)

// Config interface for dependency injection
type Config interface {
	GetServerPort() string
	GetDataDogServiceName() string
}

// Simple config struct that implements the interface
type serverConfig struct {
	port        string
	serviceName string
}

func (c *serverConfig) GetServerPort() string {
	return c.port
}

func (c *serverConfig) GetDataDogServiceName() string {
	return c.serviceName
}

// NewServerConfig creates a config from the main config
func NewServerConfig(cfg *appconfig.Config) Config {
	return &serverConfig{
		port:        cfg.Server.Port,
		serviceName: cfg.DataDog.ServiceName,
	}
}
