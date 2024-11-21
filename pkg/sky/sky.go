/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package sky

/*****************************************************************************************************************/

import (
	"time"
)

/*****************************************************************************************************************/

type Params struct {
	ExposureDuration         time.Duration // exposure duration
	MaxADU                   float64       // maximum ADU value
	BiasOffset               float64       // bias offset in units of ADU
	Gain                     float64       // gain in units of e-/ADU
	ReadNoise                float64       // read noise in units of e-/pixel
	DarkCurrent              float64       // dark current in units of e-/s/pixel
	BinningX                 int           // binning factor on the x axis in units of pixels
	BinningY                 int           // binning factor on the y axis in units of pixels
	PixelSizeX               float64       // pixel size on the x axis in units of meters
	PixelSizeY               float64       // pixel size on the y axis in units of meters
	FocalLength              float64       // focal length of the telescope in units of m
	ApertureDiameter         float64       // aperture diameter of the telescope in units of m
	SkyBackground            float64       // the sky background in units of e-/m2/arcsec2/s
	Seeing                   float64       // the perceived seeing in units of arcsec
	AverageQuantumEfficiency float64       // the average quantum efficiency of the CCD
}

/*****************************************************************************************************************/
