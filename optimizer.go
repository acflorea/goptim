package main

import (
	"github.com/acflorea/goptim/generators"
	"strconv"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"github.com/acflorea/goptim/core"
	"flag"
	"strings"
)

func main() {

	noOfExperimentsPtr := flag.Int("noOfExperiments", 10000, "Number of experiments.")
	maxAttemptsPtr := flag.Int("maxAttempts", 1000, "Maximum number of trials in an experiment")
	targetstopPtr := flag.Int("targetstop", 0, "Target stop")

	// Hyperopt specifics
	probsPtr := flag.String("probs", "1 1 1 1 1 1", "Probabilities to change each value")
	grievankPtr := flag.Int("grievank", 6, "Number of variables in Grievank function")

	flag.Parse()

	noOfExperiments := *noOfExperimentsPtr

	grievank := *grievankPtr

	var restrictions []generators.GenerationStrategy
	for i := 0; i < grievank; i++ {
		restrictions = append(restrictions, generators.NewUniform("x"+strconv.Itoa(i+1), -600.0, 600.0))
	}

	// We target a stop after targetstop attempts
	var probabilityToChange = []float32{}
	allOnes := true
	probsStr := *probsPtr
	for _, prob := range strings.Split(probsStr, " ") {
		if n, err := strconv.ParseFloat(prob, 32); err == nil {
			prob := float32(n)
			if prob != 1.0 {
				allOnes = false
			}
			probabilityToChange = append(probabilityToChange, prob)

		}
	}

	adjustSingleValue := false

	maxAttempts := *maxAttemptsPtr

	targetstop := *targetstopPtr
	if targetstop == 0 {
		targetstop = maxAttempts
	}

	Ws := []int{1, 2, 3, 4, 5, 6, 7, 8}

	// loop through algorithms ?
	algorithm := generators.Algorithms["Parametrization"]

	targetFunction := functions.Functions["F_Griewank"]

	silent := true

	vargs := make(map[string]interface{})
	vargs["verbose"] = false
	vargs["grievank"] = grievank

	for optimalSlicePercent := 0; optimalSlicePercent < 75; optimalSlicePercent++ {

		for _, W := range Ws {

			fmt.Println(fmt.Sprintf("\n\n"+
				"+++++++  optimalSlicePercent=%d, W=%d  +++++++", optimalSlicePercent, W))

			core.Optimize(
				noOfExperiments,
				restrictions,
				probabilityToChange,
				adjustSingleValue,
				float64(optimalSlicePercent),
				maxAttempts,
				targetstop,
				W,
				algorithm,
				targetFunction,
				silent,
				vargs)

		}

		if allOnes {
			// no point to test "probabilityToChange" if all probabilities are ones
			break
		}
	}

}
