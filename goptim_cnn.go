package main

import (
	"github.com/acflorea/goptim/core"
	"github.com/acflorea/goptim/generators"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"flag"
	"strconv"
	"math"
)

func main() {

	noOfExperimentsPtr := flag.Int("noOfExperiments", 1, "Number of experiments.")
	silentPtr := flag.Bool("silent", true, "Silent Mode.")
	maxAttemptsPtr := flag.Int("maxAttempts", 300, "Maximum number of trials in an experiment")

	command := flag.String("command", "python", "Target command.")
	fct := flag.String("fct", "Keras", "Target function")
	alg := flag.String("alg", "Parametrization", "Parallel random generator strategy")

	script := flag.String("script", "", "External script to run")

	test := flag.String("test", "False", "Test Mode.")

	workers := flag.Int("w", 1, "Number of goroutines")

	targetstop := flag.Int("targetstop", 300, "Target stop")

	adjustSingleValue := flag.Bool("adjustSingleValue", true, "Adjust Single Value")

	flag.Parse()

	vargs := map[string]interface{}{}
	vargs["noOfExperiments"] = *noOfExperimentsPtr
	vargs["silent"] = *silentPtr
	vargs["maxAttempts"] = *maxAttemptsPtr
	vargs["command"] = *command
	vargs["fct"] = *fct
	vargs["alg"] = *alg
	vargs["test"] = *test
	vargs["script"] = *script
	vargs["workers"] = *workers
	vargs["targetstop"] = *targetstop
	vargs["adjustSingleValue"] = *adjustSingleValue

	Optimize_cnn(vargs)

}

func Optimize_cnn(vargs map[string]interface{}) {

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

	// Generators

	// Number of convolutional layers from 3 to 6
	conv_layers_map := make(map[interface{}]float64)
	for i := 3; i <= 6; i++ {
		conv_layers_map[i] = 1.0
	}
	conv_layers := generators.NewDiscrete("conv_layers", conv_layers_map)

	// Number of maps in a convolutional layer from 8 to 512
	maps_map := make(map[interface{}]float64)
	for i := 8; i <= 512; i++ {
		maps_map[i] = 1.0
	}

	// Number of fully connected layers from 1 to 4
	full_layers_map := make(map[interface{}]float64)
	for i := 1; i <= 4; i++ {
		full_layers_map[i] = 1.0
	}
	full_layers := generators.NewDiscrete("full_layers", full_layers_map)

	// Number of neurons in fully connected layer from 5 to 2048
	neurons_map := make(map[interface{}]float64)
	for i := 5; i <= 2048; i++ {
		neurons_map[i] = 1.0
	}

	restrictions := []generators.GenerationStrategy{conv_layers, full_layers}

	for i := 1; i <= 6; i++ {
		restrictions = append(restrictions, generators.NewDiscrete("maps_"+strconv.Itoa(i), maps_map))
	}
	for i := 1; i <= 4; i++ {
		restrictions = append(restrictions, generators.NewDiscrete("neurons_"+strconv.Itoa(i), neurons_map))
	}

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
		0.0,
		maxAttempts,
		targetstop,
		W,
		algorithm,
		targetFunction,
		silent,
		vargs)

}
