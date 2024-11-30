/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package transform

/*****************************************************************************************************************/

// Affine2DParameters represents the parameters of a 2D affine transformation.
type Affine2DParameters struct {
	A, B, C float64 // Transformation for X: x' = A*x + B*y + C
	D, E, F float64 // Transformation for Y: y' = D*x + E*y + F
}

/*****************************************************************************************************************/
