/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package wcs

/*****************************************************************************************************************/

import (
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
