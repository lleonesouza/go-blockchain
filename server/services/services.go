package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/craton-api/chain/blockchain"
	"github.com/craton-api/chain/database"
	"github.com/craton-api/chain/network"
	"github.com/craton-api/chain/server/config"
	"github.com/craton-api/chain/server/utils"
)

type Service struct {
	nodeId string
	Db     *database.AdaptorsDatabase
	Chain  *blockchain.BlockChain
	Wallet *Wallet
}

const (
	dbPath = "./tmp/blocks_%s"
)

func MakeServices() (*Service, error) {
	env := config.MakeEnv()
	nodeId := env.NODE_ID
	var chain *blockchain.BlockChain
	var err error
	path := fmt.Sprintf(dbPath, nodeId)
	db, err := database.Make(path)
	if err != nil {
		return nil, err
	}
	// exists, err := db.DBExists()
	// if err != nil {
	// 	return nil, err
	// }

	w := MakeWallet()

	_, err = db.Get([]byte("lh"))
	if err == nil {
		fmt.Println("\n\n\n Continue ")
		chain = blockchain.ContinueBlockChain(nodeId, db)
	} else {
		fmt.Println("\n\n\n Creating ")
		walletAddress, err := w.Create()
		if err != nil {
			return nil, err
		}
		fmt.Println("firstAddress : ", walletAddress)
		chain, err = blockchain.InitBlockChain(walletAddress, nodeId, db)
		if err != nil {
			return nil, err
		}
	}

	service := &Service{nodeId, db, chain, w}

	service.ReindexesUTXO()

	return service, nil
}

func (service *Service) StartNode(minerAddress string) {
	env := config.MakeEnv()
	nodeId := env.NODE_ID
	path := fmt.Sprintf(dbPath, nodeId)
	db, _ := database.Make(path)

	fmt.Printf("Starting Node %s\n", service.nodeId)

	if len(minerAddress) > 0 {
		if blockchain.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	network.StartServer(service.nodeId, minerAddress, db)
}

func (service *Service) ReindexesUTXO() {
	UTXOSet := blockchain.UTXOSet{service.Chain, service.Db}
	UTXOSet.Reindexes()

	count := UTXOSet.CountTransactions()
	fmt.Printf("ReindexesUTXO Done! There are %d transactions in the UTXO set.\n", count)
}

func (service *Service) ListAddresses() ([]string, error) {
	addresses, err := service.Wallet.GetWallets()
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (service *Service) CreateWallet() (string, error) {
	address, err := service.Wallet.Create()
	if err != nil {
		return "", err
	}

	return address, nil
}

func (service *Service) PrintChain() {
	iter := service.Chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (service *Service) GetBalance(address string) int {
	if !blockchain.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}
	// chain := blockchain.ContinueBlockChain(service.nodeId)
	UTXOSet := blockchain.UTXOSet{service.Chain, service.Db}
	// defer chain.Database.Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	return balance
}

func (service *Service) Send(from, to string, amount int, mineNow bool) {
	var txs []*blockchain.Transaction
	if !blockchain.ValidateAddress(to) {
		log.Panic("Address is not Valid")
	}
	if !blockchain.ValidateAddress(from) {
		log.Panic("Address is not Valid")
	}

	UTXOSet := blockchain.UTXOSet{service.Chain, service.Db}
	// defer chain.Database.Close()

	w := service.Wallet.GetWallet(from)

	tx := blockchain.NewTransaction(w, to, amount, &UTXOSet)

	if mineNow {
		txs = append(txs, tx)
		// cbTx := blockchain.CoinbaseTx(from, "")
		// txs := []*blockchain.Transaction{cbTx, tx}
		block := service.Chain.MineBlock(txs)

		UTXOSet.Update(block)
	} else {
		network.SendTx(network.KnownNodes[0], tx)
		fmt.Println("send tx")
	}

	fmt.Println("Success!")
}

func (service *Service) Mine(txs []*blockchain.Transaction) *blockchain.Block {
	block := service.Chain.MineBlock(txs)

	UTXOSet := blockchain.UTXOSet{service.Chain, service.Db}

	UTXOSet.Update(block)

	return block
}
