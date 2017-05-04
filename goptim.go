package main

import (
	"fmt"
	"errors"
)

func main() {

	emptyMap := make(map[string]float64)
	constant(emptyMap)

	fmt.Println(Eval(constant, emptyMap))
	fmt.Println(Eval(identity, map[string]float64{
		"x": 1.2345,
	}))
	fmt.Println(Eval(x_square, map[string]float64{
		"x": 2.5,
	}))

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

// Evaluate the funtion f in the point p
func Eval(f func(map[string]float64) (float64, error), p map[string]float64) (float64, error) {
	return f(p)
}
