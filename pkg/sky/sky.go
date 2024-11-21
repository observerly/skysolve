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
	"math/rand"
	"time"

	"github.com/observerly/skysolve/pkg/astrometry"
	stats "github.com/observerly/skysolve/pkg/statistics"
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

//lint:ignore U1000 Reserved for future implementation.
func (s *SimulatedSkyImage) normaliseFieldImage(data []float64, width, height int) ([][]uint32, error) {
	// Initialize the 2D slice with the desired height.
	image := make([][]uint32, height)

	// Create a flat slice to hold all processed pixel values.
	flatImage := make([]uint32, height*width)

	// Single loop to process all pixels.
	for index := 0; index < height*width; index++ {
		// Assign the inner slice when starting a new row.
		if index%width == 0 {
			y := index / width
			image[y] = flatImage[index : index+width]
		}

		// Apply gain and bias offset.
		value := data[index]/s.Gain + s.BiasOffset

		// Clamp the value to the valid ADU range.
		if value < 0.0 {
			value = 0.0
		}
		if value > s.MaxADU {
			value = s.MaxADU
		}

		// Store the rounded value as uint32 in the flat slice.
		flatImage[index] = uint32(math.Round(value))
	}

	return image, nil
}

/*****************************************************************************************************************/

func (s *SimulatedSkyImage) GenerateBackgroundImage() ([]float64, error) {
	// Calculate the aperture area in m²:
	apertureArea := math.Pi * math.Pow(s.ApertureDiameter/2.0, 2)

	// Calculate sky background per pixel in e⁻/s/pixel:
	skyBackgroundPerPixel := s.SkyBackground * apertureArea * s.PixelScaleX * s.PixelScaleY * 3600.0 * 3600.0

	// Initialize a flat base image with zeros for the entire field:
	image := make([]float64, s.Width*s.Height)

	// Precompute background noise, which is the sum of dark current, read noise, and sky background:
	background := stats.PoissonDistributedRandomNumber(s.DarkCurrent*s.ExposureDuration) +
		stats.NormalDistributedRandomNumber(0.0, s.ReadNoise) +
		stats.PoissonDistributedRandomNumber(skyBackgroundPerPixel*s.ExposureDuration)

	// Add background to the entire image with some random noise:
	for i := range image {
		image[i] += background * rand.Float64()
	}

	return image, nil
}

/*****************************************************************************************************************/
