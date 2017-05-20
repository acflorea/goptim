package generators

import (
	"math/rand"
	"github.com/acflorea/goptim/functions"
	"math"
	"time"
)

// Generates a uniform random value between a and b
func Float64(a, b float64) (float64, float64) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	original := r.Float64()
	return original, a + (b-a)*original
}

// One dimensional restriction [lowerBound, upperBound)
type Range struct {
	LowerBound, UpperBound float64
}

type Generator interface {
	AllAvailable() (points []functions.MultidimensionalPoint)
	Next() (point functions.MultidimensionalPoint)
	HasNext() bool
}

// A multipoint generator structure
type RandomUniformGenerator struct {
	// Number of dimensions
	DimensionsNo int
	// Optional restrictions on each dimension
	// Restrictions are considered in the order they are defined (1st restriction applies to 1st dimension etc)
	Restrictions []Range
	// How many points to generate
	PointsNo int
	// The index of the last generated point
	index int
}

// Generates g.PointsNo.
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g RandomUniformGenerator) AllAvailable() (points []functions.MultidimensionalPoint) {

	points = make([]functions.MultidimensionalPoint, g.PointsNo)

	for pIdx := 0; pIdx < g.PointsNo; pIdx++ {
		values := make([]float64, g.DimensionsNo)
		for dimIdx := 0; dimIdx < g.DimensionsNo; dimIdx++ {
			lowerBound, upperBound := -math.MaxFloat32, math.MaxFloat32
			if len(g.Restrictions) > dimIdx {
				lowerBound = g.Restrictions[dimIdx].LowerBound
				upperBound = g.Restrictions[dimIdx].UpperBound
			}
			_, values[dimIdx] = Float64(lowerBound, upperBound)
		}
		points[pIdx] = functions.MultidimensionalPoint{Values: values}
	}

	g.index = g.PointsNo

	return points
}

// Generates a new point
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g RandomUniformGenerator) Next() (point functions.MultidimensionalPoint) {

	values := make([]float64, g.DimensionsNo)
	for dimIdx := 0; dimIdx < g.DimensionsNo; dimIdx++ {
		lowerBound, upperBound := -math.MaxFloat32, math.MaxFloat32
		if len(g.Restrictions) > dimIdx {
			lowerBound = g.Restrictions[dimIdx].LowerBound
			upperBound = g.Restrictions[dimIdx].UpperBound
		}
		_, values[dimIdx] = Float64(lowerBound, upperBound)
	}

	point = functions.MultidimensionalPoint{Values: values}

	g.index++

	return
}

func (g RandomUniformGenerator) HasNext() bool {
	return g.index < g.PointsNo
}