package main

import (
	"flag"
	"fmt"
	"github.com/acflorea/goptim/core"
	"github.com/acflorea/goptim/functions"
	"github.com/acflorea/goptim/generators"
	"github.com/bluele/slack"
)

// The result of one trial
type OptimizationOutput struct {
	Optim  float64
	GOptim float64
	X      functions.MultidimensionalPoint
	Trials int
}

func main() {

	slackTokenPtr := flag.String("slackToken", "", "Token to connect to Slack")
	slackChannelPtr := flag.String("slackChannel", "k7m-updates", "Token to connect to Slack")

	fileNamePtr := flag.String("fileName", "", "Name of the input file.")
	targetFolderPtr := flag.String("targetFolder", "", "The folder in which to run.")
	noOfExperimentsPtr := flag.Int("noOfExperiments", 1, "Number of experiments.")
	silentPtr := flag.Bool("silent", true, "Silent Mode.")
	maxAttemptsPtr := flag.Int("maxAttempts", 300, "Maximum number of trials in an experiment")
	fct := flag.String("fct", "F_identity", "Target function")
	alg := flag.String("alg", "SeqSplit", "Parallel random generator strategy")
	script := flag.String("script", "", "External script to run")
	command := flag.String("command", "", "External program to execute")
	workers := flag.Int("w", 8, "Number of goroutines")
	targetstop := flag.Int("targetstop", 0, "Target stop")

	useRandomSamplePtr := flag.Bool("useRandomSample", true, "Use a single random sample instead of whole target space")

	flag.Parse()

	vargs := map[string]interface{}{}
	vargs["fileName"] = *fileNamePtr
	vargs["targetFolder"] = *targetFolderPtr
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
	vargs["optimalSlicePercent"] = 100.0

	vargs["useRandomSample"] = *useRandomSamplePtr

	// Deal with the Slack API
	if *slackTokenPtr != "" {
		api := slack.New(*slackTokenPtr)

		api.JoinChannel(*slackChannelPtr)

		vargs["slackAPI"] = api
		vargs["slackChannel"] = *slackChannelPtr
	}

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
	//# 1 to 9 increment of 1
	//max_breadth = int(getValue(argumentsDict, '-w', '--max_breadth', 1))
	max_breadth_map := make(map[interface{}]float64)
	for i := 1; i <= 9; i++ {
		max_breadth_map[i] = 1.0
	}
	max_breadth := generators.NewDiscrete("max_breadth", max_breadth_map)

	//# 1 to 55 increment of 1
	//# 1 to 9 increment of 1
	//max_depth = int(getValue(argumentsDict, '-d', '--max_depth', 2))
	max_depth_map := make(map[interface{}]float64)
	for i := 1; i <= 9; i++ {
		max_depth_map[i] = 1.0
	}
	max_depth := generators.NewDiscrete("max_depth", max_depth_map)

	//# 0 to 500 increment of 0.1
	//# 0 to 20 increment of 0.1
	// attr_b := generators.NewUniform("attr_b", 0, 500)
	attr_b_map := make(map[interface{}]float64)
	for i := 0.0; i <= 20.0; {
		attr_b_map[i] = 1.0
		i = i + 0.1
	}
	attr_b := generators.NewDiscrete("attr_b", attr_b_map)

	//# 0 to 1 increment of 0.01
	//# 0 to 1 increment of 0.01
	// attr_c := generators.NewUniform("attr_c", -1, 1)
	attr_c_map := make(map[interface{}]float64)
	for i := 0.0; i <= 1.0; {
		attr_c_map[i] = 1.0
		i = i + 0.01
	}
	attr_c := generators.NewDiscrete("attr_c", attr_c_map)

	//# 0.1 to 1 increment of 0.1
	//# 0.2 to 1 increment of 0.1
	// edge_cost := generators.NewUniform("edge_cost", 0.1, 1)
	edge_cost_map := make(map[interface{}]float64)
	for i := 0.2; i <= 1.0; {
		edge_cost_map[i] = 1.0
		i = i + 0.1
	}
	edge_cost := generators.NewDiscrete("edge_cost", edge_cost_map)

	//# 1 to 55 increment of 1
	//# 2 to 10 increment of 1
	movement_factor_map := make(map[interface{}]float64)
	for i := 2; i <= 10; i++ {
		movement_factor_map[i] = 1.0
	}
	movement_factor := generators.NewDiscrete("movement_factor", movement_factor_map)

	restrictions := []generators.GenerationStrategy{max_breadth, max_depth, attr_b, attr_c, edge_cost, movement_factor}

	useRandomSample := vargs["useRandomSample"].(bool)
	if useRandomSample {
		seed := 123456
		vargs["seed"] = seed
		restrictions = append(restrictions, generators.NewUniform("seed", 0, 10000000))
	}

	//Sum of fractions for main effects 54.29%
	//	Sum of fractions for pairwise interaction effects 32.79%
	//0.49% due to interaction: X4 x X3
	//0.97% due to interaction: X5 x X4
	//1.19% due to interaction: X5 x X1
	//1.32% due to interaction: X4 x X0
	//1.48% due to interaction: X5 x X3
	//1.70% due to interaction: X5 x X0
	//1.91% due to interaction: X5 x X2
	//1.94% due to interaction: X3 x X0
	//2.26% due to interaction: X3 x X1
	//2.27% due to interaction: X2 x X0
	//2.86% due to interaction: X3 x X2
	//3.29% due to interaction: X4 x X1
	//3.30% due to interaction: X1 x X0
	//3.50% due to interaction: X4 x X2
	//3.52% due to main effect: X5
	//4.33% due to interaction: X2 x X1
	//7.76% due to main effect: X0
	//8.00% due to main effect: X3
	//9.35% due to main effect: X4
	//10.59% due to main effect: X2
	//15.06% due to main effect: X1

	// fANOVA - list them here for brevity...
	// attr_c
	var x1 = 15.06
	// edge_cost
	var x2 = 10.59
	// max_depth
	var x4 = 9.35
	// max_breadth
	var x3 = 8.0
	// attr_b
	var x0 = 7.76
	// movement_factor
	var x5 = 3.52

	// x0, x1, x2, x3, x4, x5 = 1.0, 1.0, 1.0, 1.0, 1.0, 1.0

	// ...
	var probabilityToChange = []float64{x0, x1, x2, x3, x4, x5}

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
