package functions

import "errors"

// A type alias for a function taking a variable number of parameters and returning a float
type NumericalFunction func(map[string]float64) (float64, error)

// Constant function
func F_constant(_ map[string]float64) (float64, error) {
	return 10, nil
}

// Identity function
func F_identity(x map[string]float64) (float64, error) {
	for _, value := range x {
		return value, nil
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// x^2 function
func F_x_square(x map[string]float64) (float64, error) {
	for _, value := range x {
		return value * value, nil
	}
	return 0.0, errors.New("Not a single parameter map.")
}

func Negate(f NumericalFunction) NumericalFunction {
	return func(x map[string]float64) (float64, error) {
		y, err := f(x)
		return -y, err
	}
}
