package grocksdb

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestIterator(t *testing.T) {
	db := newTestDB(t, "TestIterator", nil)
	defer db.Close()

	// insert keys
	givenKeys := [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}
	wo := NewDefaultWriteOptions()
	for _, k := range givenKeys {
		ensure.Nil(t, db.Put(wo, k, []byte("val")))
	}

	ro := NewDefaultReadOptions()
	iter := db.NewIterator(ro)
	defer iter.Close()
	var actualKeys [][]byte
	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		key := make([]byte, 4)
		copy(key, iter.Key().Data())
		actualKeys = append(actualKeys, key)
	}
	ensure.Nil(t, iter.Err())
	ensure.DeepEqual(t, actualKeys, givenKeys)
}

func TestIteratorCF(t *testing.T) {
	db, cfs, cleanup := newTestDBMultiCF(t, "TestIteratorCF", []string{"default", "c1", "c2", "c3"}, nil)
	defer cleanup()

	// insert keys
	givenKeys := [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}
	wo := NewDefaultWriteOptions()
	for _, k := range givenKeys {
		for i := range cfs {
			ensure.Nil(t, db.PutCF(wo, cfs[i], k, []byte("val")))
		}
	}

	{
		ro := NewDefaultReadOptions()
		iter := db.NewIteratorCF(ro, cfs[0])
		defer iter.Close()
		var actualKeys [][]byte
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			key := make([]byte, 4)
			copy(key, iter.Key().Data())
			actualKeys = append(actualKeys, key)
		}
		ensure.Nil(t, iter.Err())
		ensure.DeepEqual(t, actualKeys, givenKeys)
	}

	{
		ro := NewDefaultReadOptions()
		iters, err := db.NewIterators(ro, cfs)
		ensure.Nil(t, err)
		ensure.DeepEqual(t, len(iters), 4)
		defer func() {
			for i := range iters {
				iters[i].Close()
			}
		}()

		for _, iter := range iters {
			var actualKeys [][]byte
			for iter.SeekToFirst(); iter.Valid(); iter.Next() {
				key := make([]byte, 4)
				copy(key, iter.Key().Data())
				actualKeys = append(actualKeys, key)
			}
			ensure.Nil(t, iter.Err())
			ensure.DeepEqual(t, actualKeys, givenKeys)
		}
	}
}
