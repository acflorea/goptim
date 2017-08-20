package functions

import (
	"fmt"
	"log"
	"github.com/danieldk/gosvm"
)

func main() {
	problem := gosvm.NewProblem()
	problem.Add(gosvm.TrainingInstance{0, gosvm.FromDenseVector([]float64{1, 1, 1, 0, 0})})
	problem.Add(gosvm.TrainingInstance{1, gosvm.FromDenseVector([]float64{1, 0, 1, 1, 1})})

	param := gosvm.DefaultParameters()
	model, err := gosvm.TrainModel(param, problem)
	if err != nil {
		log.Fatal(err)
	}

	label := model.Predict(gosvm.FromDenseVector([]float64{1, 1, 0, 0, 0}))

	fmt.Printf("Predicted label: %f\n", label)
}
