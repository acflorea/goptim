package functions

import (
	"os/exec"
	"strings"
	"fmt"
	"strconv"
)

// Calls an external script and collects the results
func Keras(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {

	targetScript, ok := vargs["script"].(string)
	if !ok {
		panic("Missing script information! Please specify a valid location!")
	}

	command, ok := vargs["command"].(string)
	if !ok {
		panic("Missing input data! Please specify a command to execute!")
	}

	test, ok := vargs["test"].(string)
	if !ok {
		test = "False"
	}

	// Add Values to vargs
	for key, value := range p.Values {
		vargs[key] = value
	}

	cmd := exec.Command(command, targetScript, "-t"+test)

	results, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}

	averages := strings.Split(string(results), "\n")

	// accuracy is last
	accuracy, _ := strconv.ParseFloat(averages[len(averages)-1], 64)

	fmt.Print(", ", string(results))
	fmt.Println(", ", accuracy)

	return accuracy, nil
}
