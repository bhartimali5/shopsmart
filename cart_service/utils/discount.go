package utils

// CalculateDiscount returns the discounted price based on user tier.
// Returns 0 if price is negative or tier is unknown.
func CalculateDiscount(price float64, userTier string) float64 {
	if price < 0 {
		return 0
	}

	discountRates := map[string]float64{
		"bronze":   0.05,
		"silver":   0.10,
		"gold":     0.15,
		"platinum": 0.20,
	}

	rate, ok := discountRates[userTier]
	if !ok {
		return 0
	}

	return price * rate
}
