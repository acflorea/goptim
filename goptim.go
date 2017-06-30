package main

import (
	"github.com/acflorea/goptim/generators"
	"math"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"time"
	"math/rand"
	"flag"
)

func main() {

	fileNamePtr := flag.String("fileName", "", "Name of the input file.")
	noOfExperimentsPtr := flag.Int("noOfExperiments", 100, "Number of experiments.")
	silentPtr := flag.Bool("silent", true, "Silent Mode.")
	maxAttemptsPtr := flag.Int("maxAttempts", 300, "Maximum number of trials in an experiment")
	fct := flag.String("fct", "F_identity", "Target function")
	alg := flag.String("alg", "SeqSplit", "Parallel random generator strategy")

	flag.Parse()

	//functions.CrossV(1, 0.1)
	//functions.Train(1.0, 1.0/10.0)
	//functions.Test()

	vargs := map[string]interface{}{}
	vargs["fileName"] = *fileNamePtr
	vargs["noOfExperiments"] = *noOfExperimentsPtr
	vargs["silent"] = *silentPtr
	vargs["maxAttempts"] = *maxAttemptsPtr
	vargs["fct"] = *fct
	vargs["alg"] = *alg

	Optimize(vargs)
}

func Optimize(vargs map[string]interface{}) {

	fmt.Println("Optimization start!")
	fmt.Println(vargs)

	start := time.Now()

	noOfExperiments := vargs["noOfExperiments"].(int)
	silent := vargs["silent"].(bool)

	// Maximum number of attempts
	maxAttempts := vargs["maxAttempts"].(int)

	// The function we attempt to optimize
	targetFunction := functions.Functions[vargs["fct"].(string)]

	// Algorithm
	//(generators.SeqSplit seems to rule)
	algorithm := generators.Algorithms[vargs["alg"].(string)]

	// number of workers
	W := 10

	// 2^-3 to 2^10
	restrictions := []generators.Range{
		{0, 100},
		{0, 100},
	}

	match := 0
	var globalTries = 0

	for expIndex := 0; expIndex < noOfExperiments; expIndex++ {

		generator :=
			generators.NewRandomUniformGenerator(2, restrictions, maxAttempts, W, algorithm)

		// channel used by workers to communicate their results
		messages := make(chan functions.Sample, W)

		for w := 0; w < W; w++ {
			go func(w int) {
				i, p, v, gv, o := DMaximize(targetFunction, vargs, generator, maxAttempts/W, w, true)
				if !silent {
					fmt.Println("Worker ", w, " MAX --> ", i, p, v, gv, o)
				}

				messages <- functions.Sample{i, p, v, gv, o == 0}
			}(w)
		}

		// Collect results
		results := make([]functions.Sample, W)
		totalTries := 0
		optim, goptim := -math.MaxFloat64, -math.MaxFloat64
		var point functions.MultidimensionalPoint
		for i := 0; i < W; i++ {
			results[i] = <-messages
			if results[i].FullSearch {
				totalTries += maxAttempts / W
			} else {
				totalTries += results[i].Index
			}
			if optim < results[i].Value {
				optim = results[i].Value
				point = results[i].Point
			}
			if goptim < results[i].GValue {
				goptim = results[i].GValue
			}
		}

		if optim == goptim {
			match++
			globalTries += totalTries
			fmt.Println("+", totalTries, point, optim, goptim)
		} else {
			globalTries += totalTries
			fmt.Println("-", totalTries, point, optim, goptim)
		}
	}

	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Results matched on %d (%f) cases", match, float64(match)/float64(noOfExperiments)))
	avg := float64(globalTries) / float64(noOfExperiments)
	fmt.Println(fmt.Sprintf("Average number of attempts %f (%f percent faster) ", avg, (float64(maxAttempts)-avg)/float64(maxAttempts)*100))
	fmt.Println(fmt.Sprintf("Optimization took %s", elapsed))

}

// Attempts to dynamically minimize the function f
// k := n / (2 * math.E)
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
// w is thw worker index
func DMinimize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, n, w int, goAllTheWay bool) (
	index int,
	p functions.MultidimensionalPoint,
	min float64,
	gmin float64,
	optimNo int) {

	k := int(float64(n) / (2 * math.E))
	return Minimize(f, vargs, generator, k, n, w, goAllTheWay)
}

// Attempts to minimize the function f
// vargs are passed to the function
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
// gmin is the global minimum (if goAllTheWay then the algorithm continues and computes it
// for comparison purposes)
// w is the worker index
func Minimize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, k, n, w int, goAllTheWay bool) (
	index int,
	p functions.MultidimensionalPoint,
	min float64,
	gmin float64,
	optimNo int) {

	index = -1
	min = math.MaxFloat64
	gmin = math.MaxFloat64
	optimNo = 0

	minReached := false

	for i := 0; i < n; i++ {
		rndPoint := generator.Next(w)
		f_rnd, _ := f(rndPoint, vargs)

		if minReached {
			if f_rnd < gmin {
				gmin = f_rnd
			}
		} else {
			if f_rnd < min {
				index = i
				p = rndPoint
				min = f_rnd
				gmin = min

				if i > k {
					s := rand.NewSource(time.Now().UnixNano())
					tmpr := rand.New(s)
					threshold := tmpr.Float64()
					if threshold < 0.5+(0.05*float64(optimNo)) {
						minReached = true
						// Increase the number of optimum points found
						optimNo += 1
						if !goAllTheWay {
							break
						}
					} else {
						// Increase the number of optimum points found
						optimNo += 1
					}
				}
			}
		}
	}

	if !minReached {
		// The stop condition was not met
		optimNo = 0
	}

	return
}

// Dynamically Minimizes the negation of the target function
func DMaximize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, n, w int, goAllTheWay bool) (
	index int,
	p functions.MultidimensionalPoint,
	max float64,
	gmax float64,
	optimNo int) {

	index, p, max, gmax, optimNo = DMinimize(functions.Negate(f), vargs, generator, n, w, goAllTheWay)
	return index, p, -max, -gmax, optimNo
}

// Minimizes the negation of the target function
func Maximize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, k, n, w int, goAllTheWay bool) (
	index int,
	p functions.MultidimensionalPoint,
	max float64,
	gmax float64,
	optimNo int) {

	index, p, max, gmax, optimNo = Minimize(functions.Negate(f), vargs, generator, k, n, w, goAllTheWay)
	return index, p, -max, -gmax, optimNo
}
