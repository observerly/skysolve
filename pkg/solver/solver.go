/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package solver

/*****************************************************************************************************************/

import (
	"errors"
	"math"
	"sort"
	"sync"

	"github.com/observerly/iris/pkg/fits"
	"github.com/observerly/iris/pkg/photometry"
	stats "github.com/observerly/iris/pkg/statistics"

	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/geometry"
)

/*****************************************************************************************************************/

type PlateSolverCatalog int

/*****************************************************************************************************************/

const (
	GAIA PlateSolverCatalog = iota
)

/*****************************************************************************************************************/

type PlateSolver struct {
	Catalog    PlateSolverCatalog
	Stars      []photometry.Star
	Sources    []catalog.Source
	RA         float64
	Dec        float64
	PixelScale float64
	Width      int32
	Height     int32
}

/*****************************************************************************************************************/

func getCatalogSources(psc PlateSolverCatalog, eq astrometry.ICRSEquatorialCoordinate, radius float64) ([]catalog.Source, error) {
	switch psc {
	case GAIA:
		// Create a new GAIA service client:
		q := catalog.NewGAIAServiceClient()
		// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
		return q.PerformRadialSearch(eq, radius, 8)
	default:
		return nil, errors.New("unsupported catalog")
	}
}

/*****************************************************************************************************************/

type Params struct {
	RA                  float64
	Dec                 float64
	PixelScale          float64
	ExtractionThreshold float64
	Radius              float32
	Sigma               float32
}

/*****************************************************************************************************************/

func NewPlateSolver(
	psc PlateSolverCatalog,
	fit *fits.FITSImage,
	params Params,
) (*PlateSolver, error) {
	var (
		stars   []photometry.Star
		sources []catalog.Source
		err     error
	)

	ra := params.RA

	dec := params.Dec

	radius := params.Radius

	sigma := params.Sigma

	if fit == nil {
		return nil, errors.New("invalid FITS image")
	}

	// Calculate the width of the image in pixels:
	xs := int(fit.Header.Naxis1)

	// Calculate the height of the image in pixels:
	ys := int(fit.Header.Naxis2)

	// Extract the image data as a float32 array from the FITS file:
	d := fit.Data

	// Apply image preprocessing techniques to the FITS file, e.g.,
	// calculate Otsu thresholding and remove background noise:

	// Setup two wait groups for the sources lookup and the stars extractor:
	var wg sync.WaitGroup
	wg.Add(2)

	// Extract bright pixels (stars) from the image:
	go func() {
		defer wg.Done()

		// Setup the statistics object:
		stats := stats.NewStats(d, fit.ADU, xs)

		// Calculate the location and scale of the image:
		location, scale := stats.FastApproxSigmaClippedMedianAndQn()

		// Extract the image from the FITS file:
		sexp := photometry.NewStarsExtractor(d, xs, ys, radius, fit.ADU)

		// Set the threshold for the bright pixels:
		sexp.Threshold = location + scale*sigma

		// Extract the bright pixels from the image:
		stars = sexp.GetBrightPixels()

		// Sort the stars by intensity, in descending order:
		sort.Slice(stars, func(i, j int) bool {
			return stars[i].Intensity > stars[j].Intensity
		})

		// Get a minimum of 32 stars from our list of stars, e.g., the brightest 32 stars:
		minimum := int64(math.Min(params.ExtractionThreshold, float64(len(stars))))

		stars = stars[:minimum]
	}()

	// Get the sources from the catalog within the search radius, using the approximate RA and Dec coordinates
	// from the FITS HDU (header):
	go func() {
		defer wg.Done()

		sources, err = getCatalogSources(psc, astrometry.ICRSEquatorialCoordinate{
			RA:  ra,
			Dec: dec,
		}, 2)

		if err != nil {
			return
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Reset the fit.Data to an empty float32 array to preserve memory:
	d = []float32{}

	// If we encounter an error when retrieving the sources from the catalog, return the error:
	if err != nil {
		return nil, err
	}

	// Return a new PlateSolver object with the catalog, stars, sources, RA, Dec, and pixel scale:
	return &PlateSolver{
		Catalog:    psc,
		Stars:      stars,
		Sources:    sources,
		RA:         ra,
		Dec:        dec,
		PixelScale: params.PixelScale,
	}, nil
}

/*****************************************************************************************************************/

func (ps *PlateSolver) GenerateStarAsterisms() []astrometry.Asterism {
	triangles := []astrometry.Asterism{}

	n := len(ps.Stars)

	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				asterism := astrometry.Asterism{
					A: ps.Stars[i],
					B: ps.Stars[j],
					C: ps.Stars[k],
				}

				// Compute invariant features for the asterism:
				features, err := geometry.ComputeInvariantFeatures(
					float64(asterism.A.X),
					float64(asterism.A.Y),
					float64(asterism.B.X),
					float64(asterism.B.Y),
					float64(asterism.C.X),
					float64(asterism.C.Y),
				)
				if err != nil {
					continue
				}

				asterism.Features = features

				triangles = append(triangles, asterism)
			}
		}
	}

	return triangles
}

/*****************************************************************************************************************/

func (ps *PlateSolver) GenerateSourceAsterisms() []catalog.SourceAsterism {
	triangles := []catalog.SourceAsterism{}

	n := len(ps.Sources)

	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				sourceAsterism := catalog.SourceAsterism{
					A: ps.Sources[i],
					B: ps.Sources[j],
					C: ps.Sources[k],
				}

				// Compute invariant features and store them
				features, err := geometry.ComputeInvariantFeatures(
					sourceAsterism.A.RA,
					sourceAsterism.A.Dec,
					sourceAsterism.B.RA,
					sourceAsterism.B.Dec,
					sourceAsterism.C.RA,
					sourceAsterism.C.Dec,
				)
				if err != nil {
					continue
				}

				sourceAsterism.Features = features

				triangles = append(triangles, sourceAsterism)
			}
		}
	}

	return triangles
}

