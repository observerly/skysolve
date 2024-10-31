/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package astrometry

/*****************************************************************************************************************/

import (
	"github.com/observerly/iris/pkg/photometry"
	"github.com/observerly/skysolve/pkg/geometry"
)

/*****************************************************************************************************************/

type ICRSEquatorialCoordinate struct {
	RA  float64
	Dec float64
}

/*****************************************************************************************************************/

type Asterism struct {
	A        photometry.Star
	B        photometry.Star
	C        photometry.Star
	Features geometry.InvariantFeatures
}

/*****************************************************************************************************************/
