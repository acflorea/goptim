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
	slackChannelPtr := flag.String("slackChannel", "hpc39-updates", "Token to connect to Slack")

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
	//learning_rate [float] : 10^-4 to 10^-1.886 python code: 10**np.random.uniform(-4,-1.886)
	learning_rate_exp := generators.NewUniform("learning_rate_exp", -4, -1.886)

	//batchsize [int]: [32, 64, 128]
	//batchsize [int]: [32, 64, 128, 256]
	batch_size_map := make(map[interface{}]float64)
	batch_size_map[32] = 1.0
	batch_size_map[64] = 1.0
	batch_size_map[128] = 1.0
	batch_size_map[256] = 1.0
	batch_size := generators.NewDiscrete("batch_size", batch_size_map)

	//x3layers [int] : [0, 1]
	x_layers_map := make(map[interface{}]float64)
	x_layers_map[0] = 1.0
	x_layers_map[1] = 1.0
	x_layers := generators.NewDiscrete("x_layers", x_layers_map)

	//cnneurons [int]: 16 to 32 step = 4
	//cnneurons1 [int]: 16 to 64 step = 4
	cnn_neurons1_map := make(map[interface{}]float64)
	for i := 16; i <= 64; {
		cnn_neurons1_map[i] = 1.0
		i = i + 4
	}
	cnn_neurons1 := generators.NewDiscrete("cnn_neurons1", cnn_neurons1_map)

	//cnneurons2 [int]: 16 to 128 step = 4
	cnn_neurons2_map := make(map[interface{}]float64)
	for i := 16; i <= 128; {
		cnn_neurons2_map[i] = 1.0
		i = i + 4
	}
	cnn_neurons2 := generators.NewDiscrete("cnn_neurons2", cnn_neurons2_map)

	//cnneurons3 [int]: 16 to 256 step = 4
	cnn_neurons3_map := make(map[interface{}]float64)
	for i := 16; i <= 256; {
		cnn_neurons3_map[i] = 1.0
		i = i + 4
	}
	cnn_neurons3 := generators.NewDiscrete("cnn_neurons3", cnn_neurons3_map)

	//fcneurons [int]: 256 to 2048 step = 24
	fc_neurons_map := make(map[interface{}]float64)
	for i := 256; i <= 2048; {
		fc_neurons_map[i] = 1.0
		i = i + 24
	}
	fc_neurons := generators.NewDiscrete("fc_neurons", fc_neurons_map)

	//dropout1 [float]: 0 to 0.5 step = 0.1
	//dropout1 [float]: 0 to 0.5 step = 0.05
	dropout1_map := make(map[interface{}]float64)
	for i := 0.0; i <= 0.5; {
		dropout1_map[i] = 1.0
		i = i + 0.05
	}
	dropout1 := generators.NewDiscrete("dropout1", dropout1_map)

	//dropout2 [float]: 0 to 0.5 step = 0.1
	//dropout2 [float]: 0 to 0.5 step = 0.05
	dropout2_map := make(map[interface{}]float64)
	for i := 0.0; i <= 0.5; {
		dropout2_map[i] = 1.0
		i = i + 0.05
	}
	dropout2 := generators.NewDiscrete("dropout2", dropout2_map)

	//dropout3 [float]: 0 to 0.5 step = 0.05
	dropout3_map := make(map[interface{}]float64)
	for i := 0.0; i <= 0.5; {
		dropout3_map[i] = 1.0
		i = i + 0.05
	}
	dropout3 := generators.NewDiscrete("dropout3", dropout3_map)

	//dropout4 [float]: 0 to 0.5 step = 0.05
	dropout4_map := make(map[interface{}]float64)
	for i := 0.0; i <= 0.5; {
		dropout4_map[i] = 1.0
		i = i + 0.05
	}
	dropout4 := generators.NewDiscrete("dropout4", dropout4_map)

	restrictions := []generators.GenerationStrategy{learning_rate_exp, batch_size, x_layers, cnn_neurons1, cnn_neurons2, cnn_neurons3,
		fc_neurons, dropout1, dropout2, dropout3, dropout4}

	useRandomSample := vargs["useRandomSample"].(bool)
	if useRandomSample {
		seed := 123456
		vargs["seed"] = seed
		restrictions = append(restrictions, generators.NewUniform("seed", 0, 10000000))
	}

	//Sum of fractions for main effects 51.13%
	//Sum of fractions for pairwise interaction effects 17.78%
	//-4.46% due to interaction: X5 x X10
	//-3.40% due to interaction: X9 x X10
	//-2.85% due to interaction: X4 x X10
	//-2.42% due to interaction: X5 x X0
	//-2.22% due to interaction: X8 x X10
	//-0.83% due to interaction: X4 x X0
	//-0.59% due to interaction: X10 x X1
	//-0.57% due to interaction: X2 x X10
	//-0.48% due to interaction: X7 x X10
	//-0.48% due to interaction: X3 x X10
	//-0.47% due to interaction: X6 x X10
	//-0.40% due to interaction: X10 x X0
	//-0.13% due to interaction: X3 x X0
	//-0.11% due to interaction: X2 x X0
	//-0.10% due to interaction: X8 x X0
	//-0.03% due to interaction: X6 x X0
	//0.04% due to main effect: X10
	//0.04% due to interaction: X1 x X0
	//0.08% due to interaction: X7 x X0
	//0.09% due to interaction: X6 x X2
	//0.12% due to interaction: X7 x X6
	//0.17% due to interaction: X2 x X1
	//0.19% due to interaction: X4 x X2
	//0.20% due to interaction: X7 x X4
	//0.23% due to interaction: X6 x X1
	//0.25% due to interaction: X8 x X6
	//0.29% due to interaction: X9 x X1
	//0.33% due to interaction: X9 x X0
	//0.36% due to interaction: X7 x X3
	//0.39% due to interaction: X9 x X2
	//0.49% due to interaction: X7 x X1
	//0.52% due to interaction: X8 x X4
	//0.53% due to interaction: X6 x X4
	//0.55% due to interaction: X4 x X3
	//0.56% due to interaction: X3 x X2
	//0.62% due to interaction: X9 x X7
	//0.65% due to interaction: X7 x X5
	//0.69% due to main effect: X0
	//0.75% due to interaction: X3 x X1
	//0.83% due to interaction: X7 x X2
	//0.86% due to interaction: X4 x X1
	//0.92% due to interaction: X6 x X3
	//0.98% due to interaction: X8 x X3
	//0.99% due to interaction: X9 x X6
	//1.00% due to interaction: X8 x X5
	//1.03% due to interaction: X9 x X4
	//1.08% due to interaction: X8 x X1
	//1.28% due to interaction: X6 x X5
	//1.31% due to interaction: X5 x X2
	//1.53% due to interaction: X5 x X1
	//1.55% due to interaction: X5 x X4
	//1.67% due to interaction: X9 x X3
	//2.30% due to main effect: X3
	//2.44% due to interaction: X8 x X2
	//2.59% due to main effect: X2
	//2.72% due to interaction: X9 x X8
	//2.72% due to main effect: X6
	//2.89% due to interaction: X5 x X3
	//3.37% due to interaction: X8 x X7
	//3.45% due to interaction: X9 x X5
	//3.77% due to main effect: X7
	//4.13% due to main effect: X1
	//6.12% due to main effect: X8
	//6.59% due to main effect: X9
	//8.32% due to main effect: X4
	//13.85% due to main effect: X5

	//0.04% due to main effect: X10
	//0.69% due to main effect: X0
	//2.30% due to main effect: X3
	//2.59% due to main effect: X2
	//2.72% due to main effect: X6
	//3.77% due to main effect: X7
	//4.13% due to main effect: X1
	//6.12% due to main effect: X8
	//6.59% due to main effect: X9
	//8.32% due to main effect: X4
	//13.85% due to main effect: X5

	//X0-batch_size,X1-cnn_neurons1,X2-cnn_neurons2,X3-cnn_neurons3,
	//X4-dropout1,X5-dropout2,X6-dropout3,X7-dropout4,
	//X8-fc_neurons,X9-learning_rate_exp,X10-x_layers

	//13.85% due to main effect: X5 dropout2
	//8.32% due to main effect: X4 dropout1
	//6.59% due to main effect: X9 learning_rate_exp
	//6.12% due to main effect: X8 fc_neurons
	//4.13% due to main effect: X1 cnn_neurons1
	//3.77% due to main effect: X7 dropout4
	//2.72% due to main effect: X6 dropout3
	//2.59% due to main effect: X2 cnn_neurons2
	//2.30% due to main effect: X3 cnn_neurons3
	//0.69% due to main effect: X0 batch_size
	//0.04% due to main effect: X10 x_layers

	// fANOVA - list them here for brevity...
	//13.85% due to main effect: X5 dropout2
	x5 := 13.85
	//8.32% due to main effect: X4 dropout1
	x4 := 8.32
	//6.59% due to main effect: X9 learning_rate_exp
	x9 := 6.59
	//6.12% due to main effect: X8 fc_neurons
	x8 := 6.12
	//4.13% due to main effect: X1 cnn_neurons1
	x1 := 4.13
	//3.77% due to main effect: X7 dropout4
	x7 := 3.77
	//2.72% due to main effect: X6 dropout3
	x6 := 2.72
	//2.59% due to main effect: X2 cnn_neurons2
	x2 := 2.59
	//2.30% due to main effect: X3 cnn_neurons3
	x3 := 2.30
	//0.69% due to main effect: X0 batch_size
	x0 := 0.69
	//0.04% due to main effect: X10 x_layers
	x10 := 0.04

	//x0, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10 = 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0

	// ...
	var probabilityToChange = []float64{x0, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10}

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
