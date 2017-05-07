package rand_test

import (
	"testing"
	"github.com/acflorea/goptim/rand"
	"fmt"
)

func Test_Float64(t *testing.T) {

	a := 10.0
	b := 20.0
	for i := 0; i < 10; i++ {
		_, r := rand.Float64(a, b)
		if (r < a || r >= b) {
			t.Error("Invalid number generated")
		}
	}
}

func Test_RandomUniformPointsGenerator(t *testing.T) {

	generator := rand.MultiPointGenerator{
		DimensionsNo: 2,
		PointsNo:     10,
		Restrictions: []rand.Range{
			{-10, 0},
			{0, 10},
		},
	}
	generatedPoints := generator.RandomUniform()

	if (len(generatedPoints) != 10) {
		msg := fmt.Sprintf("Error generating points. "+
			"Expected (%i) but got (%i).", 10, len(generatedPoints))
		t.Error(msg)
	}

	for pIdx := 0; pIdx < len(generatedPoints); pIdx++ {
		x := generatedPoints[pIdx].Values[0]
		y := generatedPoints[pIdx].Values[1]
		if x < -10 || x >= 0 || y < 0 || y >= 10 {
			msg := fmt.Sprintf("Error generating points. Coordinate out of bounds "+
				"(x, y) = (%f, %f)", x, y)
			t.Error(msg)
		}
	}
}
