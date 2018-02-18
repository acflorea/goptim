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

	conv_layers := strconv.Itoa(vargs["conv_layers"].(int))
	full_layers := strconv.Itoa(vargs["full_layers"].(int))

	var maps [48]int
	for i := 3; i <= 50; i++ {
		maps[i-3] = vargs["maps_"+strconv.Itoa(i)].(int)
	}
	maps_str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(maps)), ","), "[]")

	var neurons [4]int
	for i := 1; i <= 4; i++ {
		neurons[i-1] = vargs["neurons_"+strconv.Itoa(i)].(int)
	}
	neurons_str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(neurons)), ","), "[]")

	params := []string{targetScript, "-c" + conv_layers, "-f" + full_layers, "-t" + test, "-n" + maps_str + "&" + neurons_str}

	cmd := exec.Command(command, params...)

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
