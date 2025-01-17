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

	// Clamp theta to [0, π] (co-latitude):
	if theta < 0 {
		theta = 0
	} else if theta > math.Pi {
		theta = math.Pi
	}

	// Convert to standard spherical angles for HEALPix, phi (longitude, [0, 2π)):
	phi := projection.Radians(eq.RA)

	// Normalize phi to [0, 2π) (longitude):
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

// GetPixelIndicesFromEquatorialRadialRegion returns a list of HEALPix pixel indices for a given equatorial
// coordinate and radius.
func (h *HealPIX) GetPixelIndicesFromEquatorialRadialRegion(
	eq astrometry.ICRSEquatorialCoordinate,
	radius float64, // in degrees
) []int {
	// Use a map to collect unique pixel indices for the radial region:
	pixelIndices := make(map[int]bool)

	// This slice will store the pixel indices for the radial region:
	pixels := make([]int, 0)

	// Number of steps to take in the radial region should be proportional to the radius to ensure we always
	// sample the radial region consistently:
	steps := int(math.Ceil(radius * 10))

	// Our initial equatorial coordinate is the center of the radial region:
	center := eq

	// We aim to take a radial sub-sample of equatorial coordinates within the given radius,
	// and convert them to HEALPix pixel indices. We perform an polar radial sample of points inside
	// the given radius, and convert them to HEALPix pixel indices.
	for i := 0; i <= steps; i++ {
		r := (float64(i) / float64(steps)) * radius

		// For the central point, we simple add the central equatorial coordinate to the map:
		if r == 0 {
			i := h.ConvertEquatorialToPixelIndex(center)

			if _, exists := pixelIndices[i]; !exists {
				pixelIndices[i] = true
				pixels = append(pixels, i)
			}

			continue
		}

		// Cycle over the azimuthal angle range to sample equatorial coordinates, in 15 degree increments:
		for az := 0.0; az <= 360.0; az += 15.0 {
			ra, dec := projection.GetEquatorialCoordinateFromPolarOffset(
				eq.RA, eq.Dec, r, az,
			)

			eq := astrometry.ICRSEquatorialCoordinate{
				RA:  ra,
				Dec: dec,
			}

			i := h.ConvertEquatorialToPixelIndex(eq)

			if _, exists := pixelIndices[i]; !exists {
				pixelIndices[i] = true
				pixels = append(pixels, i)
			}
		}
	}

	return pixels
}

/*****************************************************************************************************************/

// convertSphericalToRingIndex converts spherical coordinates (theta, phi) to a HEALPix pixel index
// using the RING indexing scheme for any NSide >= 1.
func convertSphericalToRingIndex(nside int, theta, phi float64) int {
	z := math.Cos(theta)

	za := math.Abs(z)

	// Scale φ by inverse of π/2 to get [0,4):
	φ := phi * 1.0 / (0.5 * math.Pi)

	nSideFaces := float64(nside)

	// |z| <= 2/3 cooresponds to the equatorial region in the HEALPix projection:
	if za <= 2.0/3.0 {
		// Calculate j+ and j-:
		jp := int(nSideFaces*(0.5+φ) - nSideFaces*(z*0.75))
		jm := int(nSideFaces*(0.5+φ) + nSideFaces*(z*0.75))

		// Determine the ring index for the equatorial region:
		ir := (nside + 1) + jp - jm
		kshift := 1 - (ir & 1) // 1 if ir even, else 0

		ip := (jp + jm - nside + kshift + 1) / 2
		fourN := 4 * nside

		// Determine the ring index for the equatorial region:
		ip = int(math.Mod(float64(ip), float64(fourN)))
		return 2*nside*(nside-1) + (ir-1)*fourN + ip
	}

	// Otherwise, we are in the polar region of the HEALPix projection, |z| > 2/3:
	// Calculate j+ and j-:
	jp := int((φ - math.Floor(φ)) * (nSideFaces * math.Sqrt(3.0*(1.0-za))))
	jm := int((1.0 - (φ - math.Floor(φ))) * (nSideFaces * math.Sqrt(3.0*(1.0-za))))

	// Determine the ring index for the polar region:
	ir := jp + jm + 1

	// Determine the pixel index for the polar region:
	ip := ((int(φ*float64(ir)) % (4 * ir)) + 4*ir) % (4 * ir)

	// North polar cap (z > 0) of the sphere:
	if z > 0 {
		return 2*ir*(ir-1) + ip
	}

	// South polar cap (z < 0) of the sphere:
	return 12*nside*nside - 2*ir*(ir+1) + ip
}

/*****************************************************************************************************************/

// convertSphericalToNestedIndex converts spherical coordinates (theta, phi) to a HEALPix pixel index
// using the NESTED indexing scheme for any NSide >= 1.
func convertSphericalToNestedIndex(nside int, theta, phi float64) int {
	z := math.Cos(theta)

	za := math.Abs(z)

	// Scale φ by inverse of (π/2) → range [0,4):
	φ := phi * (1.0 / (0.5 * math.Pi))

	// Convert nside to float for arithmetic:
	nSideFaces := float64(nside)

	// We'll determine which face (faceIndex) we're on, and (ix, iy) within that face:
	var faceIndex, ix, iy int

	// If |z| ≤ 2/3, we're in the equatorial region:
	if za <= 2.0/3.0 {
		// j+ and j-:
		jp := int(nSideFaces*(0.5+φ) - nSideFaces*(0.75*z))
		jm := int(nSideFaces*(0.5+φ) + nSideFaces*(0.75*z))

		// Determine faceIndex by comparing j+ and j-:
		faceP := jp / nside
		faceM := jm / nside

		switch {
		case faceP == faceM:
			// Bitwise OR with 4 for this equatorial sub-face:
			faceIndex = faceP | 4
		case faceP < faceM:
			faceIndex = faceP
		default:
			faceIndex = faceM + 8
		}

		// Local x,y coordinates on that face:
		ix = jm & (nside - 1)
		iy = (nside - 1) - (jp & (nside - 1))
	} else {
		// Otherwise, we're in a polar region (|z| > 2/3):

		// Integer and fractional parts of φ:
		φFloor := math.Floor(φ)

		// Use baseFaceIndex in [0..3], clamped:
		baseFaceIndex := int(φFloor)
		if baseFaceIndex > 3 {
			baseFaceIndex = 3
		}

		// Calculate j+ and j-:
		jp := int((φ - φFloor) * (nSideFaces * math.Sqrt(3.0*(1.0-za))))
		jm := int((1.0 - (φ - φFloor)) * (nSideFaces * math.Sqrt(3.0*(1.0-za))))

		// Clamp j+ and j- to [0, nside-1]:
		if jp >= nside {
			jp = nside - 1
		}
		if jm >= nside {
			jm = nside - 1
		}

		// North polar cap if z > 0, else south polar cap:
		if z > 0 {
			faceIndex = baseFaceIndex
			ix = (nside - 1) - jm
			iy = (nside - 1) - jp
		} else {
			faceIndex = baseFaceIndex + 8
			ix = jp
			iy = jm
		}
	}

	// Bit interleaving (xyf2nest) to build the final NESTED pixel index:
	var utab = [256]int{
		0x0000, 0x0001, 0x0004, 0x0005, 0x0010, 0x0011, 0x0014, 0x0015,
		0x0040, 0x0041, 0x0044, 0x0045, 0x0050, 0x0051, 0x0054, 0x0055,
		0x0100, 0x0101, 0x0104, 0x0105, 0x0110, 0x0111, 0x0114, 0x0115,
		0x0140, 0x0141, 0x0144, 0x0145, 0x0150, 0x0151, 0x0154, 0x0155,
		0x0400, 0x0401, 0x0404, 0x0405, 0x0410, 0x0411, 0x0414, 0x0415,
		0x0440, 0x0441, 0x0444, 0x0445, 0x0450, 0x0451, 0x0454, 0x0455,
		0x0500, 0x0501, 0x0504, 0x0505, 0x0510, 0x0511, 0x0514, 0x0515,
		0x0540, 0x0541, 0x0544, 0x0545, 0x0550, 0x0551, 0x0554, 0x0555,
		0x1000, 0x1001, 0x1004, 0x1005, 0x1010, 0x1011, 0x1014, 0x1015,
		0x1040, 0x1041, 0x1044, 0x1045, 0x1050, 0x1051, 0x1054, 0x1055,
		0x1100, 0x1101, 0x1104, 0x1105, 0x1110, 0x1111, 0x1114, 0x1115,
		0x1140, 0x1141, 0x1144, 0x1145, 0x1150, 0x1151, 0x1154, 0x1155,
		0x1400, 0x1401, 0x1404, 0x1405, 0x1410, 0x1411, 0x1414, 0x1415,
		0x1440, 0x1441, 0x1444, 0x1445, 0x1450, 0x1451, 0x1454, 0x1455,
		0x1500, 0x1501, 0x1504, 0x1505, 0x1510, 0x1511, 0x1514, 0x1515,
		0x1540, 0x1541, 0x1544, 0x1545, 0x1550, 0x1551, 0x1554, 0x1555,
		0x4000, 0x4001, 0x4004, 0x4005, 0x4010, 0x4011, 0x4014, 0x4015,
		0x4040, 0x4041, 0x4044, 0x4045, 0x4050, 0x4051, 0x4054, 0x4055,
		0x4100, 0x4101, 0x4104, 0x4105, 0x4110, 0x4111, 0x4114, 0x4115,
		0x4140, 0x4141, 0x4144, 0x4145, 0x4150, 0x4151, 0x4154, 0x4155,
		0x4400, 0x4401, 0x4404, 0x4405, 0x4410, 0x4411, 0x4414, 0x4415,
		0x4440, 0x4441, 0x4444, 0x4445, 0x4450, 0x4451, 0x4454, 0x4455,
		0x4500, 0x4501, 0x4504, 0x4505, 0x4510, 0x4511, 0x4514, 0x4515,
		0x4540, 0x4541, 0x4544, 0x4545, 0x4550, 0x4551, 0x4554, 0x4555,
		0x5000, 0x5001, 0x5004, 0x5005, 0x5010, 0x5011, 0x5014, 0x5015,
		0x5040, 0x5041, 0x5044, 0x5045, 0x5050, 0x5051, 0x5054, 0x5055,
		0x5100, 0x5101, 0x5104, 0x5105, 0x5110, 0x5111, 0x5114, 0x5115,
		0x5140, 0x5141, 0x5144, 0x5145, 0x5150, 0x5151, 0x5154, 0x5155,
		0x5400, 0x5401, 0x5404, 0x5405, 0x5410, 0x5411, 0x5414, 0x5415,
		0x5440, 0x5441, 0x5444, 0x5445, 0x5450, 0x5451, 0x5454, 0x5455,
		0x5500, 0x5501, 0x5504, 0x5505, 0x5510, 0x5511, 0x5514, 0x5515,
		0x5540, 0x5541, 0x5544, 0x5545, 0x5550, 0x5551, 0x5554, 0x5555,
	}

	// Combine faceIndex, ix, iy into the final pixel index in the NESTED scheme:
	pix := faceIndex*nside*nside +
		(utab[ix&0xff] |
			(utab[(ix>>8)&0xff] << 16) |
			(utab[iy&0xff] << 1) |
			(utab[(iy>>8)&0xff] << 17))

	// Return the final pixel index:
	return pix
}

/*****************************************************************************************************************/
