package main

import (
	"fmt"
	"errors"
	"github.com/acflorea/goptim/rand"
	"math"
)

func main() {

	emptyMap := make(map[string]float64)
	constant(emptyMap)

	i, p, v := Minimize(x_square, 25, 1000)
	fmt.Println("xSquare MIN --> ", i, p, v)

	i, p, v = Minimize(identity, 25, 1000)
	fmt.Println("identity MIN --> ", i, p, v)

}

// Constant function
func constant(_ map[string]float64) (float64, error) {
	return 10, nil
}

// Identity function
func identity(x map[string]float64) (float64, error) {
	for _, value := range x {
		return value, nil
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// x^2 function
func x_square(x map[string]float64) (float64, error) {
	for _, value := range x {
		return value * value, nil
	}
	return 0.0, errors.New("Not a single parameter map.")
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
func Minimize(f func(map[string]float64) (float64, error), k, n int) (index int, p, min float64) {

	left, right := -10.0, 10.0
	index = -1
	min = math.MaxFloat64

	for i := 0; i < k; i++ {
		_, rnd := rand.Float64(left, right)
		f_rnd, _ := Eval(f, map[string]float64{
			"x": rnd,
		})

		//fmt.Println(i, " :: ", rnd, " -> ", f_rnd)

		if (f_rnd < min) {
			index = i
			p = rnd
			min = f_rnd
		}
	}

	for i := k; i < n; i++ {
		_, rnd := rand.Float64(left, right)
		f_rnd, _ := Eval(f, map[string]float64{
			"x": rnd,
		})

		//fmt.Println(i, " :: ", rnd, " -> ", f_rnd)

		if (f_rnd < min) {
			index = i
			p = rnd
			min = f_rnd
			break
		}
	}

	return

}
