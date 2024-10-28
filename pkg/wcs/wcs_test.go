/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package wcs

/*****************************************************************************************************************/

import (
	"fmt"
	"testing"

	"github.com/observerly/skysolve/pkg/transform"
)

/*****************************************************************************************************************/

func TestNewWCS(t *testing.T) {
	wcs := NewWorldCoordinateSystem(1000, 1000, transform.Affine2DParameters{
		A: 1,
		B: 0,
		C: 0,
		D: 1,
		E: 0,
		F: 0,
	})

	if wcs.CRPIX1 != 1000 {
		t.Errorf("CRPIX1 not set correctly")
	}

	if wcs.CRPIX2 != 1000 {
		t.Errorf("CRPIX2 not set correctly")
	}

	if wcs.CRVAL1 == 0 {
		t.Errorf("CRVAL1 not set correctly")
	}

	if wcs.CRVAL2 == 0 {
		t.Errorf("CRVAL2 not set correctly")
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

func TestPixelToEquatorialCoordinate(t *testing.T) {
	wcs := WCS{
		CRPIX1: 200,
		CRPIX2: 200,
		CRVAL1: 0,
		CRVAL2: 0,
		CD1_1:  0.2,
		CD1_2:  30,
		CD2_1:  0.2,
		CD2_2:  0.2,
	}

	coordinate := wcs.PixelToEquatorialCoordinate(0, 0)

	fmt.Println(coordinate)

	if coordinate.RA != 280 {
		t.Errorf("RA not calculated correctly")
	}

	if coordinate.Dec != 80 {
		t.Errorf("Dec not calculated correctly")
	}
}

/*****************************************************************************************************************/
