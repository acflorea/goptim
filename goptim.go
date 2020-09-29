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

	optimize_hpc39(vargs)

}

func optimize_hpc39(vargs map[string]interface{}) {

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

	//learning_rate [float] : 10^-2.2218 to 10^-2.045 python code: 10**np.random.uniform(-2.2218,-2.045)
	learning_rate_exp := generators.NewUniform("learning_rate_exp", -2.2218, -2.045)

	//batchsize [int]: [32, 64, 128]
	batch_size_map := make(map[interface{}]float64)
	batch_size_map[32] = 1.0
	batch_size_map[64] = 1.0
	batch_size_map[128] = 1.0
	batch_size := generators.NewDiscrete("batch_size", batch_size_map)

	//x3layers [int] : [0, 1]
	x_layers_map := make(map[interface{}]float64)
	x_layers_map[0] = 1.0
	x_layers_map[1] = 1.0
	x_layers := generators.NewDiscrete("x_layers", x_layers_map)

	//cnneurons [int]: 16 to 32 step = 4
	cnn_neurons_map := make(map[interface{}]float64)
	for i := 16; i <= 32; {
		cnn_neurons_map[i] = 1.0
		i = i + 4
	}
	cnn_neurons := generators.NewDiscrete("cnn_neurons", cnn_neurons_map)

	//fcneurons [int]: 256 to 2048 step = 24
	fc_neurons_map := make(map[interface{}]float64)
	for i := 256; i <= 2048; {
		fc_neurons_map[i] = 1.0
		i = i + 24
	}
	fc_neurons := generators.NewDiscrete("fc_neurons", fc_neurons_map)

	//dropout1 [float]: 0 to 0.5 step = 0.1
	dropout1_map := make(map[interface{}]float64)
	for i := 0.0; i <= 0.5; {
		dropout1_map[i] = 1.0
		i = i + 0.1
	}
	dropout1 := generators.NewDiscrete("dropout1", dropout1_map)

	//dropout1 [float]: 0 to 0.5 step = 0.1
	dropout2_map := make(map[interface{}]float64)
	for i := 0.0; i <= 0.5; {
		dropout2_map[i] = 1.0
		i = i + 0.1
	}
	dropout2 := generators.NewDiscrete("dropout2", dropout2_map)

	restrictions := []generators.GenerationStrategy{learning_rate_exp, batch_size, x_layers, cnn_neurons, fc_neurons, dropout1, dropout2}

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
