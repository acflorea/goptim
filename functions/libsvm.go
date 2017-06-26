package functions

import (
	"github.com/acflorea/libsvm-go"
	"fmt"
)

// LIBSVM optimization through crossvalidation
func LIBSVM_optim(p MultidimensionalPoint, vargs map[string]string) (float64, error) {
	C := p.Values[0]
	Gamma := p.Values[1]

	return CrossV(C, Gamma, vargs), nil
}

func CrossV(C, Gamma float64, vargs map[string]string) (accuracy float64) {

	fileName, found := vargs["fileName"]
	if !found {
		fileName = "/Users/acflorea/phd/libsvm-datasets/wine/wine.scale"
	}

	quietMode := true

	param := libSvm.NewParameter() // Create a parameter object with default values
	param.KernelType = libSvm.RBF  // Use the polynomial kernel
	param.QuietMode = quietMode

	param.C = C
	param.Gamma = Gamma

	// Create a problem specification from the training data and parameter attributes
	problem, err := libSvm.NewProblem(fileName, param)

	_, acc := libSvm.CrossValidationWithAccuracies(problem, param, 10)

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

func Train(C, Gamma float64, vargs map[string]string) {

	fileName, found := vargs["fileName"]
	if !found {
		fileName = "/Users/acflorea/phd/libsvm-datasets/wine/wine.scale"
	}
	modelName, found := vargs["modelName"]
	if !found {
		modelName = "/Users/acflorea/phd/libsvm-datasets/wine/wine.model"
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

func Test(vargs map[string]string) {

	fileName, found := vargs["fileName"]
	if !found {
		fileName = "/Users/acflorea/phd/libsvm-datasets/wine/wine.scale"
	}
	modelName, found := vargs["modelName"]
	if !found {
		modelName = "/Users/acflorea/phd/libsvm-datasets/wine/wine.model"
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
