package functions

import (
	"errors"
	"strconv"
	"strings"
	"math"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

// A type alias for a multidimensional point
// (basically a map which entries have the form dimensionName -> value)
// Eg: The origin in a 3d space {"x":0, "y":0, "z":0}
type MultidimensionalPoint struct {
	// Values on each dimension
	Values []float64
	// Optional dimension labels
	Labels []string
}

// A sample
// Contains the sample index, the point itself and the value corresponding to that point
type Sample struct {
	Index      int
	Point      MultidimensionalPoint
	Value      float64
	GValue     float64
	FullSearch bool
}

// Prints a point in a friendly way
// Eg: "x0=1.000000,x1=2.000000,x2=1.500000"
func (p *MultidimensionalPoint) PrettyPrint() (desc string) {

	var dimensionsLabels = make([]string, len(p.Values))
	if i := 0; p.Labels != nil {
		// if some dimension labels are specified
		min := len(p.Labels)
		if min > len(dimensionsLabels) {
			min = len(dimensionsLabels)
		}
		// In sync labels and values
		for ; i < min; i++ {
			dimensionsLabels[i] = p.Labels[i] + "=" + FloatToString(p.Values[i])
		}
		// Extra values
		for ; i < len(dimensionsLabels); i++ {
			dimensionsLabels[i] = "x" + strconv.Itoa(i) + "=" + FloatToString(p.Values[i])
		}
		// Extra labels
		for ; i < len(p.Labels); i++ {
			dimensionsLabels = append(dimensionsLabels, p.Labels[i]+"=nil")
		}
	} else {
		for ; i < len(dimensionsLabels); i++ {
			dimensionsLabels[i] = "x" + strconv.Itoa(i) + "=" + FloatToString(p.Values[i])
		}
	}

	desc = strings.Join(dimensionsLabels[:], ",")
	return
}

// A type alias for a function taking a variable number of parameters and returning a float
type NumericalFunction func(point MultidimensionalPoint) (float64, error)

// Constant function
func F_constant(_ MultidimensionalPoint) (float64, error) {
	return 10, nil
}

// Identity function
func F_identity(x MultidimensionalPoint) (float64, error) {
	for _, value := range x.Values {
		return value, nil
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// x^2 function
func F_x_square(x MultidimensionalPoint) (float64, error) {
	for _, value := range x.Values {
		return value * value, nil
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// sin(x) function
func F_sin(x MultidimensionalPoint) (float64, error) {
	for _, value := range x.Values {
		return math.Sin(value), nil
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// sin(sqrt(sq(x)+sq(y)))/sqrt(sq(x)+sq(y))
func F_sombrero(p MultidimensionalPoint) (float64, error) {
	x := p.Values[0]
	y := p.Values[1]

	w := math.Sqrt(x*x + y*y)

	if w != 0 {
		return math.Sin(w) / w, nil
	} else {
		return 0.0, errors.New("Not a single parameter map.")
	}
}

func Negate(f NumericalFunction) NumericalFunction {
	return func(x MultidimensionalPoint) (float64, error) {
		y, err := f(x)
		return -y, err
	}
}
