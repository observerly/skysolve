/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package wcs

/*****************************************************************************************************************/

import (
	"testing"
)

/*****************************************************************************************************************/

func TestNewWCS(t *testing.T) {
	wcs := NewWorldCoordinateSystem(WCS{
		CRPIX1: 0,
		CRPIX2: 0,
		CRVAL1: 0,
		CRVAL2: 0,
		CD1_1:  0,
		CD1_2:  0,
		CD2_1:  0,
		CD2_2:  0,
	})

	if wcs.CRPIX1 != 0 {
		t.Errorf("CRPIX1 not set correctly")
	}

	if wcs.CRPIX2 != 0 {
		t.Errorf("CRPIX2 not set correctly")
	}

	if wcs.CRVAL1 != 0 {
		t.Errorf("CRVAL1 not set correctly")
	}

	if wcs.CRVAL2 != 0 {
		t.Errorf("CRVAL2 not set correctly")
	}

	if wcs.CD1_1 != 0 {
		t.Errorf("CD1_1 not set correctly")
	}

	if wcs.CD1_2 != 0 {
		t.Errorf("CD1_2 not set correctly")
	}

	if wcs.CD2_1 != 0 {
		t.Errorf("CD2_1 not set correctly")
	}

	if wcs.CD2_2 != 0 {
		t.Errorf("CD2_2 not set correctly")
	}
}

/*****************************************************************************************************************/

func TestPixelToEquatorialCoordinate(t *testing.T) {
	wcs := WCS{
		CRPIX1: 0,
		CRPIX2: 0,
		CRVAL1: 180,
		CRVAL2: 0,
		CD1_1:  1,
		CD1_2:  0,
		CD2_1:  0,
		CD2_2:  1,
	}

	coordinate := wcs.PixelToEquatorialCoordinate(0, 0)

	if coordinate.RA != 180 {
		t.Errorf("RA not calculated correctly")
	}

	if coordinate.Dec != 0 {
		t.Errorf("Dec not calculated correctly")
	}
}

/*****************************************************************************************************************/
