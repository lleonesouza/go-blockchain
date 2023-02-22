package blockchain

import (
	"bytes"
	"encoding/hex"

	"github.com/craton-api/chain/database"
)

var (
	utxoPrefix = []byte("utxo-")
)

type UTXOSet struct {
	Blockchain *BlockChain
	Db         *database.AdaptorsDatabase
}

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0

	items := u.Db.IterateByPrefix(utxoPrefix)

	for _, item := range items {
		k := bytes.TrimPrefix(item.Key, utxoPrefix)
		txID := hex.EncodeToString(k)
		outs := DeserializeOutputs(item.Value)

		for outIdx, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
			}
		}

	}

	return accumulated, unspentOuts
}

func (u UTXOSet) FindUnspentTransactions(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	items := u.Db.IterateByPrefix(utxoPrefix)

	for _, item := range items {
		outs := DeserializeOutputs(item.Value)
		for _, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (u UTXOSet) CountTransactions() int {
	counter := u.Db.CountByPrefix(utxoPrefix)

	return counter
}

func (u UTXOSet) Reindexes() {
	u.DeleteByPrefix(utxoPrefix)
	UTXO := u.Blockchain.FindUTXO()

	for txId, outs := range UTXO {
		key, err := hex.DecodeString(txId)
		Handle(err)
		key = append(utxoPrefix, key...)

		err = u.Db.Set(key, outs.Serialize())
		Handle(err)
	}
}

func (u *UTXOSet) Update(block *Block) error {
	updateInputTx := func(tx *Transaction) error {
		for _, txInput := range tx.Inputs {
			updatedOuts := TxOutputs{}

			inputId := append(utxoPrefix, txInput.ID...)
			item, err := u.Db.Get(inputId)
			if err != nil {
				return err
			}

			outs := DeserializeOutputs(item)
			if err != nil {
				return err
			}

			for index, out := range outs.Outputs {
				if index != txInput.Out {
					updatedOuts.Outputs = append(updatedOuts.Outputs, out)
				}
			}

			if len(updatedOuts.Outputs) == 0 {
				u.Db.Delete(inputId)
			} else {
				a := updatedOuts.Serialize()
				u.Db.Set(inputId, a)
			}
		}
		return nil
	}

	for _, tx := range block.Transactions {
		if !tx.IsCoinbase() {
			err := updateInputTx(tx)
			if err != nil {
				return err
			}
		}

		newOutputs := TxOutputs{}

		newOutputs.Outputs = append(newOutputs.Outputs, tx.Outputs...)
		txID := append(utxoPrefix, tx.ID...)

		bytes := newOutputs.Serialize()
		u.Db.Set(txID, bytes)
	}
	return nil
}

func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	u.Db.DeleteByPrefix(prefix)
}
