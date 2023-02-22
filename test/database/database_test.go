package database

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/craton-api/chain/database"
)

func TestDatabase(t *testing.T) {
	path := "./t_database0"
	prefix := "test"
	iterator_times := 10

	// Create Database
	db, err := database.Make(path)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	// TESTING: SET AND GET
	// loop creating and reading keys and values
	for i := 0; i < iterator_times; i++ {
		key := fmt.Sprintf("%s-%d", prefix, i)
		value := fmt.Sprintf("valuetesting%d", i)

		err = db.Set([]byte(key), []byte(value))

		if err != nil {
			t.Errorf("error: %s", err)
		}

		res, err := db.Get([]byte(key))
		if err != nil {
			t.Errorf("error: %s", err)
		}

		// Compare bytes
		if !bytes.Equal(res, []byte(value)) {
			t.Errorf("'%s' is not equal '%s'", res, value)
		}

	}

	// TESTING: IteratorByPrefix
	items := db.IterateByPrefix([]byte(prefix))
	t.Logf("ItemsByPrefix: %s\n\n", items)

	// TESTING: CountByPrefix
	count := db.CountByPrefix([]byte(prefix))
	t.Logf("CountByPrefix: %d\n\n", count)

	// TESTING: Delete

	// TESTING: DeleteByPrefix
	// set, get and compare
	for i := 0; i < iterator_times; i++ {
		key := fmt.Sprintf("%s-%d", prefix, i)
		value := fmt.Sprintf("valueToBeDeleted%d", i)

		err = db.Set([]byte(key), []byte(value))

		if err != nil {
			t.Errorf("error: %s", err)
		}

		res, err := db.Get([]byte(key))
		if err != nil {
			t.Errorf("error: %s", err)
		}

		// Compare bytes
		if !bytes.Equal(res, []byte(value)) {
			t.Errorf("'%s' is not equal '%s'", res, value)
		}
	}

	// delete
	db.DeleteByPrefix([]byte(prefix))

	// get and check if is deleted
	for i := 0; i < iterator_times; i++ {
		key := fmt.Sprintf("%s-%d", prefix, i)

		res, err := db.Get([]byte(key))

		if err == nil {
			t.Errorf("item not deleted: %s", res)
		}
	}

	// TESTING: CountByPrefix
	count = db.CountByPrefix([]byte(prefix))
	t.Logf("CountByPrefix: %d\n\n", count)
}
