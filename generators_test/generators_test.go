package generators_test

import (
	"testing"
	"github.com/acflorea/goptim/generators"
	"fmt"
	"github.com/acflorea/goptim/functions"
	"math"
	"time"
	"math/rand"
)

func Test_Float64(t *testing.T) {

	// If the generator is not specified, create a new one
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	count := 1000000
	values := make([]float64, count)

	a := 0.0
	b := 1.0
	expectedMean := (b + a) / 2
	expectedVariance := (b - a) * (b - a) / 12.0
	for i := 0; i < count; i++ {
		_, values[i] = generators.Float64(a, b, r)
	}
	mean := 0.0
	for i := 0; i < count; i++ {
		mean += values[i] / float64(count)
	}
	variance := 0.0
	for i := 0; i < count; i++ {
		variance += math.Pow((values[i] - mean), 2) / float64(count)
	}
	if (math.Abs(mean-expectedMean) > 0.001) {
		t.Error("The generated values don't look uniform - the mean is different.", mean, expectedMean)
	}
	if (math.Abs(variance-expectedVariance) > 0.001) {
		t.Error("The generated values don't look unform - the variance is different", variance, expectedVariance)
	}

}

func Test_ExpFloat64(t *testing.T) {

	// If the generator is not specified, create a new one
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	count := 1000000
	values := make([]float64, count)
	a := 100.0
	expectedMean := 1.0 / a
	expectedVariance := 1.0 / (a * a)
	for i := 0; i < count; i++ {
		_, values[i] = generators.ExpFloat64(a, r)
	}
	mean := 0.0
	for i := 0; i < count; i++ {
		mean += values[i] / float64(count)
	}
	variance := 0.0
	for i := 0; i < count; i++ {
		variance += math.Pow((values[i] - mean), 2) / float64(count)
	}
	if (math.Abs(mean-expectedMean) > 0.001) {
		t.Error("The generated values don't look exponential - the mean is different.", mean, expectedMean)
	}
	if (math.Abs(variance-expectedVariance) > 0.001) {
		t.Error("The generated values don't look exponential - the variance is different", variance, expectedVariance)
	}

}

func Test_RandomPointsGeneratorNext(t *testing.T) {

	howManyPoints := 10
	dimensionsNo := 2

	restrictions := []generators.GenerationStrategy{
		generators.NewUniform(-10, 10),
		generators.NewUniform(-10, 10),
	}

	generator :=
		generators.NewRandomGenerator(dimensionsNo, restrictions, howManyPoints, 1, generators.ManagerWorker)

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

func Test_RandomPointsGeneratorAll(t *testing.T) {

	howManyPoints := 10
	dimensionsNo := 2

	restrictions := []generators.GenerationStrategy{
		generators.NewUniform(-10, 10),
		generators.NewUniform(-10, 10),
	}

	generator :=
		generators.NewRandomGenerator(dimensionsNo, restrictions, howManyPoints, 1, generators.ManagerWorker)

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
