package functions

import (
	"os/exec"
	"fmt"
	"strconv"
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
	//Gamma, ok := vargs["gamma"].(float64)
	//if !ok {
	//	Gamma = "auto"
	//}
	Degree, ok := vargs["degree"].(int)
	if !ok {
		Degree = 3
	}
	Coef0, ok := vargs["coef0"].(float64)
	if !ok {
		Coef0 = 0.0
	}

	cmd := exec.Command("python", "/Users/aflorea/phd/optimus-prime/crossVal.py",
		fileName, "rbf", FloatToString(C), "auto", strconv.Itoa(Degree), FloatToString(Coef0))

	test, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(test))

	accuracy := -1.0

	return accuracy, nil
}
