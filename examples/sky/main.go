/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package main

/*****************************************************************************************************************/

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"time"

	stats "github.com/observerly/iris/pkg/statistics"
	"github.com/observerly/iris/pkg/utils"
	"github.com/observerly/sidera/pkg/humanize"
	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/sky"
)

/*****************************************************************************************************************/

func ZScaleNormalizeImage(data [][]uint32, median, stdDev float64, scaleFactor float64) (*image.Gray16, error) {
	// Ensure the image has at least one column:
	height := len(data)

	// Ensure the image has at least one row:
	width := len(data[0])

	// Create a new 16-bit grayscale image using the given dimensions:
	img := image.NewGray16(image.Rect(0, 0, width, height))

	// Clamp vmin to [0, 65535]:
	vmin := math.Max(0, median-stdDev*scaleFactor)

	// Clamp vmax to [0, 65535]:
	vmax := math.Min(65535, median+stdDev*scaleFactor)

	// Prevent division by zero issues:
	if vmax == vmin {
		vmax = vmin + 1.0
	}

	// Normalize the image data:
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := float64(data[y][x])

			// Clamp the value to [vmin, vmax]:
			if value < vmin {
				value = vmin
			} else if value > vmax {
				value = vmax
			}

			// Normalize to [0, 65535] range:
			normalized := uint16(65535.0 * (value - vmin) / (vmax - vmin))

			// Set pixel value in the image:
			img.SetGray16(x, y, color.Gray16{Y: normalized})
		}
	}

	return img, nil
}

/*****************************************************************************************************************/

func main() {
	// Define the x dimension of the image:
	xs := 2048

	// Define the y dimension of the image:
	ys := 2048

	// Define the center of the field of view as the Pleiades:
	eq := astrometry.ICRSEquatorialCoordinate{
		RA:  56.75101,
		Dec: 24.11678,
	}

	// 16-bit CCD sensor with a maximum ADU of 65535:
	adu := 65535.0

	// Physical pixel size in meters:
	pixelSize := 0.00054

	// Focal Length in meters:
	focalLength := 1.2

	fovX := 2 * math.Atan((float64(xs)*pixelSize)/(2*focalLength)) * (180 / math.Pi)

	fovY := 2 * math.Atan((float64(ys)*pixelSize)/(2*focalLength)) * (180 / math.Pi)

	r := math.Hypot(fovX, fovY) / 2

	params := sky.Params{
		ExposureDuration:         300 * time.Second, // 30 second exposure, suitable for faint object imaging. Adjust as needed for brighter or dimmer targets.
		MaxADU:                   adu,               // Common for 16-bit CCD sensors. Ensure it aligns with your sensor’s specifications.
		BiasOffset:               300.0,             // Typical bias levels range between 500 to 1500 ADU. Adjust based on your sensor’s characteristics.
		Gain:                     0.5,               // Typical CCD gains range from 1.0 to 2.0 e⁻/ADU. A gain of 1.2 is a reasonable default.
		ReadNoise:                1.2,               // Common read noise values are between 3 to 10 e⁻. A value of 5 e⁻ is typical for many CCD sensors.
		DarkCurrent:              0.2,               // Dark current varies based on sensor temperature; cooled sensors have lower dark currents.
		BinningX:                 1,                 // (1x1) preserves the highest possible resolution. Increase binning factors (e.g., 2x2) to improve signal-to-noise at the expense of resolution.
		BinningY:                 1,                 // (1x1) preserves the highest possible resolution. Increase binning factors (e.g., 2x2) to improve signal-to-noise at the expense of resolution.
		PixelSizeX:               pixelSize,         // Common pixel sizes for CCDs range from 4 µm to 15 µm. A pixel size of 5 µm is a typical value balancing resolution and sensitivity.
		PixelSizeY:               pixelSize,         // Common pixel sizes for CCDs range from 4 µm to 15 µm. A pixel size of 5 µm is a typical value balancing resolution and sensitivity.
		FocalLength:              focalLength,       // Adjust based on the desired field of view and image scale. An increase in the focal length will result in a narrower field of view.
		ApertureDiameter:         0.417,             // Adjust according to the telescope size you wish to simulate. An increase in aperture diameter will result in a brighter image.
		SkyBackground:            50.0,              // Sky background can vary widely based on location and observing conditions (e.g., 0 for a perfectly dark sky, up to ~1000 e⁻/m²/arcsec²/s).
		Seeing:                   1.5,               // Typical seeing ranges from 0.5 to 2.0 arcseconds. Adjust based on the observing site’s atmospheric conditions.
		AverageQuantumEfficiency: 0.93,              // High quantum efficiency (90%) is excellent. Typical QE for CCDs ranges from 60% to 90%. Adjust based on the sensor’s specifications.
	}

	// Field of View
	fmt.Printf("Field of View (FOV): %.5f° x %.5f° (Radius: ~%.4f°)\n", fovX, fovY, r)

	// Center of View
	fmt.Printf("Center of View: RA: %.5f°, Dec: %.5f°\n", eq.RA, eq.Dec)

	// Image Dimensions
	fmt.Printf("Image Dimensions: %d x %d pixels\n", xs, ys)

	// Exposure Details
	fmt.Printf("Exposure Duration: %.0f seconds\n", params.ExposureDuration.Seconds())
	fmt.Printf("Maximum ADU: %.0f\n", params.MaxADU)
	fmt.Printf("Bias Offset: %.1f ADU\n", params.BiasOffset)
	fmt.Printf("Gain: %.1f e⁻/ADU\n", params.Gain)
	fmt.Printf("Read Noise: %.1f e⁻\n", params.ReadNoise)
	fmt.Printf("Dark Current: %.1f e⁻/s\n", params.DarkCurrent)

	// Sensor Configuration
	fmt.Printf("Pixel Size: %.2e meters (X & Y)\n", params.PixelSizeX)
	fmt.Printf("Binning: %dx%d\n", params.BinningX, params.BinningY)
	fmt.Printf("Focal Length: %.2f meters\n", params.FocalLength)
	fmt.Printf("Aperture Diameter: %.1f meters\n", params.ApertureDiameter)

	// Observing Conditions
	fmt.Printf("Sky Background: %.1f e⁻/m²/arcsec²/s\n", params.SkyBackground)
	fmt.Printf("Seeing: %.1f arcseconds\n", params.Seeing)
	fmt.Printf("Average Quantum Efficiency: %.2f%%\n", params.AverageQuantumEfficiency*100)

	fmt.Println("\nGenerating simulated sky image with the above parameters...")

	// Make sure we output the image to the .output directory in the workspace root directory:
	location := "./.output"

	// Ensure the output directory exists before attempting to write the image:
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", location, err)
		return
	}

	filename := fmt.Sprintf(
		"J2000_%v%v_%d_%d.png",
		humanize.FormatDecimalToDMS(eq.RA, "%s%d%d%.2f"),
		humanize.FormatDecimalToDMS(eq.Dec, "%s%d%d%.2f"),
		xs,
		ys,
	)

	path := filepath.Join(location, filename)

	// Create a new simulated sky parameter:
	sky, err := sky.NewSimulatedSky(xs, ys, eq, params)

	if err != nil {
		fmt.Printf("Error creating simulated sky: %v", err)
		panic(err)
	}

	// Create a new GAIA service client:
	q := catalog.NewGAIAServiceClient()

	// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
	sources, err := q.PerformRadialSearch(eq, math.Ceil(r*10)/10, 1000, 13)

	if err != nil {
		fmt.Printf("Error performing radial search: %v", err)
		panic(err)
	}

	fmt.Println("Generating image...")

	// Record the start time
	startTime := time.Now()

	// Generate a new simulated sky image:
	image, err := sky.GenerateFieldImage(sources)

	if err != nil {
		fmt.Printf("Error generating image: %v", err)
		panic(err)
	}

	s := stats.NewStats(utils.Flatten2DUInt32Array(image), int32(adu), xs)

	median := s.FastMedian()

	x, y := sky.WCS.EquatorialCoordinateToPixel(eq.RA, eq.Dec)

	fmt.Println(x, y)

	// Create a normalised image using the median and standard deviation:
	img, err := ZScaleNormalizeImage(image, float64(median), float64(s.StdDev), 1.2)
	if err != nil {
		fmt.Println("Error creating image:", err)
		return
	}

	// Calculate the elapsed time
	elapsedTime := time.Since(startTime)

	fmt.Println("Elapsed time to generate image:", elapsedTime)

	// Save the image to a file on disk:
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file: %v", err)
		panic(err)
	}
	defer file.Close()

	// Encode the image as a PNG:
	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("Error encoding image: %v", err)
		panic(err)
	}

	fmt.Printf("Image saved as '%s'\n", filename)
}

/*****************************************************************************************************************/
