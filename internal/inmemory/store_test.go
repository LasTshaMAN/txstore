package inmemory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreNoTransactions(t *testing.T) {
	store := NewStore()

	doSetAndCheck(t, store, "some key", "some value")
}

func TestStoreWithTransactions(t *testing.T) {
	const (
		key1   = "key1"
		key2   = "key2"
		value1 = "value1"
		value2 = "value2"
		value3 = "value3"
	)

	t.Run("invalid commit has no effect", func(t *testing.T) {
		store := NewStore()

		store.Commit()
	})
	t.Run("invalid rollback has no effect", func(t *testing.T) {
		store := NewStore()

		store.Rollback()
	})
	t.Run("simple commit", func(t *testing.T) {
		store := NewStore()

		store.Set(key1, value1)

		store.Begin()
		doSetAndCheck(t, store, key2, value2)
		store.Commit()

		got := store.Get(key1)
		assert.Equal(t, value1, got)
		got = store.Get(key2)
		assert.Equal(t, value2, got)
	})
	t.Run("simple rollback", func(t *testing.T) {
		store := NewStore()

		store.Set(key1, value1)
		store.Begin()
		store.Rollback()

		got := store.Get(key1)
		assert.Equal(t, value1, got)
		got = store.Get(key2)
		assert.Equal(t, emptyValue, got)
	})
	t.Run("scenario with multiple transactions", func(t *testing.T) {
		store := NewStore()

		store.Begin()
		store.Set(key1, value1)
		store.Rollback()

		got := store.Get(key1)
		assert.Equal(t, emptyValue, got)

		store.Begin()
		store.Set(key1, value1)
		store.Commit()

		got = store.Get(key1)
		assert.Equal(t, value1, got)

		store.Begin()
		store.Set(key2, value2)
		store.Commit()

		got = store.Get(key2)
		assert.Equal(t, value2, got)
	})
	t.Run("scenario with nested transactions", func(t *testing.T) {
		store := NewStore()

		store.Begin()
		store.Set(key1, value1)
		got := store.Get(key1)
		assert.Equal(t, value1, got)

		store.Begin()
		store.Set(key1, value2)
		got = store.Get(key1)
		assert.Equal(t, value2, got)

		store.Begin()
		store.Set(key1, value3)
		got = store.Get(key1)
		assert.Equal(t, value3, got)

		store.Rollback()
		got = store.Get(key1)
		assert.Equal(t, value2, got)

		store.Rollback()
		got = store.Get(key1)
		assert.Equal(t, value1, got)

		store.Rollback()
		got = store.Get(key1)
		assert.Equal(t, emptyValue, got)
	})
}

// doSetAndCheck sets key to value for store and performs common sanity checks.
func doSetAndCheck(t *testing.T, store *Store, key, value string) {
	got := store.Get("common checks - unknown key")
	assert.Empty(t, got)

	store.Set(key, value)
	got = store.Get(key)
	assert.Equal(t, value, got)
	gotCnt := store.Count(value)
	assert.Equal(t, 1, gotCnt)

	store.Delete(key)
	got = store.Get(key)
	assert.Empty(t, got)
	gotCnt = store.Count(value)
	assert.Zero(t, gotCnt)

	store.Set(key, value)
	got = store.Get(key)
	assert.Equal(t, value, got)
	gotCnt = store.Count(value)
	assert.Equal(t, 1, gotCnt)

	store.Set("common checks - key2", value)
	store.Set("common checks - key3", "common checks - another value")
	gotCnt = store.Count(value)
	assert.Equal(t, 2, gotCnt)
	store.Delete("common checks - key2")
	store.Delete("common checks - key3")
}
