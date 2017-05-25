package functions

import (
	"github.com/ewalker544/libsvm-go"
	"fmt"
)

func Train() {
	param := libSvm.NewParameter() // Create a parameter object with default values
	param.KernelType = libSvm.RBF  // Use the polynomial kernel

	param.C = 10
	param.Gamma = 1 / 13

	model := libSvm.NewModel(param) // Create a model object from the parameter attributes

	// Create a problem specification from the training data and parameter attributes
	problem, err := libSvm.NewProblem("/Users/acflorea/phd/libsvm-datasets/wine/wine.scale", param)

	if err != nil {
		panic(err)
	}

	model.Train(problem) // Train the model from the problem specification

	model.Dump("/Users/acflorea/phd/libsvm-datasets/wine/wine.model")

}

func Test() {

	// Create a model object from the model file generated from training
	model := libSvm.NewModelFromFile("/Users/acflorea/phd/libsvm-datasets/wine/wine.model")

	p := libSvm.Problem{}

	p.Read("/Users/acflorea/phd/libsvm-datasets/wine/wine.scale", libSvm.NewParameter())

	p.Begin()

	TP := 0
	for i := 0; i < p.ProblemSize(); i++ {
		y, x := p.GetLine()
		yp := model.Predict(x)
		fmt.Println(y, yp)
		if y == yp {
			TP += 1.0
		}
		p.Next()
	}

	fmt.Println("Accuracy is: ", float64(TP)/float64(p.ProblemSize()))

	p.Done()

}
