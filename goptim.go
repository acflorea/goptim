package main

import (
	"github.com/acflorea/goptim/generators"
	"math"
	"github.com/acflorea/goptim/functions"
	"fmt"
	"time"
	"math/rand"
)

func main() {
	//functions.Train()
	//functions.Test()

	Optimize()
}

func Optimize() {

	start := time.Now()

	// Maximum number of attempts
	maxAttempts := 1000000

	// The function we attempt to optimize
	targetFunction := functions.F_sombrero

	// Algorithm
	algorithm := generators.Leapfrog

	// number of workers
	W := 100

	restrictions := []generators.Range{
		{-1000, 1000},
		{-1000, 1000},
	}

	generator :=
		generators.NewRandomUniformGenerator(2, restrictions, maxAttempts, W, algorithm)

	// channel used by workers to communicate their results
	messages := make(chan functions.Sample, W)

	for w := 0; w < W; w++ {
		go func(w int) {
			i, p, v, o := DMaximize(targetFunction, generator, maxAttempts/W, w)
			fmt.Println("Worker ", w, " MAX --> ", i, p, v, o)

			messages <- functions.Sample{i, p, v}
		}(w)
	}

	// Collect results
	results := make([]functions.Sample, W)
	totalTries := 0
	optim := -math.MaxFloat64
	var point functions.MultidimensionalPoint
	for i := 0; i < W; i++ {
		results[i] = <-messages
		totalTries += results[i].Index
		if optim < results[i].Value {
			optim = results[i].Value
			point = results[i].Point
		}
	}

	fmt.Println(totalTries, point, optim)

	elapsed := time.Since(start)
	fmt.Println("Optimization took %s", elapsed)

}

// Attempts to dynamically minimize the function f
// k := n / (2 * math.E)
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
// w is thw worker index
func DMinimize(f functions.NumericalFunction, generator generators.Generator, n, w int) (
	index int,
	p functions.MultidimensionalPoint,
	min float64,
	optimNo int) {

	k := int(float64(n) / (2 * math.E))
	return Minimize(f, generator, k, n, w)
}

// Attempts to minimize the function f
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
// w is the worker index
func Minimize(f functions.NumericalFunction, generator generators.Generator, k, n, w int) (
	index int,
	p functions.MultidimensionalPoint,
	min float64,
	optimNo int) {

	index = -1
	min = math.MaxFloat64
	optimNo = 0

	for i := 0; i < n; i++ {
		rndPoint := generator.Next(w)
		f_rnd, _ := f(rndPoint)

		//fmt.Println(i, " :: ", rnd, " -> ", f_rnd)

		if f_rnd < min {
			index = i
			p = rndPoint
			min = f_rnd
			if i > k {
				// Increase the number of optimum points found
				optimNo += 1
				s := rand.NewSource(time.Now().UnixNano())
				tmpr := rand.New(s)
				threshold := tmpr.Float64()
				if threshold < 0.5+0.1*float64(optimNo) {
					break
				}
			}
		}
	}

	return
}

// Dynamically Minimizes the negation of the target function
func DMaximize(f functions.NumericalFunction, generator generators.Generator, n, w int) (
	index int,
	p functions.MultidimensionalPoint,
	max float64,
	optimNo int) {

	index, p, max, optimNo = DMinimize(functions.Negate(f), generator, n, w)
	return index, p, -max, optimNo
}

// Minimizes the negation of the target function
func Maximize(f functions.NumericalFunction, generator generators.Generator, k, n, w int) (
	index int,
	p functions.MultidimensionalPoint,
	max float64,
	optimNo int) {

	index, p, max, optimNo = Minimize(functions.Negate(f), generator, k, n, w)
	return index, p, -max, optimNo
}
