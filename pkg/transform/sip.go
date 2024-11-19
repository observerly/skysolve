/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package transform

/*****************************************************************************************************************/

// SIP (Simple Imaging Polynomial) is a convention used in FITS (Flexible Image Transport System)
// headers to describe complex distortions in astronomical images. It extends the standard World
// Coordinate System (WCS) by introducing higher-order polynomial terms that account for non-linear
// optical distortions, such as those introduced by telescope optics or atmospheric effects.
// @see https://fits.gsfc.nasa.gov/registry/sip/SIP_distortion_v1_0.pdf

/*****************************************************************************************************************/

// The forward parameters are polynomial coefficients used to map from pixel coordinates to world coordinates.
type SIP2DForwardParameters struct {
	AOrder int
	APower map[string]float64
	BOrder int
	BPower map[string]float64
}

/*****************************************************************************************************************/

// The inverse paramaters are polynomial coefficients used to map from world coordinates to pixel coordinates.
type SIP2DInverseParameters struct {
	APOrder int
	APPower map[string]float64
	BPOrder int
	BPPower map[string]float64
}

/*****************************************************************************************************************/
