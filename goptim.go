package main

import (
	"fmt"
	"github.com/acflorea/goptim/rand"
	"math"
	"github.com/acflorea/goptim/functions"
)

func main() {

	emptyMap := make(map[string]float64)
	functions.F_constant(emptyMap)

	i, p, v := Minimize(functions.F_x_square, 25, 10000)
	fmt.Println("xSquare MIN --> ", i, p, v)

	i, p, v = Minimize(functions.F_identity, 25, 10000)
	fmt.Println("identity MIN --> ", i, p, v)

	i, p, v = Maximize(functions.F_x_square, 25, 10000)
	fmt.Println("xSquare MAX --> ", i, p, v)

	i, p, v = Maximize(functions.F_identity, 25, 10000)
	fmt.Println("identity MAX --> ", i, p, v)

	values := []float64{1.0, 2.0, 1.5}
	point := functions.MultidimensionalPoint{Values: values}
	fmt.Println(point.PrettyPrint())

}

// Evaluate the function f in the point p
func Eval(f func(map[string]float64) (float64, error), p map[string]float64) (float64, error) {
	return f(p)
}

// Attempts to minimize the function f
// 1st it evaluated the function in k random points and computes the minimum
// it then continues to evaluate the function (up to a total maximum of n attempts)
// The algorithm stops either if a value found at the second step is lower than the minimum
// of if n attempts have been made (in which case the 1st step minimum is reported)
func Minimize(f functions.NumericalFunction, k, n int) (index int, p, min float64) {

	left, right := -100.0, 100.0
	index = -1
	min = math.MaxFloat64

	for i := 0; i < n; i++ {
		_, rnd := rand.Float64(left, right)
		f_rnd, _ := Eval(f, map[string]float64{
			"x": rnd,
		})

		//fmt.Println(i, " :: ", rnd, " -> ", f_rnd)

		if (f_rnd < min) {
			index = i
			p = rnd
			min = f_rnd
			if i > k {
				break
			}
		}
	}

	return
}

// Minimizes the negation of the target function
func Maximize(f functions.NumericalFunction, k, n int) (index int, p, max float64) {
	index, p, max = Minimize(functions.Negate(f), k, n)
	return index, p, -max
}
