package generators

import (
	"math/rand"
	"github.com/acflorea/goptim/functions"
	"math"
	"time"
)

// Generates a uniform random value between a and b
func Float64(a, b float64, r *rand.Rand) (float64, float64) {
	if r == nil {
		// If the generator is not specified, create a new one
		source := rand.NewSource(time.Now().UnixNano())
		r = rand.New(source)
	}
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

// The random generation algorithm
type Algorithm int

// Types of parallel random generators
const (
	// A single generator, the master generates the values and pushes them to workers
	ManagerWorker Algorithm = iota
	// For each sequence of random numbers x[r]...x[r+p]...x[r+2p]... the process p takes every p-th value
	Leapfrog
	// Block allocation of data to tasks
	SeqSplit
	// Each woker has it's own parametrized generator
	Parametrization
)

// Algorithm labels
var algorithms = [...]string{
	"ManagerWorker",
	"FebrLeapfroguary",
	"SeqSplit",
	"Parametrization",
}

// A multipoint generator structure
type randomUniformGenerator struct {
	// Number of dimensions
	dimensionsNo int
	// Optional restrictions on each dimension
	// Restrictions are considered in the order they are defined (1st restriction applies to 1st dimension etc)
	restrictions []Range
	// How many points to generate
	pointsNo int
	// The level of parallelism
	cores int
	// The random generation algorithm
	algorithm Algorithm
	// The index of the last generated point
	index int
}

func NewRandomUniformGenerator(dimensionsNo int, restrictions []Range, pointsNo int, cores int, algorithm Algorithm) Generator {
	return randomUniformGenerator{
		dimensionsNo: dimensionsNo,
		restrictions: restrictions,
		pointsNo:     pointsNo,
		cores:        cores,
		algorithm:    algorithm,
		index:        0,
	}
}

// Generates g.PointsNo.
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g randomUniformGenerator) AllAvailable() (points []functions.MultidimensionalPoint) {

	points = make([]functions.MultidimensionalPoint, g.pointsNo)

	for pIdx := 0; pIdx < g.pointsNo; pIdx++ {
		values := make([]float64, g.dimensionsNo)
		for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {
			lowerBound, upperBound := -math.MaxFloat32, math.MaxFloat32
			if len(g.restrictions) > dimIdx {
				lowerBound = g.restrictions[dimIdx].LowerBound
				upperBound = g.restrictions[dimIdx].UpperBound
			}
			_, values[dimIdx] = Float64(lowerBound, upperBound, nil)
		}
		points[pIdx] = functions.MultidimensionalPoint{Values: values}
	}

	g.index = g.pointsNo

	return points
}

// Generates a new point
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g randomUniformGenerator) Next() (point functions.MultidimensionalPoint) {

	values := make([]float64, g.dimensionsNo)
	for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {
		lowerBound, upperBound := -math.MaxFloat32, math.MaxFloat32
		if len(g.restrictions) > dimIdx {
			lowerBound = g.restrictions[dimIdx].LowerBound
			upperBound = g.restrictions[dimIdx].UpperBound
		}
		_, values[dimIdx] = Float64(lowerBound, upperBound, nil)
	}

	point = functions.MultidimensionalPoint{Values: values}

	g.index++

	return
}

func (g randomUniformGenerator) HasNext() bool {
	return g.index < g.pointsNo
}
