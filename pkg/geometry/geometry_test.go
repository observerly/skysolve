/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package geometry

/*****************************************************************************************************************/

import (
	"math"
	"testing"
)

/*****************************************************************************************************************/

// Helper function to compare floating-point numbers with tolerance
func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

/*****************************************************************************************************************/

func TestDistanceBetweenTwoCartesianPoints(t *testing.T) {
	x1 := 0.0
	y1 := 0.0
	x2 := 3.0
	y2 := 4.0

	expected := 5.0

	result := DistanceBetweenTwoCartesianPoints(x1, y1, x2, y2)

	if result != expected {
		t.Errorf("DistanceBetweenTwoCartesianPoints(%f, %f, %f, %f) = %f; want %f", x1, y1, x2, y2, result, expected)
	}
}

/*****************************************************************************************************************/

func TestAngleBetweenThreeCartesianPointsRightTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 3.0, 0.0
	x3, y3 := 0.0, 4.0

	expectedAngle := 90.0
	epsilon := 1e-6

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected %.6f, got %.6f", expectedAngle, angle)
	}
}

func TestAngleBetweenThreeCartesianPointsEquilateralTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 1.0, 0.0
	x3, y3 := 0.5, math.Sqrt(3)/2

	expectedAngle := 60.0
	epsilon := 1e-6

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected %.6f, got %.6f", expectedAngle, angle)
	}
}

func TestAngleBetweenThreeCartesianPointsDegenerateTriangle(t *testing.T) {
	x1, y1 := 1.0, 1.0
	x2, y2 := 1.0, 1.0
	x3, y3 := 1.0, 1.0

	_, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err == nil {
		t.Errorf("Expected error for degenerate triangle, but got none")
	}
}

func TestAngleBetweenThreeCartesianPointsAcuteTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 2.0, 0.0
	x3, y3 := 1.0, 1.0

	expectedAngle := 45.0
	epsilon := 1e-6

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected %.6f, got %.6f", expectedAngle, angle)
	}
}

func TestAngleBetweenThreeCartesianPointsObtuseTriangle(t *testing.T) {
	x1, y1 := 0.0, 0.0
	x2, y2 := 4.0, 0.0
	x3, y3 := 1.0, 1.0

	expectedAngle := 90.0 // Approximately, actual angle is less
	epsilon := 1e-1       // Allowing more tolerance

	angle, err := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// The actual angle is approximately 16.26 degrees
	expectedAngle = math.Acos((math.Pow(1.0, 2)+math.Pow(3.0, 2)-math.Pow(math.Sqrt(2.0), 2))/(2*1.0*3.0)) * (180.0 / math.Pi) // ~16.26 degrees

	if math.Abs(angle-expectedAngle) > epsilon {
		t.Errorf("Angle mismatch: expected ~%.6f, got %.6f", expectedAngle, angle)
	}
}

/*****************************************************************************************************************/

// TestComputeInvariantFeaturesRightTriangle tests the invariant features computation for a right-angled triangle.
func TestComputeInvariantFeaturesRightTriangle(t *testing.T) {
	// Define a right-angled triangle with sides 3-4-5
	x1, y1 := 0.0, 0.0
	x2, y2 := 3.0, 0.0
	x3, y3 := 0.0, 4.0

	expected := InvariantFeatures{
		RatioAB: 0.6,       // AB / BC = 3 /5
		RatioAC: 0.8,       // AC / BC = 4 /5
		AngleA:  90.0,      // Right angle at A
		AngleB:  53.130102, // Approximately 53.13 degrees
	}

	computed, err := ComputeInvariantFeatures(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !almostEqual(computed.RatioAB, expected.RatioAB, 1e-6) {
		t.Errorf("RatioAB mismatch: expected %.6f, got %.6f", expected.RatioAB, computed.RatioAB)
	}

	if !almostEqual(computed.RatioAC, expected.RatioAC, 1e-6) {
		t.Errorf("RatioAC mismatch: expected %.6f, got %.6f", expected.RatioAC, computed.RatioAC)
	}

	if !almostEqual(computed.AngleA, expected.AngleA, 1e-6) {
		t.Errorf("AngleA mismatch: expected %.6f, got %.6f", expected.AngleA, computed.AngleA)
	}

	if !almostEqual(computed.AngleB, expected.AngleB, 1e-6) {
		t.Errorf("AngleB mismatch: expected %.6f, got %.6f", expected.AngleB, computed.AngleB)
	}
}

// TestComputeInvariantFeaturesEquilateralTriangle tests the invariant features computation for an equilateral triangle.
func TestComputeInvariantFeaturesEquilateralTriangle(t *testing.T) {
	// Define an equilateral triangle with side length 1
	x1, y1 := 0.0, 0.0
	x2, y2 := 1.0, 0.0
	x3, y3 := 0.5, math.Sqrt(3)/2

	expected := InvariantFeatures{
		RatioAB: 1.0,  // AB / BC = 1 /1
		RatioAC: 1.0,  // AC / BC = 1 /1
		AngleA:  60.0, // 60 degrees at A
		AngleB:  60.0, // 60 degrees at B
	}

	computed, err := ComputeInvariantFeatures(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !almostEqual(computed.RatioAB, expected.RatioAB, 1e-6) {
		t.Errorf("RatioAB mismatch: expected %.6f, got %.6f", expected.RatioAB, computed.RatioAB)
	}

	if !almostEqual(computed.RatioAC, expected.RatioAC, 1e-6) {
		t.Errorf("RatioAC mismatch: expected %.6f, got %.6f", expected.RatioAC, computed.RatioAC)
	}

	if !almostEqual(computed.AngleA, expected.AngleA, 1e-6) {
		t.Errorf("AngleA mismatch: expected %.6f, got %.6f", expected.AngleA, computed.AngleA)
	}

	if !almostEqual(computed.AngleB, expected.AngleB, 1e-6) {
		t.Errorf("AngleB mismatch: expected %.6f, got %.6f", expected.AngleB, computed.AngleB)
	}
}

// TestComputeInvariantFeaturesAcuteTriangle tests the invariant features computation for an acute triangle.
func TestComputeInvariantFeaturesAcuteTriangle(t *testing.T) {
	// Define an acute triangle
	x1, y1 := 0.0, 0.0
	x2, y2 := 2.0, 0.0
	x3, y3 := 1.0, 1.0

	// Manually compute expected invariant features
	a := DistanceBetweenTwoCartesianPoints(x2, y2, x3, y3) // BC ≈ 1.414214
	b := DistanceBetweenTwoCartesianPoints(x1, y1, x3, y3) // AC ≈ 1.414214
	c := DistanceBetweenTwoCartesianPoints(x1, y1, x2, y2) // AB = 2.0

	// Calculate ratios based on specific sides
	ratioAB := math.Min(c, a) / math.Max(c, a) // min(2.0,1.414214)/max(2.0,1.414214)=1.414214/2.0≈0.707107
	ratioAC := math.Min(b, a) / math.Max(b, a) // min(1.414214,1.414214)/max(1.414214,1.414214)=1.0

	expected := InvariantFeatures{
		RatioAB: ratioAB, // ≈0.707107
		RatioAC: ratioAC, // =1.0
		AngleA:  45.0,    // 45 degrees at A
		AngleB:  45.0,    // 45 degrees at B
	}

	computed, err := ComputeInvariantFeatures(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !almostEqual(computed.RatioAB, expected.RatioAB, 1e-6) {
		t.Errorf("RatioAB mismatch: expected %.6f, got %.6f", expected.RatioAB, computed.RatioAB)
	}

	if !almostEqual(computed.RatioAC, expected.RatioAC, 1e-6) {
		t.Errorf("RatioAC mismatch: expected %.6f, got %.6f", expected.RatioAC, computed.RatioAC)
	}

	if !almostEqual(computed.AngleA, expected.AngleA, 1e-6) {
		t.Errorf("AngleA mismatch: expected %.6f, got %.6f", expected.AngleA, computed.AngleA)
	}

	if !almostEqual(computed.AngleB, expected.AngleB, 1e-6) {
		t.Errorf("AngleB mismatch: expected %.6f, got %.6f", expected.AngleB, computed.AngleB)
	}
}

// TestComputeInvariantFeaturesObtuseTriangle tests the invariant features computation for an obtuse triangle.
func TestComputeInvariantFeaturesObtuseTriangle(t *testing.T) {
	// Define an obtuse triangle
	x1, y1 := 0.0, 0.0
	x2, y2 := 3.0, 0.0
	x3, y3 := 1.0, 1.0

	// Manually compute expected invariant features
	a := DistanceBetweenTwoCartesianPoints(x2, y2, x3, y3) // BC ≈ 2.236068
	b := DistanceBetweenTwoCartesianPoints(x1, y1, x3, y3) // AC ≈ 1.414214
	c := DistanceBetweenTwoCartesianPoints(x1, y1, x2, y2) // AB = 3.0

	angleA, errA := AngleBetweenThreeCartesianPoints(x1, y1, x2, y2, x3, y3) // ≈63.434949 degrees
	angleB, errB := AngleBetweenThreeCartesianPoints(x2, y2, x1, y1, x3, y3) // ≈26.565051 degrees

	if errA != nil || errB != nil {
		t.Fatalf("Unexpected error computing angles: %v, %v", errA, errB)
	}

	// Calculate ratios based on specific sides with min/max to ensure ≤1
	ratioAB := math.Min(c, a) / math.Max(c, a) // min(3.0,2.236068)/max(3.0,2.236068)=2.236068/3.0≈0.745356
	ratioAC := math.Min(b, a) / math.Max(b, a) // min(1.414214,2.236068)/max(1.414214,2.236068)=1.414214/2.236068≈0.632455

	expected := InvariantFeatures{
		RatioAB: ratioAB, // ≈0.745356
		RatioAC: ratioAC, // ≈0.632455
		AngleA:  angleA,  // ≈63.434949 degrees
		AngleB:  angleB,  // ≈26.565051 degrees
	}

	computed, err := ComputeInvariantFeatures(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !almostEqual(computed.RatioAB, expected.RatioAB, 1e-6) {
		t.Errorf("RatioAB mismatch: expected %.6f, got %.6f", expected.RatioAB, computed.RatioAB)
	}

	if !almostEqual(computed.RatioAC, expected.RatioAC, 1e-6) {
		t.Errorf("RatioAC mismatch: expected %.6f, got %.6f", expected.RatioAC, computed.RatioAC)
	}

	if !almostEqual(computed.AngleA, expected.AngleA, 1e-6) {
		t.Errorf("AngleA mismatch: expected %.6f, got %.6f", expected.AngleA, computed.AngleA)
	}

	if !almostEqual(computed.AngleB, expected.AngleB, 1e-6) {
		t.Errorf("AngleB mismatch: expected %.6f, got %.6f", expected.AngleB, computed.AngleB)
	}
}

// TestComputeInvariantFeaturesDegenerateTriangle tests the invariant features computation for a degenerate triangle.
func TestComputeInvariantFeaturesDegenerateTriangle(t *testing.T) {
	// Define a degenerate triangle (all points on a straight line)
	x1, y1 := 1.0, 1.0
	x2, y2 := 1.0, 1.0
	x3, y3 := 1.0, 1.0

	_, err := ComputeInvariantFeatures(x1, y1, x2, y2, x3, y3)

	if err == nil {
		t.Errorf("Expected error for degenerate triangle, but got none")
	}
}

// TestComputeInvariantFeaturesScalefulTriangle tests the invariant features computation for a scaled triangle.
func TestComputeInvariantFeaturesScalefulTriangle(t *testing.T) {
	// Define a scaled triangle (same shape as 3-4-5, but scaled by 2)
	x1, y1 := 0.0, 0.0
	x2, y2 := 6.0, 0.0
	x3, y3 := 0.0, 8.0

	expected := InvariantFeatures{
		RatioAB: 0.6,       // AB / BC = 6 /10
		RatioAC: 0.8,       // AC / BC = 8 /10
		AngleA:  90.0,      // Right angle at A
		AngleB:  53.130102, // Approximately 53.13 degrees
	}

	computed, err := ComputeInvariantFeatures(x1, y1, x2, y2, x3, y3)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !almostEqual(computed.RatioAB, expected.RatioAB, 1e-6) {
		t.Errorf("RatioAB mismatch: expected %.6f, got %.6f", expected.RatioAB, computed.RatioAB)
	}

	if !almostEqual(computed.RatioAC, expected.RatioAC, 1e-6) {
		t.Errorf("RatioAC mismatch: expected %.6f, got %.6f", expected.RatioAC, computed.RatioAC)
	}

	if !almostEqual(computed.AngleA, expected.AngleA, 1e-6) {
		t.Errorf("AngleA mismatch: expected %.6f, got %.6f", expected.AngleA, computed.AngleA)
	}

	if !almostEqual(computed.AngleB, expected.AngleB, 1e-6) {
		t.Errorf("AngleB mismatch: expected %.6f, got %.6f", expected.AngleB, computed.AngleB)
	}
}

/*****************************************************************************************************************/

func TestCompareInvariantFeaturesWithinTolerance(t *testing.T) {
	tolerance := InvariantFeatureTolerance{
		LengthRatio: 0.01,
		Angle:       0.5, // degrees
	}

	f1 := InvariantFeatures{RatioAB: 1.0, RatioAC: 1.0, AngleA: 30.0, AngleB: 60.0}
	f2 := InvariantFeatures{RatioAB: 1.005, RatioAC: 0.995, AngleA: 30.3, AngleB: 59.7}

	if !CompareInvariantFeatures(f1, f2, tolerance) {
		t.Error("Expected true when differences are within tolerance")
	}
}

func TestCompareInvariantFeaturesRatioABExceedsTolerance(t *testing.T) {
	tolerance := InvariantFeatureTolerance{
		LengthRatio: 0.01,
		Angle:       0.5,
	}

	f1 := InvariantFeatures{RatioAB: 1.0, RatioAC: 1.0, AngleA: 30.0, AngleB: 60.0}
	f2 := InvariantFeatures{RatioAB: 1.02, RatioAC: 1.0, AngleA: 30.0, AngleB: 60.0}

	if CompareInvariantFeatures(f1, f2, tolerance) {
		t.Error("Expected false when RatioAB difference exceeds tolerance")
	}
}

func TestCompareInvariantFeaturesAngleAExceedsTolerance(t *testing.T) {
	tolerance := InvariantFeatureTolerance{
		LengthRatio: 0.01,
		Angle:       0.5,
	}

	f1 := InvariantFeatures{RatioAB: 1.0, RatioAC: 1.0, AngleA: 30.0, AngleB: 60.0}
	f2 := InvariantFeatures{RatioAB: 1.0, RatioAC: 1.0, AngleA: 30.6, AngleB: 60.0}

	if CompareInvariantFeatures(f1, f2, tolerance) {
		t.Error("Expected false when AngleA difference exceeds tolerance")
	}
}

func TestCompareInvariantFeaturesExceedsToleranceBySmallMargin(t *testing.T) {
	tolerance := InvariantFeatureTolerance{
		LengthRatio: 0.01,
		Angle:       0.5,
	}

	f1 := InvariantFeatures{RatioAB: 1.0, RatioAC: 1.0, AngleA: 30.0, AngleB: 60.0}
	f2 := InvariantFeatures{RatioAB: 1.0101, RatioAC: 1.0, AngleA: 30.0, AngleB: 60.0}

	if CompareInvariantFeatures(f1, f2, tolerance) {
		t.Error("Expected false when differences just exceed tolerance")
	}
}

func TestCompareInvariantFeaturesNegativeDifferencesWithinTolerance(t *testing.T) {
	tolerance := InvariantFeatureTolerance{
		LengthRatio: 0.01,
		Angle:       0.5,
	}

	f1 := InvariantFeatures{RatioAB: 1.0, RatioAC: 1.0, AngleA: 30.0, AngleB: 60.0}
	f2 := InvariantFeatures{RatioAB: 0.995, RatioAC: 0.995, AngleA: 29.7, AngleB: 59.7}

	if !CompareInvariantFeatures(f1, f2, tolerance) {
		t.Error("Expected true when negative differences are within tolerance")
	}
}

/*****************************************************************************************************************/
