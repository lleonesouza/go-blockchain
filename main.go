package main

import (
	"fmt"

	"github.com/craton-api/chain/server"
	"github.com/craton-api/chain/server/config"
)

func main() {
	env := config.MakeEnv()
	port := env.PORT
	fmt.Println("\n PORT: ", env.PORT)
	port = fmt.Sprintf(":%s", env.PORT)
	_, err := server.MakeServer(port)
	if err != nil {
		panic(err)
	}

}
