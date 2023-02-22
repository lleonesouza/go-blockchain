package database

import (
	"fmt"
	"log"
)

// DBExists returns true or false based if theres the db file
type DBExists func() (bool, error)

// Get receive the key, returns the value
type Get func([]byte) ([]byte, error)

// Set receive the key and the value, returns the value
type Set func(key []byte, value []byte) error

// Close finish the database connection
type Close func()

// CountByPrefix get values by prefix
type CountByPrefix func(prefix []byte) int

// DeleteByPrefix
type DeleteByPrefix func(prefix []byte)

// Delete by key bytes
type Delete func(prefix []byte)

type IterateByPrefix func(prefix []byte) []Item

// AdaptorsDatabase Struct
type AdaptorsDatabase struct {
	// DeleteByPrefix
	DeleteByPrefix
	// DBExists returns true or false based if theres the db file
	DBExists
	// Get receive the key, returns the value
	Get
	// Set receive the key and the value, returns the value
	Set
	// Close finish the database connection
	Close
	// AdaptorsDatabase Struct
	CountByPrefix
	// Delete deletes by key
	Delete
	// IterateByPrefix
	IterateByPrefix
}

// MakeModule generate a new database adaptors
func Make(path string) (*AdaptorsDatabase, error) {
	db, err := makeDatabase(path)
	if err != nil {
		return nil, err
	}

	adaptors := AdaptorsDatabase{
		DBExists: func() (bool, error) {
			exists, err := db.dbExists(path)
			if err != nil {
				fmt.Println(err)
				return false, err
			}
			return exists, nil
		},
		Get: func(data []byte) ([]byte, error) {
			return db.get(data)
		},
		Set: func(key []byte, value []byte) error {
			return db.set(key, value)
		},
		Close: func() {
			db.close()
		},
		CountByPrefix: func(prefix []byte) int {
			return db.countByPrefix(prefix)
		},
		DeleteByPrefix: func(prefix []byte) {
			db.deleteByPrefix(prefix)
		},
		Delete: func(key []byte) {
			db.delete(key)
		},
		IterateByPrefix: func(prefix []byte) []Item {
			return db.iterateByPrefix(prefix)
		},
	}

	return &adaptors, nil
}

func handleErr(err error, functionName string) {
	if (err) != nil {
		fmt.Printf("Module Name: %s \n", "database")
		fmt.Printf("Function Name: %s \n", functionName)
		fmt.Println(err.Error())
		log.Panic(err)
	}
}
