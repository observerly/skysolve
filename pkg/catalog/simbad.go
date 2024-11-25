/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
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
