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

type Scheme int

const (
	RING Scheme = iota
	NESTED
)

/*****************************************************************************************************************/

type HealPIX struct {
	NSide                 int
	Scheme                Scheme
	Longitude             float64
	Latitude              float64
	PolarLatitudeBoundary float64
}

/*****************************************************************************************************************/

// HEALPix, i.e., the "Hierarchical Equal Area isoLatitude Pixelization", is a versatile structure for the
// pixelization of coordinates on the sphere.
// @see https://healpix.jpl.nasa.gov/html/intro.htm
// @see https://healpix.sourceforge.io/pdf/intro.pdf
func NewHealPIX(sides int, scheme Scheme) *HealPIX {
	// Ensure the NSide is a power of 2 (2^k) and greater than 0:
	if sides < 1 {
		sides = 1
	} else {
		sides = 1 << uint(math.Round(math.Log2(float64(sides))))
	}

	return &HealPIX{
		NSide:                 sides,
		Scheme:                scheme,
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

// ConvertEquatorialToPixelIndex converts equatorial coordinates (RA, Dec) to a HEALPix pixel index
// (either RING or NESTED), based on the initial HealPIX configuration.
func (h *HealPIX) ConvertEquatorialToPixelIndex(eq astrometry.ICRSEquatorialCoordinate) int {
	// Convert to standard spherical angles for HEALPix, theta (co-latitude, [0, π]):
	theta := math.Pi/2.0 - projection.Radians(eq.Dec)

	// Clamp theta to [0, π]:
	if theta < 0 {
		theta = 0
	} else if theta > math.Pi {
		theta = math.Pi
	}

	// Convert to standard spherical angles for HEALPix, phi (longitude, [0, 2π)):
	phi := projection.Radians(eq.RA)

	// Normalize phi to [0, 2π):
	if phi < 0 {
		phi += 2.0 * math.Pi
	}

	// Branch to the specific indexing scheme (RING or NESTED):
	switch h.Scheme {
	case RING:
		return convertSphericalToRingIndex(h.NSide, theta, phi)
	case NESTED:
		return convertSphericalToNestedIndex(h.NSide, theta, phi)
	default:
		return convertSphericalToRingIndex(h.NSide, theta, phi)
	}
}

/*****************************************************************************************************************/

// convertSphericalToRingIndex converts spherical coordinates (theta, phi) to a HEALPix pixel index
// using the RING indexing scheme for any NSide >= 1.
func convertSphericalToRingIndex(nside int, theta, phi float64) int {
	return 0
}

/*****************************************************************************************************************/

// convertSphericalToNestedIndex converts spherical coordinates (theta, phi) to a HEALPix pixel index
// using the NESTED indexing scheme for any NSide >= 1.
func convertSphericalToNestedIndex(nside int, theta, phi float64) int {
	return 0
}

/*****************************************************************************************************************/
