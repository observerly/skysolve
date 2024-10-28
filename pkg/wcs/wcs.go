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
	CRPIX1 float64 // Reference pixel X
	CRPIX2 float64 // Reference pixel Y
	CRVAL1 float64 // Reference RA
	CRVAL2 float64 // Reference Dec
	CD1_1  float64 // Affine transform parameter A
	CD1_2  float64 // Affine transform parameter B
	CD2_1  float64 // Affine transform parameter C
	CD2_2  float64 // Affine transform parameter D
	E      float64 // Affine translation parameter e (optional)
	F      float64 // Affine translation parameter f (optional)
}

/*****************************************************************************************************************/

func NewWorldCoordinateSystem(xc float64, yc float64, params transform.Affine2DParameters) WCS {
	// Create a new WCS object:
	wcs := WCS{
		CRPIX1: float64(xc),
		CRPIX2: float64(yc),
		CRVAL1: 0,
		CRVAL2: 0,
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
