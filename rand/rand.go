package rand

import "math/rand"

// Generates a uniform random value between a and b
func Float64(a, b float64) (float64, float64) {
	original := rand.Float64()
	return original, a + (b-a)*original
}
