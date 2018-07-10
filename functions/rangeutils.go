package functions

type TwoDPoint struct {
	X float64
	Y float64
}

// A wrapper over a slice of points that also keeps the Sum and min/max values
type TwoDPointVector struct {
	TwoDPoints []TwoDPoint
	Sum        TwoDPoint
	MinY       TwoDPoint
	MaxY       TwoDPoint
}

func (v *TwoDPointVector) Append(point TwoDPoint) {

	// Keep the points sorted by Y
	v.TwoDPoints = insertionSort(append(v.TwoDPoints, point))
	if len(v.TwoDPoints) == 1 {
		v.Sum = point
		v.MinY = point
		v.MaxY = point
	} else {
		v.Sum.X += point.X
		v.Sum.Y += point.Y
	}

}

//func longestOKsequence(twoDPoints []TwoDPoint) []TwoDPoint {
//
//	maxLength := 0
//	maxStart := 0
//	start := 0
//	length := 0
//
//	avg := 0
//	for i, point := range twoDPoints {
//
//	}
//
//	for i, point := range twoDPoints {
//		if testFunction(point) {
//
//		} else {
//
//		}
//	}
//
//	return twoDPoints
//}

func insertionSort(twoDPoints []TwoDPoint) []TwoDPoint {

	n := len(twoDPoints)

	for i := 2; i < n; i++ {
		key := twoDPoints[i].X
		j := i - 1

		/* Move elements of arr[0..i-1], that are
		   greater than key, to one position ahead
		   of their current position */
		for ; j >= 0 && twoDPoints[j].X > key; {
			twoDPoints[j+1] = twoDPoints[j]
			j = j - 1
		}
		twoDPoints[j+1] = twoDPoints[i]
	}

	return twoDPoints

}
