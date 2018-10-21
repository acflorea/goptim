package functions

import (
	"testing"
	"bitbucket.org/acflorea/goptim/functions"
	"fmt"
)

func TestTwoDPointVectorAppend(t *testing.T) {

	var twoDPointVector functions.TwoDPointVector

	point1 := functions.TwoDPoint{1.0, 3.0}
	point2 := functions.TwoDPoint{-1.0, -1.0}
	point3 := functions.TwoDPoint{3.0, 1.0}
	point4 := functions.TwoDPoint{0.0, 0.0}

	twoDPointVector.Append(point1)

	if l := len(twoDPointVector.TwoDPoints); l != 1 {
		t.Error(fmt.Sprintf("Problem appending, mismatched points length, expected 1, found %d", l))
	}

	if twoDPointVector.TwoDPoints[0] != point1 {
		t.Error("Problem with points order")
	}

	twoDPointVector.Append(point2)

	if l := len(twoDPointVector.TwoDPoints); l != 2 {
		t.Error(fmt.Sprintf("Problem appending, mismatched points length, expected 2, found %d", l))
	}

	if twoDPointVector.TwoDPoints[0] != point2 || twoDPointVector.TwoDPoints[1] != point1 {
		t.Error("Problem with points order")
	}

	twoDPointVector.Append(point3)

	if l := len(twoDPointVector.TwoDPoints); l != 3 {
		t.Error(fmt.Sprintf("Problem appending, mismatched points length, expected 3, found %d", l))
	}

	if twoDPointVector.TwoDPoints[0] != point2 || twoDPointVector.TwoDPoints[1] != point1 ||
		twoDPointVector.TwoDPoints[2] != point3 {
		t.Error("Problem with points order")
	}

	twoDPointVector.Append(point4)

	if l := len(twoDPointVector.TwoDPoints); l != 4 {
		t.Error(fmt.Sprintf("Problem appending, mismatched points length, expected 4, found %d", l))
	}

	if twoDPointVector.TwoDPoints[0] != point2 || twoDPointVector.TwoDPoints[1] != point4 ||
		twoDPointVector.TwoDPoints[2] != point1 || twoDPointVector.TwoDPoints[3] != point3 {
		t.Error("Problem with points order")
	}
}
