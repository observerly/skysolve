/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package wcs

/*****************************************************************************************************************/

import (
	"math"
	"testing"

	"github.com/observerly/skysolve/pkg/transform"
)

/*****************************************************************************************************************/

func TestNewWCS(t *testing.T) {
	affine := transform.Affine2DParameters{
		A: 1,
		B: 0,
		C: 0,
		D: 1,
		E: 0,
		F: 0,
	}

	sip := transform.SIP2DParameters{
		AOrder: 1,
		BOrder: 1,
		APower: map[string]float64{
			"A_0_0": 0,
			"A_1_0": 0,
			"A_0_1": 0,
		},
		BPower: map[string]float64{
			"B_0_0": 0,
			"B_1_0": 0,
			"B_0_1": 0,
		},
	}

	wcs := NewWorldCoordinateSystem(1000, 1000, WCSParams{
		AffineParams: affine,
		Projection:   RADEC_TAN,
		SIPParams:    sip,
	})

	if wcs.CRPIX1 != 1000 {
		t.Errorf("CRPIX1 not set correctly")
	}

	if wcs.CRPIX2 != 1000 {
		t.Errorf("CRPIX2 not set correctly")
	}

	if wcs.CRVAL1 != 0 {
		t.Errorf("CRVAL1 not set correctly, expected 0 got %v", wcs.CRVAL1)
	}

	if wcs.CRVAL2 != 0 {
		t.Errorf("CRVAL2 not set correctly, expected 0 got %v", wcs.CRVAL2)
	}

	if wcs.CUNIT1 != "deg" {
		t.Errorf("CUNIT1 not set correctly")
	}

	if wcs.CUNIT2 != "deg" {
		t.Errorf("CUNIT2 not set correctly")
	}

	if wcs.CDELT1 != 1 {
		t.Errorf("CDELT1 not set correctly")
	}

	if wcs.CDELT2 != 1 {
		t.Errorf("CDELT2 not set correctly")
	}

	if wcs.CTYPE1 != "RA---TAN" {
		t.Errorf("CTYPE1 not set correctly")
	}

	if wcs.CTYPE2 != "DEC--TAN" {
		t.Errorf("CTYPE2 not set correctly")
	}

	if wcs.CD1_1 != 1 {
		t.Errorf("CD1_1 not set correctly")
	}

	if wcs.CD1_2 != 0 {
		t.Errorf("CD1_2 not set correctly")
	}

	if wcs.CD2_1 != 0 {
		t.Errorf("CD2_1 not set correctly")
	}

	if wcs.CD2_2 != 1 {
		t.Errorf("CD2_2 not set correctly")
	}
}

/*****************************************************************************************************************/

func TestPixelToEquatorialCoordinateAtImageCenter(t *testing.T) {
	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
	}

	coordinate := wcs.PixelToEquatorialCoordinate(1024.0, 1024.0)

	if coordinate.RA != 150.0 {
		t.Errorf("RA not calculated correctly")
	}

	if coordinate.Dec != 2.0 {
		t.Errorf("Dec not calculated correctly")
	}
}

/*****************************************************************************************************************/

func TestPixelToEquatorialCoordinate(t *testing.T) {
	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
	}

	coordinate := wcs.PixelToEquatorialCoordinate(1000.0, 1000.0)

	if coordinate.RA != 150.0066666672 {
		t.Errorf("RA not calculated correctly")
	}

	if coordinate.Dec != 1.9933333328 {
		t.Errorf("Dec not calculated correctly")
	}
}

/*****************************************************************************************************************/

func TestPixelToEquatorialCoordinateWithSIPDistortionAtImageCenter(t *testing.T) {
	sip := transform.SIP2DParameters{
		AOrder: 3,
		BOrder: 3,
		APower: map[string]float64{
			"A_2_0": 2.5e-5,
			"A_1_1": 1.8e-5,
			"A_0_2": 2.2e-5,
			"A_3_0": 1.1e-6,
			"A_2_1": 5.5e-6,
			"A_1_2": 4.3e-6,
			"A_0_3": 3.2e-6,
		},
		BPower: map[string]float64{
			"B_2_0": -2.5e-5,
			"B_1_1": -1.8e-5,
			"B_0_2": -2.2e-5,
			"B_3_0": -1.1e-6,
			"B_2_1": -5.5e-6,
			"B_1_2": -4.3e-6,
			"B_0_3": -3.2e-6,
		},
	}

	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
		SIP:    sip,
	}

	coordinate := wcs.PixelToEquatorialCoordinate(1024.0, 1024.0)

	if coordinate.RA != 150.0 {
		t.Errorf("RA not calculated correctly")
	}

	if coordinate.Dec != 2.0 {
		t.Errorf("Dec not calculated correctly")
	}
}

/*****************************************************************************************************************/

func TestPixelToEquatorialCoordinateWithSIPDistortion(t *testing.T) {
	sip := transform.SIP2DParameters{
		AOrder: 3,
		BOrder: 3,
		APower: map[string]float64{
			"A_2_0": 2.5e-5,
			"A_1_1": 1.8e-5,
			"A_0_2": 2.2e-5,
			"A_3_0": 1.1e-6,
			"A_2_1": 5.5e-6,
			"A_1_2": 4.3e-6,
			"A_0_3": 3.2e-6,
		},
		BPower: map[string]float64{
			"B_2_0": -2.5e-5,
			"B_1_1": -1.8e-5,
			"B_0_2": -2.2e-5,
			"B_3_0": -1.1e-6,
			"B_2_1": -5.5e-6,
			"B_1_2": -4.3e-6,
			"B_0_3": -3.2e-6,
		},
	}

	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
		SIP:    sip,
	}

	coordinate := wcs.PixelToEquatorialCoordinate(1000.0, 1000.0)

	if coordinate.RA != 150.0067104112035 {
		t.Errorf("RA not calculated correctly")
	}

	if coordinate.Dec != 1.9933770768034995 {
		t.Errorf("Dec not calculated correctly")
	}
}

/*****************************************************************************************************************/

func TestEquatorialCoordinateToPixelAtImageCenter(t *testing.T) {
	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
	}

	x, y := wcs.EquatorialCoordinateToPixel(150.0, 2.0)

	tolerance := 0.0000001

	if math.Abs((x - 1024.0)) > tolerance {
		t.Errorf("X not calculated correctly")
	}

	if math.Abs(y-1024.0) > tolerance {
		t.Errorf("Y not calculated correctly")
	}
}

/*****************************************************************************************************************/

func TestEquatorialCoordinateToPixel(t *testing.T) {
	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
	}

	coordinate := wcs.PixelToEquatorialCoordinate(1000.0, 1000.0)

	x, y := wcs.EquatorialCoordinateToPixel(coordinate.RA, coordinate.Dec)

	tolerance := 0.0000001

	if (x - 1000.0) > tolerance {
		t.Errorf("X not calculated correctly")
	}

	if (y - 1000.0) > tolerance {
		t.Errorf("Y not calculated correctly")
	}
}

/*****************************************************************************************************************/

func TestEquatorialCoordinateToPixelWithSIPDistortionAtImageCenter(t *testing.T) {
	sip := transform.SIP2DParameters{
		AOrder: 3,
		BOrder: 3,
		APower: map[string]float64{
			"A_0_2": 9.0886e-06,
			"A_0_3": 4.8066e-09,
			"A_1_1": 4.8146e-05,
			"A_1_2": -1.7096e-07,
			"A_2_0": 2.82e-05,
			"A_2_1": 3.3336e-08,
			"A_3_0": -1.8684e-07,
		},
		BPower: map[string]float64{
			"B_0_2": 4.1248e-05,
			"B_0_3": -1.9016e-07,
			"B_1_1": 1.4761e-05,
			"B_1_2": 2.1973e-08,
			"B_2_0": -6.4708e-06,
			"B_2_1": -1.8188e-07,
			"B_3_0": 1.0084e-10,
		},
	}

	isip := transform.SIP2DInverseParameters{
		APOrder: 3,
		APPower: map[string]float64{
			"AP_0_1": 3.6698e-06,
			"AP_0_2": -9.1825e-06,
			"AP_0_3": -3.8909e-09,
			"AP_1_0": -2.0239e-05,
			"AP_1_1": -4.8946e-05,
			"AP_1_2": 1.7951e-07,
			"AP_2_0": -2.8622e-05,
			"AP_2_1": -2.9553e-08,
			"AP_3_0": 1.9119e-07,
		},
		BPPower: map[string]float64{
			"BP_0_1": -2.1339e-05,
			"BP_0_2": -4.189e-05,
			"BP_0_3": 1.9696e-07,
			"BP_1_0": 2.8502e-06,
			"BP_1_1": -1.5089e-05,
			"BP_1_2": -2.0219e-08,
			"BP_2_0": 6.4625e-06,
			"BP_2_1": 1.849e-07,
			"BP_3_0": -7.6669e-10,
		},
	}

	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
		SIP:    sip,
		ISIP:   isip,
	}

	coordinate := wcs.PixelToEquatorialCoordinate(1024.0, 1024.0)

	x, y := wcs.EquatorialCoordinateToPixel(coordinate.RA, coordinate.Dec)

	tolerance := 0.0000001

	if math.Abs((x - 1024.0)) > tolerance {
		t.Errorf("X not calculated correctly")
	}

	if math.Abs(y-1024.0) > tolerance {
		t.Errorf("Y not calculated correctly")
	}
}

/*****************************************************************************************************************/

func TestEquatorialCoordinateToPixelWithSIPDistortion(t *testing.T) {
	sip := transform.SIP2DParameters{
		AOrder: 3,
		BOrder: 3,
		APower: map[string]float64{
			"A_0_2": 9.0886e-06,
			"A_0_3": 4.8066e-09,
			"A_1_1": 4.8146e-05,
			"A_1_2": -1.7096e-07,
			"A_2_0": 2.82e-05,
			"A_2_1": 3.3336e-08,
			"A_3_0": -1.8684e-07,
		},
		BPower: map[string]float64{
			"B_0_2": 4.1248e-05,
			"B_0_3": -1.9016e-07,
			"B_1_1": 1.4761e-05,
			"B_1_2": 2.1973e-08,
			"B_2_0": -6.4708e-06,
			"B_2_1": -1.8188e-07,
			"B_3_0": 1.0084e-10,
		},
	}

	isip := transform.SIP2DInverseParameters{
		APOrder: 3,
		APPower: map[string]float64{
			"AP_0_1": 3.6698e-06,
			"AP_0_2": -9.1825e-06,
			"AP_0_3": -3.8909e-09,
			"AP_1_0": -2.0239e-05,
			"AP_1_1": -4.8946e-05,
			"AP_1_2": 1.7951e-07,
			"AP_2_0": -2.8622e-05,
			"AP_2_1": -2.9553e-08,
			"AP_3_0": 1.9119e-07,
		},
		BPPower: map[string]float64{
			"BP_0_1": -2.1339e-05,
			"BP_0_2": -4.189e-05,
			"BP_0_3": 1.9696e-07,
			"BP_1_0": 2.8502e-06,
			"BP_1_1": -1.5089e-05,
			"BP_1_2": -2.0219e-08,
			"BP_2_0": 6.4625e-06,
			"BP_2_1": 1.849e-07,
			"BP_3_0": -7.6669e-10,
		},
	}

	wcs := WCS{
		CRPIX1: 1024.0, // Assuming a 2048x2048 image, center pixel
		CRPIX2: 1024.0,
		CRVAL1: 150.0,         // Right Ascension in degrees (e.g., 10h)
		CRVAL2: 2.0,           // Declination in degrees
		CD1_1:  -0.0002777778, // -1/3600 deg/pixel (scale: ~1 arcsec/pixel)
		CD1_2:  0.0,           // No rotation
		CD2_1:  0.0,           // No rotation
		CD2_2:  0.0002777778,  // 1/3600 deg/pixel
		SIP:    sip,
		ISIP:   isip,
	}

	coordinate := wcs.PixelToEquatorialCoordinate(1000.0, 1000.0)

	x, y := wcs.EquatorialCoordinateToPixel(coordinate.RA, coordinate.Dec)

	tolerance := 0.001

	if math.Abs((x - 1000.0)) > tolerance {
		t.Errorf("X not calculated correctly")
	}

	if math.Abs(y-1000.0) > tolerance {
		t.Errorf("Y not calculated correctly")
	}
}

/*****************************************************************************************************************/
