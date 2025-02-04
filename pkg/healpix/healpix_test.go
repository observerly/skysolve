/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package healpix

/*****************************************************************************************************************/

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/observerly/skysolve/pkg/astrometry"
)

/*****************************************************************************************************************/

func TestHealpixGetNSide(t *testing.T) {
	nside := 2

	healpix := NewHealPIX(nside, RING)

	if healpix.GetNSide() != nside {
		t.Errorf("Expected NSide=%d, Got NSide=%d", nside, healpix.GetNSide())
	}
}

/*****************************************************************************************************************/

func TestHealpixGetPixelArea(t *testing.T) {
	// Define a slice of NSide values to test
	nsides := []int{128, 256, 512, 1024}

	// Define expected pixel areas based on your JSON data
	expectedPixelAreas := map[int]float64{
		128:  0.209823,
		256:  0.052456,
		512:  0.013114,
		1024: 0.003278,
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelAreaRing := hpRing.GetPixelArea()
				expectedPixelArea := expectedPixelAreas[nside]

				if math.Abs(pixelAreaRing-expectedPixelArea) > 1e-6 {
					t.Errorf("RING Scheme: NSide=%d => Expected Pixel Area=%.6f, Got Pixel Area=%.6f",
						nside, expectedPixelArea, pixelAreaRing)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixGetPixelRadialExtent tests the GetPixelRadialExtent method for various NSide values.
// It verifies that the calculated radial extent matches the expected values within a small tolerance.
func TestHealpixGetPixelRadialExtent(t *testing.T) {
	// Define a slice of NSide values to test
	nsides := []int{128, 256, 512, 1024}

	// Define expected radial extents based on HEALPix pixel area calculations
	// These values are approximate and based on the formula:
	// r = arccos(1 - A / (2π)) where A = 4π / Npixels
	// Convert the result from radians to degrees
	expectedRadialExtents := map[int]float64{
		128:  0.2584, // degrees
		256:  0.1292, // degrees
		512:  0.0646, // degrees
		1024: 0.0323, // degrees
	}

	// Define a small tolerance for floating-point comparisons
	tolerance := 1e-4 // degrees

	for _, nside := range nsides {
		// Define the test case name
		testName := fmt.Sprintf("NSide=%d,Scheme=RING", nside)

		// Run subtests for each NSide
		t.Run(testName, func(t *testing.T) {
			// Initialize HealPIX with the current NSide and RING scheme
			hpRing := NewHealPIX(nside, RING)

			// Choose a representative pixel ID (e.g., 0)
			pixelID := 0

			// Calculate the radial extent using the method under test
			calculatedRadialExtent := hpRing.GetPixelRadialExtent(pixelID)

			// Retrieve the expected radial extent
			expectedRadialExtent, exists := expectedRadialExtents[nside]
			if !exists {
				t.Fatalf("No expected radial extent defined for NSide=%d", nside)
			}

			// Calculate the absolute difference
			diff := math.Abs(calculatedRadialExtent - expectedRadialExtent)

			// Check if the difference exceeds the tolerance
			if diff > tolerance {
				t.Errorf("RING Scheme: NSide=%d, PixelID=%d => Expected Radial Extent=%.6f°, Got=%.6f° (Difference=%.6f°)",
					nside, pixelID, expectedRadialExtent, calculatedRadialExtent, diff)
			}
		})
	}
}

/*****************************************************************************************************************/

// TestHealpixNorthPole tests the North Pole coordinates across multiple NSide values using both RING and NESTED schemes.
func TestHealpixNorthPole(t *testing.T) {
	ra := 0.0
	dec := 90.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 0,
		2: 0,
		4: 0,
		8: 0,
	}

	expectedPixelsNESTED := map[int]int{
		1: 0,
		2: 3,
		4: 15,
		8: 63,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  45.0,
			Dec: 41.81031,
		},
		2: {
			RA:  45.0,
			Dec: 66.44354,
		},
		4: {
			RA:  45.0,
			Dec: 78.28415,
		},
		8: {
			RA:  45.0,
			Dec: 84.14973,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA0Dec0 tests equatorial coordinates at RA=0°, Dec=0° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA0Dec0(t *testing.T) {
	ra := 0.0
	dec := 0.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 4,
		2: 12,
		4: 72,
		8: 336,
	}

	expectedPixelsNESTED := map[int]int{
		1: 4,
		2: 19,
		4: 76,
		8: 304,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  0.0,
			Dec: 0.0,
		},
		2: {
			RA:  0.0,
			Dec: 19.47122,
		},
		4: {
			RA:  0.0,
			Dec: 9.59407,
		},
		8: {
			RA:  0.0,
			Dec: 4.78019,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA90Dec0 tests equatorial coordinates at RA=90°, Dec=0° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA90Dec0(t *testing.T) {
	ra := 90.0
	dec := 0.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 5,
		2: 22,
		4: 92,
		8: 376,
	}

	expectedPixelsNESTED := map[int]int{
		1: 5,
		2: 21,
		4: 86,
		8: 346,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  90.0,
			Dec: 0.0,
		},
		2: {
			RA:  112.5,
			Dec: 0.0, // Replace with accurate value
		},
		4: {
			RA:  101.25,
			Dec: 0.0, // Replace with accurate value
		},
		8: {
			RA:  95.625,
			Dec: 0.0, // Replace with accurate value
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA180Dec0 tests equatorial coordinates at RA=180°, Dec=0° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA180Dec0(t *testing.T) {
	ra := 180.0
	dec := 0.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 6,
		2: 24,
		4: 96,
		8: 384,
	}

	expectedPixelsNESTED := map[int]int{
		1: 6,
		2: 25,
		4: 102,
		8: 410,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  180.0,
			Dec: 0.0,
		},
		2: {
			RA:  202.5,
			Dec: 0.0, // Replace with accurate value
		},
		4: {
			RA:  191.25,
			Dec: 0.0, // Replace with accurate value
		},
		8: {
			RA:  185.625,
			Dec: 0.0, // Replace with accurate value
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA270Dec0 tests equatorial coordinates at RA=270°, Dec=0° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA270Dec0(t *testing.T) {
	ra := 270.0
	dec := 0.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 7,
		2: 26,
		4: 100,
		8: 392,
	}

	expectedPixelsNESTED := map[int]int{
		1: 7,
		2: 29,
		4: 118,
		8: 474,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  270.0,
			Dec: 0.0,
		},
		2: {
			RA:  292.5,
			Dec: 0.0,
		},
		4: {
			RA:  281.25,
			Dec: 0.0,
		},
		8: {
			RA:  275.625,
			Dec: 0.0,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA45Dec45 tests equatorial coordinates at RA=45°, Dec=45° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA45Dec45(t *testing.T) {
	ra := 45.0
	dec := 45.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 0,
		2: 0,
		4: 13,
		8: 87,
	}

	expectedPixelsNESTED := map[int]int{
		1: 0,
		2: 3,
		4: 12,
		8: 48,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  45.0,
			Dec: 41.81031,
		},
		2: {
			RA:  45.0,
			Dec: 66.44354,
		},
		4: {
			RA:  45.0,
			Dec: 54.34091,
		},
		8: {
			RA:  45.0,
			Dec: 48.14121,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA135Dec45 tests equatorial coordinates at RA=135°, Dec=45° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA135Dec45(t *testing.T) {
	ra := 135.0
	dec := 45.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 1,
		2: 1,
		4: 16,
		8: 94,
	}

	expectedPixelsNESTED := map[int]int{
		1: 1,
		2: 7,
		4: 28,
		8: 112,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  135.0,
			Dec: 41.81031,
		},
		2: {
			RA:  135.0,
			Dec: 66.44354,
		},
		4: {
			RA:  135.0,
			Dec: 54.34091,
		},
		8: {
			RA:  135.0,
			Dec: 48.14121,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA225DecNeg45 tests equatorial coordinates at RA=225°, Dec=-45° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA225DecNeg45(t *testing.T) {
	ra := 225.0
	dec := -45.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 10,
		2: 46,
		4: 175,
		8: 673,
	}

	expectedPixelsNESTED := map[int]int{
		1: 10,
		2: 40,
		4: 163,
		8: 655,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  225.0,
			Dec: -41.81031,
		},
		2: {
			RA:  225.0,
			Dec: -66.44354,
		},
		4: {
			RA:  225.0,
			Dec: -54.34091,
		},
		8: {
			RA:  225.0,
			Dec: -48.14121,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixRA315DecNeg45 tests equatorial coordinates at RA=315°, Dec=-45° across multiple NSide values using both RING and NESTED schemes.
func TestHealpixRA315DecNeg45(t *testing.T) {
	ra := 315.0
	dec := -45.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 11,
		2: 47,
		4: 178,
		8: 680,
	}

	expectedPixelsNESTED := map[int]int{
		1: 11,
		2: 44,
		4: 179,
		8: 719,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  315.0,
			Dec: -41.81031,
		},
		2: {
			RA:  315.0,
			Dec: -66.44354,
		},
		4: {
			RA:  315.0,
			Dec: -54.34091,
		},
		8: {
			RA:  315.0,
			Dec: -48.14121,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

// TestHealpixSouthPole tests the South Pole coordinates across multiple NSide values using both RING and NESTED schemes.
func TestHealpixSouthPole(t *testing.T) {
	ra := 0.0
	dec := -90.0

	coord := astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}

	// Define a slice of NSide values to test
	nsides := []int{1, 2, 4, 8}

	expectedPixelsRING := map[int]int{
		1: 8,
		2: 44,
		4: 188,
		8: 764,
	}

	expectedPixelsNESTED := map[int]int{
		1: 8,
		2: 32,
		4: 128,
		8: 512,
	}

	expectedAngles := map[int]astrometry.ICRSEquatorialCoordinate{
		1: {
			RA:  45.0,
			Dec: -41.81031,
		},
		2: {
			RA:  45.0,
			Dec: -66.44354,
		},
		4: {
			RA:  45.0,
			Dec: -78.28415,
		},
		8: {
			RA:  45.0,
			Dec: -84.14973,
		},
	}

	for _, nside := range nsides {
		// Test RING Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=RING", nside),
			func(t *testing.T) {
				hpRing := NewHealPIX(nside, RING)
				pixelRing := hpRing.ConvertEquatorialToPixelIndex(coord)
				expectedPixelRing := expectedPixelsRING[nside]

				if pixelRing != expectedPixelRing {
					t.Errorf("RING Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelRing, pixelRing)
				}

				equatorialRing := hpRing.ConvertPixelIndexToEquatorial(expectedPixelRing)

				expectedAngleRing, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in RING scheme", nside)
				}

				if math.Abs(equatorialRing.RA-expectedAngleRing.RA) > 1e-6 || math.Abs(equatorialRing.Dec-expectedAngleRing.Dec) > 1e-5 {
					t.Errorf("RING Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelRing, expectedAngleRing.RA, expectedAngleRing.Dec, equatorialRing.RA, equatorialRing.Dec)
				}
			},
		)

		// Test NESTED Scheme
		t.Run(
			fmt.Sprintf("NSide=%d,Scheme=NESTED", nside),
			func(t *testing.T) {
				hpNested := NewHealPIX(nside, NESTED)
				pixelNested := hpNested.ConvertEquatorialToPixelIndex(coord)
				expectedPixelNested := expectedPixelsNESTED[nside]

				if pixelNested != expectedPixelNested {
					t.Errorf("NESTED Scheme: NSide=%d, RA=%.1f°, Dec=%.1f° => Expected Pixel=%d, Got Pixel=%d",
						nside, ra, dec, expectedPixelNested, pixelNested)
				}

				equatorialNested := hpNested.ConvertPixelIndexToEquatorial(expectedPixelNested)

				expectedAngleNested, exists := expectedAngles[nside]
				if !exists {
					t.Fatalf("Expected angles not defined for NSide=%d in NESTED scheme", nside)
				}

				if math.Abs(equatorialNested.RA-expectedAngleNested.RA) > 1e-6 || math.Abs(equatorialNested.Dec-expectedAngleNested.Dec) > 1e-5 {
					t.Errorf("NESTED Scheme: NSide=%d, Pixel=%d => Expected RA=%.5f°, Dec=%.5f°, Got RA=%.5f°, Dec=%.5f°",
						nside, expectedPixelNested, expectedAngleNested.RA, expectedAngleNested.Dec, equatorialNested.RA, equatorialNested.Dec)
				}
			},
		)
	}
}

/*****************************************************************************************************************/

func TestGetPixelIndicesFromEquatorialRadialRegion(t *testing.T) {
	nside := 2

	healpix := NewHealPIX(nside, RING)

	eq := astrometry.ICRSEquatorialCoordinate{
		RA:  0.0,
		Dec: 0.0,
	}

	radius := 1.2

	pixelIndices := healpix.GetPixelIndicesFromEquatorialRadialRegion(eq, radius)

	expectedPixelIndices := []int{12, 20, 28, 27}

	if !reflect.DeepEqual(pixelIndices, expectedPixelIndices) {
		t.Errorf("Expected Pixel Indices=%v, Got Pixel Indices=%v", expectedPixelIndices, pixelIndices)
	}
}

/*****************************************************************************************************************/

func TestGetFaceXY(t *testing.T) {
	nside := 2

	healpix := NewHealPIX(nside, RING)

	pixelIndex := 12

	face, x, y := healpix.GetFaceXY(pixelIndex)

	expectedFace := 4
	expectedX := 1
	expectedY := 1

	if face != expectedFace || x != expectedX || y != expectedY {
		t.Errorf("Expected Face=%d, X=%d, Y=%d, Got Face=%d, X=%d, Y=%d", expectedFace, expectedX, expectedY, face, x, y)
	}
}

/*****************************************************************************************************************/

func TestGetPixelIndexFromFaceXY(t *testing.T) {
	nside := 2

	healpix := NewHealPIX(nside, RING)

	face := 4
	x := 1
	y := 1

	pixelIndex := healpix.GetPixelIndexFromFaceXY(face, x, y)

	expectedPixelIndex := 12

	if pixelIndex != expectedPixelIndex {
		t.Errorf("Expected Pixel Index=%d, Got Pixel Index=%d", expectedPixelIndex, pixelIndex)
	}
}

/*****************************************************************************************************************/

// TestGetNeighbouringPixelsNested tests the neighbour lookup in the NESTED scheme.
func TestGetNeighbouringPixelsNested(t *testing.T) {
	// Create a HealPIX instance with NSide=8 and the NESTED scheme.
	h := NewHealPIX(8, NESTED)

	// Define test cases; the expected neighbour lists are taken from your healpy output.
	testCases := []struct {
		pixel    int
		expected []int
		name     string
	}{
		{
			pixel:    0,
			expected: []int{277, 279, 2, 3, 1, 363, 362, 575},
			name:     "Pixel0",
		},
		{
			pixel:    10,
			expected: []int{287, 309, 32, 33, 11, 9, 8, 285},
			name:     "Pixel10",
		},
		{
			pixel:    50,
			expected: []int{39, 45, 56, 57, 51, 49, 48, 37},
			name:     "Pixel50",
		},
		{
			pixel:    100,
			expected: []int{97, 99, 102, 103, 101, 79, 78, 75},
			name:     "Pixel100",
		},
		{
			pixel:    200,
			expected: []int{477, 479, 202, 203, 201, 195, 194, 471},
			name:     "Pixel200",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := h.GetNeighbouringPixels(tc.pixel)
			sort.Ints(got)
			sort.Ints(tc.expected)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("NESTED: For pixel %d, expected neighbours %v, but got %v", tc.pixel, tc.expected, got)
			}
		})
	}
}

/*****************************************************************************************************************/

// TestGetNeighbouringPixelsRing tests the neighbour lookup in the RING scheme.
func TestGetNeighbouringPixelsRing(t *testing.T) {
	// Create a HealPIX instance with NSide=8 and the RING scheme.
	h := NewHealPIX(8, RING)

	// Define test cases; expected values come from your healpy output.
	testCases := []struct {
		pixel    int
		expected []int
		name     string
	}{
		{
			pixel:    0,
			expected: []int{4, 11, 3, 1, 6, 5, 13},
			name:     "Pixel0",
		},
		{
			pixel:    10,
			expected: []int{21, 20, 9, 2, 3, 11, 22, 37},
			name:     "Pixel10"},
		{
			pixel:    50,
			expected: []int{72, 71, 49, 31, 32, 51, 73, 99},
			name:     "Pixel50"},
		{
			pixel:    100,
			expected: []int{130, 99, 73, 51, 74, 101, 131, 163},
			name:     "Pixel100"},
		{
			pixel:    200,
			expected: []int{232, 199, 168, 136, 169, 201, 233, 264},
			name:     "Pixel200",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := h.GetNeighbouringPixels(tc.pixel)
			sort.Ints(got)
			sort.Ints(tc.expected)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("RING: For pixel %d, expected neighbours %v, but got %v", tc.pixel, tc.expected, got)
			}
		})
	}
}

/*****************************************************************************************************************/
