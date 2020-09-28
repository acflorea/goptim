package functions

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// This is a wrapper over the K7M optimizer
func K7M(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {

	targetScript, ok := vargs["script"].(string)
	if !ok {
		panic("Missing script information! Please specify a valid location!")
	}

	command, ok := vargs["command"].(string)
	if !ok {
		panic("Missing input data! Please specify a command to execute!")
	}

	targetFolder, ok := vargs["targetFolder"].(string)
	if !ok {
		panic("Missing input data! Please specify a target folder!")
	}

	// Add Values to vargs
	for key, value := range p.Values {
		vargs[key] = value
	}

	// unixOptions = "w:d:b:c:e:m:"
	// gnuOptions = ["max_breadth=", "max_depth=", "attr_b=", "attr_c=", "edge_cost=", "movement_factor="]

	max_breadth := strconv.Itoa(vargs["max_breadth"].(int))
	max_depth := strconv.Itoa(vargs["max_depth"].(int))
	attr_b := fmt.Sprintf("%f", vargs["attr_b"].(float64))
	attr_c := fmt.Sprintf("%f", vargs["attr_c"].(float64))
	edge_cost := fmt.Sprintf("%f", vargs["edge_cost"].(float64))
	movement_factor := strconv.Itoa(vargs["movement_factor"].(int))

	params := []string{targetFolder, targetScript, "-w" + max_breadth, "-d" + max_depth, "-b" + attr_b, "-c" + attr_c, "-e", edge_cost, "-m", movement_factor}

	cmd := exec.Command(command, params...)

	result, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}

	// target
	rs := strings.Split(string(result), "\n")
	target, _ := strconv.ParseFloat(string(rs[len(rs)-1]), 64)

	fmt.Println(fmt.Sprintf("%f for %v", target, params))

	return -target, nil
}
