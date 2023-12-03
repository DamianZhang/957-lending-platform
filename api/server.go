package api

import (
	"fmt"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/token"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/gofiber/fiber/v2"
)

// Server serves HTTP requests for our lending platform
type Server struct {
	config          util.Config
	app             *fiber.App
	borrowerService service.BorrowerService
	tokenMaker      token.Maker
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config util.Config, borrowerService service.BorrowerService) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		borrowerService: borrowerService,
		tokenMaker:      tokenMaker,
	}

	server.setUpRoutes()
	return server, nil
}

func (server *Server) setUpRoutes() {
	app := fiber.New()

	borrowerHandler := NewBorrowerHandler(server.config, server.borrowerService, server.tokenMaker)
	borrowerHandler.Route(app)

	server.app = app
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.app.Listen(address)
}
