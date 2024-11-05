/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package solver

/*****************************************************************************************************************/

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/observerly/iris/pkg/fits"
	"github.com/observerly/iris/pkg/photometry"
	stats "github.com/observerly/iris/pkg/statistics"
	iutils "github.com/observerly/iris/pkg/utils"

	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/geometry"
	"github.com/observerly/skysolve/pkg/matrix"
	"github.com/observerly/skysolve/pkg/projection"
	"github.com/observerly/skysolve/pkg/transform"
	"github.com/observerly/skysolve/pkg/utils"
	"github.com/observerly/skysolve/pkg/wcs"
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

/*****************************************************************************************************************/

type Match struct {
	Star   photometry.Star
	Source catalog.Source
}

/*****************************************************************************************************************/

func (ps *PlateSolver) MatchAsterismsWithCatalog(
	asterism astrometry.Asterism,
	sourceAsterism catalog.SourceAsterism,
	tolerance geometry.InvariantFeatureTolerance,
) ([]Match, error) {
	// Define the stars for the asterism:
	stars := []photometry.Star{asterism.A, asterism.B, asterism.C}

	// Define the sources for the source asterism:
	sources := []catalog.Source{sourceAsterism.A, sourceAsterism.B, sourceAsterism.C}

	// Define the reference Right Ascension coordinates for the image:
	CRVAL1 := ps.RA

	// Define the reference declination coordinates for the image:
	CRVAL2 := ps.Dec

	projectedSources := make([]struct{ x, y float64 }, 3)
	for i, source := range sources {
		// Convert the equatorial coordinates to gnomic coordinates:
		x, y := projection.ConvertEquatorialToGnomic(source.RA, source.Dec, CRVAL1, CRVAL2)
		// Store the projected coordinates:
		projectedSources[i] = struct{ x, y float64 }{x, y}
	}

	// Define permutations of indices (total 6 permutations for 3 elements) to rearrange the sources:
	// This ensures we try all possible combinations of sources to match the asterism:
	permutations := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, perm := range permutations {
		// Rearrange the sources according to the permutation order:
		mappedSources := []catalog.Source{
			sources[perm[0]],
			sources[perm[1]],
			sources[perm[2]],
		}

		// Use projected coordinates for the rearranged sources:
		x1, y1 := projectedSources[perm[0]].x, projectedSources[perm[0]].y
		x2, y2 := projectedSources[perm[1]].x, projectedSources[perm[1]].y
		x3, y3 := projectedSources[perm[2]].x, projectedSources[perm[2]].y

		// Compute invariant features for the rearranged sources:
		features, err := geometry.ComputeInvariantFeatures(
			x1, y1,
			x2, y2,
			x3, y3,
		)

		if err != nil {
			continue
		}

		// Compare the features with the asterism's features
		if geometry.CompareInvariantFeatures(asterism.Features, features, tolerance) {
			// If they match, create the matches:
			matches := []Match{
				{Star: stars[0], Source: mappedSources[0]},
				{Star: stars[1], Source: mappedSources[1]},
				{Star: stars[2], Source: mappedSources[2]},
			}
			return matches, nil
		}
	}

	// If no match is found, return an error
	return nil, errors.New("no match found between asterism and source asterism")
}

/*****************************************************************************************************************/

func (ps *PlateSolver) FindSourceMatches(tolerance geometry.InvariantFeatureTolerance) ([]Match, error) {
	var (
		asterisms       []astrometry.Asterism
		sourceAsterisms []catalog.SourceAsterism
	)

	// Setup two wait groups for the sources lookup and the stars extractor:
	var wg sync.WaitGroup
	wg.Add(2)

	// Get the asterisms or triangulated stars from the image:
	go func() {
		defer wg.Done()
		asterisms = ps.GenerateStarAsterisms()
	}()

	// Get the asterisms or triangulated sources from the catalog:
	go func() {
		defer wg.Done()
		sourceAsterisms = ps.GenerateSourceAsterisms()
	}()

	wg.Wait()

	// Define the precision for quantization, 4 seems to be a good value between being too
	// specific and not specific enough:
	precision := 4

	// Create a map to index source triangles by their quantized invariant features:
	sourceTriangleIndex := make(map[string][]catalog.SourceAsterism)

	// Compute invariant features for source triangles and index them:
	for _, source := range sourceAsterisms {
		key := utils.QuantizeFeatures(source.Features, precision)
		sourceTriangleIndex[key] = append(sourceTriangleIndex[key], source)
	}

	matches := make([]Match, 0)

	for i, asterism := range asterisms {
		// Quantize the asterism's features to create a key:
		key := utils.QuantizeFeatures(asterism.Features, precision)

		// Get the source triangles with the same invariant features, e.g. by looking for matching source
		// triangles in the index:
		if sources, found := sourceTriangleIndex[key]; found {
			for _, source := range sources {
				// Attempt to match the individual stars to sources in the catalog, using the star triangle and source triangle:
				match, err := ps.MatchAsterismsWithCatalog(asterism, source, tolerance)

				if err == nil && match != nil {
					// Add the matches to the list of matches:
					matches = append(matches, match...)
				}
			}
		}

		// If we have more than 32 matches, we can stop looking for more, this seems to be an optimal
		// amount of matches to plate solve to a high degree of accuracy, this covers more than 30% of
		// the total source count.
		if len(matches) >= 32 && i < 40 {
			break
		}
	}

	return matches, nil
}

/*****************************************************************************************************************/

// solveForAffineParameters fits an affine transformation matrix to the matches.
//
//lint:ignore U1000 Reserved for future implementation.
func (ps *PlateSolver) solveForAffineParameters(
	a [][]float64,
	b []float64,
	n int,
) (*transform.Affine2DParameters, error) {
	var (
		B      *matrix.Matrix
		aT     *matrix.Matrix
		aTaInv *matrix.Matrix
		err    error
	)

	// WaitGroup for creating A and B concurrently:
	var wg sync.WaitGroup
	wg.Add(2)

	errorChannel := make(chan error, 2)

	go func() {
		defer wg.Done()

		// Convert A to a matrix:
		A, err := matrix.NewFromSlice(iutils.Flatten2DFloat64Array(a), n, 6)

		if err != nil {
			errorChannel <- fmt.Errorf("failed to create A matrix: %v", err)
			return
		}

		// Compute A^T (transpose of A):
		aT, err = A.Transpose()

		if err != nil {
			errorChannel <- fmt.Errorf("failed to compute A^T: %v", err)
			return
		}

		// Compute A^T * A (matrix multiplication):
		aTa, err := aT.Multiply(A)
		if err != nil {
			errorChannel <- fmt.Errorf("failed to compute A^T * A: %v", err)
			return
		}

		// Compute the inverse of A^T * A (matrix inversion):
		aTaInv, err = aTa.Invert()
		if err != nil {
			errorChannel <- fmt.Errorf("failed to invert A^T * A: %v", err)
			return
		}
	}()

	go func() {
		defer wg.Done()

		// Convert B to a matrix:
		B, err = matrix.NewFromSlice(b, n, 1)

		if err != nil {
			errorChannel <- fmt.Errorf("failed to create B matrix: %v", err)
			return
		}
	}()

	// Wait for both goroutines to finish and close the error channel:
	wg.Wait()
	// Close the error channel to signal that no more errors will be sent:
	close(errorChannel)

	// If an error occured, return the error:
	if err := <-errorChannel; err != nil {
		return nil, err
	}

	// Compute A^T * B (matrix multiplication):
	aTb, err := aT.Multiply(B)
	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T * B: %v", err)
	}

	if aTaInv == nil || aTb == nil {
		return nil, errors.New("failed to compute affine transformation matrix parameters")
	}

	// Compute the affine transformation matrix parameters using the least squares method:
	params := make([]float64, 6)

	// Calculate the affine transformation matrix parameters:
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			params[i] += aTaInv.Value[i*6+j] * aTb.Value[j]
		}
	}

	affineParams := transform.Affine2DParameters{
		A: params[0],
		B: params[1],
		C: params[3],
		D: params[4],
		E: params[2],
		F: params[5],
	}

	return &affineParams, nil
}

/*****************************************************************************************************************/

// solveForSIPParameters fits higher-order SIP polynomials to the residuals after the affine transformation.
//
//lint:ignore U1000 Reserved for future implementation.
func (ps *PlateSolver) solveForSIPParameters(
	aRA [][]float64,
	aDec [][]float64,
	bRA []float64,
	bDec []float64,
	n int,
	sipOrder int,
) (*transform.SIP2DParameters, error) {
	// Calculate the number of terms in the SIP polynomial:
	numTerms := (sipOrder + 1) * (sipOrder + 2) / 2

	// Convert SIP design matrices and B vectors to matrices
	aSIP_RA, err := matrix.NewFromSlice(iutils.Flatten2DFloat64Array(aRA), n, numTerms)
	if err != nil {
		return nil, fmt.Errorf("failed to create SIP RA matrix: %v", err)
	}

	bSIP_RA, err := matrix.NewFromSlice(bRA, n, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create SIP RA vector: %v", err)
	}

	aSIP_Dec, err := matrix.NewFromSlice(iutils.Flatten2DFloat64Array(aDec), n, numTerms)
	if err != nil {
		return nil, fmt.Errorf("failed to create SIP Dec matrix: %v", err)
	}

	bSIP_Dec, err := matrix.NewFromSlice(bDec, n, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create SIP Dec vector: %v", err)
	}

	aSIPT_RA, err := aSIP_RA.Transpose()
	if err != nil {
		return nil, fmt.Errorf("failed to transpose SIP RA matrix: %v", err)
	}

	aTaSIP_RA, err := aSIPT_RA.Multiply(aSIP_RA)
	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T * A for SIP RA: %v", err)
	}

	aTbSIP_RA, err := aSIPT_RA.Multiply(bSIP_RA)
	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T * B for SIP RA: %v", err)
	}

	aTaInvSIP_RA, err := aTaSIP_RA.Invert()
	if err != nil {
		return nil, fmt.Errorf("failed to invert A^T * A for SIP RA: %v", err)
	}

	sipParamsRA := make([]float64, numTerms)
	for i := 0; i < numTerms; i++ {
		for j := 0; j < numTerms; j++ {
			sipParamsRA[i] += aTaInvSIP_RA.Value[i*numTerms+j] * aTbSIP_RA.Value[j]
		}
	}

	// Solve for SIP Dec Parameters
	aSIPT_Dec, err := aSIP_Dec.Transpose()
	if err != nil {
		return nil, fmt.Errorf("failed to transpose SIP Dec matrix: %v", err)
	}

	aTaSIP_Dec, err := aSIPT_Dec.Multiply(aSIP_Dec)
	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T * A for SIP Dec: %v", err)
	}

	aTbSIP_Dec, err := aSIPT_Dec.Multiply(bSIP_Dec)
	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T * B for SIP Dec: %v", err)
	}

	aTaInvSIP_Dec, err := aTaSIP_Dec.Invert()
	if err != nil {
		return nil, fmt.Errorf("failed to invert A^T * A for SIP Dec: %v", err)
	}

	sipParamsDec := make([]float64, numTerms)
	for i := 0; i < numTerms; i++ {
		for j := 0; j < numTerms; j++ {
			sipParamsDec[i] += aTaInvSIP_Dec.Value[i*numTerms+j] * aTbSIP_Dec.Value[j]
		}
	}

	// Map SIP coefficients to FITS term keys for RA and Dec:
	sipTermKeysA := utils.GeneratePolynomialTermKeys("A", sipOrder)
	sipTermKeysB := utils.GeneratePolynomialTermKeys("B", sipOrder)

	if len(sipTermKeysA) != numTerms || len(sipTermKeysB) != numTerms {
		return nil, fmt.Errorf("incorrect number of SIP term keys: got %d for A and %d for B, expected %d each", len(sipTermKeysA), len(sipTermKeysB), numTerms)
	}

	aPowerMap := make(map[string]float64)
	bPowerMap := make(map[string]float64)

	for idx, term := range sipTermKeysA {
		aPowerMap[term] = sipParamsRA[idx]
	}

	for idx, term := range sipTermKeysB {
		bPowerMap[term] = sipParamsDec[idx]
	}

	sipParams := transform.SIP2DParameters{
		AOrder: sipOrder,
		BOrder: sipOrder,
		APower: aPowerMap,
		BPower: bPowerMap,
	}

	return &sipParams, nil
}

/*****************************************************************************************************************/

func (ps *PlateSolver) Solve(tolerance geometry.InvariantFeatureTolerance) (*wcs.WCS, error) {
	matches, err := ps.FindSourceMatches(tolerance)

	if err != nil {
		return nil, err
	}

	if len(matches) < 3 {
		return nil, errors.New("insufficient matches to perform plate solving")
	}

	n := 2 * len(matches)
	A := make([][]float64, n)
	B := make([]float64, n)

	// Calculate the affine transformation matrix from the matches:
	for i, match := range matches {
		x := float64(match.Star.X)
		y := float64(match.Star.Y)

		ra := match.Source.RA
		dec := match.Source.Dec

		// Create the matrix A and vector B for the least squares method:
		A[2*i] = []float64{x, y, 1, 0, 0, 0}
		B[2*i] = ra
		A[2*i+1] = []float64{0, 0, 0, x, y, 1}
		B[2*i+1] = dec
	}

	// Convert A to a matrix:
	a, err := matrix.NewFromSlice(iutils.Flatten2DFloat64Array(A), n, 6)
	if err != nil {
		return nil, err
	}

	// Convert B to a matrix:
	b, err := matrix.NewFromSlice(B, n, 1)
	if err != nil {
		return nil, err
	}

	// Compute A^T
	aT, err := a.Transpose()

	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T: %v", err)
	}

	// Compute A^T * A
	aTa, err := aT.Multiply(a)
	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T * A: %v", err)
	}

	// Compute A^T * B
	aTb, err := aT.Multiply(b)
	if err != nil {
		return nil, fmt.Errorf("failed to compute A^T * B: %v", err)
	}

	// Compute the inverse of A^T * A
	aTaInv, err := aTa.Invert()
	if err != nil {
		return nil, fmt.Errorf("failed to invert A^T * A: %v", err)
	}

	// Compute the affine transformation matrix parameters using the least squares method:
	params := make([]float64, 6)
	if aTaInv == nil || aTb == nil {
		return nil, errors.New("failed to compute affine transformation matrix parameters")
	}

	// Calculate the affine transformation matrix parameters:
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			params[i] += aTaInv.Value[i*6+j] * aTb.Value[j]
		}
	}

	// Calculate the x-coordinate of the center of the image:
	x := float64(ps.Width) / 2

	// Calculate the y-coordinate of the center of the image:
	y := float64(ps.Height) / 2

	// Create the affine parameters:
	affineParams := transform.Affine2DParameters{
		A: params[0],
		B: params[1],
		C: params[3],
		D: params[4],
		E: params[2],
		F: params[5],
	}

	// Create the SIP parameters:
	sipParams := transform.SIP2DParameters{}

	// Now that we have the affine parameters, we can calculate the actual RA and dec coordinate
	// for the center of the image:
	t := wcs.NewWorldCoordinateSystem(
		x,
		y,
		wcs.WCSParams{
			Projection:   wcs.RADEC_TAN,
			AffineParams: affineParams,
			SIPParams:    sipParams,
		},
	)

	return &t, nil
}

/*****************************************************************************************************************/
