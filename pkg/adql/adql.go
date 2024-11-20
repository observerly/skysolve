/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package adql

/*****************************************************************************************************************/

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

/*****************************************************************************************************************/

type TapResponse struct {
	Data [][]interface{} `json:"data"`
}

/*****************************************************************************************************************/

// TapADQLClient represents a client for the TAP ADQL service.
type TapClient struct {
	URI     string
	Client  *http.Client
	Timeout time.Duration
	Headers map[string]string
}

/*****************************************************************************************************************/

// NewTapClient initializes a new generic TAP ADQL client with optional configurations.
func NewTapClient(serviceURL url.URL, timeout time.Duration, headers map[string]string) *TapClient {
	client := &http.Client{
		Timeout: timeout,
	}

	return &TapClient{
		URI:     serviceURL.String(),
		Client:  client,
		Timeout: timeout,
		Headers: headers,
	}
}

/*****************************************************************************************************************/

// BuildADQLQuery constructs an ADQL query using a provided template and data.
func (t *TapClient) BuildADQLQuery(templateStr string, data interface{}) (string, error) {
	// Parse the ADQL template:
	tmpl, err := template.New("adql").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse ADQL template: %w", err)
	}

	// Execute the ADQL template and write the result to a buffer:
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute ADQL template: %w", err)
	}

	// Return the constructed ADQL query:
	return buf.String(), nil
}

/*****************************************************************************************************************/
