/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package healpix

/*****************************************************************************************************************/

import (
	"math"

	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/projection"
)

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

// ConvertEquatorialToCartesian converts equatorial coordinates (RA, Dec) to cartesian coordinates (x, y)
// using the HEALPix projection, see (https://healpix.sourceforge.io/) for further detail.
// The HEALPix projection is a hybrid projection that uses the interrupted Collignon projection for the
// polar regions and the Lambert-cylindrical closer to the equator.
func (h *HealPIX) ConvertEquatorialToCartesian(
	eq astrometry.ICRSEquatorialCoordinate,
) (x, y float64) {
	z := math.Sin(projection.Radians(eq.Dec))

	// Closer to the equator, we use the Lambert cylindrical projection:
	if math.Abs(z) <= h.PolarLatitudeBoundary {
		return projection.ConvertEquatorialToLambertCylindricalCartesian(eq, z)
	}

	// Closer to the polar regions, we use the interrupted Collignon projection:
	return projection.ConvertEquatorialToInterruptedCollignonCartesian(eq, z)
}

/*****************************************************************************************************************/
