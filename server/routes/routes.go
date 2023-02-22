package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/craton-api/chain/server/handlers"
)

func MakeRoutes(e *echo.Echo) (*echo.Echo, error) {
	h, err := handlers.MakeHandlers()
	if err != nil {
		return nil, err
	}

	e.GET("/health", h.Health)
	// e.GET("/chain", h.GetBlocks)

	e.GET("/balance/:address", h.GetBalance)

	e.GET("/reindexes", h.ReindexesUTXO)

	e.GET("/blocks", h.GetBlocks)

	e.POST("/send", h.Send)

	return e, nil
}
