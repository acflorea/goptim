package functions

import (
	"os/exec"
	"fmt"
)

// Calls an external script and collects the results
func Script(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {

	fileName, ok := vargs["fileName"].(string)
	if !ok {
		panic("Missing input data! Please specify a fileName!")
	}

	cmd := exec.Command("python", "/Users/aflorea/phd/optimus-prime/random.py", fileName)

	test, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(test))

	accuracy := -1.0

	return accuracy, nil
}
