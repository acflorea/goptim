package rand_test

import (
	"testing"
	"github.com/acflorea/goptim/rand"
)

func TestFloat64(t *testing.T) {

	a := 10.0
	b := 20.0
	for i := 0; i < 10; i++ {
		_, r := rand.Float64(a, b)
		if (r < a || r >= b) {
			t.Error("Invelid number generated")
		}
	}
}
