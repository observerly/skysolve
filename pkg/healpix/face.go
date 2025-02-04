/*****************************************************************************************************************/

//  @author     Michael Roberts <michael@observerly.com>
//  @package    @observerly/skysolve
//  @license    Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package healpix

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
