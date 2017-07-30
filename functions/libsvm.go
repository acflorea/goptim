package functions

import (
	"github.com/acflorea/libsvm-go"
	"fmt"
	"encoding/json"
)

// LIBSVM optimization through crossvalidation
func LIBSVM_optim(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {

	// Add Values to vargs
	for key, value := range p.Values {
		vargs[key] = value
	}

	accuracy, _, _ := CrossV(vargs)

	return accuracy, nil
}

func CrossV(vargs map[string]interface{}) (accuracy float64, all, TPs int) {

	quietMode := true
	fileName, ok := vargs["fileName"].(string)
	if !ok {
		panic("Missing input data! Please specify a fileName!")
	}

	param := libSvm.NewParameter() // Create a parameter object with default values
	param.QuietMode = quietMode

	kernel, ok := vargs["kernel"].(int)
	if ok {
		param.KernelType = kernel
	}
	C, ok := vargs["C"].(float64)
	if ok {
		param.C = C
	}
	Gamma, ok := vargs["gamma"].(float64)
	if ok {
		param.Gamma = Gamma
	}
	Degree, ok := vargs["degree"].(int)
	if ok {
		param.Degree = Degree
	}
	Coef0, ok := vargs["coef0"].(float64)
	if ok {
		param.Coef0 = Coef0
	}

	// Create a problem specification from the training data and parameter attributes
	problem, err := libSvm.NewProblem(fileName, param)

	_, acc, confusion := libSvm.CrossValidationWithAccuracies(problem, param, 10)

	accuracy = 0
	for i := 0; i < len(acc); i++ {
		if !quietMode {
			fmt.Println("Accuracy for fold ", i, " is ", acc[i])
		}
		accuracy += acc[i] / float64(len(acc))
	}
	//if !quietMode {
	fmt.Println("Accuracy is ", accuracy)
	jsonedCM, err := json.Marshal(confusion)
	fmt.Println("Confusion Matrix ", string(jsonedCM))
	//}

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
