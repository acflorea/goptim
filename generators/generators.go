package generators

import (
	"math/rand"
	"github.com/acflorea/goptim/functions"
	"math"
	"time"
	"sync"
)

var mutex sync.Mutex

// Generates a uniform random value between a and b
func Float64(a, b float64, r *rand.Rand) (float64, float64) {
	if r == nil {
		// If the generator is not specified, create a new one
		source := rand.NewSource(time.Now().UnixNano())
		r = rand.New(source)
	}

	// Lock - avoid data races TODO - find a way to implement this using channels
	mutex.Lock()
	defer mutex.Unlock()

	original := r.Float64()
	return original, a + (b-a)*original
}

// One dimensional restriction [lowerBound, upperBound)
type Range struct {
	LowerBound, UpperBound float64
}

type Generator interface {
	AllAvailable(w int) (points []functions.MultidimensionalPoint)
	Next(w int) (point functions.MultidimensionalPoint)
	HasNext(w int) bool
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
	index []int
	// internal random generator(s)
	rs []*rand.Rand
}

func NewRandomUniformGenerator(dimensionsNo int, restrictions []Range, pointsNo int, cores int, algorithm Algorithm) Generator {
	generator := randomUniformGenerator{
		dimensionsNo: dimensionsNo,
		restrictions: restrictions,
		pointsNo:     pointsNo,
		cores:        cores,
		algorithm:    algorithm,
		index:        make([]int, cores),
	}

	// Init generator
	now := time.Now().UnixNano()

	rs := make([]*rand.Rand, cores, cores)

	switch algorithm {
	case ManagerWorker:
		// Same generator for all workers
		source := rand.NewSource(now)
		r := rand.New(source)
		for i := 0; i < cores; i++ {
			rs[i] = r
		}
	case Leapfrog:
		for i := 0; i < cores; i++ {
			source := rand.NewSource(now)
			rs[i] = rand.New(source)
		}
	case SeqSplit:
		for i := 0; i < cores; i++ {
			source := rand.NewSource(now)
			rs[i] = rand.New(source)
			// Advance the generator
			for j := 0; j < i*pointsNo/cores; j++ {
				rs[i].Float64()
			}
		}
	case Parametrization:
		for i := 0; i < cores; i++ {
			// Different seed meaning different sequences
			source := rand.NewSource(now - int64(i))
			rs[i] = rand.New(source)
		}
	}

	generator.rs = rs

	return generator
}

// Generates g.PointsNo.
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g randomUniformGenerator) AllAvailable(w int) (points []functions.MultidimensionalPoint) {

	points = make([]functions.MultidimensionalPoint, g.pointsNo)

	for pIdx := 0; pIdx < g.pointsNo; pIdx++ {
		values := make([]float64, g.dimensionsNo)
		for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {
			lowerBound, upperBound := -math.MaxFloat32, math.MaxFloat32
			if len(g.restrictions) > dimIdx {
				lowerBound = g.restrictions[dimIdx].LowerBound
				upperBound = g.restrictions[dimIdx].UpperBound
			}

			if g.algorithm == Leapfrog {
				if g.index[w] == 0 {
					// Set the counter in place
					for i := 0; i < w; i++ {
						g.rs[w].Float64()
					}
				} else {
					// Jump "cores" positions
					for i := 0; i < g.cores; i++ {
						g.rs[w].Float64()
					}
				}
			}

			_, values[dimIdx] = Float64(lowerBound, upperBound, g.rs[w])
		}
		points[pIdx] = functions.MultidimensionalPoint{Values: values}
	}

	g.index[w] = g.pointsNo

	return points
}

// Generates a new point
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g randomUniformGenerator) Next(w int) (point functions.MultidimensionalPoint) {

	values := make([]float64, g.dimensionsNo)
	for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {
		lowerBound, upperBound := -math.MaxFloat32, math.MaxFloat32
		if len(g.restrictions) > dimIdx {
			lowerBound = g.restrictions[dimIdx].LowerBound
			upperBound = g.restrictions[dimIdx].UpperBound
		}

		if g.algorithm == Leapfrog {
			if g.index[w] == 0 {
				// Set the counter in place
				for i := 0; i < w; i++ {
					g.rs[w].Float64()
				}
			} else {
				// Jump "cores" positions
				for i := 0; i < g.cores; i++ {
					g.rs[w].Float64()
				}
			}
		}

		_, values[dimIdx] = Float64(lowerBound, upperBound, g.rs[w])
	}

	point = functions.MultidimensionalPoint{Values: values}

	g.index[w]++

	return
}

func (g randomUniformGenerator) HasNext(w int) bool {
	return g.index[w] < g.pointsNo
}
