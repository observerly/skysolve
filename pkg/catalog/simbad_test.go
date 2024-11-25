/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

import (
	"fmt"
	"testing"

	"github.com/observerly/skysolve/pkg/astrometry"
)

/*****************************************************************************************************************/

func TestSIMBADQueryExecutedSuccessfully(t *testing.T) {
	var q = NewSIMBADServiceClient()

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

		fmt.Println(star.Designation)
	}

	// The SIMBAD catalog is expected to return a maximum of 100 stars for this query:
	if len(stars) > 100 {
		t.Errorf("Too many stars returned")
	}
}

/*****************************************************************************************************************/
