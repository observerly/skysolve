/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package fov

import "math"

/*****************************************************************************************************************/

type PixelScale struct {
	X float64 // Pixel size in the x direction (in degrees)
	Y float64 // Pixel size in the y direction (in degrees)
}

/*****************************************************************************************************************/

func GetRadialExtent(
	xs float64,
	ys float64,
	pixelScale PixelScale,
) float64 {
	// Calculate the field of view in the x direction (in degrees):
	xr := pixelScale.X * xs

	// Calculate the field of view in the y direction (in degrees):
	yr := pixelScale.Y * ys

	r := math.Min(xr, yr)

	// Calculate the radial field of view (in degrees):
	return math.Sqrt(r*r + r*r)
}

/*****************************************************************************************************************/
