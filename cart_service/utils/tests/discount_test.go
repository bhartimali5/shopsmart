package tests

import (
	"testing"

	"example.com/rest-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateDiscount(t *testing.T) {

	t.Run("ValidTiers", func(t *testing.T) {
		t.Run("ShouldReturn15PercentDiscount_WhenUserTierIsGold", func(t *testing.T) {
			// Given a price of 100 and user tier gold
			price := 100.0
			tier := "gold"

			// When CalculateDiscount is called
			result := utils.CalculateDiscount(price, tier)

			// Then the discount should be 15% of the price
			assert.Equal(t, 15.0, result)
		})

		t.Run("ShouldReturn20PercentDiscount_WhenUserTierIsPlatinum", func(t *testing.T) {
			// Arrange
			price := 200.0
			tier := "platinum"

			// Act
			result := utils.CalculateDiscount(price, tier)

			// Assert
			assert.Equal(t, 40.0, result)
		})
	})

	t.Run("EdgeCases", func(t *testing.T) {
		t.Run("ShouldReturnZero_WhenPriceIsNegative", func(t *testing.T) {
			result := utils.CalculateDiscount(-50, "silver")
			assert.Equal(t, 0.0, result)
		})

		t.Run("ShouldReturnZero_WhenPriceIsZero", func(t *testing.T) {
			result := utils.CalculateDiscount(0, "bronze")
			assert.Equal(t, 0.0, result)
		})

		t.Run("ShouldReturnZeroDiscount_WhenUserTierIsUnknown", func(t *testing.T) {
			result := utils.CalculateDiscount(100, "vip")
			assert.Equal(t, 0.0, result)
		})
	})
}
