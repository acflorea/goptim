package functions

import (
	"testing"
	"github.com/acflorea/goptim/functions"
	"fmt"
)

func TestTwoDPointVectorAppend(t *testing.T) {

	var twoDPointVector functions.TwoDPointVector

	point11 := functions.TwoDPoint{1.0, 3.0}
	point21 := functions.TwoDPoint{2.0, -1.0}
	point13 := functions.TwoDPoint{1.0, 1.0}

	twoDPointVector.Append(point11)

	if l := len(twoDPointVector.TwoDPoints); l != 1 {
		t.Error(fmt.Sprintf("Problem appending, mismatched points length, expected 1, found %d", l))
	}

	if twoDPointVector.TwoDPoints[0] != point11 {
		t.Error("Problem with points order")
	}

	twoDPointVector.Append(point21)

	if l := len(twoDPointVector.TwoDPoints); l != 2 {
		t.Error(fmt.Sprintf("Problem appending, mismatched points length, expected 2, found %d", l))
	}

	if twoDPointVector.TwoDPoints[0] != point11 {
		t.Error("Problem with points order")
	}

	twoDPointVector.Append(point13)

	if l := len(twoDPointVector.TwoDPoints); l != 3 {
		t.Error(fmt.Sprintf("Problem appending, mismatched points length, expected 3, found %d", l))
	}

	if twoDPointVector.TwoDPoints[0] != point11 {
		t.Error("Problem with points order")
	}

}
