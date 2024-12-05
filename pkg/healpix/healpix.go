/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package healpix

/*****************************************************************************************************************/

type HealPIX struct {
	Longitude             float64
	Latitude              float64
	PolarLatitudeBoundary float64
}

/*****************************************************************************************************************/

// HEALPix, i.e., the "Hierarchical Equal Area isoLatitude Pixelization", is a versatile structure for the
// pixelization of coordinates on the sphere.
func NewHealPIX() *HealPIX {
	return &HealPIX{
		Longitude:             180.0,
		Latitude:              0.0,
		PolarLatitudeBoundary: 2.0 / 3.0, // in radians (approximately 38.1972 degrees)
	}
}

/*****************************************************************************************************************/
