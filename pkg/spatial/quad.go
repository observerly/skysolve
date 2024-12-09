/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package spatial

/*****************************************************************************************************************/

import (
	"github.com/observerly/skysolve/pkg/quad"
)

/*****************************************************************************************************************/

// Match holds the matched Quad and the distance between the generated Quad and the matched Quad:
type QuadMatch struct {
	Quad     quad.Quad
	Distance float64
}

/*****************************************************************************************************************/
