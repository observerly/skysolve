/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
	"errors"

	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/geometry"
)

/*****************************************************************************************************************/

type Catalog int

/*****************************************************************************************************************/

const (
	GAIA Catalog = iota
	SIMBAD
)

/*****************************************************************************************************************/

type Source struct {
	UID                       string  `json:"uid" gaia:"source_id" simbad:"uid"`                   // Source ID (unique)
	Designation               string  `json:"designation" gaia:"designation" simbad:"designation"` // Source Designation
	RA                        float64 `json:"ra" gaia:"ra" simbad:"ra"`                            // Right Ascension (in degrees)
	Dec                       float64 `json:"dec" gaia:"dec" simbad:"dec"`                         // Declination (in degrees)
	ProperMotionRA            float64 `json:"pmra" gaia:"pmra" simbad:"pmra"`                      // Proper Motion in RA (in mas/yr)
	ProperMotionDec           float64 `json:"pmdec" gaia:"pmdec" simbad:"pmdec"`                   // Proper Motion in Dec (in mas/yr)
	Parallax                  float64 `json:"parallax" gaia:"parallax" simbad:"parallax"`          // Parallax (in mas)
	PhotometricGMeanFlux      float64 `json:"flux" gaia:"phot_g_mean_flux" simbad:"flux"`          // G-band Mean Flux (in e-/s)
	PhotometricGMeanMagnitude float64 `json:"magnitude" gaia:"phot_g_mean_mag" simbad:"magnitude"` // G-band Mean Magnitude (in mag)
}

/*****************************************************************************************************************/

type SourceAsterism struct {
	A        Source
	B        Source
	C        Source
	Features geometry.InvariantFeatures
}

/*****************************************************************************************************************/

type CatalogService struct {
	Catalog   Catalog
	Limit     int
	Threshold float64
}

/*****************************************************************************************************************/

type Params struct {
	RA        float64
	Dec       float64
	Radius    float64
	Limit     int
	Threshold float64
}

/*****************************************************************************************************************/

func NewCatalogService(
	catalog Catalog,
	params Params,
) *CatalogService {
	return &CatalogService{
		Catalog:   catalog,
		Limit:     params.Limit,
		Threshold: params.Threshold,
	}
}

/*****************************************************************************************************************/

func (c *CatalogService) PerformRadialSearch(
	eq astrometry.ICRSEquatorialCoordinate,
	radius float64,
) ([]Source, error) {
	switch c.Catalog {
	case GAIA:
		// Create a new GAIA service client:
		q := NewGAIAServiceClient()
		// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
		return q.PerformRadialSearch(eq, radius, c.Limit, c.Threshold)
	case SIMBAD:
		// Create a new SIMBAD service client:
		q := NewSIMBADServiceClient()
		// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
		return q.PerformRadialSearch(eq, radius, c.Limit, c.Threshold)
	default:
		return nil, errors.New("unsupported catalog")
	}
}

/*****************************************************************************************************************/
