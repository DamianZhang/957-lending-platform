package api

import (
	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/gofiber/fiber/v2"
)

// Server serves HTTP requests for our lending platform
type Server struct {
	borrowerService service.BorrowerService
	app             *fiber.App
}

// NewServer creates a new HTTP server and set up routing
func NewServer(borrowerService service.BorrowerService) (*Server, error) {
	server := &Server{
		borrowerService: borrowerService,
	}

	server.setUpRoutes()
	return server, nil
}

func (server *Server) setUpRoutes() {
	app := fiber.New()

	borrowerHandler := NewBorrowerHandler(server.borrowerService)
	borrowerHandler.Route(app)

	server.app = app
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.app.Listen(address)
}

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}

	GeneralResponse struct {
		Data interface{} `json:"data"`
	}
)
