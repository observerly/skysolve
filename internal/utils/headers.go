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

func ResolveOrExtractDecFromHeaders(value float32, header fits.FITSHeader) (float32, error) {
	// First, pick a candidate Dec (v):
	v := value

	// If the candidate Dec (v) is NaN, try to get it from the header:
	if math.IsNaN(float64(v)) {
		dec, exists := header.Floats["DEC"]
		if !exists {
			return float32(math.NaN()), fmt.Errorf("dec header not found in the supplied FITS file")
		}
		v = dec.Value
	}

	// Validate the candidate Dec (v) is a valid float32:
	if math.IsNaN(float64(v)) {
		return float32(math.NaN()), fmt.Errorf("dec value needs to be a valid float32")
	}

	// Validate the candidate Dec (v) is within the range [-90, 90]:
	if v < -90 || v > 90 {
		return float32(math.NaN()), fmt.Errorf("dec value is out of range: %f", v)
	}

	// Return the candidate Dec (v):
	return v, nil
}

/*****************************************************************************************************************/

func ExtractImageWidthFromHeaders(header fits.FITSHeader) (int32, error) {
	// Attempt to get the width header from the FITS file:
	width, exists := header.Ints["NAXIS1"]
	if !exists {
		return -1, fmt.Errorf("width header not found in the supplied FITS file")
	}

	// Validate the width is within the range [0, ∞):
	if width.Value <= 0 || width.Value == math.MaxInt32 {
		return -1, fmt.Errorf("width value is out of range: %v", width.Value)
	}

	// Return the width:
	return width.Value, nil
}

/*****************************************************************************************************************/

func ExtractImageHeightFromHeaders(header fits.FITSHeader) (int32, error) {
	// Attempt to get the height header from the FITS file:
	height, exists := header.Ints["NAXIS2"]
	if !exists {
		return -1, fmt.Errorf("height header not found in the supplied FITS file")
	}

	// Validate the height is within the range [0, ∞):
	if height.Value <= 0 || height.Value == math.MaxInt32 {
		return -1, fmt.Errorf("height value is out of range: %v", height.Value)
	}

	// Return the height:
	return height.Value, nil
}

/*****************************************************************************************************************/
