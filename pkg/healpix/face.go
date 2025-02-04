/*****************************************************************************************************************/

//  @author     Michael Roberts <michael@observerly.com>
//  @package    @observerly/skysolve
//  @license    Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package healpix

/*****************************************************************************************************************/

import (
	"fmt"
)

/*****************************************************************************************************************/

const (
	BasePixelsPerRow int = 4
	BasePixelRows    int = 3
)

/*****************************************************************************************************************/

// Face represents the properties of a base pixel (or “face”) of a HEALPix map.
type Face struct {
	faceId       int          // the index of the face (0 to 11), arranged in rings around the map
	row          int          // the row in which the face resides (0, 1, or 2)
	southVertexY int          // the y coordinate of the southernmost vertex (typically 2, 3, or 4)
	southVertexX int          // the x coordinate of the southernmost vertex (typically between 0 and 7)
	neighbors    map[byte]int // a map of neighboring faces; the key packs the (x,y) offset into a byte
}

/*****************************************************************************************************************/

var faces []Face

/*****************************************************************************************************************/

// Predefined neighbor information for each face in the map:
var neighbors = [][]int{
	{8, 4, 3, 5, 0, 3, 1, 1},
	{9, 5, 0, 6, 1, 0, 2, 2},
	{10, 6, 1, 7, 2, 1, 3, 3},
	{11, 7, 2, 8, 3, 2, 0, 0},

	{11, 7, 8, 4, 3, 5, 0},
	{8, 4, 9, 5, 0, 6, 1},
	{9, 5, 10, 6, 1, 7, 2},
	{10, 6, 11, 7, 2, 4, 3},

	{11, 11, 9, 8, 4, 9, 5, 0},
	{8, 8, 10, 9, 5, 10, 6, 1},
	{9, 9, 11, 10, 6, 11, 7, 2},
	{10, 10, 8, 11, 7, 8, 4, 3},
}

/*****************************************************************************************************************/

// init precomputes all the properties for each of the 12 faces.
func init() {
	faces = make([]Face, 12)
	for i := 0; i < 12; i++ {
		row := i / BasePixelsPerRow
		col := i % BasePixelsPerRow
		faceSouthY := row + 2
		faceSouthX := 2*col - (row % 2) + 1

		faces[i] = Face{
			faceId:       i,
			row:          row,
			southVertexY: faceSouthY,
			southVertexX: faceSouthX,
			neighbors:    make(map[byte]int, 6),
		}

		// nind indexes into neighbors[i] for the current face:
		nind := 0
		// Loop over directional offsets (dx, dy) from -1 to 1 in both dimensions:
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				// Skip certain combinations according to the row:
				if row < 2 && dx == 1 && dy == 1 {
					continue
				}

				// Skip certain combinations according to the row:
				if row > 0 && dx == -1 && dy == -1 {
					continue
				}

				// Pack the directional offsets into a key, e.g., (dx+1) and (dy+1) each yield a value
				// in {0,1,2}; the dy part is shifted left by 2:
				key := byte(dx+1) | (byte(dy+1) << 2)
				faces[i].neighbors[key] = neighbors[i][nind]
				nind++
			}
		}
	}
}

/*****************************************************************************************************************/

// NewFace returns the Face with the specified face index (0 to 11):
func NewFace(index int) Face {
	return faces[index]
}

/*****************************************************************************************************************/

// Return the face id of the neighbor of the current face in the given direction. Each of x and y expects one of
// three values: -1, 0, or 1.
// In this scheme, x and y have 0 at the bottom most vertex of the face (which are arrayed as diamonds conceptually),
// and each increases outward along the diamond boundary, x along the left side and y along the right side.
// So x,y == -1,-1 is the directly southern neighbor; x,y == 1,1 is the directly northern neighbor, and x,y == -1,1
// is the northwestern neighbor.
func (f Face) GetNeighbour(xOffset int, yOffset int) (int, error) {
	xpack := byte(xOffset + 1)
	ypack := byte(yOffset+1) << 2

	index := xpack | ypack

	if n, ok := f.neighbors[byte(index)]; ok {
		return n, nil
	}

	return -1, fmt.Errorf("no neighbor in direction %d,%d", xOffset, yOffset)
}

/*****************************************************************************************************************/
