/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package spatial

/*****************************************************************************************************************/

import (
	"errors"

	"github.com/observerly/skysolve/pkg/quad"
	"gonum.org/v1/gonum/spatial/vptree"
)

/*****************************************************************************************************************/

// Match holds the matched Quad and the distance between the generated Quad and the matched Quad:
type QuadMatch struct {
	Quad     quad.Quad
	Distance float64
}

/*****************************************************************************************************************/

type QuadMatcher struct {
	Tree *vptree.Tree
}

/*****************************************************************************************************************/

// NewMatcher initializes the Matcher with a list of source quads and maxUses.
func NewQuadMatcher(quads []quad.Quad) (*QuadMatcher, error) {
	// Convert []quad.Quad to []vptree.Comparable
	comparables := make([]vptree.Comparable, len(quads))

	for i, q := range quads {
		comparables[i] = q
	}

	// Initialize the VP-Tree with effort=2 (can be adjusted)
	tree, err := vptree.New(comparables, 1, nil) // effort=2, src=nil for default randomness
	if err != nil {
		return nil, err
	}

	return &QuadMatcher{
		Tree: tree,
	}, nil
}

/*****************************************************************************************************************/

// MatchQuad finds the nearest source Quad to the generated Quad within maxDistance.
// It ensures that no star in the matched Quad exceeds MaxUses.
func (m *QuadMatcher) MatchQuad(q quad.Quad, tolerance float64) (*QuadMatch, error) {
	// Query the VP-Tree for the nearest neighbor
	nearest, distance := m.Tree.Nearest(q)

	if distance > tolerance {
		return nil, errors.New("no match found within the specified distance")
	}

	matchedQuad, ok := nearest.(quad.Quad)

	if !ok {
		return nil, errors.New("matched element is not of type Quad")
	}

	// Create a copy of the matchedQuad to avoid modifying the original source Quad
	qc := matchedQuad

	// Set the corresponding equatorial coordinates for the matched Quad:
	qc.A.RA = q.A.RA
	qc.A.Dec = q.A.Dec

	qc.B.RA = q.B.RA
	qc.B.Dec = q.B.Dec

	qc.C.RA = q.C.RA
	qc.C.Dec = q.C.Dec

	qc.D.RA = q.D.RA
	qc.D.Dec = q.D.Dec

	// Ensure we set the designations of the matched Quad:
	qc.A.Designation = q.A.Designation
	qc.B.Designation = q.B.Designation
	qc.C.Designation = q.C.Designation
	qc.D.Designation = q.D.Designation

	return &QuadMatch{
		Quad:     qc,
		Distance: distance,
	}, nil
}

/*****************************************************************************************************************/

// MatchQuads finds matches for all generated quads.
// Returns a slice of Match containing successful matches.
func (m *QuadMatcher) MatchQuads(quads []quad.Quad, tolerance float64) ([]QuadMatch, error) {
	matches := []QuadMatch{}

	for _, q := range quads {
		match, err := m.MatchQuad(q, tolerance)
		if err != nil {
			// Handle quads with no matches or exceeded usage as needed, e.g., skip or log
			continue
		}

		matches = append(matches, *match)
	}

	return matches, nil
}

/*****************************************************************************************************************/
