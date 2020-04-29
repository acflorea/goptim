package main

import (
	"flag"
	"fmt"
	"github.com/acflorea/goptim/core"
	"github.com/acflorea/goptim/functions"
	"github.com/acflorea/goptim/generators"
	"math"
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
	noOfExperimentsPtr := flag.Int("noOfExperiments", 1, "Number of experiments.")
	silentPtr := flag.Bool("silent", true, "Silent Mode.")
	maxAttemptsPtr := flag.Int("maxAttempts", 300, "Maximum number of trials in an experiment")
	fct := flag.String("fct", "F_identity", "Target function")
	alg := flag.String("alg", "SeqSplit", "Parallel random generator strategy")
	script := flag.String("script", "", "External script to run")
	command := flag.String("command", "", "External program to execute")
	workers := flag.Int("w", 8, "Number of goroutines")
	targetstop := flag.Int("targetstop", 0, "Target stop")

	flag.Parse()

	vargs := map[string]interface{}{}
	vargs["fileName"] = *fileNamePtr
	vargs["noOfExperiments"] = *noOfExperimentsPtr
	vargs["silent"] = *silentPtr
	vargs["maxAttempts"] = *maxAttemptsPtr
	vargs["fct"] = *fct
	vargs["alg"] = *alg
	vargs["script"] = *script
	vargs["command"] = *command
	vargs["workers"] = *workers
	vargs["targetstop"] = *targetstop

	vargs["adjustSingleValue"] = false
	vargs["optimalSlicePercent"] = 1.0

	optimize_k7m(vargs)

}

func optimize_k7m(vargs map[string]interface{}) {

	fmt.Println("Optimization start!")
	fmt.Println(vargs)

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
	W := vargs["workers"].(int)

	// We target a stop after targetstop attempts
	targetstop := vargs["targetstop"].(int)
	if targetstop == 0 {
		targetstop = maxAttempts
	}

	// if this is true a single value changes for each step
	// otherwise the values are changing according to their probabilities
	adjustSingleValue := vargs["adjustSingleValue"].(bool)

	optimalSlicePercent := vargs["optimalSlicePercent"].(float64)

	// Generators

	//# 1 to 55 increment of 1
	//max_breadth = int(getValue(argumentsDict, '-w', '--max_breadth', 1))
	max_breadth_map := make(map[interface{}]float64)
	for i := 1; i <= 55; i++ {
		max_breadth_map[i] = 1.0
	}
	max_breadth := generators.NewDiscrete("max_breadth", max_breadth_map)

	//# 2 to 55 increment of 1
	//max_depth = int(getValue(argumentsDict, '-d', '--max_depth', 2))
	max_depth_map := make(map[interface{}]float64)
	for i := 2; i <= 55; i++ {
		max_depth_map[i] = 1.0
	}
	max_depth := generators.NewDiscrete("max_depth", max_depth_map)

	//# -50 to 50 increment of 0.1
	//attr_b = float(getValue(argumentsDict, '-b', '--attr_b', -50))
	attr_b := generators.NewUniform("attr_b", -50, 50)

	//# -1 to 1 increment of 0.1
	//attr_c = float(getValue(argumentsDict, '-c', '--attr_c', -1))
	attr_c := generators.NewUniform("attr_c", -1, 1)

	//# 0.1 to 1 increment of 0.1
	//edge_cost = float(getValue(argumentsDict, '-e', '--edge_cost', 0.1))
	edge_cost := generators.NewUniform("edge_cost", 0.1, 1)

	//# 1 to 55 increment of 1
	//movement_factor = int(getValue(argumentsDict, '-m', '--movement_factor', 1))
	movement_factor_map := make(map[interface{}]float64)
	for i := 1; i <= 55; i++ {
		movement_factor_map[i] = 1.0
	}
	movement_factor := generators.NewDiscrete("movement_factor", movement_factor_map)

	restrictions := []generators.GenerationStrategy{max_breadth, max_depth, attr_b, attr_c, edge_cost, movement_factor}

	//7.40% due to main effect: X0
	//11.85% due to main effect: X1
	//0.51% due to main effect: X2
	//0.79% due to main effect: X3
	//1.62% due to main effect: X4
	//0.73% due to main effect: X5
	//2.26% due to main effect: X6
	//1.26% due to main effect: X7
	//26.28% due to main effect: X8
	//0.87% due to main effect: X9
	//3.22% due to main effect: X10
	//1.75% due to main effect: X11

	// fANOVA - list them here for brevity...
	var x8 = 26.28 * math.Log(2043.0) // neurons_1 5...2048
	var x1 = 11.85 * math.Log(3.0)    // full_layers 1...4 !!!
	var x0 = 7.40 * math.Log(3.0)     // conv_layers  3...6 !!!
	var x10 = 3.22 * math.Log(2043.0) // neurons_3 5...2048
	var x6 = 2.26 * math.Log(504)     // maps_5 8...512
	var x11 = 1.75 * math.Log(2043.0) // neurons_4 5...2048
	var x4 = 1.62 * math.Log(504)     // maps_3 8...512
	var x7 = 1.26 * math.Log(504)     // maps_6 8...512
	var x9 = 0.87 * math.Log(2043.0)  // neurons_2 5...2048
	var x3 = 0.79 * math.Log(504)     // maps_2 5...2048
	var x5 = 0.73 * math.Log(504)     // maps_4 5...2048
	var x2 = 0.51 * math.Log(504)     // maps_1 5...2048

	// conv_layers, full_layers, maps_1, maps_2, maps_3, maps_4, maps_5, maps_6, [neurons_1], neurons_2, neurons_3, neurons_4
	var probabilityToChange = []float64{x0, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11}

	core.Optimize(
		noOfExperiments,
		restrictions,
		probabilityToChange,
		adjustSingleValue,
		optimalSlicePercent,
		maxAttempts,
		targetstop,
		W,
		algorithm,
		targetFunction,
		silent,
		vargs)

}
