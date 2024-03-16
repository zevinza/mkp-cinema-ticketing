package lib

// ProfitPercent calculate profit percent from purchase price and sellprice
func ProfitPercent(purchase, sell float64) float64 {
	return (sell - purchase) / purchase * 100
}
