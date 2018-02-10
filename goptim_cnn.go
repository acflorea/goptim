package main

import (
	"github.com/acflorea/goptim/core"
	"github.com/acflorea/goptim/generators"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"flag"
)

func main() {

	noOfExperimentsPtr := flag.Int("noOfExperiments", 1, "Number of experiments.")
	silentPtr := flag.Bool("silent", true, "Silent Mode.")
	maxAttemptsPtr := flag.Int("maxAttempts", 300, "Maximum number of trials in an experiment")

	command := flag.String("command", "python", "Target command.")
	fct := flag.String("fct", "Keras", "Target function")
	alg := flag.String("alg", "SeqSplit", "Parallel random generator strategy")

	script := flag.String("script", "", "External script to run")

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

	restrictions := []generators.GenerationStrategy{}

	core.Optimize(noOfExperiments, restrictions, maxAttempts, targetstop, W, algorithm, targetFunction, silent, vargs)

}
