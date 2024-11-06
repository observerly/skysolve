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

type CoordinateProjectionType int

const (
	RADEC_TAN CoordinateProjectionType = iota
	RADEC_TANSIP
)

/*****************************************************************************************************************/

type CTypeP struct {
	CType1 string
	CType2 string
}

/*****************************************************************************************************************/

func (c CoordinateProjectionType) ToCTypes() CTypeP {
	switch c {
	case RADEC_TAN:
		return CTypeP{
			CType1: "RA---TAN",
			CType2: "DEC--TAN",
		}
	case RADEC_TANSIP:
		return CTypeP{
			CType1: "RA---TAN-SIP",
			CType2: "DEC--TAN-SIP",
		}
	default:
		return CTypeP{
			CType1: "RA---TAN",
			CType2: "DEC--TAN",
		}
	}
}

/*****************************************************************************************************************/

type WCSParams struct {
	Projection   CoordinateProjectionType     // Projection type e.g., "TAN", or "TAN-SIP"
	AffineParams transform.Affine2DParameters // Affine transformation parameters
	SIPParams    transform.SIP2DParameters    // SIP transformation (distortion) coefficients
}

/*****************************************************************************************************************/

type WCS struct {
	WCAXES int                       `hdu:"WCAXES" default:"2"`        // Number of world coordinate axes
	CRPIX1 float64                   `hdu:"CRPIX1"`                    // Reference pixel X
	CRPIX2 float64                   `hdu:"CRPIX2"`                    // Reference pixel Y
	CRVAL1 float64                   `hdu:"CRVAL1" default:"0.0"`      // Reference RA (example default, often specific to image)
	CRVAL2 float64                   `hdu:"CRVAL2" default:"0.0"`      // Reference Dec (example default, often specific to image)
	CTYPE1 string                    `hdu:"CTYPE1" default:"RA---TAN"` // Coordinate type for axis 1, typically RA with TAN projection
	CTYPE2 string                    `hdu:"CTYPE2" default:"DEC--TAN"` // Coordinate type for axis 2, typically DEC with TAN projection
	CDELT1 float64                   `hdu:"CDELT1"`                    // Coordinate increment for axis 1 (no default)
	CDELT2 float64                   `hdu:"CDELT2"`                    // Coordinate increment for axis 2 (no default)
	CUNIT1 string                    `hdu:"CUNIT1" default:"deg"`      // Coordinate unit for axis 1, defaulted to degrees
	CUNIT2 string                    `hdu:"CUNIT2" default:"deg"`      // Coordinate unit for axis 2, defaulted to degrees
	CD1_1  float64                   `hdu:"CD1_1"`                     // Affine transform parameter A (no default)
	CD1_2  float64                   `hdu:"CD1_2"`                     // Affine transform parameter B (no default)
	CD2_1  float64                   `hdu:"CD2_1"`                     // Affine transform parameter C (no default)
	CD2_2  float64                   `hdu:"CD2_2"`                     // Affine transform parameter D (no default)
	E      float64                   `hdu:"E"`                         // Affine translation parameter e (optional, no default)
	F      float64                   `hdu:"F"`                         // Affine translation parameter f (optional, no default)
	SIP    transform.SIP2DParameters // SIP transformation (distortion) coefficients
}

/*****************************************************************************************************************/

func NewWorldCoordinateSystem(xc float64, yc float64, params WCSParams) WCS {
	// Get the coordinate projection types, e.g., "RA---TAN", or "RA---TAN-SIP":
	ctypes := params.Projection.ToCTypes()

	// Create a new WCS object:
	wcs := WCS{
		WCAXES: 2, // We always assume two world coordinate axes, RA and Dec.
		CRPIX1: float64(xc),
		CRPIX2: float64(yc),
		CRVAL1: 0,
		CRVAL2: 0,
		CUNIT1: "deg", // We always assume degrees.
		CUNIT2: "deg", // We always assume degrees.
		CTYPE1: ctypes.CType1,
		CTYPE2: ctypes.CType2,
		CD1_1:  params.AffineParams.A,
		CD1_2:  params.AffineParams.B,
		CD2_1:  params.AffineParams.C,
		CD2_2:  params.AffineParams.D,
		E:      params.AffineParams.E,
		F:      params.AffineParams.F,
		SIP:    params.SIPParams,
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
/*****************************************************************************************************************/

// Helper function to parse SIP term keys
func parseSIPTerm(term, prefix string) (i int, j int, err error) {
	parts := strings.Split(term, "_")

	if len(parts) != 3 || parts[0] != prefix {
		return 0, 0, fmt.Errorf("invalid SIP term format: %s", term)
	}

	_, err = fmt.Sscanf(parts[1]+" "+parts[2], "%d %d", &i, &j)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse SIP term: %s", term)
	}

	return i, j, nil
}

/*****************************************************************************************************************/

func (wcs *WCS) PixelToEquatorialCoordinate(
	x, y float64,
) (coordinate astrometry.ICRSEquatorialCoordinate) {
	// Compute the offsets from the reference pixel
	deltaX := x - wcs.CRPIX1 // Offset in X
	deltaY := y - wcs.CRPIX2 // Offset in Y

	// Calculate the reference equatorial coordinate for the right ascension:
	ra := wcs.CD1_1*deltaX + wcs.CD1_2*deltaY + wcs.E

	// Correct for large values of RA:
	ra = math.Mod(ra, 360)

	// Correct for negative values of RA:
	if ra < 0 {
		ra += 360
	}

	// Calculate the reference equatorial coordinate for the declination:
	dec := wcs.CD2_1*deltaX + wcs.CD2_2*deltaY + wcs.F

	// Correct for large values of declination:
	dec = math.Mod(dec, 90)

	return astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}
}

/*****************************************************************************************************************/
