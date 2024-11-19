/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package transform

/*****************************************************************************************************************/

import "testing"

/*****************************************************************************************************************/

func TestSIP2DForwardParameters(t *testing.T) {
	sip := SIP2DForwardParameters{
		AOrder: 1,
		BOrder: 1,
		APower: map[string]float64{
			"0_0": 1,
			"1_0": 0,
			"0_1": 0,
		},
		BPower: map[string]float64{
			"0_0": 1,
			"1_0": 0,
			"0_1": 0,
		},
	}

	if sip.AOrder != 1 {
		t.Errorf("AOrder not set correctly")
	}

	if sip.BOrder != 1 {
		t.Errorf("BOrder not set correctly")
	}

	if sip.APower["0_0"] != 1 {
		t.Errorf("APower[0_0] not set correctly")
	}

	if sip.APower["1_0"] != 0 {
		t.Errorf("APower[1_0] not set correctly")
	}

	if sip.APower["0_1"] != 0 {
		t.Errorf("APower[0_1] not set correctly")
	}

	if sip.BPower["0_0"] != 1 {
		t.Errorf("BPower[0_0] not set correctly")
	}

	if sip.BPower["1_0"] != 0 {
		t.Errorf("BPower[1_0] not set correctly")
	}

	if sip.BPower["0_1"] != 0 {
		t.Errorf("BPower[0_1] not set correctly")
	}
}

/*****************************************************************************************************************/
