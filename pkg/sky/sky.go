/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package sky

/*****************************************************************************************************************/

import (
	"fmt"
	"math"
	"time"

	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/transform"
	"github.com/observerly/skysolve/pkg/wcs"
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

type SimulatedSkyImage struct {
	RA                       float64
	Dec                      float64
	WCS                      wcs.WCS
	Width                    int
	Height                   int
	ExposureDuration         float64
	MaxADU                   float64
	BiasOffset               float64
	Gain                     float64
	ReadNoise                float64
	DarkCurrent              float64
	BinningX                 int
	BinningY                 int
	PixelSizeX               float64
	PixelSizeY               float64
	PixelScaleX              float64
	PixelScaleY              float64
	FocalLength              float64
	ApertureDiameter         float64
	SkyBackground            float64
	Seeing                   float64
	AverageQuantumEfficiency float64
}

/*****************************************************************************************************************/

func NewSimulatedSky(xs int, ys int, eq astrometry.ICRSEquatorialCoordinate, params Params) (*SimulatedSkyImage, error) {
	// Check that the image dimensions are positive (realistic):
	if xs <= 0 || ys <= 0 {
		return nil, fmt.Errorf("image dimensions must be positive")
	}

	// Check that the pixel sizes are positive:
	if params.PixelSizeX <= 0 || params.PixelSizeY <= 0 {
		return nil, fmt.Errorf("pixel sizes must be positive")
	}

	// Check that we have a realistic seeing value:
	if params.Seeing <= 0 {
		return nil, fmt.Errorf("seeing (FWHM) must be positive")
	}

	// Check that the exposure time is positive:
	if params.ExposureDuration <= time.Duration(0) {
		return nil, fmt.Errorf("exposure time must be positive")
	}

	// Calculate pixel scale in degrees per pixel from the pixel size and focal length:
	pixelScaleX := (params.PixelSizeX / params.FocalLength) * (180 / math.Pi)

	// Calculate pixel scale in degrees per pixel from the pixel size and focal length:
	pixelScaleY := (params.PixelSizeY / params.FocalLength) * (180 / math.Pi)

	wcsParams := wcs.WCSParams{
		Projection: wcs.RADEC_TAN,
		AffineParams: transform.Affine2DParameters{
			A: pixelScaleX,
			B: 0,
			C: 0,
			D: -pixelScaleY,
			E: eq.RA,
			F: eq.Dec,
		},
	}

	// Create a new WCS object, centered at the center of the image:
	wcs := wcs.NewWorldCoordinateSystem(float64(xs)/2, float64(ys)/2, wcsParams)

	// Return a new SimulatedSkyImage object:
	return &SimulatedSkyImage{
		RA:                       eq.RA,
		Dec:                      eq.Dec,
		WCS:                      wcs,
		Width:                    xs,
		Height:                   ys,
		ExposureDuration:         params.ExposureDuration.Seconds(),
		MaxADU:                   params.MaxADU,
		BiasOffset:               params.BiasOffset,
		Gain:                     params.Gain,
		ReadNoise:                params.ReadNoise,
		DarkCurrent:              params.DarkCurrent,
		BinningX:                 params.BinningX,
		BinningY:                 params.BinningY,
		PixelScaleX:              pixelScaleX,
		PixelScaleY:              pixelScaleY,
		FocalLength:              params.FocalLength,
		ApertureDiameter:         params.ApertureDiameter,
		SkyBackground:            params.SkyBackground,
		Seeing:                   params.Seeing,
		AverageQuantumEfficiency: params.AverageQuantumEfficiency,
	}, nil
}

/*****************************************************************************************************************/