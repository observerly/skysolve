/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
	"github.com/observerly/skysolve/pkg/geometry"
)

/*****************************************************************************************************************/

type Source struct {
	UID                       string  `json:"uid" gaia:"source_id"`              // Source ID (unique)
	Designation               string  `json:"designation" gaia:"designation"`    // Source Designation
	RA                        float64 `json:"ra" gaia:"ra"`                      // Right Ascension (in degrees)
	Dec                       float64 `json:"dec" gaia:"dec"`                    // Declination (in degrees)
	ProperMotionRA            float64 `json:"pmra" gaia:"pmra"`                  // Proper Motion in RA (in mas/yr)
	ProperMotionDec           float64 `json:"pmdec" gaia:"pmdec"`                // Proper Motion in Dec (in mas/yr)
	Parallax                  float64 `json:"parallax" gaia:"parallax"`          // Parallax (in mas)
	PhotometricGMeanFlux      float64 `json:"flux" gaia:"phot_g_mean_mag"`       // Mean Flux (in e-/s)
	PhotometricGMeanMagnitude float64 `json:"magnitude" gaia:"phot_g_mean_flux"` // Mean Magnitude (in mag)
}

/*****************************************************************************************************************/

type SourceAsterism struct {
	A        Source
	B        Source
	C        Source
	Features geometry.InvariantFeatures
}

/*****************************************************************************************************************/
