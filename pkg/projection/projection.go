/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package projection

/*****************************************************************************************************************/

import (
	"math"

	"github.com/observerly/skysolve/pkg/astrometry"
)

/*****************************************************************************************************************/

var RAD2DEG = 180 / math.Pi

/*****************************************************************************************************************/

var DEG2RAD = math.Pi / 180

/*****************************************************************************************************************/

func Radians(degrees float64) float64 {
	return degrees * DEG2RAD
}

/*****************************************************************************************************************/

func Degrees(radians float64) float64 {
	return radians * RAD2DEG
}

/*****************************************************************************************************************/

func ConvertEquatorialToGnomic(ra, dec, ra0, dec0 float64) (x, y float64) {
	// Threshold to determine if cosalt1 is effectively zero:
	const epsilon = 1e-10

	// Convert all coordinates from degrees to radians:
	ra = Radians(ra)
	dec = Radians(dec)
	ra0 = Radians(ra0)
	dec0 = Radians(dec0)

	// Gnomonic projection formula:
	cosalt1 := math.Sin(dec0)*math.Sin(dec) + math.Cos(dec0)*math.Cos(dec)*math.Cos(ra-ra0)

	// Check for division by zero:
	if cosalt1 < epsilon {
		return 0, 0
	}

	// Calculate the x coordinate:
	x = math.Cos(dec) * math.Sin(ra-ra0) / cosalt1

	// Calculate the y coordinate:
	y = (math.Cos(dec0)*math.Sin(dec) - math.Sin(dec0)*math.Cos(dec)*math.Cos(ra-ra0)) / cosalt1

	return x, y
}

/*****************************************************************************************************************/

func ConvertEquatorialToLambertCylindricalCartesian(eq astrometry.ICRSEquatorialCoordinate, z float64) (x, y float64) {
	// Calculate the y coordinate:
	y = 3 * (math.Pi / 8) * z

	// Return the x and y coordinates in degrees:
	return eq.RA, Degrees(y)
}

/*****************************************************************************************************************/

func ConvertEquatorialToInterruptedCollignonCartesian(eq astrometry.ICRSEquatorialCoordinate, z float64) (x, y float64) {
	sigma := math.Sqrt(3 * (1 - math.Abs(z)))

	// Calculate the y coordinate:
	y = (math.Pi / 4) * (2 - sigma)

	facetX := (math.Pi / 4) * (2 * math.Floor(2+(2*Radians(eq.RA)/math.Pi)-3))
	// Calculate the x coordinate:
	x = facetX + sigma*(Radians(eq.Dec)-facetX)
	// Return the x and y coordinates in degrees:
	return Degrees(x), Degrees(math.Copysign(y, eq.Dec))
}

/*****************************************************************************************************************/

func GetEquatorialCoordinateFromPolarOffset(ra, dec, z, θ float64) (x, y float64) {
	// Convert all coordinates from degrees to radians:
	ra0 := Radians(ra)
	dec0 := Radians(dec)
	r := Radians(z)
	θ = Radians(θ)

	// Calculate the declination for the polar offset:
	dec = math.Asin(
		math.Sin(dec0)*math.Cos(r) +
			math.Cos(dec0)*math.Sin(r)*math.Cos(θ),
	)

	// Calculate the right ascension for the polar offset:
	ra = ra0 + math.Atan2(
		math.Sin(θ)*math.Sin(r)*math.Cos(dec0),
		math.Cos(r)-math.Sin(dec0)*math.Sin(dec),
	)

	// Normalize the right ascension to the range [0, 2π) if it is less than 0:
	for ra < 0 {
		ra += 2.0 * math.Pi
	}

	// Normalize the right ascension to the range [0, 2π) if it is greater than 2π:
	for ra >= 2.0*math.Pi {
		ra -= 2.0 * math.Pi
	}
	// Return the right ascension and declination in degrees:
	return Degrees(ra), Degrees(dec)
}

/*****************************************************************************************************************/
