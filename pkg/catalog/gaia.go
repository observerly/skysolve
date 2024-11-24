/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
	"fmt"
	"net/url"
	"time"

	"github.com/observerly/skysolve/pkg/adql"
	"github.com/observerly/skysolve/pkg/astrometry"
)

/*****************************************************************************************************************/

type GAIAQuery struct {
	RA        float64 // right ascension (in degrees)
	Dec       float64 // right ascension (in degrees)
	Radius    float64 // search radius (in degrees)
	Limit     int     // maximum number of records to return
	Threshold float64 // limiting magnitude
}

/*****************************************************************************************************************/

type GAIAServiceClient struct {
	*adql.TapClient
	Query GAIAQuery
}

/*****************************************************************************************************************/

// Gaia DR3 service handler. The five-parameter astrometric solution, positions on the sky (α, δ),
// parallaxes, and proper motions, are given for around 1.46 billion sources, with a limiting magnitude
// of G = 21.
func NewGAIAServiceClient() *GAIAServiceClient {
	// https://gea.esac.esa.int/tap-server/tap/sync
	url := url.URL{
		Scheme: "https",
		Host:   "gea.esac.esa.int",
		Path:   "/tap-server/tap/sync",
	}

	headers := map[string]string{
		// Default content type for TAP services
		"Content-Type": "application/x-www-form-urlencoded",
		// Ensure we are good citizens and identify ourselves:
		"X-Requested-By": "@observerly/skysolve",
	}

	client := adql.NewTapClient(url, 60*time.Second, headers)

	return &GAIAServiceClient{
		TapClient: client,
		Query:     GAIAQuery{},
	}
}

/*****************************************************************************************************************/

const gaiaRecord = `source_id, designation, ra, dec, pmra, pmdec, parallax, phot_g_mean_flux, phot_g_mean_mag`

/*****************************************************************************************************************/

func (g *GAIAServiceClient) PerformRadialSearch(eq astrometry.ICRSEquatorialCoordinate, radius float64, limit int, threshold float64) ([]Source, error) {
	// Define the ADQL query template for the GAIA TAP service:
	// @see https://gea.esac.esa.int/archive/documentation/GDR2/Gaia_archive/chap_datamodel/
	// N.B. (use only gold standard data, e.g., photometry processing mode (byte) i.e., phot_proc_mode = '0'):
	const gaiaADQLTemplate = `
		SELECT TOP {{.Limit}} {{.Record}}
		FROM gaiadr2.gaia_source
		WHERE CONTAINS(
			POINT('ICRS', ra, dec),
			CIRCLE('ICRS', {{.RA}}, {{.Dec}}, {{.Radius}})
		) = 1 
		AND phot_g_mean_mag < {{.Threshold}}
		AND phot_rp_mean_flux IS NOT NULL
		ORDER BY phot_g_mean_mag ASC;
	`

	// Set the query parameters:
	g.Query.RA = eq.RA
	g.Query.Dec = eq.Dec
	g.Query.Radius = radius
	g.Query.Limit = limit
	g.Query.Threshold = threshold

	// Construct the ADQL query from the template:
	adqlQuery, err := g.BuildADQLQuery(gaiaADQLTemplate, struct {
		Record    string
		RA        float64
		Dec       float64
		Radius    float64
		Limit     int
		Threshold float64
	}{
		Record:    gaiaRecord,
		RA:        g.Query.RA,
		Dec:       g.Query.Dec,
		Radius:    g.Query.Radius,
		Limit:     g.Query.Limit,
		Threshold: g.Query.Threshold,
	})
	if err != nil {
		return nil, err
	}

	// Execute the query and get the response:
	tapResponse, err := g.ExecuteADQLQuery(adqlQuery)
	if err != nil {
		return nil, err
	}

	var stars []Source

	// Helper functions defined locally within the method to convert an unknown interface{} to float64:
	toFloat64 := func(val interface{}) (float64, bool) {
		v, ok := val.(float64)
		return v, ok
	}

	for _, record := range tapResponse.Data {
		// Create a new Source struct from the record:
		// Initialize a new Source struct
		var star Source

		// Assign UID and Designation using fmt.Sprintf to handle various types
		star.UID = fmt.Sprintf("%v", record[0])
		star.Designation = fmt.Sprintf("%v", record[1])

		// Safely assign RA and assig a default value if not a float64:
		if ra, ok := toFloat64(record[2]); ok {
			star.RA = ra
		} else {
			// Handle unexpected type or assign a default value:
			star.RA = 0.0
		}

		// Safely assign Dec and assign a default value if not a float64:
		if dec, ok := toFloat64(record[3]); ok {
			star.Dec = dec
		} else {
			// Handle unexpected type or assign a default value:
			star.Dec = 0.0
		}

		// Safely assign ProperMotionRA if not nil:
		if record[4] != nil {
			if pmra, ok := toFloat64(record[4]); ok {
				star.ProperMotionRA = pmra
			}
		}

		// Safely assign ProperMotionDec if not nil:
		if record[5] != nil {
			if pmdec, ok := toFloat64(record[4]); ok {
				star.ProperMotionDec = pmdec
			}
		}

		// Safely assign Parallax if not nil:
		if record[6] != nil {
			if parallax, ok := toFloat64(record[6]); ok {
				star.Parallax = parallax
			}
		}

		// Safely assign PhotometricGMeanFlux if not nil:
		if record[7] != nil {
			if flux, ok := toFloat64(record[7]); ok {
				star.PhotometricGMeanFlux = flux
			}
		}

		// Safely assign PhotometricGMeanMagnitude if not nil:
		if record[8] != nil {
			if magnitude, ok := toFloat64(record[8]); ok {
				star.PhotometricGMeanMagnitude = magnitude
			}
		}

		// Append to the stars slice of Source structs:
		stars = append(stars, star)
	}

	return stars, nil
}

/*****************************************************************************************************************/
