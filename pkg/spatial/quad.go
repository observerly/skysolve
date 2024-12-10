/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package spatial

/*****************************************************************************************************************/

import (
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
