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
