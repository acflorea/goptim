package functions_test

import (
	"testing"
	"github.com/acflorea/goptim/functions"
	"fmt"
)

func TestF_constant(t *testing.T) {

	x := 12.34
	expected_y := 10.0

	y, err := functions.F_constant(map[string]float64{
		"x": x,
	})
	if err != nil {
		msg := fmt.Sprintf("F_constant(%f) raised an error ", x)
		t.Error(msg, err)
	}
	if err != nil || y != expected_y {
		msg := fmt.Sprintf("F_constant(%f) returned (%f). Expected is (%f).", x, y, expected_y)
		t.Error(msg)
	}

}

func TestF_identity(t *testing.T) {

	x := 12.34
	expected_y := x

	y, err := functions.F_identity(map[string]float64{
		"x": x,
	})
	if err != nil {
		msg := fmt.Sprintf("F_identity(%f) raised an error ", x)
		t.Error(msg, err)
	}
	if err != nil || y != expected_y {
		msg := fmt.Sprintf("F_identity(%f) returned (%f). Expected is (%f).", x, y, expected_y)
		t.Error(msg)
	}

}

func TestF_x_square(t *testing.T) {

	x := 12.34
	expected_y := 152.2756

	y, err := functions.F_x_square(map[string]float64{
		"x": x,
	})
	if err != nil {
		msg := fmt.Sprintf("F_x_square(%f) raised an error ", x)
		t.Error(msg, err)
	}
	if err != nil || y != expected_y {
		msg := fmt.Sprintf("F_x_square(%f) returned (%f). Expected is (%f).", x, y, expected_y)
		t.Error(msg)
	}

}
