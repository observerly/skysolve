/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package adql

/*****************************************************************************************************************/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// ExecuteQuery sends the ADQL query to the TAP service and returns the parsed response.
func (t *TapClient) ExecuteADQLQuery(adqlQuery string) (*TapResponse, error) {
	formData := url.Values{}
	formData.Set("REQUEST", "doQuery")
	formData.Set("LANG", "ADQL")
	formData.Set("FORMAT", "json")
	formData.Set("QUERY", adqlQuery)

	req, err := http.NewRequest("POST", t.URI, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set the content type to form encoded data:
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the content length:
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(formData.Encode())))

	// Set any additional headers:
	for key, value := range t.Headers {
		req.Header.Set(key, value)
	}

	// Perform the HTTP request:
	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body:
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check the response status code:
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TAP query failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the JSON response:
	var tapResp TapResponse
	err = json.Unmarshal(bodyBytes, &tapResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Return the parsed response:
	return &tapResp, nil
}

/*****************************************************************************************************************/
