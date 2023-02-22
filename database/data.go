package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgraph-io/badger"
)

type database struct {
	badger *badger.DB
}

type Item struct {
	Key   []byte
	Value []byte
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

func makeDatabase(path string) (*database, error) {
	opts := badger.DefaultOptions

	opts.ValueDir = path
	opts.Dir = path
	opts.ValueDir = path

	db, err := badger.Open(opts)
	if err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(path, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				database := database{db}
				return &database, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	}

	fmt.Printf("\nDB Path: %s\n", path)

	database := database{db}

	return &database, nil
}

func (db *database) dbExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (db *database) DbExists(path string) (bool, error) {
	_, err := os.Stat(path + "/MANIFEST")
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	exists := os.IsNotExist(err)
	if exists {
		return true, nil
	}
	return false, nil

}

func (db *database) set(key []byte, value []byte) error {
	err := db.badger.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err != nil {
			return err
		}
		err = txn.Set(key, value)
		return err
	})

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (db *database) get(key []byte) ([]byte, error) {
	var value []byte

	err := db.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		value, err = item.Value()

		return err
	})

	if err != nil {

		return nil, err
	} else {
		return value, nil
	}
}

func (db *database) deleteByPrefix(prefix []byte) {
	deleteKeys := func(keysForDelete [][]byte) error {

		badgerTxn := func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}

			return nil
		}

		if err := db.badger.Update(badgerTxn); err != nil {
			return err
		}

		return nil
	}

	collectSize := 100000
	db.badger.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					log.Panic(err)
				}
				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}
		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

func (db *database) iterateByPrefix(prefix []byte) []Item {
	var items []Item

	err := db.badger.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(k)
			handleErr(err, "data.go:145")

			items = append(items, Item{k, v})
		}
		return nil
	})
	handleErr(err, "data.go:151")

	return items
}

func (db *database) countByPrefix(prefix []byte) int {
	counter := 0
	err := db.badger.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			counter++
		}

		return nil
	})

	handleErr(err, "data.go:170")
	return counter
}

// func retry() {

// }

// func openDB() {

// }

func (db *database) delete(key []byte) {
}

// func (db *database) update(key []byte) {
// }

func (db *database) close() {
	db.badger.Close()
}
