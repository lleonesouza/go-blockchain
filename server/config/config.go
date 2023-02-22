package config

import "os"

type Env struct {
	WALLET_URL   string
	NODE_ID      string
	PORT         string
	GENESIS_DATA string
}

func MakeEnv() Env {
	var wallet_url string
	var node_id string
	var port string
	var genesis_data string

	if os.Getenv("WALLET_URL") != "" {
		wallet_url = os.Getenv("WALLET_URL")
	} else {
		wallet_url = "http://localhost:5000/v1"
	}
	if os.Getenv("NODE_ID") != "" {
		node_id = os.Getenv("NODE_ID")
	} else {
		node_id = "00000node_id"
	}
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "5001"
	}
	if os.Getenv("GENESIS_DATA") != "" {
		genesis_data = os.Getenv("GENESIS_DATA")
	} else {
		genesis_data = "g3n3s4s d4t4"
	}

	return Env{
		WALLET_URL:   wallet_url,
		NODE_ID:      node_id,
		PORT:         port,
		GENESIS_DATA: genesis_data,
	}
}
