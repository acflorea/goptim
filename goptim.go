package main

import (
	"github.com/acflorea/goptim/core"
	"github.com/acflorea/goptim/generators"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"flag"
	"strings"
	"strconv"
)

func main() {

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

	// Spark specifics
	sparkMaster := flag.String("sparkMaster", "local[*]", "Spark master")
	sparkSubmit := flag.String("sparkSubmit", "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/bin/spark-submit", "Location of the Spark submit script")
	targetJar := flag.String("targetJar", "/Users/acflorea/phd/columbugus/target/scala-2.10/columbugus-assembly-2.3.1.jar", "Location of the job jar")
	mainClass := flag.String("mainClass", "dr.acf.recc.ReccomenderBackbone", "The target class to execute")
	configFile := flag.String("configFile", "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/columbugus-conf/netbeans.conf", "Config file for the spark job")
	fsRoot := flag.String("fsRoot", "/Users/acflorea/phd/columbugus_data/netbeans_final_test", "Location of data files")

	flag.Parse()

	//functions.CrossV(1, 0.1)
	//functions.Train(1.0, 1.0/10.0)
	//functions.Test()

	vargs := map[string]interface{}{}
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

	Optimize(vargs)

}

func Optimize(vargs map[string]interface{}) {

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

	// We target a stop after targetstop attempts
	var probabilityToChange = []float32{}
	probsStr := vargs["probs"].(string)
	for _, prob := range strings.Split(probsStr, " ") {
		if n, err := strconv.ParseFloat(prob, 32); err == nil {
			probabilityToChange = append(probabilityToChange, float32(n))
		}
	}

	optSlicePercent := vargs["optSlicePercent"].(float64)

	//2^-3 to 2^10
	//restrictions := []generators.GenerationStrategy{
	//	generators.NewUniform("C", math.Pow(2, -2), math.Pow(2, 15)),
	//	generators.NewUniform("gamma", math.Pow(2, -15), math.Pow(2, 3)),
	//}

	////{"linear", "polynomial", "rbf", "sigmoid"}
	//restrictions := []generators.GenerationStrategy{
	//	generators.NewDiscrete("kernel", map[interface{}]float64{
	//		libSvm.LINEAR: 1.0, // 0
	//		libSvm.POLY:   1.0, // 1
	//		libSvm.RBF:    1.0, // 2
	//	}),
	//	generators.NewExponential("C", 10.0),
	//	generators.NewExponential("gamma", 10.0),
	//	generators.NewDiscrete("degree", map[interface{}]float64{
	//		2: 1.0,
	//		3: 1.0,
	//		4: 1.0,
	//		5: 1.0,
	//	}),
	//	generators.NewUniform("coef0", 0.0, 1.0),
	//	generators.NewUniform("categoryScalingFactor", 1.0, 100.0),
	//	generators.NewUniform("productScalingFactor", 1.0, 100.0),
	//}

	//onetoonehundred := map[interface{}]float64{}
	//for i := 1; i <= 1000; i++ {
	//	onetoonehundred[float64(i)] = 1.0
	//}
	//
	//restrictions := []generators.GenerationStrategy{
	//	generators.NewDiscrete("x", onetoonehundred),
	//}

	//{"linear", "polynomial", "rbf", "sigmoid"}
	restrictions := []generators.GenerationStrategy{
		//generators.NewUniform("x", -3000.0, 3000.0),
		//generators.NewUniform("y", -3000.0, 3000.0),
		//generators.NewUniform("z", -3000.0, 3000.0),
		generators.NewUniform("x", -3.0, 3.0),
		generators.NewUniform("y", -3.0, 3.0),
		generators.NewUniform("z", -3.0, 3.0),
	}

	//81.72% due to main effect: X0
	//8.23% due to main effect: X1
	//0.90% due to main effect: X2


	// heart
	//65.08% due to main effect: X0
	//9.80% due to main effect: X2
	//7.24% due to interaction: X1 x X0
	//7.07% due to interaction: X2 x X0
	//5.17% due to main effect: X1
	//3.50% due to interaction: X2 x X1
	//65.08 5.17 9.80

	//var probabilityToChange = []float32{81.72, 8.23, 0.9}
	//var probabilityToChange = []float32{10.01, 18.76, 60.22}
	// 3*x-2*y+z
	//var probabilityToChange = []float32{60.52, 26.97, 3.09}
	//var probabilityToChange = []float32{3, 2, 1}
	//var probabilityToChange = []float32{0.85, 0.1, 0.1}
	//var probabilityToChange = []float32{1, 1, 1}
	// if this is true a single value changes for each step
	// otherwise the values are changing according to their probabilities
	var adjustSingleValue = false

	core.Optimize(noOfExperiments, restrictions, probabilityToChange, adjustSingleValue, optSlicePercent, maxAttempts, targetstop, W, algorithm, targetFunction, silent, vargs)

}
