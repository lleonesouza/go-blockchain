package blockchain

import (
	"github.com/craton-api/chain/database"
)

type BlockChainIterator struct {
	CurrentHash []byte
	Db          *database.AdaptorsDatabase
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Db}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	encodedBlock, err := iter.Db.Get(iter.CurrentHash)
	block := Deserialize(encodedBlock)

	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
