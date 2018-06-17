package main

import (
	"github.com/acflorea/goptim/core"
	"github.com/acflorea/goptim/generators"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"flag"
	"strings"
	"strconv"
	"time"
)

func main() {

	verbose := flag.Bool("verbose", false, "Talkative or not")

	fileNamePtr := flag.String("fileName", "", "Name of the input file.")
	noOfExperimentsPtr := flag.Int("noOfExperiments", 1, "Number of experiments.")
	silentPtr := flag.Bool("silent", true, "Silent Mode.")
	maxAttemptsPtr := flag.Int("maxAttempts", 300, "Maximum number of trials in an experiment")
	fct := flag.String("fct", "F_identity", "Target function")
	alg := flag.String("alg", "SeqSplit", "Parallel random generator strategy")
	script := flag.String("script", "", "External script to run")
	workers := flag.Int("w", 1, "Number of goroutines")
	targetstop := flag.Int("targetstop", 0, "Target stop")

	// Hyperopt specifics
	probs := flag.String("probs", "", "Probabilities to change each value")
	optSlicePercent := flag.Float64("optSlicePercent", 0.0, "Slice of results considered to be optimal")
	grievank := flag.Int("grievank", 3, "Number of variables in Grievank function")

	// Spark specifics
	sparkMaster := flag.String("sparkMaster", "local[*]", "Spark master")
	sparkSubmit := flag.String("sparkSubmit", "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/bin/spark-submit", "Location of the Spark submit script")
	targetJar := flag.String("targetJar", "/Users/acflorea/phd/columbugus/target/scala-2.10/columbugus-assembly-2.3.1.jar", "Location of the job jar")
	mainClass := flag.String("mainClass", "dr.acf.recc.ReccomenderBackbone", "The target class to execute")
	configFile := flag.String("configFile", "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/columbugus-conf/netbeans.conf", "Config file for the spark job")
	fsRoot := flag.String("fsRoot", "/Users/acflorea/phd/columbugus_data/netbeans_final_test", "Location of data files")

	flag.Parse()

	vargs := map[string]interface{}{}

	vargs["verbose"] = *verbose

	vargs["fileName"] = *fileNamePtr
	vargs["noOfExperiments"] = *noOfExperimentsPtr
	vargs["silent"] = *silentPtr
	vargs["maxAttempts"] = *maxAttemptsPtr
	vargs["fct"] = *fct
	vargs["alg"] = *alg
	vargs["script"] = *script
	vargs["workers"] = *workers
	vargs["targetstop"] = *targetstop

	// Spark specifics
	vargs["sparkMaster"] = *sparkMaster
	vargs["sparkSubmit"] = *sparkSubmit
	vargs["targetJar"] = *targetJar
	vargs["mainClass"] = *mainClass
	vargs["configFile"] = *configFile
	vargs["fsRoot"] = *fsRoot

	// Hyperopt specifics
	vargs["probs"] = *probs
	vargs["optSlicePercent"] = *optSlicePercent
	vargs["grievank"] = *grievank

	fmt.Println(Optimize(vargs))

}

func Optimize(vargs map[string]interface{}) map[string]interface{} {

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

	// 2
	// 91.72% due to main effect: X1
	// 2.44% due to main effect: X0
	// "2.44 91.72"

	// 3
	// 77.57% due to main effect: X2
	// 16.88% due to main effect: X1
	// 0.04% due to main effect: X0
	// "0.04 16.88 77.57"

	// 4
	// 62.96% due to main effect: X3
	// 23.14% due to main effect: X2
	// 2.31% due to main effect: X1
	// 0.05% due to main effect: X0
	// "0.05 2.31 23.14 62.96"

	// We target a stop after targetstop attempts
	var probabilityToChange = []float32{}
	probsStr := vargs["probs"].(string)
	for _, prob := range strings.Split(probsStr, " ") {
		if n, err := strconv.ParseFloat(prob, 32); err == nil {
			probabilityToChange = append(probabilityToChange, float32(n))
		}
	}

	optSlicePercent := vargs["optSlicePercent"].(float64)

	grievank := vargs["grievank"].(int)

	var restrictions []generators.GenerationStrategy

	for i := 0; i < grievank; i++ {
		restrictions = append(restrictions, generators.NewUniform("x"+strconv.Itoa(i+1), -600.0, 600.0))
	}

	// if this is true a single value changes for each step
	// otherwise the values are changing according to their probabilities
	var adjustSingleValue = false

	// if silent add a progress "bar"
	if silent {
		go func() {
			for {
				fmt.Print("...")
				time.Sleep(1 * time.Second)
			}
		}()
	}

	return core.Optimize(
		noOfExperiments,
		restrictions,
		probabilityToChange,
		adjustSingleValue,
		optSlicePercent,
		maxAttempts,
		targetstop,
		W,
		algorithm,
		targetFunction,
		silent,
		vargs)

}
