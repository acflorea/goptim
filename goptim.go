package main

import (
	"github.com/acflorea/goptim/generators"
	"math"
	"github.com/acflorea/goptim/functions"
	"fmt"
)

func main() {

	generator := generators.RandomUniformGenerator{
		DimensionsNo: 2,
		PointsNo:     300,
		Restrictions: []generators.Range{
			{0, 120},
			{0, 120},
		},
	}

	i, p, v := Maximize(functions.SparkIt, generator, 25, 300)
	fmt.Println("SparkIt MAX --> ", i, p, v)

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

// Minimizes the negation of the target function
func Maximize(f functions.NumericalFunction, generator generators.Generator, k, n int) (
	index int,
	p functions.MultidimensionalPoint,
	max float64) {

	index, p, max = Minimize(functions.Negate(f), generator, k, n)
	return index, p, -max
}
