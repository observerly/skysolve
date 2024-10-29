/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package utils

/*****************************************************************************************************************/

import "testing"

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
