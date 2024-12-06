/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package quad

/*****************************************************************************************************************/

import (
	"fmt"
	"math"

	"github.com/observerly/skysolve/pkg/geometry"
	"github.com/observerly/skysolve/pkg/star"
)

/*****************************************************************************************************************/

var NORMALISATION_ANGLE = math.Pi / 4

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

// Match represents a correspondence between an extracted star (from our image) and a
// catalog source star in pixel coordinates. This could be a false positive match, but
// we aim to minimize these by applying statistical methods for ruling out near or false
// matches.
type QuadMatch struct {
	StarQaud   Quad
	SourceQuad Quad
}

/*****************************************************************************************************************/

// DetermineAB determines which points are A and B based on the criteria that A and B are the two points
// with the largest distance between all of the points in the quad.
// C is then the point that is closest to A in the x dimension, e.g., Cx < Dx.
func DetermineABCD(a, b, c, d star.Star) (star.Star, star.Star, star.Star, star.Star) {
	stars := []star.Star{a, b, c, d}
	maximum := -1.0
	var A, B star.Star

	// Find the pair with the maximum distance in the quad:
	for i := 0; i < len(stars); i++ {
		for j := i + 1; j < len(stars); j++ {
			distance := geometry.DistanceBetweenTwoCartesianPoints(stars[i].X, stars[i].Y, stars[j].X, stars[j].Y)

			if distance > maximum {
				maximum = distance
				// Assign A and B based on X coordinate such that Ax < Bx:
				if stars[i].X < stars[j].X {
					A, B = stars[i], stars[j]
				} else {
					A, B = stars[j], stars[i]
				}
			}
		}
	}

	// Assign C and D as the remaining stars in the quad:
	var remaining []star.Star
	for _, s := range stars {
		if s != A && s != B {
			remaining = append(remaining, s)
		}
	}

	// Ensure Cx < Dx by swapping if necessary:
	if remaining[0].X < remaining[1].X {
		return A, B, remaining[0], remaining[1]
	}

	return A, B, remaining[1], remaining[0]
}

/*****************************************************************************************************************/

// NormalizeToAB normalizes the Quad such that point A maps to (0,0) and point B maps to (1,1).
func NormalizeToAB(a, b, c, d star.Star) (star.Star, star.Star, star.Star, star.Star, error) {
	Ax, Ay := 0.0, 0.0
	Bx, By := b.X-a.X, b.Y-a.Y
	Cx, Cy := c.X-a.X, c.Y-a.Y
	Dx, Dy := d.X-a.X, d.Y-a.Y

	// Step 2: Calculate the rotation angle to align A->B with y=x
	rotationAngle := NORMALISATION_ANGLE - math.Atan2(By, Bx)

	cosA := math.Cos(rotationAngle)
	sinA := math.Sin(rotationAngle)

	rAx, rAy := Ax*cosA-Ay*sinA, Ax*sinA+Ay*cosA
	rBx, rBy := Bx*cosA-By*sinA, Bx*sinA+By*cosA
	rCx, rCy := Cx*cosA-Cy*sinA, Cx*sinA+Cy*cosA
	rDx, rDy := Dx*cosA-Dy*sinA, Dx*sinA+Dy*cosA

	// Step 4: Calculate scale based on rotated B.x (which equals rotated B.y)
	scale := rBx // Since after rotation, rBx == rBy

	// Prevent division by zero
	if scale == 0 {
		scale = 1
	}

	a.X = rAx / scale
	a.Y = rAy / scale

	b.X = rBx / scale
	b.Y = rBy / scale

	c.X = rCx / scale
	c.Y = rCy / scale

	d.X = rDx / scale
	d.Y = rDy / scale

	if c.X+d.X > 1 {
		return a, b, c, d, fmt.Errorf("quad invalid: Cx + Dx > 1, which breaks normalisation symmetry")
	}

	return a, b, c, d, nil
}

/*****************************************************************************************************************/
