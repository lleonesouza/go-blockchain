package server

import (
	"github.com/labstack/echo/v4"
	"github.com/craton-api/chain/server/routes"
)

type Server struct {
	Echo *echo.Echo
}

func MakeServer(port string) (*Server, error) {
	e := echo.New()

	// Create Routes
	e, err := routes.MakeRoutes(e)
	if err != nil {
		return nil, err
	}

	// Start Server
	err = e.Start(port)
	if err != nil {
		return nil, err
	}

	return &Server{Echo: e}, nil
}
