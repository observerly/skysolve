/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package catalog

/*****************************************************************************************************************/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/observerly/skysolve/pkg/astrometry"
)

/*****************************************************************************************************************/

type GAIAQuery struct {
	RA        float64 // right ascension (in degrees)
	Dec       float64 // right ascension (in degrees)
	Radius    float64 // search radius (in degrees)
	Limit     float64 // maximum number of records to return
	Threshold float64 // limiting magnitude
}

/*****************************************************************************************************************/

type GAIAServiceClient struct {
	URI    string
	Query  GAIAQuery
	Client *http.Client
}

/*****************************************************************************************************************/

// GAIAResponse represents the JSON structure returned by GAIA TAP service.
type GAIAResponse struct {
	// Assuming GAIA TAP returns a "data" field containing an array of records.
	Data [][]interface{} `json:"data"`
}

/*****************************************************************************************************************/

// Gaia DR3 service handler. The five-parameter astrometric solution, positions on the sky (α, δ),
// parallaxes, and proper motions, are given for around 1.46 billion sources, with a limiting magnitude
// of G = 21.
func NewGAIAServiceClient() *GAIAServiceClient {
	// Create a custom dialer with a timeout of 5 seconds:
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Create a custom transport with the dialer and a TLS handshake timeout of 1 second:
	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 1 * time.Second,
	}

	// Create a custom HTTP client with the transport and overall timeout:
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return &GAIAServiceClient{
		URI:    "https://gea.esac.esa.int/tap-server/tap/sync",
		Query:  GAIAQuery{},
		Client: client,
	}
}

/*****************************************************************************************************************/

const record = `source_id, designation, ra, dec, pmra, pmdec, parallax, phot_rp_mean_flux, phot_g_mean_mag`

/*****************************************************************************************************************/

func (g *GAIAServiceClient) Build() (string, error) {
	// Define the ADQL query template for the GAIA TAP service:
	// @see https://gea.esac.esa.int/archive/documentation/GDR2/Gaia_archive/chap_datamodel/
	// N.B. (use only gold standard data, e.g., photometry processing mode (byte) i.e., phot_proc_mode = '0'):
	const queryTemplate = `
		SELECT TOP {{.Limit}} {{.Record}}
		FROM gaiadr2.gaia_source
		WHERE CONTAINS(
			POINT('ICRS', ra, dec),
			CIRCLE('ICRS', {{.RA}}, {{.Dec}}, {{.Radius}})
		) = 1 
		AND phot_g_mean_mag < {{.Threshold}}
		AND phot_rp_mean_flux IS NOT NULL
		ORDER BY phot_g_mean_mag DESC;
	`

	// Parse the ADQL query template:
	tmpl, err := template.New("adql").Parse(queryTemplate)
	if err != nil {
		return "", err
	}

	// Set the template data from the query parameters:
	// Data to populate the template
	data := struct {
		Record    string
		RA        float64
		Dec       float64
		Radius    float64
		Limit     float64
		Threshold float64
	}{
		Record:    record,
		RA:        g.Query.RA,
		Dec:       g.Query.Dec,
		Radius:    g.Query.Radius,
		Limit:     g.Query.Limit,
		Threshold: g.Query.Threshold,
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

func (g *GAIAServiceClient) PerformRadialSearch(eq astrometry.ICRSEquatorialCoordinate, radius float64, limit float64, threshold float64) ([]Source, error) {
	// Set the query parameters:
	g.Query.RA = eq.RA
	g.Query.Dec = eq.Dec
	g.Query.Radius = radius
	g.Query.Limit = limit
	g.Query.Threshold = threshold

	// Construct the ADQL query from the template:
	adqlQuery, err := g.Build()
	if err != nil {
		return nil, err
	}

	// Prepare the POST form data for the HTTP request:
	formData := url.Values{}
	formData.Set("REQUEST", "doQuery")
	formData.Set("LANG", "ADQL")
	formData.Set("FORMAT", "json")
	formData.Set("QUERY", adqlQuery)

	// Create a new POST request with the form data:
	req, err := http.NewRequest("POST", g.URI, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set the headers to handle the form data:
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request using the custom HTTP client:
	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice:
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors and return the response body if not OK:
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GAIA TAP query failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the JSON response into a GAIAResponse struct:
	var gaia GAIAResponse
	err = json.Unmarshal(bodyBytes, &gaia)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	var stars []Source

	for _, record := range gaia.Data {
		// Create a new Source struct from the record:
		star := Source{
			UID:         fmt.Sprintf("%v", record[0]),
			Designation: fmt.Sprintf("%v", record[1]),
			RA:          record[2].(float64),
			Dec:         record[3].(float64),
		}

		if record[4] != nil {
			star.ProperMotionRA = record[4].(float64)
		}

		if record[5] != nil {
			star.ProperMotionDec = record[5].(float64)
		}

		if record[6] != nil {
			star.Parallax = record[6].(float64)
		}

		if record[7] != nil {
			star.PhotometricGMeanFlux = record[7].(float64)
		}

		if record[8] != nil {
			star.PhotometricGMeanMagnitude = record[8].(float64)
		}

		// Append to the stars slice of Source structs:
		stars = append(stars, star)
	}

	return stars, nil
}

/*****************************************************************************************************************/
