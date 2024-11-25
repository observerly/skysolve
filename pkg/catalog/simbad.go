/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
	"net/url"
	"time"

	"github.com/observerly/skysolve/pkg/adql"
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
