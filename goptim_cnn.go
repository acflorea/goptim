package main

import (
	"github.com/acflorea/goptim/core"
	"github.com/acflorea/goptim/generators"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"flag"
	"strconv"
)

func main() {

	noOfExperimentsPtr := flag.Int("noOfExperiments", 1, "Number of experiments.")
	silentPtr := flag.Bool("silent", true, "Silent Mode.")
	maxAttemptsPtr := flag.Int("maxAttempts", 300, "Maximum number of trials in an experiment")

	command := flag.String("command", "python", "Target command.")
	fct := flag.String("fct", "Keras", "Target function")
	alg := flag.String("alg", "SeqSplit", "Parallel random generator strategy")

	script := flag.String("script", "", "External script to run")

	test := flag.String("test", "False", "Test Mode.")

	workers := flag.Int("w", 1, "Number of goroutines")

	targetstop := flag.Int("targetstop", 300, "Target stop")

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

	// Generators

	// Number of convolutional layers from 3 to 50/15 (after 15 seems to crash.... no more dimensions)
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

	for i := 3; i <= 50; i++ {
		restrictions = append(restrictions, generators.NewDiscrete("maps_"+strconv.Itoa(i), maps_map))
	}
	for i := 1; i <= 4; i++ {
		restrictions = append(restrictions, generators.NewDiscrete("neurons_"+strconv.Itoa(i), neurons_map))
	}

	core.Optimize(noOfExperiments, restrictions, maxAttempts, targetstop, W, algorithm, targetFunction, silent, vargs)

}
