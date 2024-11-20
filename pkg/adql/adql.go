/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package adql

/*****************************************************************************************************************/

import (
	"net/http"
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
