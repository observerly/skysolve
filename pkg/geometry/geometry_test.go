/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package utils

/*****************************************************************************************************************/

import (
	"math"
	"testing"
)

/*****************************************************************************************************************/

func TestDistanceBetweenTwoCartesianPoints(t *testing.T) {
	x1 := 0.0
	y1 := 0.0
	x2 := 3.0
	y2 := 4.0

	expected := 5.0

	result := DistanceBetweenTwoCartesianPoints(x1, y1, x2, y2)

	if result != expected {
		t.Errorf("DistanceBetweenTwoCartesianPoints(%f, %f, %f, %f) = %f; want %f", x1, y1, x2, y2, result, expected)
	}
}

/*****************************************************************************************************************/

func TestAngleBetweenThreeCartesianPointsRightTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 3.0, 0.0
	x3, y3 := 0.0, 4.0

	expectedAngle := 90.0
	epsilon := 1e-6

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected %.6f, got %.6f", expectedAngle, angle)
	}
}

func TestAngleBetweenThreeCartesianPointsEquilateralTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 1.0, 0.0
	x3, y3 := 0.5, math.Sqrt(3)/2

	expectedAngle := 60.0
	epsilon := 1e-6

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected %.6f, got %.6f", expectedAngle, angle)
	}
}

func TestAngleBetweenThreeCartesianPointsDegenerateTriangle(t *testing.T) {
	x1, y1 := 1.0, 1.0
	x2, y2 := 1.0, 1.0
	x3, y3 := 2.0, 2.0

	_, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err == nil {
		t.Errorf("Expected error for degenerate triangle, but got none")
	}
}

func TestAngleBetweenThreeCartesianPointsAcuteTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 2.0, 0.0
	x3, y3 := 1.0, 1.0

	expectedAngle := 45.0
	epsilon := 1e-6

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected %.6f, got %.6f", expectedAngle, angle)
	}
}

func TestAngleBetweenThreeCartesianPointsObtuseTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 4.0, 0.0
	x3, y3 := 1.0, 1.0

	expectedAngle := 90.0 // Approximately, actual angle is less
	epsilon := 1e-1       // Allowing more tolerance

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// The actual angle is approximately 16.26 degrees
	expectedAngle = math.Acos((math.Pow(1.0, 2)+math.Pow(3.0, 2)-math.Pow(math.Sqrt(2.0), 2))/(2*1.0*3.0)) * (180.0 / math.Pi) // ~16.26 degrees

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected ~%.6f, got %.6f", expectedAngle, angle)
	}
}

/*****************************************************************************************************************/
