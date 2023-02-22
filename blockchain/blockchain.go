package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/craton-api/chain/database"
)

const (
	dbPath      = "./tmp/blocks_%s"
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Db       *database.AdaptorsDatabase
}

func ContinueBlockChain(nodeId string, db *database.AdaptorsDatabase) *BlockChain {
	lastHash, err := db.Get([]byte("lh"))
	Handle(err)

	chain := BlockChain{lastHash, db}

	return &chain
}

func InitBlockChain(address, nodeId string, db *database.AdaptorsDatabase) (*BlockChain, error) {

	cbtx := CoinbaseTx(address, genesisData)
	genesis := Genesis(cbtx)
	fmt.Println("Genesis created")

	err := db.Set(genesis.Hash, genesis.Serialize())
	if err != nil {
		return nil, err
	}

	err = db.Set([]byte("lh"), genesis.Hash)
	if err != nil {
		return nil, err
	}

	lastHash := genesis.Hash

	blockchain := BlockChain{lastHash, db}
	return &blockchain, nil
}

func (chain *BlockChain) GetLastHash() ([]byte, error) {
	lasthash, err := chain.Db.Get([]byte("lh"))
	if err != nil {
		return nil, err
	}
	return lasthash, nil
}

func (chain *BlockChain) SetLastHash(data []byte) error {
	err := chain.Db.Set([]byte("lh"), data)
	if err != nil {
		return err
	}
	chain.LastHash = data
	return nil
}

func (chain *BlockChain) AddBlock(block *Block) {
	_, err := chain.Db.Get(block.Hash)
	if err == nil {
		return // block already exists
	}

	err = chain.Db.Set(block.Hash, block.Serialize())
	Handle(err)

	lastBlockData, err := chain.Db.Get([]byte("lh"))
	Handle(err)
	lastBlock := Deserialize(lastBlockData)

	if block.Height > lastBlock.Height {
		err := chain.Db.Set([]byte("lh"), block.Hash)
		Handle(err)
		chain.LastHash = block.Hash
	}
}

func (chain *BlockChain) GetBestHeight() (int, error) {
	lastBlockData, err := chain.Db.Get(chain.LastHash)
	if err != nil {
		return 0, err
	}

	lastBlock := *Deserialize(lastBlockData)
	return lastBlock.Height, nil

}

func (chain *BlockChain) GetBlock(blockHash []byte) (Block, error) {
	blockData, err := chain.Db.Get(blockHash)
	if err != nil {
		return Block{}, err
	}
	block := *Deserialize(blockData)

	return block, nil
}

func (chain *BlockChain) GetBlockHashes() [][]byte {
	var blocks [][]byte

	iter := chain.Iterator()

	for {
		block := iter.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}

func (chain *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if !chain.VerifyTransaction(tx) {
			log.Panic("Invalid Transaction")
		}
	}

	lastBlockData, err := chain.Db.Get(chain.LastHash)
	Handle(err)
	lastBlock := Deserialize(lastBlockData)
	lastHeight = lastBlock.Height
	lastHash = lastBlock.Hash

	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)

	err = chain.Db.Set(newBlock.Hash, newBlock.Serialize())
	Handle(err)
	err = chain.Db.Set([]byte("lh"), newBlock.Hash)
	Handle(err)

	chain.LastHash = newBlock.Hash
	return newBlock
}

func (chain *BlockChain) FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return UTXO
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
