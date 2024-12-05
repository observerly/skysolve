/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package projection

/*****************************************************************************************************************/

import (
	"math"
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
