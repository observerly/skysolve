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
	"time"

	"github.com/fogleman/gg"
	"github.com/observerly/iris/pkg/fits"
	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/fov"
	"github.com/observerly/skysolve/pkg/solve"
)

/*****************************************************************************************************************/

func main() {
	// Attempt to open the file from the given filepath:
	file, err := os.Open("./samples/Rosetta_Nebula_[Ha]_Monochrome_M_300s_2024-11-26T17_20_00Z.fits")
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
		return
	}

	// Defer closing the file:
	defer file.Close()

	// Assume an image of 2x2 pixels with 16-bit depth, and no offset:
	fit := fits.NewFITSImage(2, 0, 0, 65535)

	// Read in our exposure data into the image:
	err = fit.Read(file)
	if err != nil {
		fmt.Printf("failed to read file: %v", err)
		return
	}

	// Attempt to get the RA header from the FITS file:
	ra, exists := fit.Header.Floats["RA"]
	if !exists {
		fmt.Printf("ra header not found")
		return
	}

	fmt.Println("RA:", ra)

	// Attempt to get the Dec header from the FITS file:
	dec, exists := fit.Header.Floats["DEC"]
	if !exists {
		fmt.Printf("dec header not found")
		return
	}

	fmt.Println("Dec:", dec)

	pixelScaleX := 0.000540 // 0.000540 degrees per pixel in the x-axis

	pixelScaleY := 0.000540 // 0.000540 degrees per pixel in the y-axis

	radius := fov.GetRadialExtent(float64(fit.Header.Naxis1), float64(fit.Header.Naxis2), fov.PixelScale{X: pixelScaleX, Y: pixelScaleY})

	// Create a new SIMBAD service client:
	service := catalog.NewCatalogService(catalog.SIMBAD, catalog.Params{
		Limit:     100, // Limit the number of records to 48
		Threshold: 16,  // Limiting Magntiude, filter out any stars that are magnitude 16 or above (fainter)
	})

	fmt.Println("Radius:", radius)

	// Attempt to create a new PlateSolver:
	solver, err := solve.NewPlateSolver(solve.Params{
		Data:                fit.Data,               // The exposure data from the fits image
		Width:               int(fit.Header.Naxis1), // The width of the image
		Height:              int(fit.Header.Naxis2), // The height of the image
		PixelScaleX:         pixelScaleX,            // The pixel scale in the x-axis
		PixelScaleY:         pixelScaleY,            // The pixel scale in the y-axis
		ADU:                 fit.ADU,                // The analog-to-digital unit of the image
		ExtractionThreshold: 16,                     // Extract a minimum of 20 of the brightest stars
		Radius:              16,                     // 16 pixels radius for the star extraction
		Sigma:               2.5,                    // 8 pixels sigma for the Gaussian kernel
	})

	if err != nil {
		fmt.Printf("there was an error while creating the plate solver: %v", err)
		return
	}

	// Define the tolerances for the solver, we can adjust these as needed:
	tolerance := solve.ToleranceParams{
		QuadTolerance:           0.02,
		EuclidianPixelTolerance: 10,
	}

	// Whilst we have no matches, and whilst we are within 1 degree of the initial { ra, dec } guess, keep solving:
	eq := astrometry.ICRSEquatorialCoordinate{
		RA:  98.6,
		Dec: 2.5,
	}

	// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
	sources, err := service.PerformRadialSearch(eq, radius)
	if err != nil {
		fmt.Printf("there was an error while performing the SIMBAD radial search: %v", err)
		return
	}

	// Record the start time
	startTime := time.Now()

	solver.Sources = append(solver.Sources, sources...)

	fmt.Println("Number of Sources:", len(solver.Sources))

	wcs, matches, err := solver.Solve(tolerance, 3)

	fmt.Println("Number of Matches:", len(matches))

	// Calculate the elapsed time
	elapsedTime := time.Since(startTime)

	fmt.Println(elapsedTime)

	if err != nil {
		fmt.Println("an error occured while plate solving:", err)
	}

	if wcs == nil {
		fmt.Println("no WCS solution found")
	}

	if wcs != nil {
		fmt.Println(wcs)

		eq = wcs.PixelToEquatorialCoordinate(float64(solver.Width)/2, float64(solver.Height)/2)

		fmt.Println(float64(solver.Width)/2, float64(solver.Height)/2)

		fmt.Println(
			"Expected RA 98.2746 (deg), got:", eq.RA,
			"Expected Dec 2.637267 (deg), got:", eq.Dec,
		)

		eq = wcs.PixelToEquatorialCoordinate(578.234176636, 485.578643799)

		fmt.Println(
			"Expected RA 98.647 (deg), got:", eq.RA,
			"Expected Dec 2.537 (deg), got:", eq.Dec,
		)

		eq = wcs.PixelToEquatorialCoordinate(0, 0)

		fmt.Println("Top Left Corner:", eq)

		eq = wcs.PixelToEquatorialCoordinate(float64(solver.Width), float64(solver.Height))

		fmt.Println("Bottom Right Corner:", eq)
	}

	height := int(solver.Height)

	width := int(solver.Width)

	// // 2. Validate Image Dimensions
	if width <= 0 || height <= 0 {
		fmt.Println("invalid image dimensions: width and height must be greater than zero")
		return
	}

	expectedDataLength := width * height
	if len(solver.Data) < expectedDataLength {
		fmt.Printf("invalid data length: expected at least %d pixels, got %d\n", expectedDataLength, len(solver.Data))
		return
	}

	// a. Find min and max values for normalization
	minVal, maxVal := solver.Data[0], solver.Data[0]
	for _, pixel := range solver.Data {
		if pixel < minVal {
			minVal = pixel
		}
		if pixel > maxVal {
			maxVal = pixel
		}
	}
	if maxVal == minVal {
		maxVal = minVal + 1 // Prevent division by zero
	}

	// Debugging Statements
	fmt.Printf("Pixel Value Range: min=%f, max=%f\n", minVal, maxVal)

	// Create a new grayscale image
	imgGray := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*width + x
			if index >= len(solver.Data) {
				// This should not happen due to earlier validation
				break
			}
			normalized := (solver.Data[index] - minVal) / (maxVal - minVal)
			if math.IsNaN(float64(normalized)) || math.IsInf(float64(normalized), 0) {
				normalized = 0
			}
			gray := uint8(math.Round(float64(normalized) * 255))
			imgGray.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	fmt.Println("Image Dimensions:", width, height)

	fmt.Println("Initializing Drawing Context...")

	// // c. Initialize the drawing context
	dc := gg.NewContext(width, height)

	fmt.Println("Drawing Image...")

	// d. Draw the grayscale image onto the context
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := imgGray.GrayAt(x, y).Y
			dc.SetRGB(float64(gray)/255, float64(gray)/255, float64(gray)/255)
			dc.SetPixel(x, y)
		}
	}

	// fmt.Println("Drawing Stars...")

	for _, star := range solver.Stars {
		dc.SetColor(color.RGBA{R: 241, G: 245, B: 249, A: 255})
		dc.DrawCircle(float64(star.X), float64(star.Y), 16.0)
		dc.SetLineWidth(2)
		dc.Stroke()
	}

	for _, match := range matches {
		dc.SetColor(color.RGBA{R: 129, G: 140, B: 248, A: 255})
		dc.DrawCircle(float64(match.Quad.A.X), float64(match.Quad.A.Y), 20.0)
		dc.DrawCircle(float64(match.Quad.B.X), float64(match.Quad.B.Y), 20.0)
		dc.DrawCircle(float64(match.Quad.C.X), float64(match.Quad.C.Y), 20.0)
		dc.DrawCircle(float64(match.Quad.D.X), float64(match.Quad.D.Y), 20.0)
		dc.SetLineWidth(2)
		dc.Stroke()

		dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		dc.DrawString(match.Quad.A.Designation, float64(match.Quad.A.X), float64(match.Quad.A.Y)-30)
		dc.DrawString(match.Quad.B.Designation, float64(match.Quad.B.X), float64(match.Quad.B.Y)-30)
		dc.DrawString(match.Quad.C.Designation, float64(match.Quad.C.X), float64(match.Quad.C.Y)-30)
		dc.DrawString(match.Quad.D.Designation, float64(match.Quad.D.X), float64(match.Quad.D.Y)-30)

		dc.SetLineWidth(10)
		dc.Stroke()
	}

	fmt.Println("Saving Image...")

	// Save the annotated image as PNG:
	outputPath := "samples/Rosetta_Nebula_[Ha]_Monochrome_M_300s_2024-11-26T17_20_00Z.png"

	outputFile, err := os.Create(outputPath)

	if err != nil {
		fmt.Println("error creating output file:", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, dc.Image())
	if err != nil {
		fmt.Println("error encoding PNG:", err)
		return
	}

	fmt.Printf("Annotated image saved as %s\n", outputPath)
}

/*****************************************************************************************************************/
