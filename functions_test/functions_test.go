package functions_test

import (
	"testing"
	"github.com/acflorea/goptim/functions"
	"fmt"
)

func Test_pretyPrintPointNoLabels(t *testing.T) {
	values := []float64{1.0, 2.0, 1.5}

	point := functions.MultidimensionalPoint{Values: values}

	pretyPrint := point.PrettyPrint()
	expected := "x0=1.000000,x1=2.000000,x2=1.500000"
	if pretyPrint != expected {
		msg := fmt.Sprintf("Error describing point. Expected (%s) but got (%s).", expected, pretyPrint)
		t.Error(msg)
	}
}

func Test_pretyPrintPointPartialLabels(t *testing.T) {
	values := []float64{1.0, 2.0, 1.5}
	labels := []string{"XX0", "XX1"}

	point := functions.MultidimensionalPoint{Values: values, Labels: labels}

	pretyPrint := point.PrettyPrint()
	expected := "XX0=1.000000,XX1=2.000000,x2=1.500000"
	if pretyPrint != expected {
		msg := fmt.Sprintf("Error describing point. Expected (%s) but got (%s).", expected, pretyPrint)
		t.Error(msg)
	}
}

func Test_pretyPrintPointFullLabels(t *testing.T) {
	values := []float64{1.0, 2.0, 1.5}
	labels := []string{"XX0", "XX1", "XX2"}

	point := functions.MultidimensionalPoint{Values: values, Labels: labels}

	pretyPrint := point.PrettyPrint()
	expected := "XX0=1.000000,XX1=2.000000,XX2=1.500000"
	if pretyPrint != expected {
		msg := fmt.Sprintf("Error describing point. Expected (%s) but got (%s).", expected, pretyPrint)
		t.Error(msg)
	}
}

func Test_pretyPrintPointExtraLabels(t *testing.T) {
	values := []float64{1.0, 2.0, 1.5}
	labels := []string{"XX0", "XX1", "XX2", "XX3"}

	point := functions.MultidimensionalPoint{Values: values, Labels: labels}

	pretyPrint := point.PrettyPrint()
	expected := "XX0=1.000000,XX1=2.000000,XX2=1.500000,XX3=nil"
	if pretyPrint != expected {
		msg := fmt.Sprintf("Error describing point. Expected (%s) but got (%s).", expected, pretyPrint)
		t.Error(msg)
	}
}

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
