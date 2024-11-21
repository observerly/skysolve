/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package solver

/*****************************************************************************************************************/

import (
	"math"
	"os"
	"testing"
	"time"

	"github.com/observerly/iris/pkg/fits"
	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/geometry"
)

/*****************************************************************************************************************/

var sources []catalog.Source

/*****************************************************************************************************************/

func GetCatalogSourcesForTesting(eq astrometry.ICRSEquatorialCoordinate, radius float64, limit int, threshold float64) ([]catalog.Source, error) {
	// If there are already sources, don't fetch them again:
	if len(sources) > 0 {
		return sources, nil
	}

	// Create a new GAIA service client:
	q := catalog.NewGAIAServiceClient()

	// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
	sources, err := q.PerformRadialSearch(eq, radius, limit, threshold)

	if err != nil {
		return nil, err
	}

	return sources, nil
}

/*****************************************************************************************************************/

func TestSolverOnMatches(t *testing.T) {
	// Attempt to open the file from the given filepath:
	file, err := os.Open("../../samples/noise16.fits")

	if err != nil {
		t.Errorf("NewFITSImageFromReader() os.Open(): %v", err)
	}

	// Defer closing the file:
	defer file.Close()

	// Assume an image of 2x2 pixels with 16-bit depth, and no offset:
	fit := fits.NewFITSImage(2, 0, 0, 65535)

	// Read in our exposure data into the image:
	err = fit.Read(file)

	if err != nil {
		t.Errorf("Read() error: %v", err)
	}

	// Attempt to get the RA header from the FITS file:
	ra, exists := fit.Header.Floats["RA"]
	if !exists {
		t.Errorf("ra header not found")
	}

	// Attempt to get the Dec header from the FITS file:
	dec, exists := fit.Header.Floats["DEC"]
	if !exists {
		t.Errorf("dec header not found")
	}

	eq := astrometry.ICRSEquatorialCoordinate{
		RA:  float64(ra.Value),
		Dec: float64(dec.Value),
	}

	radius := 2.0

	// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
	sources, err := GetCatalogSourcesForTesting(eq, radius, 100, 8)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Attempt to create a new PlateSolver:
	solver, err := NewPlateSolver(fit, Params{
		RA:                  float64(ra.Value),  // The appoximate RA of the center of the image
		Dec:                 float64(dec.Value), // The appoximate Dec of the center of the image
		PixelScale:          2.061 / 3600.0,     // 2.061 arcseconds per pixel (0.0005725 degrees)
		ExtractionThreshold: 50,                 // Extract a minimum of 80 of the brightest stars
		Radius:              16,                 // 16 pixels radius for the star extraction
		Sigma:               8,                  // 8 pixels sigma for the Gaussian kernel
		Sources:             sources,
	})

	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Define the tolerances for the solver, we can adjust these as needed:
	tolerances := geometry.InvariantFeatureTolerance{
		LengthRatio: 0.025, // 5% tolerance in side length ratios
		Angle:       0.5,   // 1 degree tolerance in angles
	}

	// Record the start time
	startTime := time.Now()
	wcs, err := solver.Solve(tolerances, 2)

	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Calculate the elapsed time
	elapsedTime := time.Since(startTime)

	// Calculate the reference equatorial coordinate:
	eq = wcs.PixelToEquatorialCoordinate(578.231147766, 485.620500565)

	// We cross-reference here with calibration data from the astrometry.net API:
	if math.Abs(eq.RA-98.6467) > 0.001 {
		t.Errorf("RA not set correctly")
	}

	// We cross-reference here with calibration data from the astrometry.net API:
	if math.Abs(eq.Dec-2.5375) > 0.001 {
		t.Errorf("Dec not set correctly")
	}

	// Ensure that the solver executed in a reasonable amount of time:
	if elapsedTime.Seconds() > 0.5 {
		t.Errorf("plate solver took too long to execute")
	}

	t.Logf("solver.Solve(tolerances) completed in %v", elapsedTime)

	t.Logf("RA: %v, Dec: %v", eq.RA, eq.Dec)

	t.Logf(wcs.CTYPE1)

	t.Logf(wcs.CTYPE2)
}

/*****************************************************************************************************************/

func BenchmarkSolverOnMatches(b *testing.B) {
	// Attempt to open the file from the given filepath:
	file, err := os.Open("../../samples/noise16.fits")

	if err != nil {
		b.Errorf("NewFITSImageFromReader() os.Open(): %v", err)
	}

	// Defer closing the file:
	defer file.Close()

	// Assume an image of 2x2 pixels with 16-bit depth, and no offset:
	fit := fits.NewFITSImage(2, 0, 0, 65535)

	// Read in our exposure data into the image:
	err = fit.Read(file)

	if err != nil {
		b.Errorf("Read() error: %v", err)
	}

	// Attempt to get the RA header from the FITS file:
	ra, exists := fit.Header.Floats["RA"]
	if !exists {
		b.Errorf("ra header not found")
	}

	// Attempt to get the Dec header from the FITS file:
	dec, exists := fit.Header.Floats["DEC"]
	if !exists {
		b.Errorf("dec header not found")
	}

	eq := astrometry.ICRSEquatorialCoordinate{
		RA:  float64(ra.Value),
		Dec: float64(dec.Value),
	}

	radius := 2.0

	// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
	sources, err := GetCatalogSourcesForTesting(eq, radius, 100, 8)

	if err != nil {
		b.Errorf("error: %v", err)
	}

	// Attempt to create a new PlateSolver:
	solver, err := NewPlateSolver(fit, Params{
		RA:                  float64(ra.Value),  // The appoximate RA of the center of the image
		Dec:                 float64(dec.Value), // The appoximate Dec of the center of the image
		PixelScale:          2.061 / 3600.0,     // 2.061 arcseconds per pixel (0.0005725 degrees)
		ExtractionThreshold: 50,                 // Extract a minimum of 80 of the brightest stars
		Radius:              16,                 // 16 pixels radius for the star extraction
		Sigma:               8,                  // 8 pixels sigma for the Gaussian kernel
		Sources:             sources,
	})

	if err != nil {
		b.Errorf("error: %v", err)
	}

	// Define the tolerances for the solver, we can adjust these as needed:
	tolerances := geometry.InvariantFeatureTolerance{
		LengthRatio: 0.025, // 5% tolerance in side length ratios
		Angle:       0.5,   // 1 degree tolerance in angles
	}

	// Reset the timer to exclude setup time from the benchmark.
	b.ResetTimer()

	// Perform the solve operation in parallel:
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// It's important to ensure that each iteration is independent.
			// If Solve modifies the solver state, consider creating a new solver inside the loop.

			// Perform the solve operation in parallel:
			_, err = solver.Solve(tolerances, 2)
			if err != nil {
				b.Fatalf("Solver.Solve failed: %v", err)
			}
		}
	})
}

/*****************************************************************************************************************/
