/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package utils

/*****************************************************************************************************************/

import (
	"math"
	"testing"

	"github.com/observerly/iris/pkg/fits"
)

/*****************************************************************************************************************/

// newFITSHeaderWithRA is a helper that returns a fits.FITSHeader
// populated with an RA value.
func newFITSHeaderWithRA(ra float32) fits.FITSHeader {
	return fits.FITSHeader{
		Floats: map[string]fits.FITSHeaderFloat{
			"RA": {
				Value:   ra,
				Comment: "Right Ascension",
			},
		},
	}
}

/*****************************************************************************************************************/

// newFITSHeaderNoRA is a helper that returns a fits.FITSHeader
// with no RA entry in the Floats map.
func newFITSHeaderNoRA() fits.FITSHeader {
	return fits.FITSHeader{
		Floats: map[string]fits.FITSHeaderFloat{},
	}
}

/*****************************************************************************************************************/

func TestRAValueIsNotNaNAndWithinRange(t *testing.T) {
	value := float32(180.0)
	header := newFITSHeaderNoRA()

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != value {
		t.Errorf("Expected %f, got %f", value, got)
	}
}

func TestRAValueIsNotNaNAndOutOfRangeNegative(t *testing.T) {
	value := float32(-1.0)
	header := newFITSHeaderNoRA()

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for negative RA, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range RA, got %f", got)
	}
}

func TestRAValueIsNotNaNAndOutOfRangePositive(t *testing.T) {
	value := float32(361.0)
	header := newFITSHeaderNoRA()

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for RA > 360, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range RA, got %f", got)
	}
}

func TestRAValueIsNaNAndRAHeaderFoundAndValid(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithRA(123.45)

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err != nil {
		t.Fatalf("Unexpected error when RA is pulled from header: %v", err)
	}
	expected := float32(123.45)
	if got != expected {
		t.Errorf("Expected RA = %f, got %f", expected, got)
	}
}

func TestRAValueIsNaNAndRAHeaderMissing(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderNoRA()

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error when RA header is missing, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN when RA header is missing, got %f", got)
	}
}

func TestRAValueIsNaNAndRAHeaderOutOfRangeNegative(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithRA(-5.0)

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for negative RA, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range RA, got %f", got)
	}
}

func TestRAValueIsNaNAndRAHeaderOutOfRangePositive(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithRA(400.0)

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for RA > 360, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range RA, got %f", got)
	}
}

func TestRAValueIsNaNAndRAHeaderNaN(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithRA(float32(math.NaN()))

	got, err := ResolveOrExtractRAFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for RA=NaN in header, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN when RA=NaN in header, got %f", got)
	}
}

/*****************************************************************************************************************/
