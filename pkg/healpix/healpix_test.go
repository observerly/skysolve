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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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

	// Define expected pixel indices for RING and NESTED schemes based on your JSON data
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
