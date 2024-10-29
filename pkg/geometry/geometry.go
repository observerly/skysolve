/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package geometry

/*****************************************************************************************************************/

import (
	"errors"
	"math"
)

/*****************************************************************************************************************/

func DistanceBetweenTwoCartesianPoints(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}

/*****************************************************************************************************************/

func AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3 float64) (float64, error) {
	a := DistanceBetweenTwoCartesianPoints(x2, y2, x3, y3) // Side opposite to point A (x1, y1)
	b := DistanceBetweenTwoCartesianPoints(x1, y1, x3, y3) // Side opposite to point B (x2, y2)
	c := DistanceBetweenTwoCartesianPoints(x1, y1, x2, y2) // Side opposite to point C (x3, y3)

	// Check for degenerate triangle (i.e. collinear points):
	if a == 0 || b == 0 || c == 0 {
		return 0, errors.New("degenerate triangle with zero-length sides")
	}

	// From the Law of Cosines, we can calculate the numerator of the arc-cosine:
	n := (math.Pow(b, 2) + math.Pow(c, 2) - math.Pow(a, 2))

	// From the Law of Cosines, we can calculate the denominator of the arc-cosine:
	d := 2 * b * c

	if d == 0 {
		return 0, errors.New("division by zero")
	}

	// Calculate the angle between the three points:
	return math.Acos(n/d) * 180 / math.Pi, nil
}

/*****************************************************************************************************************/

type InvariantFeatures struct {
	RatioAB float64
	RatioAC float64
	AngleA  float64
	AngleB  float64
}

/*****************************************************************************************************************/

func ComputeInvariantFeatures(x1, y1, x2, y2, x3, y3 float64) (InvariantFeatures, error) {
	// Compute side lengths of the triangle:
	a := DistanceBetweenTwoCartesianPoints(x2, y2, x3, y3) // BC
	b := DistanceBetweenTwoCartesianPoints(x1, y1, x3, y3) // AC
	c := DistanceBetweenTwoCartesianPoints(x1, y1, x2, y2) // AB

	// Check for degenerate triangle (i.e. collinear points):
	if a == 0 || b == 0 || c == 0 {
		return InvariantFeatures{}, errors.New("degenerate triangle with zero-length sides")
	}

	// Compute the angle A which is opposite to side a:
	angleA, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		return InvariantFeatures{}, err
	}

	// Compute the angle B which is opposite to side b:
	angleB, err := AngleBetweenThreeCartesianPoints(x2, y2, x1, y1, x3, y3)
	if err != nil {
		return InvariantFeatures{}, err
	}

	// Calculate ratios based on specific sides without normalization
	ratioAB := math.Min(c, a) / math.Max(c, a)
	ratioAC := math.Min(b, a) / math.Max(b, a)

	return InvariantFeatures{
		RatioAB: ratioAB,
		RatioAC: ratioAC,
		AngleA:  angleA,
		AngleB:  angleB,
	}, nil
}

/*****************************************************************************************************************/
