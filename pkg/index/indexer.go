/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package index

/*****************************************************************************************************************/

import (
	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/healpix"
	"github.com/observerly/skysolve/pkg/quad"
	"github.com/observerly/skysolve/pkg/solve"
	"github.com/observerly/skysolve/pkg/star"
)

/*****************************************************************************************************************/

type Indexer struct {
	Catalog catalog.CatalogService
	HealPIX healpix.HealPIX
}

/*****************************************************************************************************************/

func NewIndexer(
	healpix healpix.HealPIX,
	catalog catalog.CatalogService,
) *Indexer {
	return &Indexer{
		Catalog: catalog,
		HealPIX: healpix,
	}
}

/*****************************************************************************************************************/

func (i *Indexer) GenerateStarsForPixel(pixel int) ([]star.Star, error) {
	// Convert the pixel index to the pixel's equatorial coordinate:
	eq := i.HealPIX.ConvertPixelIndexToEquatorial(pixel)

	// Get the radial extent of the pixel:
	radius := i.HealPIX.GetPixelRadialExtent(pixel)

	// Get the sources within the pixel's radial extent:
	sources, err := i.Catalog.PerformRadialSearch(eq, radius)

	// If we encounter an error, return it:
	if err != nil {
		return nil, err
	}

	// Convert the sources to stars:
	var stars []star.Star

	for _, source := range sources {
		// For each source, we just need to sense check that the source is within the pixel:
		eq := astrometry.ICRSEquatorialCoordinate{
			RA:  source.RA,
			Dec: source.Dec,
		}

		pixelIndex := i.HealPIX.ConvertEquatorialToPixelIndex(eq)

		// If the source is not within the pixel, skip it:
		if pixelIndex != pixel {
			continue
		}

		stars = append(stars, star.Star{
			Designation: source.Designation,
			X:           source.RA,
			Y:           source.Dec,
			RA:          source.RA,
			Dec:         source.Dec,
			Intensity:   source.PhotometricGMeanFlux,
		})
	}

	// We should have at least 5 stars, if we have more just slice the first 5:
	if len(stars) > 5 {
		stars = stars[:5]
	}

	return stars, nil
}

/*****************************************************************************************************************/

func (i *Indexer) GenerateQuadsForPixel(pixel int) ([]quad.Quad, error) {
	// Convert the pixel index to the pixel's equatorial coordinate:
	eq := i.HealPIX.ConvertPixelIndexToEquatorial(pixel)

	// Get the radial extent of the pixel:
	radius := i.HealPIX.GetPixelRadialExtent(pixel)

	// Get the sources within the pixel's radial extent:
	sources, err := i.Catalog.PerformRadialSearch(eq, radius)

	// If we encounter an error, return it:
	if err != nil {
		return nil, err
	}

	// Convert the sources to stars:
	stars := make([]star.Star, len(sources))

	for _, source := range sources {
		// For each source, we just need to sense check that the source is within the pixel:
		eq := astrometry.ICRSEquatorialCoordinate{
			RA:  source.RA,
			Dec: source.Dec,
		}

		pixelIndex := i.HealPIX.ConvertEquatorialToPixelIndex(eq)

		// If the source is not within the pixel, skip it:
		if pixelIndex != pixel {
			continue
		}

		stars = append(stars, star.Star{
			Designation: source.Designation,
			X:           source.RA,
			Y:           source.Dec,
			RA:          source.RA,
			Dec:         source.Dec,
			Intensity:   source.PhotometricGMeanFlux,
		})
	}

	// We should have at least 5 sources to generate a quad:
	if len(stars) < 5 {
		return nil, nil
	}

	// Indexing the sources here involves storing the sources locally for future reference:
	quads, err := solve.GenerateEuclidianStarQuads(stars, 5)

	// If we encounter an error, return it:
	if err != nil {
		return nil, err
	}

	return quads, nil
}

/*****************************************************************************************************************/
