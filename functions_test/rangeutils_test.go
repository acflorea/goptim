package functions

import (
	"testing"
	"github.com/acflorea/goptim/functions"
)

func TestTwoDPointVector(t *testing.T) {

	var twoDPointVector functions.TwoDPointVector

	point11 := functions.TwoDPoint{1.0, 1.0}
	point21 := functions.TwoDPoint{2.0, 1.0}
	point13 := functions.TwoDPoint{1.0, 3.0}

	twoDPointVector.Append(point11)


	twoDPointVector.Append(point21)
	twoDPointVector.Append(point13)

}
