package functions

import (
	"github.com/acflorea/libsvm-go"
	"fmt"
)

// LIBSVM optimization through crossvalidation
func LIBSVM_optim(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {
	C, ok := p.Values["C"].(float64)
	if !ok {
		C = 0.01
	}
	Gamma, ok := p.Values["gamma"].(float64)
	if !ok {
		Gamma = 0.0
	}
	kernel, ok := p.Values["kernel"].(int)
	if !ok {
		kernel = libSvm.RBF
	}

	accuracy, _, _ := CrossV(kernel, C, Gamma, vargs)

	return accuracy, nil
}

func CrossV(kernel int, C, Gamma float64, vargs map[string]interface{}) (accuracy float64, all, TPs int) {

	fileName, found := vargs["fileName"].(string)
	if !found {
		panic("Missing input data! Please specify a fileName!")
	}

	quietMode := true

	param := libSvm.NewParameter() // Create a parameter object with default values
	param.QuietMode = quietMode

	// Parameters
	param.KernelType = kernel
	param.C = C
	param.Gamma = Gamma

	// Create a problem specification from the training data and parameter attributes
	problem, err := libSvm.NewProblem(fileName, param)

	_, acc, all, TPs := libSvm.CrossValidationWithAccuracies(problem, param, 10)

	accuracy = 0
	for i := 0; i < len(acc); i++ {
		if !quietMode {
			fmt.Println("Accuracy for fold ", i, " is ", acc[i])
		}
		accuracy += acc[i] / float64(len(acc))
	}

	if err != nil {
		panic(err)
	}

	return
}

func Train(C, Gamma float64, vargs map[string]interface{}) {

	fileName, found := vargs["fileName"].(string)
	if !found {
		panic("Missing input data! Please specify a fileName!")
	}
	modelName, found := vargs["modelName"].(string)
	if !found {
		panic("Missing model data! Please specify a modelName!")
	}

	param := libSvm.NewParameter() // Create a parameter object with default values
	param.KernelType = libSvm.RBF  // Use the polynomial kernel

	param.C = C
	param.Gamma = Gamma

	model := libSvm.NewModel(param) // Create a model object from the parameter attributes

	// Create a problem specification from the training data and parameter attributes
	problem, err := libSvm.NewProblem(fileName, param)

	if err != nil {
		panic(err)
	}

	model.Train(problem) // Train the model from the problem specification

	model.Dump(modelName)

}

func Test(vargs map[string]interface{}) {

	fileName, found := vargs["fileName"].(string)
	if !found {
		panic("Missing input data! Please specify a fileName!")
	}
	modelName, found := vargs["modelName"].(string)
	if !found {
		panic("Missing model data! Please specify a modelName!")
	}

	// Create a model object from the model file generated from training
	model := libSvm.NewModelFromFile(modelName)

	p := libSvm.Problem{}

	p.Read(fileName, libSvm.NewParameter())

	p.Begin()

	TP := 0
	for i := 0; i < p.ProblemSize(); i++ {
		y, x := p.GetLine()
		yp := model.Predict(x)
		//fmt.Println(y, yp)
		if y == yp {
			TP += 1.0
		}
		p.Next()
	}

	fmt.Println("Accuracy is: ", float64(TP)/float64(p.ProblemSize()))

	p.Done()

}
