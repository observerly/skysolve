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
	"github.com/observerly/skysolve/pkg/catalog"
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
			A: -pixelScaleX,
			B: 0,
			C: eq.RA,
			D: 0,
			E: pixelScaleY,
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

// generateMoffatProfile generates a flattened Moffat profile and returns it along with its bounds and dimensions.
// This function is independent of the SimulatedSkyImage struct.
func generateMoffatProfile(
	x0, y0 float64,
	xMin, xMax, yMin, yMax int, flux float64,
	beta float64, precisionX, precisionY float64,
) ([]float64, int, int) {
	width := xMax - xMin + 1
	height := yMax - yMin + 1
	totalPixels := width * height
	profile := make([]float64, totalPixels)

	totalIntensity := 0.0

	// Single loop iterating over all pixels in the grid:
	for idx := 0; idx < totalPixels; idx++ {
		// Map idx to y and x coordinates:
		yIdx := idx / width
		xIdx := idx % width

		// Calculate the pixel coordinates:
		y := yMin + yIdx
		x := xMin + xIdx

		// Calculate the distance from the center of the profile:
		dy := float64(y) - y0 + 0.5
		dx := float64(x) - x0 + 0.5

		// Compute r using precomputed inverses:
		r := (dx*dx)*precisionX + (dy*dy)*precisionY

		// Calculate the intensity using the Moffat profile:
		intensity := math.Exp(-beta * math.Log(1.0+r))

		// Assign intensity to profile and accumulate sum of intensities:
		profile[idx] = intensity

		// Accumulate total intensity:
		totalIntensity += intensity
	}

	// Normalize and scale by total flux:
	scaleFactor := flux / totalIntensity
	for i := 0; i < totalPixels; i++ {
		profile[i] *= scaleFactor
	}

	return profile, width, height
}

/*****************************************************************************************************************/

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

// GenerateFieldImage generates the star field image by placing sources onto the image data.
func (s *SimulatedSkyImage) GenerateFieldImage(sources []catalog.Source) ([][]uint32, error) {
	// Initialize a flat base image with zeros for the entire field:
	image, err := s.GenerateBackgroundImage()

	// If we failed to generate the background image, return the error:
	if err != nil {
		return nil, err
	}

	// Calculate the aperture area in m²:
	apertureArea := math.Pi * math.Pow(s.ApertureDiameter/2.0, 2)

	// Precompute flux density factor for sources:
	fluxDensity := s.AverageQuantumEfficiency * s.ExposureDuration * apertureArea

	// Calculate FWHM in pixels in the x dimension for the Moffat profile:
	fwhmPixelsX := s.Seeing / (s.PixelScaleX * 3600.0)

	// Calculate FWHM in pixels in the y dimension for the Moffat profile:
	fwhmPixelsY := s.Seeing / (s.PixelScaleY * 3600.0)

	// Calculate the precesion factor in the x dimension for the Moffat profile:
	precisionX := math.Pow(fwhmPixelsX, -2)

	// Calculate the precesion factor in the y dimension for the Moffat profile:
	precisionY := math.Pow(fwhmPixelsY, -2)

	// Calculate beta for the Moffat profile using the flux of the source, such that brighter sources have a smaller beta:
	beta := 3.0

	// Add stars to the image
	for _, source := range sources {
		// Calculate flux in e⁻ for this source:
		e := source.PhotometricGMeanFlux * fluxDensity * math.Pow(10, -0.4*source.PhotometricGMeanMagnitude)

		scale := float64(((s.Width + s.Height) / 2)) * math.Pow(10, -0.2*source.PhotometricGMeanMagnitude)

		// Calculate render radius in the x dimension for the Moffat profile:
		renderRadiusX := fwhmPixelsX * scale

		// Calculate render radius in the y dimension for the Moffat profile:
		renderRadiusY := fwhmPixelsY * scale

		// Transform source RA and Dec to pixel coordinates
		x0, y0 := s.WCS.EquatorialCoordinateToPixel(source.RA, source.Dec)

		// Skip sources outside the image frame (RA, Dec out of bounds):
		if x0 < 0 || x0 >= float64(s.Width) || y0 < 0 || y0 >= float64(s.Height) {
			continue
		}

		// Determine the grid bounds for the individual source:
		xMin := int(math.Max(0, x0-renderRadiusX))
		xMax := int(math.Min(float64(s.Width-1), x0+renderRadiusX))
		yMin := int(math.Max(0, y0-renderRadiusY))
		yMax := int(math.Min(float64(s.Height-1), y0+renderRadiusY))

		// Generate Moffat profile for the source:
		profile, width, height := generateMoffatProfile(
			x0, y0, xMin, xMax, yMin, yMax, e, beta, precisionX, precisionY,
		)

		// Add profile to the image with clipping:
		for idx := 0; idx < width*height; idx++ {
			yIdx := idx / width
			xIdx := idx % width

			imageY := yIdx + yMin
			imageX := xIdx + xMin

			if imageY >= 0 && imageY < s.Height && imageX >= 0 && imageX < s.Width {
				image[imageY*s.Width+imageX] += profile[idx]
			}
		}
	}

	// Normalize and convert to 2D image:
	return s.normaliseFieldImage(image, s.Width, s.Height)
}

/*****************************************************************************************************************/
