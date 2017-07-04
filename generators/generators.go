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
	r = initGenerator(r)

	original := r.Float64()
	return original, a + (b-a)*original
}

// Generates a random value from an exponential distribution with rate lambda
func ExpFloat64(lambda float64, r *rand.Rand) (float64, float64) {
	r = initGenerator(r)

	original := r.Float64()
	// x = log(1-u)/(−λ)
	return original, math.Log(1-original) / (-lambda)
}

func initGenerator(r *rand.Rand) *rand.Rand {
	if r == nil {
		// If the generator is not specified, create a new one
		source := rand.NewSource(time.Now().UnixNano())
		r = rand.New(source)
	}
	// Lock - avoid data races TODO - find a way to implement this using channels
	mutex.Lock()
	defer mutex.Unlock()
	return r
}

// One dimensional restriction [lowerBound, upperBound)
type GenerationStrategy struct {
	Distribution           Distribution
	LowerBound, UpperBound float64
}

//
func NewUniform(a, b float64) (GenerationStrategy) {
	return GenerationStrategy{
		Uniform, a, b,
	}
}

//
func NewExponential(lambda float64) (GenerationStrategy) {
	return GenerationStrategy{
		Exponential, lambda, 0.0,
	}
}

type Generator interface {
	AllAvailable(w int) (points []functions.MultidimensionalPoint)
	Next(w int) (point functions.MultidimensionalPoint)
	HasNext(w int) bool
}

// The distributions
type Distribution int

// Map with distribution by name
var Distributions = map[string]Distribution{
	"Uniform":     Uniform,
	"Exponential": Exponential,
}

// Types of distributions
const (
	// Uniform
	Uniform Distribution = iota
	// Exponential
	Exponential
)

// Algorithm labels
var distributions = [...]string{
	"Uiform",
	"Exponential",
}

// The random generation algorithm
type Algorithm int

// Map with functions by name
var Algorithms = map[string]Algorithm{
	"ManagerWorker":   ManagerWorker,
	"Leapfrog":        Leapfrog,
	"SeqSplit":        SeqSplit,
	"Parametrization": Parametrization,
}

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
type randomGenerator struct {
	// Number of dimensions
	dimensionsNo int
	// Optional restrictions on each dimension
	// Restrictions are considered in the order they are defined (1st restriction applies to 1st dimension etc)
	restrictions []GenerationStrategy
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

func NewRandomGenerator(dimensionsNo int, restrictions []GenerationStrategy, pointsNo int, cores int, algorithm Algorithm) Generator {
	generator := randomGenerator{
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
func (g randomGenerator) AllAvailable(w int) (points []functions.MultidimensionalPoint) {

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
			//Float64(lowerBound, upperBound, g.rs[w])
		}
		points[pIdx] = functions.MultidimensionalPoint{Values: values}
	}

	g.index[w] = g.pointsNo

	return points
}

// Generates a new point
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g randomGenerator) Next(w int) (point functions.MultidimensionalPoint) {

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

func (g randomGenerator) HasNext(w int) bool {
	return g.index[w] < g.pointsNo
}
