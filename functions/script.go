package functions

import (
	"os/exec"
	"strconv"
	"strings"
	"fmt"
)

// Calls an external script and collects the results
func Script(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {

	fileName, ok := vargs["fileName"].(string)
	if !ok {
		panic("Missing input data! Please specify a fileName!")
	}

	// Add Values to vargs
	for key, value := range p.Values {
		vargs[key] = value
	}

	//kernel, ok := vargs["kernel"].(int)
	//if !ok {
	//	kernel = "linear"
	//}
	C, ok := vargs["C"].(float64)
	if !ok {
		C = 1
	}
	_Gamma, ok := vargs["gamma"].(float64)
	Gamma := "auto"
	if ok {
		Gamma = FloatToString(_Gamma)
	}
	Degree, ok := vargs["degree"].(int)
	if !ok {
		Degree = 3
	}
	Coef0, ok := vargs["coef0"].(float64)
	if !ok {
		Coef0 = 0.0
	}

	fmt.Println("python", "/Users/aflorea/phd/optimus-prime/crossVal.py",
		fileName, "rbf", FloatToString(C), Gamma, strconv.Itoa(Degree), FloatToString(Coef0))

	cmd := exec.Command("python", "/Users/aflorea/phd/optimus-prime/crossVal.py",
		fileName, "rbf", FloatToString(C), Gamma, strconv.Itoa(Degree), FloatToString(Coef0))

	results, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	averages := strings.Split(string(results), ",")

	accuracy := 0.0
	for _, value := range averages {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}

		accuracy = accuracy + parsed
	}
	accuracy = accuracy / 10.0

	return accuracy, nil
}
