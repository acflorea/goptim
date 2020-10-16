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
	//batchsize [int]: [32, 64, 128], 256
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
	//cnneurons [int]: 16 to 64 step = 4
	cnn_neurons_map := make(map[interface{}]float64)
	for i := 16; i <= 64; {
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

	restrictions := []generators.GenerationStrategy{learning_rate_exp, batch_size, x_layers, cnn_neurons, fc_neurons, dropout1, dropout2}

	useRandomSample := vargs["useRandomSample"].(bool)
	if useRandomSample {
		seed := 123456
		vargs["seed"] = seed
		restrictions = append(restrictions, generators.NewUniform("seed", 0, 10000000))
	}

	//Sum of fractions for main effects 53.39%
	//	Sum of fractions for pairwise interaction effects 34.69%
	//-1.96% due to interaction: X2 x X0
	//-0.77% due to interaction: X4 x X2
	//-0.37% due to interaction: X6 x X2
	//-0.05% due to interaction: X5 x X2
	//0.02% due to interaction: X2 x X1
	//0.04% due to interaction: X3 x X2
	//0.53% due to interaction: X3 x X1
	//0.75% due to interaction: X6 x X4
	//0.80% due to main effect: X2
	//0.84% due to interaction: X6 x X3
	//1.03% due to interaction: X6 x X1
	//1.11% due to interaction: X3 x X0
	//1.12% due to interaction: X6 x X0
	//1.13% due to interaction: X5 x X1
	//1.36% due to interaction: X6 x X5
	//1.56% due to main effect: X6
	//2.22% due to interaction: X4 x X1
	//2.76% due to interaction: X5 x X3
	//3.70% due to interaction: X5 x X4
	//4.33% due to main effect: X1
	//4.50% due to interaction: X4 x X3
	//4.52% due to interaction: X1 x X0
	//5.01% due to main effect: X3
	//5.42% due to interaction: X4 x X0
	//6.77% due to interaction: X5 x X0
	//11.54% due to main effect: X4
	//12.54% due to main effect: X5
	//17.63% due to main effect: X0

	//learning_rate_exp, batch_size, x_layers, cnn_neurons, fc_neurons, dropout1, dropout2
	//17.63% due to main effect: X0 - learning_rate_exp
	//12.54% due to main effect: X5 - dropout1
	//11.54% due to main effect: X4 - fc_neurons
	//5.01% due to main effect: X3 - cnn_neurons
	//4.33% due to main effect: X1 - batch_size
	//1.56% due to main effect: X6 - dropout2
	//0.80% due to main effect: X2 - x_layers

	// fANOVA - list them here for brevity...
	// learning_rate_exp
	var x0 = 17.63
	// dropout1
	var x5 = 12.54
	// fc_neurons
	var x4 = 11.54
	// cnn_neurons
	var x3 = 5.01
	// batch_size
	var x1 = 4.33
	// dropout2
	var x6 = 1.56
	// x_layers
	var x2 = 0.8

	// x0, x1, x2, x3, x4, x5 = 1.0, 1.0, 1.0, 1.0, 1.0, 1.0

	// ...
	var probabilityToChange = []float64{x0, x1, x2, x3, x4, x5, x6}

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
