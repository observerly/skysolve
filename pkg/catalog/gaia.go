/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"text/template"

	"github.com/observerly/skysolve/pkg/astrometry"
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

// Gaia DR3 service handler. The five-parameter astrometric solution, positions on the sky (α, δ),
// parallaxes, and proper motions, are given for around 1.46 billion sources, with a limiting magnitude
// of G = 21.
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

func (g *GAIAServiceClient) PerformRadialSearch(eq astrometry.ICRSEquatorialCoordinate, radius float64, limit float64) ([]Source, error) {
	// Set the query parameters:
	g.Query.RA = eq.RA
	g.Query.Dec = eq.Dec
	g.Query.Radius = radius
	g.Query.Limit = limit

	// Construct the ADQL query from the template:
	adqlQuery, err := g.Build()
	if err != nil {
		return nil, err
	}

	// Prepare the POST form data for the HTTP request:
	formData := url.Values{}
	formData.Set("REQUEST", "doQuery")
	formData.Set("LANG", "ADQL")
	formData.Set("FORMAT", "csv")
	formData.Set("QUERY", adqlQuery)

	// Send the HTTP request to the GAIA TAP service:
	resp, err := http.PostForm(g.URI, formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body:
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Check for HTTP errors and return the response body if not OK:
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GAIA TAP query failed: %s", string(bodyBytes))
	}

	// Parse the CSV data from the response body:
	records, err := csv.NewReader(bytes.NewReader(bodyBytes)).ReadAll()
	if err != nil {
		return nil, err
	}

	// Unmarshal the CSV data into our struct array:
	var stars []Source

	// Iterate over the records and extract the source star data, skipping the header row:
	for _, record := range records[1:] {
		ra, err := strconv.ParseFloat(fmt.Sprintf("%v", record[2]), 64)
		if err != nil {
			continue
		}

		dec, err := strconv.ParseFloat(fmt.Sprintf("%v", record[3]), 64)
		if err != nil {
			continue
		}

		pmra, err := strconv.ParseFloat(fmt.Sprintf("%v", record[4]), 64)
		if err != nil {
			continue
		}

		pmdec, err := strconv.ParseFloat(fmt.Sprintf("%v", record[5]), 64)
		if err != nil {
			continue
		}

		parallax, err := strconv.ParseFloat(fmt.Sprintf("%v", record[6]), 64)
		if err != nil {
			continue
		}

		flux, err := strconv.ParseFloat(fmt.Sprintf("%v", record[7]), 64)
		if err != nil {
			continue
		}

		mag, err := strconv.ParseFloat(fmt.Sprintf("%v", record[8]), 64)
		if err != nil {
			continue
		}

		// Create a new source star object:
		star := Source{
			UID:                       record[0],
			Designation:               record[1],
			RA:                        ra,
			Dec:                       dec,
			ProperMotionRA:            pmra,
			ProperMotionDec:           pmdec,
			Parallax:                  parallax,
			PhotometricGMeanFlux:      flux,
			PhotometricGMeanMagnitude: mag,
		}

		// Append the source star to the array:
		stars = append(stars, star)
	}

	// Convert string fields to float64 for RA, Dec, and Magnitude:
	for i, star := range stars {
		ra, err := strconv.ParseFloat(fmt.Sprintf("%v", star.RA), 64)
		if err != nil {
			continue
		}

		dec, err := strconv.ParseFloat(fmt.Sprintf("%v", star.Dec), 64)
		if err != nil {
			continue
		}

		mag, err := strconv.ParseFloat(fmt.Sprintf("%v", star.PhotometricGMeanMagnitude), 64)
		if err != nil {
			continue
		}

		stars[i].RA = ra
		stars[i].Dec = dec
		stars[i].PhotometricGMeanMagnitude = mag
	}

	// Return the extracted source star data:
	return stars, nil
}

/*****************************************************************************************************************/
