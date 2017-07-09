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
	Label        string
	Distribution Distribution
	// Lambda for exponential distribution
	Lambda float64
	// Lower and Upper bounds for uniform distribution
	LowerBound, UpperBound float64
	// Map of value->probability for discrete distribution
	Values map[interface{}]float64
}

// Generates values uniform distributed between a and b
func NewUniform(label string, a, b float64) (GenerationStrategy) {
	return GenerationStrategy{
		label, Uniform, 0.0, a, b, nil,
	}
}

// Generates values exponentially distributed with parameter lambda
func NewExponential(label string, lambda float64) (GenerationStrategy) {
	return GenerationStrategy{
		label, Exponential, lambda, 0.0, 0.0, nil,
	}
}

// Generates values exponentially distributed with parameter lambda
func NewDiscrete(label string, values map[interface{}]float64) (GenerationStrategy) {

	// normalize the values so the sum gives one
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	if sum == 1.0 {
		return GenerationStrategy{
			label, Discrete, 1.0, 0.0, 0.0, values,
		}
	} else {
		factor := 1.0 / sum
		nValues := make(map[interface{}]float64)
		for key, value := range values {
			nValues[key] = value * factor
		}
		return GenerationStrategy{
			label, Discrete, 1.0, 0.0, 0.0, nValues,
		}
	}

}

type Generator interface {
	AllAvailable(w int) (points []functions.MultidimensionalPoint)
	Next(w int) (point functions.MultidimensionalPoint)
	HasNext(w int) bool
}

// A multipoint generator structure
type randomGenerator struct {
	// Number of dimensions
	dimensionsNo int
	// The generation strategy on each dimension
	// GenerationStrategies are considered in the order they are defined (1st strategy applies to 1st dimension etc)
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
		values := make(map[string]interface{})
		labels := make([]string, g.dimensionsNo)
		for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {
			lowerBound, upperBound, lambda := -math.MaxFloat32, math.MaxFloat32, 1.0
			distribution := Uniform
			var samples map[interface{}]float64
			if len(g.restrictions) > dimIdx {
				lowerBound = g.restrictions[dimIdx].LowerBound
				upperBound = g.restrictions[dimIdx].UpperBound
				lambda = g.restrictions[dimIdx].Lambda
				distribution = g.restrictions[dimIdx].Distribution
				samples = g.restrictions[dimIdx].Values
				labels[dimIdx] = g.restrictions[dimIdx].Label
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

			switch distribution {
			case Uniform:
				_, values[labels[dimIdx]] = Float64(lowerBound, upperBound, g.rs[w])
			case Exponential:
				_, values[labels[dimIdx]] = ExpFloat64(lambda, g.rs[w])
			case Discrete:
				raw := g.rs[w].Float64()
				sum := 0.0
				for key, value := range samples {
					sum += value
					if raw <= sum {
						values[labels[dimIdx]] = key
						break
					}
				}
			}
		}
		points[pIdx] = functions.MultidimensionalPoint{Values: values}
	}

	g.index[w] = g.pointsNo

	return points
}

// Generates a new point
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g randomGenerator) Next(w int) (point functions.MultidimensionalPoint) {

	values := make(map[string]interface{})
	labels := make([]string, g.dimensionsNo)
	for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {
		lowerBound, upperBound, lambda := -math.MaxFloat32, math.MaxFloat32, 1.0
		distribution := Uniform
		var samples map[interface{}]float64
		if len(g.restrictions) > dimIdx {
			lowerBound = g.restrictions[dimIdx].LowerBound
			upperBound = g.restrictions[dimIdx].UpperBound
			lambda = g.restrictions[dimIdx].Lambda
			distribution = g.restrictions[dimIdx].Distribution
			samples = g.restrictions[dimIdx].Values
			labels[dimIdx] = g.restrictions[dimIdx].Label
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

		switch distribution {
		case Uniform:
			_, values[labels[dimIdx]] = Float64(lowerBound, upperBound, g.rs[w])
		case Exponential:
			_, values[labels[dimIdx]] = ExpFloat64(lambda, g.rs[w])
		case Discrete:
			raw := g.rs[w].Float64()
			sum := 0.0
			for key, value := range samples {
				sum += value
				if raw <= sum {
					values[labels[dimIdx]] = key
					break
				}
			}
		}

	}

	point = functions.MultidimensionalPoint{Values: values}

	g.index[w]++

	return
}

func (g randomGenerator) HasNext(w int) bool {
	return g.index[w] < g.pointsNo
}
