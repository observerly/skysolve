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
	"gonum.org/v1/gonum/spatial/vptree"
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

// NewQuad creates a new Quad from four points.
func NewQuad(a, b, c, d star.Star, precision int) (Quad, error) {
	// We need to determine which is A and which is B, given our criteria, and then determine
	// which is C and which is D based on the x dimension.
	A, B, C, D := DetermineABCD(a, b, c, d)

	// Once we have determined A, B, C and D, we can normalised according to coordinate space such that
	// A is found at (0,0) and B is then found at (1,1).
	a, b, c, d, err := NormalizeToAB(A, B, C, D)

	if err != nil {
		return Quad{}, err
	}

	// Set up the quad to contain the original points and set the normalised points:
	q := Quad{
		A:           A,
		B:           B,
		C:           C,
		D:           D,
		NormalisedA: a,
		NormalisedB: b,
		NormalisedC: c,
		NormalisedD: d,
		Precision:   precision,
	}

	// Generate the hash code for the quad, once we have the normalised points:
	q.Hash = [4]float64{q.NormalisedC.X, q.NormalisedC.Y, q.NormalisedD.X, q.NormalisedD.Y}

	return q, nil
}

/*****************************************************************************************************************/

// Distance calculates the Euclidean distance between two quads based on their Hash fields.
// This method satisfies the vptree.Comparable interface.
func (q Quad) Distance(compare vptree.Comparable) float64 {
	o, ok := compare.(Quad)

	if !ok {
		panic("vptree: incompatible type for distance calculation")
	}

	// Calculate squared differences for C and D:
	dxC := q.NormalisedC.X - o.NormalisedC.X
	dyC := q.NormalisedC.Y - o.NormalisedC.Y
	dxD := q.NormalisedD.X - o.NormalisedD.X
	dyD := q.NormalisedD.Y - o.NormalisedD.Y

	// Compute Euclidean distance in 4D space
	return (math.Hypot(dxC, dyC) + math.Hypot(dxD, dyD)) / 2
}

/*****************************************************************************************************************/

func (q *Quad) EucliadianPixelCenter() (float64, float64) {
	// Get the center between the four points, A, B, C and D:
	x := (q.A.X + q.B.X + q.C.X + q.D.X) / 4
	// Get the center between the four points, A, B, C and D:
	y := (q.A.Y + q.B.Y + q.C.Y + q.D.Y) / 4
	return x, y
}

/*****************************************************************************************************************/

// GenerateHashCode generates a precise hash code for the quad based on projections.
// It mirrors the functionality of the Python `quad_hash` function.
func (q *Quad) GenerateHashCode() [4]float64 {
	return [4]float64{q.NormalisedC.X, q.NormalisedC.Y, q.NormalisedD.X, q.NormalisedD.Y}
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

	// If Cx + Dx > 1, then the quad is not symmetric (and thus not invariant under rotation):
	if c.X+d.X > 1 {
		return a, b, c, d, fmt.Errorf("quad invalid: Cx + Dx > 1, which makes the normalisation asymmetric")
	}

	// If either C or D are not within the unit circle, then the quad is invalid:
	if !IsWithinUnitCircle(c.X, c.Y) && !IsWithinUnitCircle(d.X, d.Y) {
		return a, b, c, d, fmt.Errorf("quad invalid: C or D is not within the unit circle")
	}

	return a, b, c, d, nil
}

/*****************************************************************************************************************/

// IsWithinUnitCircle checks if a star is within the unit circle centered at (0.5, 0.5):
func IsWithinUnitCircle(x float64, y float64) bool {
	centerX, centerY := 0.5, 0.5
	radius := math.Sqrt2 / 2
	dist := math.Hypot(x-centerX, y-centerY)
	return dist <= radius+1e-6 // Adding a small epsilon to account for floating-point precision
}

/*****************************************************************************************************************/
