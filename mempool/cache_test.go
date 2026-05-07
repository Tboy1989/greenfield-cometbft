package mempool

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLRUTxCacheGetList verifies that GetList returns the backing linked list.
func TestLRUTxCacheGetList(t *testing.T) {
	cache := NewLRUTxCache(10)
	require.NotNil(t, cache.GetList())
	require.Equal(t, cache.list, cache.GetList())
}

// TestLRUTxCacheReset verifies that Reset clears both the map and list.
func TestLRUTxCacheReset(t *testing.T) {
	cache := NewLRUTxCache(10)
	tx := []byte("reset_test_tx")
	cache.Push(tx)
	require.Equal(t, 1, len(cache.cacheMap))

	cache.Reset()
	require.Equal(t, 0, len(cache.cacheMap))
	require.Equal(t, 0, cache.GetList().Len())
}

// TestLRUTxCacheHas verifies that Has returns true for cached and false for missing txs.
func TestLRUTxCacheHas(t *testing.T) {
	cache := NewLRUTxCache(10)
	tx := []byte("has_test_tx")

	assert.False(t, cache.Has(tx))
	cache.Push(tx)
	assert.True(t, cache.Has(tx))

	cache.Remove(tx)
	assert.False(t, cache.Has(tx))
}

// TestNopTxCache verifies that NopTxCache satisfies the TxCache interface and
// always behaves as a no-op (Push always returns true, Has always returns false).
func TestNopTxCache(t *testing.T) {
	var c TxCache = NopTxCache{}
	tx := []byte("nop_tx")

	// Reset is a no-op – just make sure it doesn't panic.
	c.Reset()

	// Push always returns true (treat every tx as new).
	assert.True(t, c.Push(tx))
	assert.True(t, c.Push(tx))

	// Has always returns false.
	assert.False(t, c.Has(tx))

	// Remove is a no-op – just make sure it doesn't panic.
	c.Remove(tx)
}

func TestCacheRemove(t *testing.T) {
	cache := NewLRUTxCache(100)
	numTxs := 10

	txs := make([][]byte, numTxs)
	for i := 0; i < numTxs; i++ {
		// probability of collision is 2**-256
		txBytes := make([]byte, 32)
		_, err := rand.Read(txBytes)
		require.NoError(t, err)

		txs[i] = txBytes
		cache.Push(txBytes)

		// make sure its added to both the linked list and the map
		require.Equal(t, i+1, len(cache.cacheMap))
		require.Equal(t, i+1, cache.list.Len())
	}

	for i := 0; i < numTxs; i++ {
		cache.Remove(txs[i])
		// make sure its removed from both the map and the linked list
		require.Equal(t, numTxs-(i+1), len(cache.cacheMap))
		require.Equal(t, numTxs-(i+1), cache.list.Len())
	}
}
