/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package fov

/*****************************************************************************************************************/

import (
	"math"
	"testing"
)

/*****************************************************************************************************************/

// Helper function to compare floating-point numbers with a tolerance.
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

/*****************************************************************************************************************/

func TestGetRadialExtentEqualXsYsEqualPixelScales(t *testing.T) {
	xs := 100.0
	ys := 100.0
	pixelScale := PixelScale{X: 0.1, Y: 0.1}
	want := math.Sqrt(2) * 10 // ~14.1421356237
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentXsGreaterThanYsDifferentPixelScales(t *testing.T) {
	xs := 200.0
	ys := 100.0
	pixelScale := PixelScale{X: 0.05, Y: 0.1}
	want := math.Sqrt(2) * 10 // ~14.1421356237
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentXsLessThanYsDifferentPixelScales(t *testing.T) {
	xs := 50.0
	ys := 200.0
	pixelScale := PixelScale{X: 0.2, Y: 0.05}
	want := math.Sqrt(2) * 10 // ~14.1421356237
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentDifferentScaledMinimum(t *testing.T) {
	xs := 100.0
	ys := 100.0
	pixelScale := PixelScale{X: 0.2, Y: 0.1}
	want := math.Sqrt(2) * 10 // min(20, 10) * sqrt(2) = ~14.1421356237
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentZeroXs(t *testing.T) {
	xs := 0.0
	ys := 100.0
	pixelScale := PixelScale{X: 0.1, Y: 0.1}
	want := 0.0 // min(0, 10) * sqrt(2) = 0
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentZeroYs(t *testing.T) {
	xs := 100.0
	ys := 0.0
	pixelScale := PixelScale{X: 0.1, Y: 0.1}
	want := 0.0 // min(10, 0) * sqrt(2) = 0
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentZeroPixelScaleX(t *testing.T) {
	xs := 100.0
	ys := 100.0
	pixelScale := PixelScale{X: 0.0, Y: 0.1}
	want := 0.0 // min(0, 10) * sqrt(2) = 0
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentZeroPixelScaleY(t *testing.T) {
	xs := 100.0
	ys := 100.0
	pixelScale := PixelScale{X: 0.1, Y: 0.0}
	want := 0.0 // min(10, 0) * sqrt(2) = 0
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-6
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, want)
	}
}

func TestGetRadialExtentNegativeXs(t *testing.T) {
	xs := -100.0
	ys := 100.0
	pixelScale := PixelScale{X: 0.1, Y: 0.1}

	got := GetRadialExtent(xs, ys, pixelScale)
	// Depending on intended behavior, you might want to take absolute values
	// Here, we'll follow the original function's logic
	tolerance := 1e-6
	expected := math.Sqrt(2) * 10 // ~14.1421356237
	if !floatEquals(got, expected, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, expected)
	}
}

func TestGetRadialExtentNegativeYs(t *testing.T) {
	xs := 100.0
	ys := -100.0
	pixelScale := PixelScale{X: 0.1, Y: 0.1}

	got := GetRadialExtent(xs, ys, pixelScale)
	// Depending on intended behavior, you might want to take absolute values
	// Here, we'll follow the original function's logic
	tolerance := 1e-6
	expected := math.Sqrt(2) * 10 // ~14.1421356237
	if !floatEquals(got, expected, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, expected)
	}
}

func TestGetRadialExtentBothXsAndYsNegative(t *testing.T) {
	xs := -100.0
	ys := -100.0
	pixelScale := PixelScale{X: 0.1, Y: 0.1}

	got := GetRadialExtent(xs, ys, pixelScale)
	// Depending on intended behavior, you might want to take absolute values
	// Here, we'll follow the original function's logic
	tolerance := 1e-6
	expected := math.Sqrt(2) * 10 // ~14.1421356237
	if !floatEquals(got, expected, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want %v", xs, ys, pixelScale, got, expected)
	}
}

func TestGetRadialExtentFloatingPointPrecision(t *testing.T) {
	xs := 123.456
	ys := 78.912
	pixelScale := PixelScale{X: 0.345, Y: 0.678}
	// Calculate min(pixelScale.X * xs, pixelScale.Y * ys)
	minScaled := math.Min(0.345*123.456, 0.678*78.912) // min(42.59232, 53.526) = 42.59232
	want := math.Sqrt(2) * minScaled                   // ~60.2463
	got := GetRadialExtent(xs, ys, pixelScale)
	tolerance := 1e-4
	if !floatEquals(got, want, tolerance) {
		t.Errorf("GetRadialExtent(%v, %v, %+v) = %v; want ~%v", xs, ys, pixelScale, got, want)
	}
}

/*****************************************************************************************************************/
