/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
	"bytes"
	"text/template"
)

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

const record = `source_id, designation, ra, dec, pmra, pmdec, parallax, phot_g_mean_flux, phot_g_mean_mag`

/*****************************************************************************************************************/

func (g *GAIAServiceClient) Build() (string, error) {
	// Define the ADQL query template for the GAIA TAP service:
	// @see https://gea.esac.esa.int/archive/documentation/GDR3/Gaia_archive/chap_datamodel/
	// N.B. (use only gold standard data, e.g., photometry processing mode (byte) i.e., phot_proc_mode = '0'):
	const queryTemplate = `
		SELECT {{.Record}}
		FROM gaiadr3.gaia_source
		WHERE CONTAINS(
			POINT('ICRS', ra, dec),
			CIRCLE('ICRS', {{.RA}}, {{.Dec}}, {{.Radius}})
		) = 1 AND phot_g_mean_mag < {{.Limit}} AND phot_proc_mode = '0'
	`

	// Parse the ADQL query template:
	tmpl, err := template.New("adql").Parse(queryTemplate)
	if err != nil {
		return "", err
	}

	// Set the template data from the query parameters:
	// Data to populate the template
	data := struct {
		Record string
		RA     float64
		Dec    float64
		Radius float64
		Limit  float64
	}{
		Record: record,
		RA:     g.Query.RA,
		Dec:    g.Query.Dec,
		Radius: g.Query.Radius,
		Limit:  g.Query.Limit,
	}

	// Execute the template with the data:
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	// Return the ADQL query string:
	return buf.String(), nil
}

/*****************************************************************************************************************/
