package main

import (
	"github.com/acflorea/goptim/generators"
	"math"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"time"
	"math/rand"
	"flag"
	"github.com/acflorea/libsvm-go"
)

// The result of one trial
type OptimizationOutput struct {
	Optim  float64
	GOptim float64
	X      functions.MultidimensionalPoint
	Trials int
}

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
	W := 5

	//2^-3 to 2^10
	//restrictions := []generators.GenerationStrategy{
	//	generators.NewUniform("C", math.Pow(2, -2), math.Pow(2, 15)),
	//	generators.NewUniform("gamma", math.Pow(2, -15), math.Pow(2, 3)),
	//}

	//{"linear", "polynomial", "rbf", "sigmoid"}
	restrictions := []generators.GenerationStrategy{
		generators.NewDiscrete("kernel", map[interface{}]float64{
			libSvm.LINEAR:  1.0, // 0
			libSvm.POLY:    1.0, // 1
			libSvm.RBF:     1.0, // 2
			libSvm.SIGMOID: 1.0, // 3
		}),
		generators.NewExponential("C", 10),
		generators.NewExponential("gamma", 10),
		generators.NewDiscrete("degree", map[interface{}]float64{
			2: 1.0,
			3: 1.0,
			4: 1.0,
			5: 1.0,
		}),
		generators.NewUniform("coef0", -10, 10),
	}

	//onetoonehundred := map[interface{}]float64{}
	//for i := 1; i <= 1000; i++ {
	//	onetoonehundred[float64(i)] = 1.0
	//}
	//
	//restrictions := []generators.GenerationStrategy{
	//	generators.NewDiscrete("x", onetoonehundred),
	//}

	match := 0
	var globalTries = 0

	OptResults := make([]OptimizationOutput, noOfExperiments)

	for expIndex := 0; expIndex < noOfExperiments; expIndex++ {

		generator :=
			generators.NewRandom(restrictions, maxAttempts, W, algorithm)

		// channel used by workers to communicate their results
		resultsChans := make(chan functions.Sample, W)

		for w := 0; w < W; w++ {

			localvargs := map[string]interface{}{}
			for k, v := range vargs {
				localvargs[k] = v
			}

			go func(w int) {
				i, p, v, gv, o := DMaximize(targetFunction, localvargs, generator, maxAttempts/W, w, true)
				if !silent {
					fmt.Println("Worker ", w, " MAX --> ", i, p, v, gv, o)
				}

				resultsChans <- functions.Sample{i, p, v, gv, o == 0}
			}(w)
		}

		// Collect results
		results := make([]functions.Sample, W)
		totalTries := 0
		optim, goptim := -math.MaxFloat64, -math.MaxFloat64
		var point functions.MultidimensionalPoint
		for i := 0; i < W; i++ {
			results[i] = <-resultsChans
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

		OptResults[expIndex] = OptimizationOutput{optim, goptim, point, totalTries}

		if optim == goptim {
			match++
			globalTries += totalTries
			fmt.Println("+", expIndex, match, totalTries, point.PrettyPrint(), optim, goptim)
		} else {
			globalTries += totalTries
			fmt.Println("-", expIndex, match, totalTries, point.PrettyPrint(), optim, goptim)
		}
	}

	best, gbest, avg, std := 0.0, 0.0, 0.0, 0.0
	for expIndex := 0; expIndex < noOfExperiments; expIndex++ {
		avg += OptResults[expIndex].GOptim / float64(noOfExperiments)
		if best < OptResults[expIndex].Optim {
			best = OptResults[expIndex].Optim
		}
		if gbest < OptResults[expIndex].GOptim {
			gbest = OptResults[expIndex].GOptim
		}
	}
	for expIndex := 0; expIndex < noOfExperiments; expIndex++ {
		std += (OptResults[expIndex].GOptim - avg) * (OptResults[expIndex].GOptim - avg) / float64(noOfExperiments)
	}
	std = math.Sqrt(std)

	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Results matched on %d (%f) cases", match, float64(match)/float64(noOfExperiments)))
	avgTrials := float64(globalTries) / float64(noOfExperiments)
	fmt.Println(fmt.Sprintf("Average number of attempts %f (%f percent faster) ", avgTrials,
		(float64(maxAttempts)-avgTrials)/float64(maxAttempts)*100))
	fmt.Println(fmt.Sprintf("Optimisation best and global best results are %f, %f", best, gbest))
	fmt.Println(fmt.Sprintf("Optimisation average result and standard deviation are %f, %f", avg, std))
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

	k := int(math.Max(1, float64(n)/(2*math.E)))
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
					if accept(optimNo) {
						//if acceptAll() {
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

func acceptAll() bool {
	return true
}

func accept(optimNo int) bool {
	s := rand.NewSource(time.Now().UnixNano())
	return rand.New(s).Float64() < 0.5+(0.1*float64(optimNo))
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
