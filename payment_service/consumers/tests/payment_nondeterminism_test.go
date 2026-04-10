package tests

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FLAWED TEST — non-deterministic
// rand.Intn(10) != 0 means 90% success, 10% failure
// This test will randomly fail ~10% of the time with no code change
func TestPaymentStatus_Flawed_ShouldSucceed_WhenRandomIsNonZero(t *testing.T) {
	success := rand.Intn(10) != 0 // hidden random dependency — no seed set

	// This assertion will fail ~10% of the time
	assert.True(t, success, "payment should succeed")
}

// ROOT CAUSE: The test calls rand.Intn() without a fixed seed, so each run
// produces a different value, making the outcome unpredictable and the test unreliable.

// FIXED TEST — deterministic via injected randomiser
// Instead of calling rand directly, accept a rand.Source so tests can control the output.
type randomiser interface {
	Intn(n int) int
}

func resolvePaymentStatus(r randomiser) string {
	if r.Intn(10) != 0 {
		return "succeeded"
	}
	return "failed"
}

func TestPaymentStatus_Fixed_ShouldSucceed_WhenSeedProducesNonZero(t *testing.T) {
	// Seed 1 produces Intn(10) = 5 on first call — always succeeds
	r := rand.New(rand.NewSource(1))

	status := resolvePaymentStatus(r)

	assert.Equal(t, "succeeded", status)
}

func TestPaymentStatus_Fixed_ShouldFail_WhenSeedProducesZero(t *testing.T) {
	// Find a seed that produces 0 on first Intn(10) call
	// Seed 10 produces 0 — always fails
	r := rand.New(rand.NewSource(10))
	// burn first value which is non-zero for this seed
	for r.Intn(10) != 0 {
		r = rand.New(rand.NewSource(10))
		break
	}
	// Use a mock instead for full control
	r2 := rand.New(rand.NewSource(1))
	_ = r2 // showing the pattern; use mockRandomiser below for guaranteed zero

	mock := &mockRandomiser{value: 0} // always returns 0 → always "failed"
	status := resolvePaymentStatus(mock)

	assert.Equal(t, "failed", status)
}

// mockRandomiser gives full control over random output in tests
type mockRandomiser struct {
	value int
}

func (m *mockRandomiser) Intn(_ int) int {
	return m.value
}
