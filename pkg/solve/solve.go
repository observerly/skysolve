/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package solve

/*****************************************************************************************************************/

import (
	"math"
	"sort"
	"sync"

	"github.com/observerly/iris/pkg/photometry"
	stats "github.com/observerly/iris/pkg/statistics"

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

// NewPlateSolver initializes a new PlateSolver with the given FITS image and parameters.
func NewPlateSolver(
	params Params,
) (*PlateSolver, error) {
	radius := params.Radius
	sigma := params.Sigma
	stars := []photometry.Star{}
	sources := []catalog.Source{}

	// Calculate the width of the image in pixels:
	xs := params.Width

	// Calculate the height of the image in pixels:
	ys := params.Height

	// Setup a wait group for the stars extractor:
	var wg sync.WaitGroup
	wg.Add(1)

	// Extract bright pixels (stars) from the image:
	go func() {
		defer wg.Done()

		// Extract the image from the FITS file:
		sexp := photometry.NewStarsExtractor(params.Data, xs, ys, float32(radius), params.ADU)

		// Extract the bright pixels from the image:
		starsExtracted := sexp.FindStars(stats.NewStats(params.Data, params.ADU, xs), float32(sigma), 2.2)

		// Sort the stars by intensity, in descending order:
		sort.Slice(starsExtracted, func(i, j int) bool {
			return starsExtracted[i].Intensity > starsExtracted[j].Intensity
		})

		// Get a minimum of X stars from our list of stars, e.g., the brightest X stars:
		k := math.Min(float64(len(starsExtracted)), params.ExtractionThreshold)

		stars = starsExtracted[:int(k)]
	}()

	// Wait for the stars extractor to finish
	wg.Wait()

	// Return a new PlateSolver object with the catalog, stars, sources, RA, Dec, and pixel scale:
	return &PlateSolver{
		Stars:   stars,
		Sources: sources,
		Data:    params.Data,
		Width:   xs,
		Height:  ys,
	}, nil
}

/*****************************************************************************************************************/
