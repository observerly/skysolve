/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package wcs

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
}

/*****************************************************************************************************************/

func NewWorldCoordinateSystem(wcs WCS) WCS {
	return wcs
}

/*****************************************************************************************************************/
