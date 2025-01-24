/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package utils

/*****************************************************************************************************************/

import (
	"fmt"
	"math"

	"github.com/observerly/iris/pkg/fits"
)

/*****************************************************************************************************************/

func ResolveOrExtractRAFromHeaders(value float32, header fits.FITSHeader) (float32, error) {
	// First, pick a candidate RA (v):
	v := value

	// If the candidate RA (v) is NaN, try to get it from the header:
	if math.IsNaN(float64(v)) {
		ra, exists := header.Floats["RA"]
		if !exists {
			return float32(math.NaN()), fmt.Errorf("ra header not found in the supplied FITS file")
		}
		v = ra.Value
	}

	// Validate the candidate RA (v) is a valid float32:
	if math.IsNaN(float64(v)) {
		return float32(math.NaN()), fmt.Errorf("ra value needs to be a valid float32")
	}

	// Validate the candidate RA (v) is within the range [0, 360]:
	if v < 0 || v > 360 {
		return float32(math.NaN()), fmt.Errorf("ra value is out of range: %f", v)
	}

	// Return the candidate RA (v):
	return v, nil
}

/*****************************************************************************************************************/
