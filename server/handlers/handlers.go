package handlers

import (
	"fmt"
	"net/http"

	service "github.com/craton-api/chain/server/services"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	services *service.Service
}

func MakeHandlers() (*Handlers, error) {
	services, err := service.MakeServices()
	if err != nil {
		return nil, err
	}

	return &Handlers{services}, nil
}

type TransactionCreateInput struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Price string `json:"price"`
}

func (h *Handlers) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK!")
}

func (h *Handlers) ReindexesUTXO(c echo.Context) error {
	h.services.ReindexesUTXO()

	return c.JSON(http.StatusOK, "OK")
}

func (h *Handlers) GetBlocks(c echo.Context) error {
	h.services.PrintChain()

	return c.JSON(http.StatusOK, "OK")
}

func (h *Handlers) GetBalance(c echo.Context) error {
	address := c.Param("address")
	balance := 0

	response := fmt.Sprintf("Balance of %s: %d\n", address, balance)

	h.services.GetBalance(address)

	return c.String(http.StatusOK, response)
}

func (h *Handlers) Send(c echo.Context) error {
	firstAddress := "16LB24uYqoXXFtgvq7khkS3QkxLKu6xgGc"

	h.services.Send(firstAddress, "to", 50, false)

	return c.JSON(http.StatusOK, "OK")
}

func (h *Handlers) StartNode(c echo.Context) error {
	address := ""
	h.services.StartNode(address)

	return c.JSON(http.StatusOK, "OK")
}
