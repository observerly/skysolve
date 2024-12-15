/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package solve

/*****************************************************************************************************************/

import (
	"github.com/observerly/iris/pkg/photometry"

	"github.com/observerly/skysolve/pkg/catalog"
)

/*****************************************************************************************************************/

type PlateSolver struct {
	Stars       []photometry.Star
	Sources     []catalog.Source
	Data        []float32
	RA          float64
	Dec         float64
	Width       int
	Height      int
	PixelScaleX float64
	PixelScaleY float64
}

/*****************************************************************************************************************/

type Params struct {
	Data                []float32
	Width               int
	Height              int
	PixelScaleX         float64
	PixelScaleY         float64
	ADU                 int32
	ExtractionThreshold float64
	Radius              float64
	Sigma               float64
}

/*****************************************************************************************************************/
