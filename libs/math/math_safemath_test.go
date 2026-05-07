package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- math.go ----

func TestMaxInt64(t *testing.T) {
	assert.Equal(t, int64(5), MaxInt64(3, 5))
	assert.Equal(t, int64(5), MaxInt64(5, 3))
	assert.Equal(t, int64(5), MaxInt64(5, 5))
	assert.Equal(t, int64(0), MaxInt64(0, -1))
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 7, MaxInt(7, 2))
	assert.Equal(t, 7, MaxInt(2, 7))
	assert.Equal(t, 7, MaxInt(7, 7))
	assert.Equal(t, 0, MaxInt(0, -1))
}

func TestMinInt64(t *testing.T) {
	assert.Equal(t, int64(3), MinInt64(3, 5))
	assert.Equal(t, int64(3), MinInt64(5, 3))
	assert.Equal(t, int64(3), MinInt64(3, 3))
	assert.Equal(t, int64(-1), MinInt64(0, -1))
}

func TestMinInt(t *testing.T) {
	assert.Equal(t, 2, MinInt(7, 2))
	assert.Equal(t, 2, MinInt(2, 7))
	assert.Equal(t, 2, MinInt(2, 2))
	assert.Equal(t, -1, MinInt(0, -1))
}

// ---- safemath.go ----

func TestSafeAddInt32(t *testing.T) {
	// Normal addition.
	assert.Equal(t, int32(7), SafeAddInt32(3, 4))
	assert.Equal(t, int32(-3), SafeAddInt32(-1, -2))

	// Positive overflow should panic.
	require.Panics(t, func() {
		SafeAddInt32(math.MaxInt32, 1)
	})

	// Negative overflow should panic.
	require.Panics(t, func() {
		SafeAddInt32(math.MinInt32, -1)
	})
}

func TestSafeSubInt32(t *testing.T) {
	// Normal subtraction.
	assert.Equal(t, int32(1), SafeSubInt32(3, 2))
	assert.Equal(t, int32(-5), SafeSubInt32(-3, 2))

	// Subtracting a positive value from MinInt32 overflows.
	require.Panics(t, func() {
		SafeSubInt32(math.MinInt32, 1)
	})

	// Subtracting a negative value from MaxInt32 overflows.
	require.Panics(t, func() {
		SafeSubInt32(math.MaxInt32, -1)
	})
}

func TestSafeConvertInt32(t *testing.T) {
	// In-range values.
	assert.Equal(t, int32(42), SafeConvertInt32(42))
	assert.Equal(t, int32(-42), SafeConvertInt32(-42))

	// Values at the exact boundaries.
	assert.Equal(t, int32(math.MaxInt32), SafeConvertInt32(math.MaxInt32))
	assert.Equal(t, int32(math.MinInt32), SafeConvertInt32(math.MinInt32))

	// Overflow above MaxInt32.
	require.Panics(t, func() {
		SafeConvertInt32(int64(math.MaxInt32) + 1)
	})

	// Overflow below MinInt32.
	require.Panics(t, func() {
		SafeConvertInt32(int64(math.MinInt32) - 1)
	})
}

func TestSafeConvertUint8(t *testing.T) {
	// Normal value.
	v, err := SafeConvertUint8(200)
	require.NoError(t, err)
	assert.Equal(t, uint8(200), v)

	// Boundary value.
	v, err = SafeConvertUint8(math.MaxUint8)
	require.NoError(t, err)
	assert.Equal(t, uint8(math.MaxUint8), v)

	// Zero is valid.
	v, err = SafeConvertUint8(0)
	require.NoError(t, err)
	assert.Equal(t, uint8(0), v)

	// Overflow above MaxUint8.
	_, err = SafeConvertUint8(int64(math.MaxUint8) + 1)
	require.ErrorIs(t, err, ErrOverflowUint8)

	// Negative value.
	_, err = SafeConvertUint8(-1)
	require.ErrorIs(t, err, ErrOverflowUint8)
}

func TestSafeConvertInt8(t *testing.T) {
	// Normal value.
	v, err := SafeConvertInt8(100)
	require.NoError(t, err)
	assert.Equal(t, int8(100), v)

	// Negative value.
	v, err = SafeConvertInt8(-100)
	require.NoError(t, err)
	assert.Equal(t, int8(-100), v)

	// Boundaries.
	v, err = SafeConvertInt8(math.MaxInt8)
	require.NoError(t, err)
	assert.Equal(t, int8(math.MaxInt8), v)

	v, err = SafeConvertInt8(math.MinInt8)
	require.NoError(t, err)
	assert.Equal(t, int8(math.MinInt8), v)

	// Overflow above MaxInt8.
	_, err = SafeConvertInt8(int64(math.MaxInt8) + 1)
	require.ErrorIs(t, err, ErrOverflowInt8)

	// Overflow below MinInt8.
	_, err = SafeConvertInt8(int64(math.MinInt8) - 1)
	require.ErrorIs(t, err, ErrOverflowInt8)
}

// ---- fraction.go ----

func TestFractionString(t *testing.T) {
	fr := Fraction{Numerator: 2, Denominator: 3}
	assert.Equal(t, "2/3", fr.String())

	fr = Fraction{Numerator: 0, Denominator: 1}
	assert.Equal(t, "0/1", fr.String())
}
