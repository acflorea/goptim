package main

import (
	"github.com/acflorea/goptim/generators"
	"math"
	"github.com/acflorea/goptim/functions"
	"fmt"
)

func main() {

	// Maximum number of attempts
	maxAttepts := 3000

	// The function we attempt to optimize
	targetFunction := functions.F_identity

	generator := generators.RandomUniformGenerator{
		DimensionsNo: 2,
		PointsNo:     maxAttepts,
		Restrictions: []generators.Range{
			{0, 120},
			{0, 120},
		},
	}

	// number of workers
	W := 100
	// channel used by workers to communicate their results
	messages := make(chan functions.Sample)

	for w := 0; w < W; w++ {
		go func(w int) {
			i, p, v := DMaximize(targetFunction, generator, maxAttepts/W)
			// fmt.Println("Worker ", w, " SparkIt MAX --> ", i, p, v)

			messages <- functions.Sample{i, p, v}
		}(w)
	}

	// Collect results
	results := make([]functions.Sample, W)
	totalTries := 0
	optim := -math.MaxFloat64
	for i := 0; i < W; i++ {
		results[i] = <-messages
		totalTries += results[i].Index
		if optim < results[i].Value {
			optim = results[i].Value
		}
	}

	fmt.Println(totalTries, optim)

}

// Attempts to dynamically minimize the function f
// k := n / (2 * math.E)
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
func DMinimize(f functions.NumericalFunction, generator generators.Generator, n int) (
	index int,
	p functions.MultidimensionalPoint,
	min float64) {

	k := int(float64(n) / (2 * math.E))
	return Minimize(f, generator, k, n)
}

// Attempts to minimize the function f
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
func Minimize(f functions.NumericalFunction, generator generators.Generator, k, n int) (
	index int,
	p functions.MultidimensionalPoint,
	min float64) {

	index = -1
	min = math.MaxFloat64

	for i := 0; i < n; i++ {
		rndPoint := generator.Next()
		f_rnd, _ := f(rndPoint)

		//fmt.Println(i, " :: ", rnd, " -> ", f_rnd)

		if (f_rnd < min) {
			index = i
			p = rndPoint
			min = f_rnd
			if i > k {
				break
			}
		}
	}

	return
}

// Dynamically Minimizes the negation of the target function
func DMaximize(f functions.NumericalFunction, generator generators.Generator, n int) (
	index int,
	p functions.MultidimensionalPoint,
	max float64) {

	index, p, max = DMinimize(functions.Negate(f), generator, n)
	return index, p, -max
}

// Minimizes the negation of the target function
func Maximize(f functions.NumericalFunction, generator generators.Generator, k, n int) (
	index int,
	p functions.MultidimensionalPoint,
	max float64) {

	index, p, max = Minimize(functions.Negate(f), generator, k, n)
	return index, p, -max
}
