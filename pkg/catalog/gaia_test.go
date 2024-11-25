/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

import (
	"math"
	"testing"

	"github.com/observerly/skysolve/pkg/astrometry"
)

/*****************************************************************************************************************/

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

/*****************************************************************************************************************/

func IsWithinICRSPolarRadius(ra, dec, r float64) bool {
	// Clamp cosine  to the valid range [-1, 1] to prevent NaN from math.Acos and
	// calculate the angular distance in radians:
	d := math.Acos(math.Max(-1.0, math.Min(1.0, math.Cos(radians(dec))*math.Cos(radians(ra)))))
	// Determine if the angular distance is within the radius R
	return d <= radians(r)
}

/*****************************************************************************************************************/

func TestGAIAQueryExecutedSuccessfully(t *testing.T) {
	var q = NewGAIAServiceClient()

	stars, err := q.PerformRadialSearch(astrometry.ICRSEquatorialCoordinate{
		RA:  0,
		Dec: 0,
	}, 2.5, 100, 10)

	if err != nil {
		t.Errorf("Failed to execute query: %v", err)
	}

	if len(stars) == 0 {
		t.Errorf("No stars returned")
	}

	for _, star := range stars {
		// Test that the star is within the search radius:
		if !IsWithinICRSPolarRadius(star.RA, star.Dec, 2.5) {
			t.Errorf("Star is not within the search radius")
		}
	}

	// The GAIA catalog is expected to return a maximum of 100 stars for this query:
	if len(stars) > 100 {
		t.Errorf("Too many stars returned")
	}
}

/*****************************************************************************************************************/
