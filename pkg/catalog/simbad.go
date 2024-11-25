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
	"strings"
	"time"

	"github.com/observerly/skysolve/pkg/adql"
	"github.com/observerly/skysolve/pkg/astrometry"
)

/*****************************************************************************************************************/

type SIMBADQuery struct {
	RA        float64 // right ascension (in degrees)
	Dec       float64 // right ascension (in degrees)
	Radius    float64 // search radius (in degrees)
	Limit     int     // maximum number of records to return
	Threshold float64 // limiting magnitude
}

/*****************************************************************************************************************/

type SIMBADServiceClient struct {
	*adql.TapClient
	Query SIMBADQuery
}

/*****************************************************************************************************************/

func NewSIMBADServiceClient() *SIMBADServiceClient {
	// https://simbad.unistra.fr/simbad/sim-tap/sync
	url := url.URL{
		Scheme: "https",
		Host:   "simbad.unistra.fr",
		Path:   "/simbad/sim-tap/sync",
	}

	headers := map[string]string{
		// Default content type for TAP services
		"Content-Type": "application/x-www-form-urlencoded",
		// Ensure we are good citizens and identify ourselves:
		"X-Requested-By": "@observerly/skysolve",
	}

	client := adql.NewTapClient(url, 60*time.Second, headers)

	return &SIMBADServiceClient{
		TapClient: client,
		Query:     SIMBADQuery{},
	}
}

/*****************************************************************************************************************/

const simbadRecord = "basic.oid AS uid, basic.main_id AS designation, basic.ra AS ra, basic.dec AS dec, basic.pmra AS pmra, basic.pmdec AS pmdec, basic.plx_value AS parallax, flux.flux AS flux, allfluxes.G AS magnitude"

/*****************************************************************************************************************/

func (s *SIMBADServiceClient) PerformRadialSearch(eq astrometry.ICRSEquatorialCoordinate, radius float64, limit int, threshold float64) ([]Source, error) {
	// Define the ADQL query template for the SIMBAD TAP service:
	// @see https://simbad.u-strasbg.fr/Pages/guide/sim-q.htx
	const simbadADQLTemplate = `
		SELECT TOP {{.Limit}} {{.Record}}
		FROM basic
		LEFT JOIN flux 
			ON basic.oid = flux.oidref 
			AND flux.filter = 'G'
		LEFT JOIN allfluxes 
			ON basic.oid = allfluxes.oidref
		WHERE CONTAINS(
			POINT('ICRS', basic.ra, basic.dec),
			CIRCLE('ICRS', {{.RA}}, {{.Dec}}, {{.Radius}})
		) = 1
		ORDER BY magnitude ASC;
	`

	// Set the query parameters:
	s.Query.RA = eq.RA
	s.Query.Dec = eq.Dec
	s.Query.Radius = radius
	s.Query.Limit = limit
	s.Query.Threshold = threshold

	// Construct the ADQL query from the template:
	adqlQuery, err := s.BuildADQLQuery(simbadADQLTemplate, struct {
		Record    string
		RA        float64
		Dec       float64
		Radius    float64
		Limit     int
		Threshold float64
	}{
		Record:    simbadRecord,
		RA:        s.Query.RA,
		Dec:       s.Query.Dec,
		Radius:    s.Query.Radius,
		Limit:     s.Query.Limit,
		Threshold: s.Query.Threshold,
	})
	if err != nil {
		return nil, err
	}

	// Execute the query and get the response:
	tapResponse, err := s.ExecuteADQLQuery(adqlQuery)

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

		// Assign UID and Designation using fmt.Sprintf to handle various types:
		star.UID = fmt.Sprintf("%v", record[0])
		star.Designation = strings.Join(strings.Fields(fmt.Sprintf("%v", record[1])), " ")

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
