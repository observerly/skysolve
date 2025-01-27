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

// newFITSHeaderWithDec is a helper that returns a fits.FITSHeader
// populated with an Dec value.
func newFITSHeaderWithDec(dec float32) fits.FITSHeader {
	return fits.FITSHeader{
		Floats: map[string]fits.FITSHeaderFloat{
			"DEC": {
				Value:   dec,
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

// newFITSHeaderNoDec is a helper that returns a fits.FITSHeader
// with no Dec entry in the Floats map.
func newFITSHeaderNoDec() fits.FITSHeader {
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

func TestDecValueIsNotNaNAndWithinRange(t *testing.T) {
	value := float32(45.0) // well within -90..+90
	header := newFITSHeaderNoDec()

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != value {
		t.Errorf("Expected %f, got %f", value, got)
	}
}

func TestDecValueIsNotNaNAndOutOfRangeNegative(t *testing.T) {
	value := float32(-91.0)
	header := newFITSHeaderNoDec()

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for Dec < -90, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range Dec, got %f", got)
	}
}

func TestDecValueIsNotNaNAndOutOfRangePositive(t *testing.T) {
	value := float32(91.0)
	header := newFITSHeaderNoDec()

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for Dec > +90, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range Dec, got %f", got)
	}
}

func TestDecValueIsNaNAndDecHeaderFoundAndValid(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithDec(30.0) // valid Dec

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err != nil {
		t.Fatalf("Unexpected error when Dec is pulled from header: %v", err)
	}
	expected := float32(30.0)
	if got != expected {
		t.Errorf("Expected Dec = %f, got %f", expected, got)
	}
}

func TestDecValueIsNaNAndDecHeaderMissing(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderNoDec()

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error when Dec header is missing, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN when Dec header is missing, got %f", got)
	}
}

func TestDecValueIsNaNAndDecHeaderOutOfRangeNegative(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithDec(-95.0)

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for Dec < -90, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range Dec, got %f", got)
	}
}

func TestDecValueIsNaNAndDecHeaderOutOfRangePositive(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithDec(100.0)

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for Dec > +90, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN for out-of-range Dec, got %f", got)
	}
}

func TestDecValueIsNaNAndDecHeaderNaN(t *testing.T) {
	value := float32(math.NaN())
	header := newFITSHeaderWithDec(float32(math.NaN()))

	got, err := ResolveOrExtractDecFromHeaders(value, header)
	if err == nil {
		t.Fatalf("Expected an error for Dec=NaN in header, but got none")
	}
	if !math.IsNaN(float64(got)) {
		t.Errorf("Expected NaN when Dec=NaN in header, got %f", got)
	}
}

/*****************************************************************************************************************/

func TestWidthIsPresentAndValid(t *testing.T) {
	header := fits.FITSHeader{
		Ints: map[string]fits.FITSHeaderInt{
			"NAXIS1": {
				Value:   1024,
				Comment: "Width of the image",
			},
		},
	}

	got, err := ExtractImageWidthFromHeaders(header)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := int32(1024)
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestWidthIsMissingFromHeader(t *testing.T) {
	header := fits.FITSHeader{
		Ints: map[string]fits.FITSHeaderInt{},
	}

	got, err := ExtractImageWidthFromHeaders(header)
	if err == nil {
		t.Fatalf("Expected an error for missing width, but got none")
	}
	if got != -1 {
		t.Errorf("Expected -1 for missing width, got %d", got)
	}
}

func TestWidthIsOutOfRangeZeroOrNegative(t *testing.T) {
	header := fits.FITSHeader{
		Ints: map[string]fits.FITSHeaderInt{
			"NAXIS1": {
				Value:   -100,
				Comment: "Width of the image",
			},
		},
	}

	got, err := ExtractImageWidthFromHeaders(header)
	if err == nil {
		t.Fatalf("Expected an error for width <= 0, but got none")
	}
	if got != -1 {
		t.Errorf("Expected -1 for invalid width, got %d", got)
	}
}

func TestWidthIsOutOfRangeMaxInt32(t *testing.T) {
	header := fits.FITSHeader{
		Ints: map[string]fits.FITSHeaderInt{
			"NAXIS1": {
				Value:   math.MaxInt32,
				Comment: "Width of the image",
			},
		},
	}

	got, err := ExtractImageWidthFromHeaders(header)
	if err == nil {
		t.Fatalf("Expected an error for width == math.MaxInt32, but got none")
	}
	if got != -1 {
		t.Errorf("Expected -1 for invalid width, got %d", got)
	}
}

/*****************************************************************************************************************/
