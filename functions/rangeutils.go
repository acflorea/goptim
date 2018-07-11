package functions

type TwoDPoint struct {
	X float64
	Y float64
}

// A wrapper over a slice of points that also keeps the Sum and min/max values
type TwoDPointVector struct {
	// Ordered by x
	TwoDPoints []TwoDPoint
	Sum        TwoDPoint
	MinY       TwoDPoint
	MaxY       TwoDPoint
}

// Adds a point to the vector, recomputes Min, Max and Sum
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
		if v.MinY.Y <= point.Y {
			v.MinY = point
		}
		if v.MaxY.Y >= point.Y {
			v.MaxY = point
		}
	}

}

// Inserts the point in it's appropriate position
func insertionSort(twoDPoints []TwoDPoint) []TwoDPoint {

	n := len(twoDPoints)

	for i := 1; i < n; i++ {
		elem := twoDPoints[i]
		key := twoDPoints[i].X
		j := i - 1

		/* Move elements of arr[0..i-1], that are
		   greater than key, to one position ahead
		   of their current position */
		for ; j >= 0 && twoDPoints[j].X > key; {
			twoDPoints[j+1] = twoDPoints[j]
			j = j - 1
		}
		twoDPoints[j+1] = elem
	}

	return twoDPoints

}
