/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package projection

/*****************************************************************************************************************/

import (
	"math"
	"testing"
)

/*****************************************************************************************************************/

// Helper function to compare two float64 numbers within a tolerance
func floatEquals(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

/*****************************************************************************************************************/

// TestConvertEquatorialToGnomicStandardCase tests a standard projection scenario
func TestConvertEquatorialToGnomicStandardCase(t *testing.T) {
	ra := 10.0   // degrees
	dec := 20.0  // degrees
	ra0 := 10.0  // degrees
	dec0 := 20.0 // degrees

	expectedX := 0.0
	expectedY := 0.0

	x, y := ConvertEquatorialToGnomic(ra, dec, ra0, dec0)

	if !floatEquals(x, expectedX, 1e-6) || !floatEquals(y, expectedY, 1e-6) {
		t.Errorf("Standard Case Failed: Expected (%f, %f), Got (%f, %f)", expectedX, expectedY, x, y)
	}
}

// TestConvertEquatorialToGnomicZeroDivision tests the edge case where cosalt1 is effectively zero
func TestConvertEquatorialToGnomicZeroDivision(t *testing.T) {
	// Choose ra and dec such that cosalt1 ≈ 0
	// For example, when dec = 90 degrees (North Pole) and dec0 = 0 degrees
	ra := 0.0   // degrees
	dec := 90.0 // degrees (North Pole)
	ra0 := 0.0  // degrees
	dec0 := 0.0 // degrees

	expectedX := 0.0
	expectedY := 0.0

	x, y := ConvertEquatorialToGnomic(ra, dec, ra0, dec0)

	if !floatEquals(x, expectedX, 1e-6) || !floatEquals(y, expectedY, 1e-6) {
		t.Errorf("Zero Division Case Failed: Expected (%f, %f), Got (%f, %f)", expectedX, expectedY, x, y)
	}
}

// TestConvertEquatorialToGnomicSameCoordinates tests when input coordinates are the same as reference
func TestConvertEquatorialToGnomicSameCoordinates(t *testing.T) {
	ra := 150.0   // degrees
	dec := -30.0  // degrees
	ra0 := 150.0  // degrees
	dec0 := -30.0 // degrees

	expectedX := 0.0
	expectedY := 0.0

	x, y := ConvertEquatorialToGnomic(ra, dec, ra0, dec0)

	if !floatEquals(x, expectedX, 1e-6) || !floatEquals(y, expectedY, 1e-6) {
		t.Errorf("Same Coordinates Case Failed: Expected (%f, %f), Got (%f, %f)", expectedX, expectedY, x, y)
	}
}

// TestConvertEquatorialToGnomicNorthPole tests projection at the North Celestial Pole
func TestConvertEquatorialToGnomicNorthPole(t *testing.T) {
	ra := 0.0    // degrees
	dec := 90.0  // degrees (North Pole)
	ra0 := 180.0 // degrees
	dec0 := 0.0  // degrees

	expectedX := 0.0
	expectedY := 0.0 // Since cosalt1 ≈ 0, expect (0, 0)

	x, y := ConvertEquatorialToGnomic(ra, dec, ra0, dec0)

	if !floatEquals(x, expectedX, 1e-6) {
		t.Errorf("North Pole Projection X Failed: Expected %f, Got %f", expectedX, x)
	}
	if !floatEquals(y, expectedY, 1e-6) {
		t.Errorf("North Pole Projection Y Failed: Expected %f, Got %f", expectedY, y)
	}
}

// TestConvertEquatorialToGnomicSouthPole tests projection at the South Celestial Pole
func TestConvertEquatorialToGnomicSouthPole(t *testing.T) {
	ra := 0.0    // degrees
	dec := -90.0 // degrees (South Pole)
	ra0 := 180.0 // degrees
	dec0 := 0.0  // degrees

	expectedX := 0.0
	expectedY := 0.0 // Since cosalt1 ≈ 0, expect (0, 0)

	x, y := ConvertEquatorialToGnomic(ra, dec, ra0, dec0)

	if !floatEquals(x, expectedX, 1e-6) {
		t.Errorf("South Pole Projection X Failed: Expected %f, Got %f", expectedX, x)
	}
	if !floatEquals(y, expectedY, 1e-6) {
		t.Errorf("South Pole Projection Y Failed: Expected %f, Got %f", expectedY, y)
	}
}

// Additional Test: Projection with a point 45 degrees away from the reference point
func TestConvertEquatorialToGnomicFortyFiveDegreesOffset(t *testing.T) {
	ra := 10.0   // degrees
	dec := 20.0  // degrees
	ra0 := 15.0  // degrees
	dec0 := 25.0 // degrees

	// Manually calculate expected x and y using the projection formula
	raRad := ra * math.Pi / 180
	decRad := dec * math.Pi / 180
	ra0Rad := ra0 * math.Pi / 180
	dec0Rad := dec0 * math.Pi / 180

	cosalt1 := math.Sin(dec0Rad)*math.Sin(decRad) + math.Cos(dec0Rad)*math.Cos(decRad)*math.Cos(raRad-ra0Rad)
	expectedX := math.Cos(decRad) * math.Sin(raRad-ra0Rad) / cosalt1
	expectedY := (math.Cos(dec0Rad)*math.Sin(decRad) - math.Sin(dec0Rad)*math.Cos(decRad)*math.Cos(raRad-ra0Rad)) / cosalt1

	x, y := ConvertEquatorialToGnomic(ra, dec, ra0, dec0)

	if !floatEquals(x, expectedX, 1e-6) || !floatEquals(y, expectedY, 1e-6) {
		t.Errorf("Forty-Five Degrees Offset Case Failed: Expected (%f, %f), Got (%f, %f)", expectedX, expectedY, x, y)
	}
}

/*****************************************************************************************************************/
