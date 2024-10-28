/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package transform

/*****************************************************************************************************************/

import "testing"

/*****************************************************************************************************************/

func TestAffine2DParameters(t *testing.T) {
	affine := Affine2DParameters{
		A: 1,
		B: 0,
		C: 0,
		D: 1,
		E: 0,
		F: 0,
	}

	if affine.A != 1 {
		t.Errorf("A not set correctly")
	}

	if affine.B != 0 {
		t.Errorf("B not set correctly")
	}

	if affine.C != 0 {
		t.Errorf("C not set correctly")
	}

	if affine.D != 1 {
		t.Errorf("D not set correctly")
	}

	if affine.E != 0 {
		t.Errorf("E not set correctly")
	}

	if affine.F != 0 {
		t.Errorf("F not set correctly")
	}
}

/*****************************************************************************************************************/
