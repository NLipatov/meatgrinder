package domain

func CalculateDamage(base float64, resist float64) float64 {
	return base * (1.0 - resist)
}
