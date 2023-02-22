package tests

import (
	"testing"

	service "github.com/craton-api/chain/server/services"
)

func TestBlockchain(t *testing.T) {
	s, err := service.MakeServices()
	if err != nil {
		t.Error("Error : ", err)
	}

	s.ReindexesUTXO()

	// First Address
	// addresses, err := s.ListAddresses()
	// if err != nil {
	// 	t.Error("Error : ", err)
	// }
	a := "1CwbKruzuJmVaJea5VWMoF8GN5PgLXdjmC"

	balance := s.GetBalance(a)
	t.Logf("Address1 balance: %d", balance)

	// Create a second wallet
	address2, err := s.CreateWallet()
	if err != nil {
		t.Error("Error : ", err)
	}

	balance = s.GetBalance(address2)
	t.Logf("Address2 balance: %d", balance)

	// Print the Chain
	s.PrintChain()

	// Send to his address
	s.Send(a, address2, 1, true)

}
