/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package quad

/*****************************************************************************************************************/

import (
	"github.com/observerly/skysolve/pkg/star"
)

/*****************************************************************************************************************/

// Quad represents a quadrilateral formed by four cartesian points in Euclidean space.
type Quad struct {
	A           star.Star  // The original value of quad point A (at 0,0)
	B           star.Star  // The original value of quad point B (at 1,1)
	C           star.Star  // The original value of quad point C (at cx, cy)
	D           star.Star  // The original value of quad point D (at dx, dy)
	NormalisedA star.Star  // The normalised value of quad point A in Euclidean space
	NormalisedB star.Star  // The normalised value of quad point B in Euclidean space
	NormalisedC star.Star  // The normalised value of quad point C in Euclidean space
	NormalisedD star.Star  // The normalised value of quad point D in Euclidean space
	Hash        [4]float64 // An exactly precise hash for the quad, representing Cx, Cy, Dx, Dy
	Precision   int        // The precision of the hash code (default is 3, which is 3 decimal places)
}

/*****************************************************************************************************************/
