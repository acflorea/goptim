package functions

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// This is a wrapper over the K7M optimizer
func HPC39(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {

	targetScript, ok := vargs["script"].(string)
	if !ok {
		panic("Missing script information! Please specify a valid location!")
	}

	command, ok := vargs["command"].(string)
	if !ok {
		panic("Missing input data! Please specify a command to execute!")
	}

	//targetFolder, ok := vargs["targetFolder"].(string)
	//if !ok {
	//	panic("Missing input data! Please specify a target folder!")
	//}

	// Add Values to vargs
	for key, value := range p.Values {
		vargs[key] = value
	}

	//learning_rate_exp, batch_size, x_layers, cnn_neurons, fc_neurons, dropout1, dropout2

	learning_rate_exp := vargs["learning_rate_exp"].(float64)
	learning_rate := fmt.Sprintf("%f", math.Pow(10, learning_rate_exp))
	batch_size := strconv.Itoa(vargs["batch_size"].(int))
	x_layers := strconv.Itoa(vargs["x_layers"].(int))
	cnn_neurons := strconv.Itoa(vargs["cnn_neurons"].(int))
	fc_neurons := strconv.Itoa(vargs["fc_neurons"].(int))
	dropout1 := fmt.Sprintf("%f", vargs["dropout1"].(float64))
	dropout2 := fmt.Sprintf("%f", vargs["dropout2"].(float64))

	params := []string{targetScript, learning_rate, batch_size, x_layers, cnn_neurons, fc_neurons, dropout1, dropout2}

	cmd := exec.Command(command, params...)

	result, err := cmd.Output()

	if err != nil {
		fmt.Print("ERR :: ")
		fmt.Println(err)
	}

	fmt.Println(string(result))

	// target
	rs := strings.Split(string(result), "\n")
	target, _ := strconv.ParseFloat(string(rs[len(rs)-2])[25:], 64)

	fmt.Println(fmt.Sprintf("%f for %v", target, params))

	return target, nil
}
