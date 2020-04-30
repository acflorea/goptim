package core

import (
	"fmt"
	"github.com/acflorea/goptim/functions"
	"github.com/acflorea/goptim/generators"
	"github.com/bluele/slack"
	"log"
	"math"
	"math/rand"
	"time"
)

// The result of one trial
type OptimizationOutput struct {
	Optim  float64
	GOptim float64
	X      functions.MultidimensionalPoint
	Trials int
}

func Optimize(noOfExperiments int,
	restrictions []generators.GenerationStrategy,
	probabilityToChange []float64,
	adjustSingleValue bool,
	optimalSlicePercent float64,
	maxAttempts int,
	targetstop int,
	W int,
	algorithm generators.Algorithm,
	targetFunction functions.NumericalFunction,
	silent bool,
	vargs map[string]interface{}) map[string]interface{} {

	start := time.Now()

	match := 0
	early := 0
	var globalTries = 0

	OptResults := make([]OptimizationOutput, noOfExperiments)

	for expIndex := 0; expIndex < noOfExperiments; expIndex++ {

		tuningTrials := int(math.Max(1, float64(targetstop)/(math.E)))
		//tuningTrials := maxAttempts
		generator :=
			generators.NewRandom(restrictions, probabilityToChange, adjustSingleValue, optimalSlicePercent, maxAttempts, tuningTrials, W, algorithm)

		// channel used by workers to communicate their results
		resultsChans := make(chan functions.Sample, W)

		for w := 0; w < W; w++ {

			localvargs := map[string]interface{}{}
			for k, v := range vargs {
				localvargs[k] = v
			}

			go func(w int, ch chan functions.Sample) {

				// Add the worker id to the args map
				localvargs["workerId"] = w

				i, p, v, gv, o := DMaximize(targetFunction, localvargs, generator, targetstop/W, maxAttempts/W, w, true)
				if !silent {
					fmt.Println("Worker ", w, " MAX --> ", i, p, v, gv, o)
				}

				ch <- functions.Sample{i, p, v, gv, o == 0}
			}(w, resultsChans)
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

		globalTries += totalTries

		if totalTries < maxAttempts {
			early++
			if optim == goptim {
				match++
			}
		}

		if !silent {
			if optim == goptim {
				fmt.Println("+", expIndex, match, totalTries, point.PrettyPrint(), optim, goptim)
			} else {
				fmt.Println("-", expIndex, match, totalTries, point.PrettyPrint(), optim, goptim)
			}
		}
	}

	fmt.Println()
	fmt.Println("Optimisation done. Computing results")

	best, gbest, avg, gavg, std, gstd := 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	for expIndex := 0; expIndex < noOfExperiments; expIndex++ {
		// fmt.Print(OptResults[expIndex].GOptim, ",")
		avg += OptResults[expIndex].Optim / float64(noOfExperiments)
		gavg += OptResults[expIndex].GOptim / float64(noOfExperiments)
		if expIndex == 0 || best < OptResults[expIndex].Optim {
			best = OptResults[expIndex].Optim
		}
		if expIndex == 0 || gbest < OptResults[expIndex].GOptim {
			gbest = OptResults[expIndex].GOptim
		}
	}
	for expIndex := 0; expIndex < noOfExperiments; expIndex++ {
		std += (OptResults[expIndex].Optim - avg) * (OptResults[expIndex].Optim - avg) / float64(noOfExperiments)
		gstd += (OptResults[expIndex].Optim - gavg) * (OptResults[expIndex].Optim - gavg) / float64(noOfExperiments)
	}
	std = math.Sqrt(std)
	gstd = math.Sqrt(gstd)

	elapsed := time.Since(start)
	earlyStopPercent := float64(early) / float64(noOfExperiments)

	fmt.Println()
	fmt.Println(fmt.Sprintf("Early stop in %d (%f) cases", early, earlyStopPercent))

	matchStopPercent := float64(match) / float64(early)
	matchPercent := float64(match) / float64(noOfExperiments)
	fmt.Println(fmt.Sprintf("Results matched on %d (%f) early stoping cases meaning %f of the total trials",
		match, matchStopPercent, matchPercent))

	avgTrials := float64(globalTries) / float64(noOfExperiments)
	fmt.Println(fmt.Sprintf("Average number of attempts %f (%f percent faster) ", avgTrials,
		(float64(maxAttempts)-avgTrials)/float64(maxAttempts)*100))

	fmt.Println(fmt.Sprintf("Optimisation best and global best results are %f, %f", best, gbest))

	fmt.Println(fmt.Sprintf("(ES) Optimisation average result and standard deviation are %f, %f", avg, std))

	fmt.Println(fmt.Sprintf("(GB) Optimisation average result and standard deviation are %f, %f", gavg, gstd))

	fmt.Println(fmt.Sprintf("Optimization took %s", elapsed))

	results := make(map[string]interface{})
	results["earlyStopPercent"] = earlyStopPercent
	results["matchStopPercent"] = matchStopPercent
	results["matchPercent"] = matchPercent
	results["avg"] = avg
	results["std"] = std
	results["optimalSlicePercent"] = optimalSlicePercent

	fmt.Println("[optimalSlicePercent, earlyStopPercent, matchStopPercent, matchPercent, avg, std]")
	fmt.Println(fmt.Sprintf("[%f, %f, %f, %f, %f, %f]",
		optimalSlicePercent, earlyStopPercent, matchStopPercent, matchPercent, avg, std))

	return results
}

// Attempts to dynamically minimize the function f
// k := n / (2 * math.E)
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
// w is thw worker index
func DMinimize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, n, N, w int, goAllTheWay bool) (
	index int,
	p functions.MultidimensionalPoint,
	min float64,
	gmin float64,
	optimNo int) {

	k := int(math.Max(1, float64(n)/math.E))
	return Minimize(f, vargs, generator, k, N, w, goAllTheWay)
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
func Minimize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, k, N, w int, goAllTheWay bool) (
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

	api, slackEnabled := vargs["slackAPI"].(*slack.Slack)
	slackChannel, ok := vargs["slackChannel"].(string)
	if !ok {
		slackChannel = "goptim-updates"
	}
	if slackEnabled {
		err := api.ChatPostMessage(slackChannel, fmt.Sprintf("[w=%d] Optimization Start", w), nil)
		if err != nil {
			log.Println("Problem connecting to Slack ", err)
		}

		defer func() {
			err = api.ChatPostMessage(slackChannel, fmt.Sprintf("[w=%d] Optimization Stop", w), nil)
			if err != nil {
				log.Println("Problem connecting to Slack ", err)
			}
		}()
	}

	state := generators.GeneratorState{
		[]functions.MultidimensionalPoint{},
		[]functions.TwoDPointVector{},
		[]float64{},
		functions.MultidimensionalPoint{}}

	for i := 0; i < N; i++ {

		rndPoint, newState := generator.Next(w, state)
		f_rnd, _ := f(rndPoint, vargs)
		centroid := newState.Centroid

		if slackEnabled {
			err := api.ChatPostMessage(slackChannel, fmt.Sprintf("[w=%d] %s", w, functions.FloatToString(f_rnd)+" :: "+rndPoint.PrettyPrint()), nil)
			if err != nil {
				log.Println("Problem connecting to Slack ", err)
			}
		}

		if i == 0 {
			// in case the centroid was not initialized
			centroid = rndPoint
		}

		if minReached {
			if f_rnd < gmin {
				gmin = f_rnd
			}
		} else {
			if f_rnd < min {
				// the centroid is tmpOpt
				centroid = rndPoint

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

		state = generators.GeneratorState{
			newState.GeneratedPoints,
			newState.Statistics,
			append(newState.Output, f_rnd),
			centroid}
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
	return rand.New(s).Float64() < 0.1+(0.1*float64(optimNo))
}

// Dynamically Minimizes the negation of the target function
func DMaximize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, n, N, w int, goAllTheWay bool) (
	index int,
	p functions.MultidimensionalPoint,
	max float64,
	gmax float64,
	optimNo int) {

	index, p, max, gmax, optimNo = DMinimize(functions.Negate(f), vargs, generator, n, N, w, goAllTheWay)
	return index, p, -max, -gmax, optimNo
}

// Minimizes the negation of the target function
func Maximize(f functions.NumericalFunction, vargs map[string]interface{}, generator generators.Generator, k, n, N, w int, goAllTheWay bool) (
	index int,
	p functions.MultidimensionalPoint,
	max float64,
	gmax float64,
	optimNo int) {

	index, p, max, gmax, optimNo = Minimize(functions.Negate(f), vargs, generator, k, N, w, goAllTheWay)
	return index, p, -max, -gmax, optimNo
}
