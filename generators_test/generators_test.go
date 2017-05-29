package generators_test

import (
	"testing"
	"github.com/acflorea/goptim/generators"
	"fmt"
	"github.com/acflorea/goptim/functions"
)

func Test_Float64(t *testing.T) {

	a := 10.0
	b := 20.0
	for i := 0; i < 10; i++ {
		_, r := generators.Float64(a, b, nil)
		if (r < a || r >= b) {
			t.Error("Invalid number generated")
		}
	}
}

func Test_RandomUniformPointsGeneratorNext(t *testing.T) {

	howManyPoints := 10
	dimensionsNo := 2

	restrictions := []generators.Range{
		{-10, 10},
		{-10, 10},
	}

	generator :=
		generators.NewRandomUniformGenerator(dimensionsNo, restrictions, howManyPoints, 1, generators.ManagerWorker)

	generatedPoints := make([]functions.MultidimensionalPoint, howManyPoints)
	for pIdx := 0; generator.HasNext(0); pIdx++ {
		generatedPoints[pIdx] = generator.Next(0)
	}

	if len(generatedPoints) != howManyPoints {
		msg := fmt.Sprintf("Error generating points. "+
			"Expected (%i) but got (%i).", howManyPoints, len(generatedPoints))
		t.Error(msg)
	}

	for pIdx := 0; pIdx < len(generatedPoints); pIdx++ {
		x := generatedPoints[pIdx].Values[0]
		y := generatedPoints[pIdx].Values[1]
		if x < -10 || x >= 10 || y < -10 || y >= 10 {
			msg := fmt.Sprintf("Error generating points. Coordinate out of bounds "+
				"(x, y) = (%f, %f)", x, y)
			t.Error(msg)
		}
	}
}

func Test_RandomUniformPointsGeneratorAll(t *testing.T) {

	howManyPoints := 10
	dimensionsNo := 2

	restrictions := []generators.Range{
		{-10, 10},
		{-10, 10},
	}

	generator :=
		generators.NewRandomUniformGenerator(dimensionsNo, restrictions, howManyPoints, 1, generators.ManagerWorker)

	generatedPoints := generator.AllAvailable(0)

	if len(generatedPoints) != howManyPoints {
		msg := fmt.Sprintf("Error generating points. "+
			"Expected (%i) but got (%i).", howManyPoints, len(generatedPoints))
		t.Error(msg)
	}

	for pIdx := 0; pIdx < len(generatedPoints); pIdx++ {
		x := generatedPoints[pIdx].Values[0]
		y := generatedPoints[pIdx].Values[1]
		if x < -10 || x >= 10 || y < -10 || y >= 10 {
			msg := fmt.Sprintf("Error generating points. Coordinate out of bounds "+
				"(x, y) = (%f, %f)", x, y)
			t.Error(msg)
		}
	}
}
