/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package wcs

/*****************************************************************************************************************/

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/spatial"
	"github.com/observerly/skysolve/pkg/transform"
	"gonum.org/v1/gonum/mat"
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
	Projection       CoordinateProjectionType         // Projection type e.g., "TAN", or "TAN-SIP"
	AffineParams     transform.Affine2DParameters     // Affine transformation parameters
	ReferenceX       float64                          // Reference X coordinate
	ReferenceY       float64                          // Reference Y coordinate
	SIPForwardParams transform.SIP2DForwardParameters // SIP forward transformation (distortion) coefficients, x, y to RA, Dec
	SIPInverseParams transform.SIP2DInverseParameters // SIP inverse transformation (distortion) coefficients RA, Dec to x, y
}

/*****************************************************************************************************************/

type WCS struct {
	WCAXES int                              `json:"wcaxes" hdu:"WCAXES" default:"2"`        // Number of world coordinate axes
	CRPIX1 float64                          `json:"crpix1" hdu:"CRPIX1"`                    // Reference pixel X
	CRPIX2 float64                          `json:"crpix2" hdu:"CRPIX2"`                    // Reference pixel Y
	CRVAL1 float64                          `json:"crval1" hdu:"CRVAL1" default:"0.0"`      // Reference RA (example default, often specific to image)
	CRVAL2 float64                          `json:"crval2" hdu:"CRVAL2" default:"0.0"`      // Reference Dec (example default, often specific to image)
	CTYPE1 string                           `json:"ctype1" hdu:"CTYPE1" default:"RA---TAN"` // Coordinate type for axis 1, typically RA with TAN projection
	CTYPE2 string                           `json:"ctype2" hdu:"CTYPE2" default:"DEC--TAN"` // Coordinate type for axis 2, typically DEC with TAN projection
	CDELT1 float64                          `json:"cdelt1" hdu:"CDELT1"`                    // Coordinate increment for axis 1 (no default)
	CDELT2 float64                          `json:"cdelt2" hdu:"CDELT2"`                    // Coordinate increment for axis 2 (no default)
	CUNIT1 string                           `json:"cunit1" hdu:"CUNIT1" default:"deg"`      // Coordinate unit for axis 1, defaulted to degrees
	CUNIT2 string                           `json:"cunit2" hdu:"CUNIT2" default:"deg"`      // Coordinate unit for axis 2, defaulted to degrees
	CD1_1  float64                          `json:"cd1_1" hdu:"CD1_1"`                      // Affine transform parameter A (no default)
	CD1_2  float64                          `json:"cd1_2" hdu:"CD1_2"`                      // Affine transform parameter B (no default)
	CD2_1  float64                          `json:"cd2_1" hdu:"CD2_1"`                      // Affine transform parameter C (no default)
	CD2_2  float64                          `json:"cd2_2" hdu:"CD2_2"`                      // Affine transform parameter D (no default)
	E      float64                          `json:"e" hdu:"E"`                              // Affine translation parameter e (optional, no default)
	F      float64                          `json:"f" hdu:"F"`                              // Affine translation parameter f (optional, no default)
	FSIP   transform.SIP2DForwardParameters `json:"fsip" hdu:"FSIP"`                        // SIP forward transformation (distortion) coefficients
	ISIP   transform.SIP2DInverseParameters `json:"isip" hdu:"ISIP"`                        // SIP inverse transformation (distortion) coefficients
}

/*****************************************************************************************************************/

// NewWorldCoordinateSystem creates a new WCS object with correctly mapped affine parameters.
func NewWorldCoordinateSystem(xc float64, yc float64, params WCSParams) WCS {
	// Get the coordinate projection types, e.g., "RA---TAN", or "RA---TAN-SIP":
	ctypes := params.Projection.ToCTypes()

	// Set CD matrix elements
	CD1_1 := params.AffineParams.A
	CD1_2 := params.AffineParams.B
	CD2_1 := params.AffineParams.D
	CD2_2 := params.AffineParams.E

	// Calculate CRVAL1 and CRVAL2 based on the affine transformation at the reference pixel (xc, yc) being at (0,0):
	crval1 := params.AffineParams.C
	crval2 := params.AffineParams.F

	// Create a new WCS object with correctly mapped CD matrix and CRVALs
	wcs := WCS{
		WCAXES: 2, // We always assume two world coordinate axes, RA and Dec.
		CRPIX1: xc,
		CRPIX2: yc,
		CRVAL1: crval1,
		CRVAL2: crval2,
		CUNIT1: "deg", // Degrees
		CUNIT2: "deg", // Degrees
		CTYPE1: ctypes.CType1,
		CTYPE2: ctypes.CType2,
		CD1_1:  CD1_1,
		CD1_2:  CD1_2,
		CD2_1:  CD2_1,
		CD2_2:  CD2_2,
		FSIP:   params.SIPForwardParams,
		ISIP:   params.SIPInverseParams,
	}

	// Calculate the coordinate increment for axis 1 (CDELT1)
	wcs.CDELT1 = -math.Sqrt(wcs.CD1_1*wcs.CD1_1 + wcs.CD2_1*wcs.CD2_1)

	// Calculate the coordinate increment for axis 2 (CDELT2)
	wcs.CDELT2 = math.Sqrt(wcs.CD1_2*wcs.CD1_2 + wcs.CD2_2*wcs.CD2_2)

	// If a reference pixel is provided, then re-calibrate the WCS object to the reference pixel
	if params.ReferenceX != 0 && params.ReferenceY != 0 {
		eq := wcs.PixelToEquatorialCoordinate(params.ReferenceX, params.ReferenceY)

		wcs.CRVAL1 = eq.RA
		wcs.CRVAL2 = eq.Dec

		wcs.CRPIX1 = params.ReferenceX
		wcs.CRPIX2 = params.ReferenceY
	}

	return wcs
}

/*****************************************************************************************************************/

func (wcs *WCS) SolveForCentroid() (coordinate astrometry.ICRSEquatorialCoordinate) {
	return wcs.PixelToEquatorialCoordinate(wcs.CRPIX1, wcs.CRPIX2)
}

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

	// Compute non-linear SIP distortion corrections A and B:
	A := 0.0
	B := 0.0

	// Apply A polynomial corrections:
	for term, coeff := range wcs.FSIP.APower {
		i, j, err := parseSIPTerm(term, "A")
		if err != nil {
			continue
		}
		A += coeff * math.Pow(deltaX, float64(i)) * math.Pow(deltaY, float64(j))
	}

	// Apply B polynomial corrections:
	for term, coeff := range wcs.FSIP.BPower {
		i, j, err := parseSIPTerm(term, "B")
		if err != nil {
			continue
		}
		B += coeff * math.Pow(deltaX, float64(i)) * math.Pow(deltaY, float64(j))
	}

	// Apply forward SIP transformation to correct for non-linear distortions:
	deltaX += A
	deltaY += B

	// Calculate the reference equatorial coordinate for the right ascension:
	ra := wcs.CD1_1*deltaX + wcs.CD1_2*deltaY + wcs.CRVAL1

	// Correct for large values of RA:
	ra = math.Mod(ra, 360)

	// Correct for negative values of RA:
	if ra < 0 {
		ra += 360
	}

	// Calculate the reference equatorial coordinate for the declination:
	dec := wcs.CD2_1*deltaX + wcs.CD2_2*deltaY + wcs.CRVAL2

	// Correct for large values of declination:
	dec = math.Mod(dec, 90)

	return astrometry.ICRSEquatorialCoordinate{
		RA:  ra,
		Dec: dec,
	}
}

/*****************************************************************************************************************/

// without iterative removal of SIP distortions.
func (wcs *WCS) EquatorialCoordinateToPixel(
	ra, dec float64,
) (x, y float64) {
	// Find the determinant of the CD matrix, for the inverse CD matrix:
	det := wcs.CD1_1*wcs.CD2_2 - wcs.CD1_2*wcs.CD2_1

	// If the determinant is zero, then it is considered singular and the inverse matrix is defaulted:
	invCD1_1 := 0.0
	invCD1_2 := 0.0
	invC := 0.0
	invCD2_1 := 0.0
	invCD2_2 := 0.0
	invF := 0.0

	// If it is non-zero, then compute the inverse CD matrix:
	if det != 0 {
		invCD1_1 = wcs.CD2_2 / det
		invCD1_2 = -wcs.CD1_2 / det
		invC = (wcs.CD1_2*wcs.CRVAL2 - wcs.CRVAL1*wcs.CD2_2) / det
		invCD2_1 = -wcs.CD2_1 / det
		invCD2_2 = wcs.CD1_1 / det
		invF = (wcs.CRVAL1*wcs.CD2_1 - wcs.CD1_1*wcs.CRVAL2) / det
	}

	// Apply the inverse CD matrix to get initial deltaX and deltaY
	deltaX := invCD1_1*ra + invCD1_2*dec + invC
	deltaY := invCD2_1*ra + invCD2_2*dec + invF

	// Compute non-linear SIP distortion corrections A and B:
	A := 0.0
	B := 0.0

	// Apply inverse A polynomial corrections
	for term, coeff := range wcs.ISIP.APPower {
		i, j, err := parseSIPTerm(term, "AP")
		if err != nil {
			continue
		}
		A += coeff * math.Pow(deltaX, float64(i)) * math.Pow(deltaY, float64(j))
	}

	// Apply inverse B polynomial corrections
	for term, coeff := range wcs.ISIP.BPPower {
		i, j, err := parseSIPTerm(term, "BP")
		if err != nil {
			continue
		}
		B += coeff * math.Pow(deltaX, float64(i)) * math.Pow(deltaY, float64(j))
	}

	// Apply backward SIP transformation to correct for non-linear distortions:
	deltaX += A
	deltaY += B

	// Add the reference pixel coordinates to obtain final pixel positions
	x = deltaX + wcs.CRPIX1
	y = deltaY + wcs.CRPIX2

	return x, y
}

/*****************************************************************************************************************/

// Each match provides two point correspondences: C and D
type PointPair struct {
	X, Y    float64 // Generated Quad NormalisedC and NormalisedD
	RA, Dec float64 // Source Quad SourceRA and SourceDec
}

/*****************************************************************************************************************/

// ComputeAffineTransformation computes the affine transformation parameters based on matched quads.
// It returns the affine parameters and an error if the computation fails.
func ComputeAffineTransformation(matches []spatial.QuadMatch) (transform.Affine2DParameters, float64, float64, error) {
	var pairs []PointPair

	// Iterate over each match to extract all four point correspondences:
	for _, match := range matches {
		// Extract Point A:
		pairs = append(pairs, PointPair{
			X:   match.Quad.A.X,   // Generated Quad X
			Y:   match.Quad.A.Y,   // Generated Quad Y
			RA:  match.Quad.A.RA,  // Source Quad RA
			Dec: match.Quad.A.Dec, // Source Quad Dec
		})

		// Extract Point B:
		pairs = append(pairs, PointPair{
			X:   match.Quad.B.X,
			Y:   match.Quad.B.Y,
			RA:  match.Quad.B.RA,
			Dec: match.Quad.B.Dec,
		})

		// Extract Point C:
		pairs = append(pairs, PointPair{
			X:   match.Quad.C.X,
			Y:   match.Quad.C.Y,
			RA:  match.Quad.C.RA,
			Dec: match.Quad.C.Dec,
		})

		// Extract Point D:
		pairs = append(pairs, PointPair{
			X:   match.Quad.D.X,
			Y:   match.Quad.D.Y,
			RA:  match.Quad.D.RA,
			Dec: match.Quad.D.Dec,
		})
	}

	n := len(pairs)
	if n < 2 { // Need at least three point correspondences for affine transformation:
		return transform.Affine2DParameters{}, math.Inf(1), math.Inf(1), errors.New("not enough point correspondences to compute affine transformation")
	}

	// Thus, for N points, we have 2N equations and 6 unknowns (a, b, c, d, e, f):
	A := mat.NewDense(2*n, 6, nil)
	bVec := mat.NewVecDense(2*n, nil)

	for i, pair := range pairs {
		// First equation: RA = a*X + b*Y + c:
		A.Set(2*i, 0, pair.X) // a
		A.Set(2*i, 1, pair.Y) // b
		A.Set(2*i, 2, 1.0)    // c
		A.Set(2*i, 3, 0.0)    // d
		A.Set(2*i, 4, 0.0)    // e
		A.Set(2*i, 5, 0.0)    // f
		bVec.SetVec(2*i, pair.RA)

		// Second equation: Dec = d*X + e*Y + f:
		A.Set(2*i+1, 0, 0.0)    // a
		A.Set(2*i+1, 1, 0.0)    // b
		A.Set(2*i+1, 2, 0.0)    // c
		A.Set(2*i+1, 3, pair.X) // d
		A.Set(2*i+1, 4, pair.Y) // e
		A.Set(2*i+1, 5, 1.0)    // f
		bVec.SetVec(2*i+1, pair.Dec)
	}

	// Solve the least squares problem: A * params = b:
	var qr mat.QR
	qr.Factorize(A)

	var params mat.VecDense
	err := qr.SolveVecTo(&params, false, bVec)
	if err != nil {
		return transform.Affine2DParameters{}, math.Inf(1), math.Inf(1), fmt.Errorf("failed to solve affine transformation: %v", err)
	}

	return transform.Affine2DParameters{
		A: params.AtVec(0),
		B: params.AtVec(1),
		C: params.AtVec(2),
		D: params.AtVec(3),
		E: params.AtVec(4),
		F: params.AtVec(5),
	}, 0, 0, nil
}

/*****************************************************************************************************************/
