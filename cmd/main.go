/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/observerly/iris/pkg/fits"
	"github.com/observerly/skysolve/pkg/geometry"
	"github.com/observerly/skysolve/pkg/solver"
)

/*****************************************************************************************************************/

func main() {
	// Attempt to open the file from the given filepath:
	file, err := os.Open("./samples/noise16.fits")

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

	// Attempt to get the Dec header from the FITS file:
	dec, exists := fit.Header.Floats["DEC"]
	if !exists {
		fmt.Printf("dec header not found")
		return
	}

	// Attempt to create a new PlateSolver:
	solver, err := solver.NewPlateSolver(solver.GAIA, fit, solver.Params{
		RA:                  float64(ra.Value),  // The appoximate RA of the center of the image
		Dec:                 float64(dec.Value), // The appoximate Dec of the center of the image
		PixelScale:          2.061 / 3600.0,     // 2.061 arcseconds per pixel (0.0005725 degrees)
		ExtractionThreshold: 80,                 // Extract a minimum of 80 of the brightest stars
		Radius:              16,                 // 16 pixels radius for the star extraction
		Sigma:               8,                  // 8 pixels sigma for the Gaussian kernel
	})

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	// Define the tolerances for the solver, we can adjust these as needed:
	tolerances := geometry.InvariantFeatureTolerance{
		LengthRatio: 0.025, // 5% tolerance in side length ratios
		Angle:       0.5,   // 1 degree tolerance in angles
	}

	// Record the start time
	startTime := time.Now()
	wcs, err := solver.Solve(tolerances)

	if err != nil {
		fmt.Printf("an error occured while plate solving: %v", err)
		return
	}

	// Calculate the elapsed time
	elapsedTime := time.Since(startTime)

	fmt.Println(elapsedTime)

	// Calculate the reference equatorial coordinate:
	eq := wcs.PixelToEquatorialCoordinate(578.231147766, 485.620500565)

	fmt.Println(eq)
}

/*****************************************************************************************************************/
