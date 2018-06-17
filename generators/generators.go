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

type GeneratorState struct {
	// points generated so far
	GeneratedPoints []functions.MultidimensionalPoint
	// values for those points
	Output []float64
	// centroid
	Centroid functions.MultidimensionalPoint
}

type Generator interface {
	AllAvailable(w int) (points []functions.MultidimensionalPoint)
	Next(w int, initialState GeneratorState) (point functions.MultidimensionalPoint, state GeneratorState)
	HasNext(w int) bool
	Improvement(state GeneratorState) bool
}

// A multipoint generator structure
type randomGenerator struct {
	// Number of dimensions
	dimensionsNo int
	// The generation strategy on each dimension
	// GenerationStrategies are considered in the order they are defined (1st strategy applies to 1st dimension etc)
	restrictions []GenerationStrategy
	// probability to change for each dimension
	// the probability to change for each dimension
	probabilityToChange []float32
	// the probability to change for each dimension - reversed
	reverse_probabilityToChange []float32
	// change a single value per step
	adjustSingleValue bool
	// optimalSlicePercent - the slice of results that are considered in the optimal range
	optimalSlicePercent float64
	// How many points to generate in total
	pointsNo int
	// Minimum number of point to generate
	minPointsNo int
	// The level of parallelism
	cores int
	// The random generation algorithm
	algorithm Algorithm
	// The index of the last generated point
	index []int
	// internal random generator(s)
	rs []*rand.Rand
}

func NewRandom(restrictions []GenerationStrategy, probabilityToChange []float32, adjustSingleValue bool, optimalSlicePercent float64, pointsNo int, minPointsNo int, cores int, algorithm Algorithm) Generator {

	// adjust the probabilityToChange values to sum up to 1.0
	// normalize the values so the sum gives one
	sum := float32(0.0)
	for _, value := range probabilityToChange {
		sum += value
	}
	if sum != 1.0 {
		factor := 1.0 / sum
		for key, value := range probabilityToChange {
			probabilityToChange[key] = value * factor
		}
	}

	// compute the reverse probabilityToChange
	reverse_probabilityToChange := []float32{}
	sum = float32(len(probabilityToChange) - 1)
	for _, value := range probabilityToChange {
		reverse_probabilityToChange = append(reverse_probabilityToChange, (1.0-value)/sum)
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

	generator := randomGenerator{
		dimensionsNo:                len(restrictions),
		restrictions:                restrictions,
		probabilityToChange:         probabilityToChange,
		reverse_probabilityToChange: reverse_probabilityToChange,
		adjustSingleValue:           adjustSingleValue,
		optimalSlicePercent:         optimalSlicePercent,
		pointsNo:                    pointsNo,
		cores:                       cores,
		algorithm:                   algorithm,
		minPointsNo:                 minPointsNo,
		index:                       make([]int, cores),
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

			lowerBound, upperBound, lambda, distribution, samples, label := getRestrictionsPerDimension(g, dimIdx)
			labels[dimIdx] = label

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

// Check if improvement was made
func (g randomGenerator) Improvement(state GeneratorState) bool {
	previousOutputLength := len(state.Output)

	// TODO - Store this in the state
	min, max := math.MaxFloat32, -math.MaxFloat32
	for _, value := range state.Output {
		if min > value {
			min = value
		}
		if max < value {
			max = value
		}
	}

	boundary := min + (max-min)/100.0*g.optimalSlicePercent

	if state.Output[previousOutputLength-1] < boundary {
		return true
	}

	return false
}

// Generates a new point
// Each point is a collection of g.DimensionsNo uniform random values bounded to g.Restrictions
func (g randomGenerator) Next(w int, initialState GeneratorState) (point functions.MultidimensionalPoint, state GeneratorState) {

	values := make(map[string]interface{})
	labels := make([]string, g.dimensionsNo)
	currentIndex := g.index[w]

	state = initialState
	if len(state.GeneratedPoints) > 0 {
		// we have state (either we have generated some numbers or this is the provided initial state)

		previousPoint := state.GeneratedPoints[len(state.GeneratedPoints)-1]

		// check if previous point was an improvement
		var wasAnImprovement = false
		if currentIndex >= g.minPointsNo/g.cores {
			wasAnImprovement = g.Improvement(state)
		}

		var probabilities = g.probabilityToChange
		if wasAnImprovement {
			// we reverse probabilities if the value was an improvement
			probabilities = g.reverse_probabilityToChange
		}

		// Each value changes with this probability
		globalProbabilityToChange := g.rs[w].Float32()
		indexToChange := -1

		if g.adjustSingleValue {

			// Identify which value should change
			sum := float32(0.0)
			for key, value := range probabilities {
				sum += value
				if globalProbabilityToChange <= sum {
					indexToChange = key
					break
				}
			}

		} else {
			change := false
			for !change {
				// check if at least one dimension changes
				for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {
					if g.probabilityToChange[dimIdx] >= globalProbabilityToChange {
						change = true
						break
					}
				}
				// if nothing changes regenerate the probability
				if !change {
					globalProbabilityToChange = g.rs[w].Float32()
				}
			}

		}

		for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {

			// attempt to retrieve individual probability to change
			var probabilityToChange float32 = 0.0

			// we are still in the tuning phase
			if currentIndex < g.minPointsNo/g.cores {
				// so we change all the values
				probabilityToChange = 1.0
			} else {

				if g.adjustSingleValue {
					// we are in the case where a single value changes
					if dimIdx == indexToChange {
						// and this is the one
						probabilityToChange = 1.0
					} else {
						// and this is not the one
						probabilityToChange = 0.0
					}
				} else {
					// we are in the case where multiple values change...
					// thy to get the probability to change for each one
					if len(g.probabilityToChange) > dimIdx {
						probabilityToChange = probabilities[dimIdx]
					} else {
						// if the probability is not explicit, consider it 1.0
						probabilityToChange = 1.0
					}
				}
			}

			lowerBound, upperBound, lambda, distribution, samples, label := getRestrictionsPerDimension(g, dimIdx)
			labels[dimIdx] = label

			if probabilityToChange >= globalProbabilityToChange {
				// change
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

			} else {
				// preserve
				values[labels[dimIdx]] = previousPoint.Values[labels[dimIdx]]
				// values[labels[dimIdx]] = state.Centroid.Values[labels[dimIdx]]
			}

		}

		point = functions.MultidimensionalPoint{Values: values}

	} else {

		// 1st attempt, no state given

		for dimIdx := 0; dimIdx < g.dimensionsNo; dimIdx++ {

			lowerBound, upperBound, lambda, distribution, samples, label := getRestrictionsPerDimension(g, dimIdx)
			labels[dimIdx] = label

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

	}

	state.GeneratedPoints = append(state.GeneratedPoints, point)

	g.index[w]++

	return
}

func getRestrictionsPerDimension(g randomGenerator, dimIdx int) (float64, float64, float64, Distribution, map[interface{}]float64, string) {
	lowerBound, upperBound, lambda := -math.MaxFloat32, math.MaxFloat32, 1.0
	distribution := Uniform
	var samples map[interface{}]float64
	label := ""
	if len(g.restrictions) > dimIdx {
		lowerBound = g.restrictions[dimIdx].LowerBound
		upperBound = g.restrictions[dimIdx].UpperBound
		lambda = g.restrictions[dimIdx].Lambda
		distribution = g.restrictions[dimIdx].Distribution
		samples = g.restrictions[dimIdx].Values
		label = g.restrictions[dimIdx].Label
	}
	return lowerBound, upperBound, lambda, distribution, samples, label
}

func (g randomGenerator) HasNext(w int) bool {
	return g.index[w] < g.pointsNo
}
