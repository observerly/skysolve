/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/
/*****************************************************************************************************************/

type GAIAQuery struct {
	RA     float64 // right ascension (in degrees)
	Dec    float64 // right ascension (in degrees)
	Radius float64 // search radius (in degrees)
	Limit  float64 // limiting magnitude
}

/*****************************************************************************************************************/

type GAIAServiceClient struct {
	URI   string
	Query GAIAQuery
}

/*****************************************************************************************************************/

func NewGAIAServiceClient() *GAIAServiceClient {
	return &GAIAServiceClient{
		URI:   "https://gea.esac.esa.int/tap-server/tap/sync",
		Query: GAIAQuery{},
	}
}

/*****************************************************************************************************************/
