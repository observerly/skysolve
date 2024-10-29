/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package wcs

/*****************************************************************************************************************/

import (
	"math"

	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/transform"
)

/*****************************************************************************************************************/

type WCS struct {
	CRPIX1 float64 `hdu:"CRPIX1"`                    // Reference pixel X
	CRPIX2 float64 `hdu:"CRPIX2"`                    // Reference pixel Y
	CRVAL1 float64 `hdu:"CRVAL1" default:"0.0"`      // Reference RA (example default, often specific to image)
	CRVAL2 float64 `hdu:"CRVAL2" default:"0.0"`      // Reference Dec (example default, often specific to image)
	CTYPE1 string  `hdu:"CTYPE1" default:"RA---TAN"` // Coordinate type for axis 1, typically RA with TAN projection
	CTYPE2 string  `hdu:"CTYPE2" default:"DEC--TAN"` // Coordinate type for axis 2, typically DEC with TAN projection
	CDELT1 float64 `hdu:"CDELT1"`                    // Coordinate increment for axis 1 (no default)
	CDELT2 float64 `hdu:"CDELT2"`                    // Coordinate increment for axis 2 (no default)
	CUNIT1 string  `hdu:"CUNIT1" default:"deg"`      // Coordinate unit for axis 1, defaulted to degrees
	CUNIT2 string  `hdu:"CUNIT2" default:"deg"`      // Coordinate unit for axis 2, defaulted to degrees
	CD1_1  float64 `hdu:"CD1_1"`                     // Affine transform parameter A (no default)
	CD1_2  float64 `hdu:"CD1_2"`                     // Affine transform parameter B (no default)
	CD2_1  float64 `hdu:"CD2_1"`                     // Affine transform parameter C (no default)
	CD2_2  float64 `hdu:"CD2_2"`                     // Affine transform parameter D (no default)
	E      float64 `hdu:"E"`                         // Affine translation parameter e (optional, no default)
	F      float64 `hdu:"F"`                         // Affine translation parameter f (optional, no default)
}

/*****************************************************************************************************************/

func NewWorldCoordinateSystem(xc float64, yc float64, params transform.Affine2DParameters) WCS {
	// Create a new WCS object:
	wcs := WCS{
		CRPIX1: float64(xc),
		CRPIX2: float64(yc),
		CRVAL1: 0,
		CRVAL2: 0,
		CUNIT1: "deg",      // We always assume degrees.
		CUNIT2: "deg",      // We always assume degrees.
		CTYPE1: "RA---TAN", // We always assume a tangential projection.
		CTYPE2: "DEC--TAN", // We always assume a tangential projection.
		CD1_1:  params.A,
		CD1_2:  params.B,
		CD2_1:  params.C,
		CD2_2:  params.D,
		E:      params.E,
		F:      params.F,
	}

	// Calculate the reference equatorial coordinate:
	eq := wcs.SolveForCentroid()

	// Set the reference equatorial coordinate:
	wcs.CRVAL1 = eq.RA
	wcs.CRVAL2 = eq.Dec

	// Calculate the coordinate increment for axis 1:
	wcs.CDELT1 = math.Sqrt(wcs.CD1_1*wcs.CD1_1 + wcs.CD2_1*wcs.CD2_1)

	// Calculate the coordinate increment for axis 2:
	wcs.CDELT2 = math.Sqrt(wcs.CD1_2*wcs.CD1_2 + wcs.CD2_2*wcs.CD2_2)

	return wcs
}

/*****************************************************************************************************************/

func (wcs *WCS) SolveForCentroid() (coordinate astrometry.ICRSEquatorialCoordinate) {
	return wcs.PixelToEquatorialCoordinate(wcs.CRPIX1, wcs.CRPIX2)
}

/*****************************************************************************************************************/

func (wcs *WCS) PixelToEquatorialCoordinate(
	x, y float64,
) (coordinate astrometry.ICRSEquatorialCoordinate) {
	// Calculate the reference equatorial coordinate for the right ascension:
	ra := wcs.CD1_1*wcs.CRPIX1 + wcs.CD1_2*wcs.CRPIX2 + wcs.E

	// Correct for large values of RA:
	if ra > 360 {
		ra = math.Mod(ra, 360)
	}

	// Correct for negative values of RA:
	if ra < 0 {
		ra += 360
	}

	// Calculate the reference equatorial coordinate for the declination:
	dec := wcs.CD2_1*wcs.CRPIX1 + wcs.CD2_2*wcs.CRPIX2 + wcs.F

	// Correct for large values of declination:
	dec = math.Mod(dec, 90)

	return astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}
}

/*****************************************************************************************************************/
