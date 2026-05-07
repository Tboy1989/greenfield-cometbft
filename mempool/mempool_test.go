package mempool

import (
	"errors"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPreCheckMaxBytes verifies that PreCheckMaxBytes rejects transactions that
// exceed the given size limit and accepts those within it.
func TestPreCheckMaxBytes(t *testing.T) {
	maxBytes := int64(100)
	checkFn := PreCheckMaxBytes(maxBytes)

	// A small transaction should pass.
	smallTx := types.Tx(make([]byte, 10))
	require.NoError(t, checkFn(smallTx))

	// A transaction whose encoded proto size exceeds maxBytes should fail.
	// types.ComputeProtoSizeForTxs adds a small overhead, so use a tx that is
	// just over the limit to reliably trigger the error.
	largeTx := types.Tx(make([]byte, 200))
	err := checkFn(largeTx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tx size is too big")
}

// TestPostCheckMaxGas verifies PostCheckMaxGas behaviour for the various
// edge cases: unlimited gas, negative GasWanted, gas within limit and gas
// exceeding the limit.
func TestPostCheckMaxGas(t *testing.T) {
	tx := types.Tx("dummy")

	// maxGas == -1 means unlimited – any response should pass.
	unlimitedFn := PostCheckMaxGas(-1)
	require.NoError(t, unlimitedFn(tx, &abci.ResponseCheckTx{GasWanted: 1_000_000}))

	limitedFn := PostCheckMaxGas(100)

	// GasWanted within limit.
	require.NoError(t, limitedFn(tx, &abci.ResponseCheckTx{GasWanted: 50}))

	// GasWanted exactly at limit.
	require.NoError(t, limitedFn(tx, &abci.ResponseCheckTx{GasWanted: 100}))

	// GasWanted exceeds limit.
	err := limitedFn(tx, &abci.ResponseCheckTx{GasWanted: 101})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gas wanted")

	// Negative GasWanted should be rejected.
	err = limitedFn(tx, &abci.ResponseCheckTx{GasWanted: -1})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "negative")
}

// TestErrTxTooLarge verifies the error message formatting of ErrTxTooLarge.
func TestErrTxTooLarge(t *testing.T) {
	err := ErrTxTooLarge{Max: 1024, Actual: 2048}
	msg := err.Error()
	assert.Contains(t, msg, "1024")
	assert.Contains(t, msg, "2048")
}

// TestErrMempoolIsFull verifies the error message formatting of ErrMempoolIsFull.
func TestErrMempoolIsFull(t *testing.T) {
	err := ErrMempoolIsFull{
		NumTxs:      10,
		MaxTxs:      100,
		TxsBytes:    512,
		MaxTxsBytes: 1024,
	}
	msg := err.Error()
	assert.Contains(t, msg, "10")
	assert.Contains(t, msg, "100")
	assert.Contains(t, msg, "512")
	assert.Contains(t, msg, "1024")
}

// TestErrPreCheck verifies that ErrPreCheck.Error() returns the wrapped reason
// and that IsPreCheckError correctly identifies it.
func TestErrPreCheck(t *testing.T) {
	inner := errors.New("transaction too large")
	preCheckErr := ErrPreCheck{Reason: inner}

	// Error message should be the inner error's message.
	assert.Equal(t, inner.Error(), preCheckErr.Error())

	// IsPreCheckError should return true for ErrPreCheck.
	assert.True(t, IsPreCheckError(preCheckErr))

	// IsPreCheckError should return false for a plain error.
	assert.False(t, IsPreCheckError(inner))
}
