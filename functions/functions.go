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
	Values []interface{}
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
			dimensionsLabels[i] = p.Labels[i] + "=" + p.Values[i].(string)
		}
		// Extra values
		for ; i < len(dimensionsLabels); i++ {
			dimensionsLabels[i] = "x" + strconv.Itoa(i) + "=" + p.Values[i].(string)
		}
		// Extra labels
		for ; i < len(p.Labels); i++ {
			dimensionsLabels = append(dimensionsLabels, p.Labels[i]+"=nil")
		}
	} else {
		for ; i < len(dimensionsLabels); i++ {
			dimensionsLabels[i] = "x" + strconv.Itoa(i) + "=" + p.Values[i].(string)
		}
	}

	desc = strings.Join(dimensionsLabels[:], ",")
	return
}

// A type alias for a function taking a variable number of parameters and returning a float
type NumericalFunction func(point MultidimensionalPoint, vargs map[string]interface{}) (float64, error)

// Constant function
func F_constant(_ MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {
	return 10, nil
}

// Identity function
func F_identity(x MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {
	for _, value := range x.Values {
		if v, ok := value.(float64); ok {
			return v, nil
		}
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// x^2 function
func F_x_square(x MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {
	for _, value := range x.Values {
		if v, ok := value.(float64); ok {
			return v * v, nil
		}
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// x^2*sin(x) function
func F_x_square_sin(x MultidimensionalPoint) (float64, error) {
	for _, value := range x.Values {
		if v, ok := value.(float64); ok {
			return v * v * math.Sin(v), nil
		}
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// sin(x) function
func F_sin(x MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {
	for _, value := range x.Values {
		if v, ok := value.(float64); ok {
			return math.Sin(v), nil
		}
	}
	return 0.0, errors.New("Not a single parameter map.")
}

// sin(sqrt(sq(x)+sq(y)))/sqrt(sq(x)+sq(y))
func F_sombrero(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {
	x, okx := p.Values[0].(float64)
	y, oky := p.Values[1].(float64)

	if okx && oky {
		w := math.Sqrt(x*x + y*y)

		if w != 0 {
			return math.Sin(w) / w, nil
		} else {
			return 0.0, errors.New("Not a single parameter map.")
		}
	} else {
		return 0.0, errors.New("Conversion failure")
	}
}

func Negate(f NumericalFunction) NumericalFunction {
	return func(x MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {
		y, err := f(x, vargs)
		return -y, err
	}
}
