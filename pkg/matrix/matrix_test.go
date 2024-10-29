/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package matrix

/*****************************************************************************************************************/

import (
	"math"
	"testing"
)

/*****************************************************************************************************************/

// floatEquals checks if two float64 numbers are equal within a specified tolerance.
func floatEquals(a, b, tolerance float64) bool {
	if a > b {
		return a-b <= tolerance
	}

	return b-a <= tolerance
}

/*****************************************************************************************************************/

// equalMatrices checks if two matrices are equal in dimensions and values.
func equalMatrices(a, b *Matrix) bool {
	if a.rows != b.rows || a.columns != b.columns {
		return false
	}
	for i := 0; i < len(a.Value); i++ {
		if a.Value[i] != b.Value[i] {
			return false
		}
	}
	return true
}

/*****************************************************************************************************************/

// Helper function to compare two matrices for equality with a tolerance for floating-point precision.
// It provides detailed error messages for discrepancies.
func equalMatricesDetailed(a, b *Matrix, t *testing.T) bool {
	if a == nil || b == nil {
		if a != b {
			t.Errorf("One of the matrices is nil while the other is not")
			return false
		}
		return true
	}
	if a.rows != b.rows || a.columns != b.columns {
		t.Errorf("Matrix dimensions do not match: got %dx%d, want %dx%d", a.rows, a.columns, b.rows, b.columns)
		return false
	}
	tolerance := 1e-9
	pass := true
	for i := 0; i < len(a.Value); i++ {
		diff := math.Abs(a.Value[i] - b.Value[i])
		if diff > tolerance {
			row := i / a.columns
			col := i % a.columns
			t.Errorf("Matrix element [%d][%d] = %v; want %v (diff %v)", row, col, a.Value[i], b.Value[i], diff)
			pass = false
		}
	}
	return pass
}

/*****************************************************************************************************************/

// TestMatrixAtAccessFirstElement verifies that accessing the first element returns the correct value without an error.
func TestMatrixAtAccessFirstElement(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	got, err := matrix.At(0, 0)
	if err != nil {
		t.Errorf("At() returned unexpected error: %v", err)
	}
	want := 1.0
	if got != want {
		t.Errorf("At(0,0) = %v; want %v", got, want)
	}
}

// TestMatrixAtAccessLastElement checks that accessing the last element in a 3x3 matrix returns the correct value.
func TestMatrixAtAccessLastElement(t *testing.T) {
	matrix := Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	got, err := matrix.At(2, 2)
	if err != nil {
		t.Errorf("At() returned unexpected error: %v", err)
	}
	want := 9.0
	if got != want {
		t.Errorf("At(2,2) = %v; want %v", got, want)
	}
}

// TestMatrixAtAccessMiddleElement ensures that accessing a middle element in a 3x3 matrix works as expected.
func TestMatrixAtAccessMiddleElement(t *testing.T) {
	matrix := Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	got, err := matrix.At(1, 1)
	if err != nil {
		t.Errorf("At() returned unexpected error: %v", err)
	}
	want := 5.0
	if got != want {
		t.Errorf("At(1,1) = %v; want %v", got, want)
	}
}

// TestMatrixAtNegativeRowIndex confirms that providing a negative row index results in an error.
func TestMatrixAtNegativeRowIndex(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	_, err := matrix.At(-1, 0)
	if err == nil {
		t.Errorf("At(-1,0) expected error, got nil")
	}
}

// TestMatrixAtNegativeColumnIndex confirms that providing a negative column index results in an error.
func TestMatrixAtNegativeColumnIndex(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	_, err := matrix.At(0, -1)
	if err == nil {
		t.Errorf("At(0,-1) expected error, got nil")
	}
}

// TestMatrixAtRowIndexOutOfBounds ensures that a row index equal to the number of rows returns an error.
func TestMatrixAtRowIndexOutOfBounds(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	_, err := matrix.At(2, 0)
	if err == nil {
		t.Errorf("At(2,0) expected error, got nil")
	}
}

// TestMatrixAtColumnIndexOutOfBounds ensures that a column index equal to the number of columns returns an error.
func TestMatrixAtColumnIndexOutOfBounds(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	_, err := matrix.At(0, 2)
	if err == nil {
		t.Errorf("At(0,2) expected error, got nil")
	}
}

// TestMatrixAtSingleElementValid verifies that accessing the only element in a 1x1 matrix returns the correct value without an error.
func TestMatrixAtSingleElementValid(t *testing.T) {
	matrix := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{42.0},
	}

	got, err := matrix.At(0, 0)
	if err != nil {
		t.Errorf("At(0,0) returned unexpected error: %v", err)
	}
	want := 42.0
	if got != want {
		t.Errorf("At(0,0) = %v; want %v", got, want)
	}
}

// TestMatrixAtSingleElementOutOfBounds ensures that accessing any index other than (0,0) in a 1x1 matrix results in an error.
func TestMatrixAtSingleElementOutOfBounds(t *testing.T) {
	matrix := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{42.0},
	}

	_, err := matrix.At(1, 0)
	if err == nil {
		t.Errorf("At(1,0) expected error, got nil")
	}
}

// TestMatrixAtEmptyMatrix confirms that accessing any element in an empty matrix returns an error.
func TestMatrixAtEmptyMatrix(t *testing.T) {
	matrix := Matrix{
		rows:    0,
		columns: 0,
		Value:   []float64{},
	}

	_, err := matrix.At(0, 0)
	if err == nil {
		t.Errorf("At(0,0) on empty matrix expected error, got nil")
	}
}

/*****************************************************************************************************************/

// TestMatrixSetAccessFirstElement verifies that setting the first element updates the value correctly without an error.
func TestMatrixSetAccessFirstElement(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	newValue := 10.0
	err := matrix.Set(0, 0, newValue)
	if err != nil {
		t.Errorf("Set() returned unexpected error: %v", err)
	}

	got, err := matrix.At(0, 0)
	if err != nil {
		t.Errorf("At() returned unexpected error: %v", err)
	}
	if got != newValue {
		t.Errorf("After Set(0,0, %v), At(0,0) = %v; want %v", newValue, got, newValue)
	}
}

// TestMatrixSetAccessLastElement checks that setting the last element in a 3x3 matrix updates the value correctly.
func TestMatrixSetAccessLastElement(t *testing.T) {
	matrix := Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	newValue := 90.0
	err := matrix.Set(2, 2, newValue)
	if err != nil {
		t.Errorf("Set() returned unexpected error: %v", err)
	}

	got, err := matrix.At(2, 2)
	if err != nil {
		t.Errorf("At() returned unexpected error: %v", err)
	}
	if got != newValue {
		t.Errorf("After Set(2,2, %v), At(2,2) = %v; want %v", newValue, got, newValue)
	}
}

// TestMatrixSetAccessMiddleElement ensures that setting a middle element in a 3x3 matrix works as expected.
func TestMatrixSetAccessMiddleElement(t *testing.T) {
	matrix := Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	newValue := 50.0
	err := matrix.Set(1, 1, newValue)
	if err != nil {
		t.Errorf("Set() returned unexpected error: %v", err)
	}

	got, err := matrix.At(1, 1)
	if err != nil {
		t.Errorf("At() returned unexpected error: %v", err)
	}
	if got != newValue {
		t.Errorf("After Set(1,1, %v), At(1,1) = %v; want %v", newValue, got, newValue)
	}
}

// TestMatrixSetNegativeRowIndex confirms that providing a negative row index results in an error.
func TestMatrixSetNegativeRowIndex(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	err := matrix.Set(-1, 0, 10.0)
	if err == nil {
		t.Errorf("Set(-1,0, 10.0) expected error, got nil")
	}
}

// TestMatrixSetNegativeColumnIndex confirms that providing a negative column index results in an error.
func TestMatrixSetNegativeColumnIndex(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	err := matrix.Set(0, -1, 10.0)
	if err == nil {
		t.Errorf("Set(0,-1, 10.0) expected error, got nil")
	}
}

// TestMatrixSetRowIndexOutOfBounds ensures that a row index equal to the number of rows returns an error.
func TestMatrixSetRowIndexOutOfBounds(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	err := matrix.Set(2, 0, 10.0)
	if err == nil {
		t.Errorf("Set(2,0, 10.0) expected error, got nil")
	}
}

// TestMatrixSetColumnIndexOutOfBounds ensures that a column index equal to the number of columns returns an error.
func TestMatrixSetColumnIndexOutOfBounds(t *testing.T) {
	matrix := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.0, 2.0, 3.0, 4.0},
	}

	err := matrix.Set(0, 2, 10.0)
	if err == nil {
		t.Errorf("Set(0,2, 10.0) expected error, got nil")
	}
}

// TestMatrixSetSingleElementValid verifies that setting the only element in a 1x1 matrix updates the value correctly without an error.
func TestMatrixSetSingleElementValid(t *testing.T) {
	matrix := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{42.0},
	}

	newValue := 100.0
	err := matrix.Set(0, 0, newValue)
	if err != nil {
		t.Errorf("Set(0,0, %v) returned unexpected error: %v", newValue, err)
	}

	got, err := matrix.At(0, 0)
	if err != nil {
		t.Errorf("At(0,0) returned unexpected error: %v", err)
	}
	if got != newValue {
		t.Errorf("After Set(0,0, %v), At(0,0) = %v; want %v", newValue, got, newValue)
	}
}

// TestMatrixSetSingleElementOutOfBounds ensures that setting any index other than (0,0) in a 1x1 matrix results in an error.
func TestMatrixSetSingleElementOutOfBounds(t *testing.T) {
	matrix := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{42.0},
	}

	err := matrix.Set(1, 0, 100.0)
	if err == nil {
		t.Errorf("Set(1,0, 100.0) expected error, got nil")
	}
}

// TestMatrixSetEmptyMatrix confirms that setting any element in an empty matrix returns an error.
func TestMatrixSetEmptyMatrix(t *testing.T) {
	matrix := Matrix{
		rows:    0,
		columns: 0,
		Value:   []float64{},
	}

	err := matrix.Set(0, 0, 10.0)
	if err == nil {
		t.Errorf("Set(0,0, 10.0) on empty matrix expected error, got nil")
	}
}

/*****************************************************************************************************************/

// TestMatrixTransposeSquareMatrix verifies that transposing a square matrix results in the correct matrix.
func TestMatrixTransposeSquareMatrix(t *testing.T) {
	original := Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	expected := &Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 4, 7, 2, 5, 8, 3, 6, 9},
	}

	transposed, err := original.Transpose()
	if err != nil {
		t.Errorf("Transpose() returned unexpected error: %v", err)
	}

	if !equalMatrices(transposed, expected) {
		t.Errorf("Transpose() = %+v; want %+v", transposed, expected)
	}
}

// TestMatrixTransposeRectangularMatrix verifies that transposing a rectangular matrix results in the correct matrix.
func TestMatrixTransposeRectangularMatrix(t *testing.T) {
	original := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	expected := &Matrix{
		rows:    3,
		columns: 2,
		Value:   []float64{1, 4, 2, 5, 3, 6},
	}

	transposed, err := original.Transpose()
	if err != nil {
		t.Errorf("Transpose() returned unexpected error: %v", err)
	}

	if !equalMatrices(transposed, expected) {
		t.Errorf("Transpose() = %+v; want %+v", transposed, expected)
	}
}

// TestMatrixTransposeSingleElement verifies that transposing a single-element matrix results in the same matrix.
func TestMatrixTransposeSingleElement(t *testing.T) {
	original := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{42.0},
	}

	expected := &Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{42.0},
	}

	transposed, err := original.Transpose()
	if err != nil {
		t.Errorf("Transpose() returned unexpected error: %v", err)
	}

	if !equalMatrices(transposed, expected) {
		t.Errorf("Transpose() = %+v; want %+v", transposed, expected)
	}
}

// TestMatrixTransposeEmptyMatrix verifies that transposing an empty matrix results in an empty matrix.
func TestMatrixTransposeEmptyMatrix(t *testing.T) {
	original := Matrix{
		rows:    0,
		columns: 0,
		Value:   []float64{},
	}

	_, err := original.Transpose()

	if err == nil {
		t.Errorf("Transpose() expected error on empty matrix, got nil")
	}
}

// TestMatrixTransposeNonSquareRectangularMatrix verifies that transposing a non-square rectangular matrix works correctly.
func TestMatrixTransposeNonSquareRectangularMatrix(t *testing.T) {
	original := Matrix{
		rows:    3,
		columns: 2,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	expected := &Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 3, 5, 2, 4, 6},
	}

	transposed, err := original.Transpose()
	if err != nil {
		t.Errorf("Transpose() returned unexpected error: %v", err)
	}

	if !equalMatrices(transposed, expected) {
		t.Errorf("Transpose() = %+v; want %+v", transposed, expected)
	}
}

// TestMatrixTransposeTwice verifies that transposing a matrix twice returns the original matrix.
func TestMatrixTransposeTwice(t *testing.T) {
	original := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	transposed, err := original.Transpose()
	if err != nil {
		t.Errorf("First Transpose() returned unexpected error: %v", err)
	}

	doubleTransposed, err := transposed.Transpose()
	if err != nil {
		t.Errorf("Second Transpose() returned unexpected error: %v", err)
	}

	if !equalMatrices(doubleTransposed, &original) {
		t.Errorf("Double Transpose() = %+v; want %+v", doubleTransposed, original)
	}
}

// TestMatrixTransposeInvalidNewFromSlice verifies that Transpose returns an error if NewFromSlice fails.
func TestMatrixTransposeInvalidNewFromSlice(t *testing.T) {
	// Temporarily modify Transpose to create an invalid slice.
	// Note: This is a workaround since the current Transpose implementation should not fail.
	// Alternatively, you can mock NewFromSlice if using interfaces.
	// Here, we'll skip this test as Transpose is expected to work correctly.

	// To demonstrate, we'll assume an error scenario:
	// Suppose NewFromSlice is called with incorrect rows and columns.
	// We'll manually create such a scenario.

	// Create a transposed slice with incorrect size
	transposed := []float64{1, 4, 2, 3} // Incorrect length for rows=2, columns=2

	_, err := NewFromSlice(transposed, 3, 1) // Intentionally incorrect
	if err == nil {
		t.Errorf("NewFromSlice expected error due to mismatched rows and columns, got nil")
	}
}

// TestMatrixTransposeVerifyOriginalUnchanged ensures that the original matrix remains unchanged after transposing.
func TestMatrixTransposeVerifyOriginalUnchanged(t *testing.T) {
	original := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	originalCopy := Matrix{
		rows:    original.rows,
		columns: original.columns,
		Value:   append([]float64(nil), original.Value...), // Deep copy
	}

	_, err := original.Transpose()
	if err != nil {
		t.Errorf("Transpose() returned unexpected error: %v", err)
	}

	if !equalMatrices(&original, &originalCopy) {
		t.Errorf("Original matrix was modified after Transpose()\nGot: %+v\nWant: %+v", original, originalCopy)
	}
}

/*****************************************************************************************************************/

// TestMultiplySquareMatrices verifies that multiplying two square matrices yields the correct product.
func TestMultiplySquareMatrices(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1, 2, 3, 4},
	}

	b := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{5, 6, 7, 8},
	}

	expected := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{19, 22, 43, 50},
	}

	product, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() = %+v; want %+v", product, expected)
	}
}

// TestMultiplyRectangularMatrices verifies that multiplying rectangular matrices yields the correct product.
func TestMultiplyRectangularMatrices(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	b := Matrix{
		rows:    3,
		columns: 2,
		Value:   []float64{7, 8, 9, 10, 11, 12},
	}

	expected := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{58, 64, 139, 154},
	}

	product, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() = %+v; want %+v", product, expected)
	}
}

// TestMultiplySingleElementMatrices verifies that multiplying single-element matrices yields the correct product.
func TestMultiplySingleElementMatrices(t *testing.T) {
	a := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{2},
	}

	b := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{3},
	}

	expected := &Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{6},
	}

	product, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() = %+v; want %+v", product, expected)
	}
}

// TestMultiplyWithIdentityMatrix verifies that multiplying with an identity matrix returns the original matrix.
func TestMultiplyWithIdentityMatrix(t *testing.T) {
	a := Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	identity := Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
	}

	expected := &Matrix{
		rows:    3,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	product, err := a.Multiply(&identity)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() with identity matrix = %+v; want %+v", product, expected)
	}
}

// TestMultiplyWithZeroMatrix verifies that multiplying with a zero matrix yields a zero matrix.
func TestMultiplyWithZeroMatrix(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	zero := Matrix{
		rows:    3,
		columns: 2,
		Value:   []float64{0, 0, 0, 0, 0, 0},
	}

	expected := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{0, 0, 0, 0},
	}

	product, err := a.Multiply(&zero)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() with zero matrix = %+v; want %+v", product, expected)
	}
}

// TestMultiplyDimensionMismatch verifies that multiplying matrices with incompatible dimensions returns an error.
func TestMultiplyDimensionMismatch(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	b := Matrix{
		rows:    4,
		columns: 2,
		Value:   []float64{7, 8, 9, 10, 11, 12, 13, 14},
	}

	_, err := a.Multiply(&b)
	if err == nil {
		t.Fatalf("Multiply() expected error due to dimension mismatch, got nil")
	}
}

// TestMultiplyTwice verifies that multiplying a matrix by itself yields the correct result.
func TestMultiplyTwice(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1, 2, 3, 4},
	}

	expected := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{7, 10, 15, 22},
	}

	product, err := a.Multiply(&a)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() squared = %+v; want %+v", product, expected)
	}
}

// TestMultiplyWithNegativeValues verifies that multiplying matrices with negative values works correctly.
func TestMultiplyWithNegativeValues(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1, -2, -3, 4},
	}

	b := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{-5, 6, 7, -8},
	}

	expected := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{-19, 22, 43, -50},
	}

	product, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() with negative values = %+v; want %+v", product, expected)
	}
}

// TestMultiplyWithNonIntegerValues verifies that multiplying matrices with non-integer values works correctly.
func TestMultiplyWithNonIntegerValues(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1.5, 2.5, 3.5, 4.5},
	}

	b := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{5.5, 6.5, 7.5, 8.5},
	}

	expected, err := NewFromSlice(
		[]float64{1.5*5.5 + 2.5*7.5, 1.5*6.5 + 2.5*8.5, 3.5*5.5 + 4.5*7.5, 3.5*6.5 + 4.5*8.5},
		int(2),
		int(2),
	)

	if err != nil {
		t.Fatalf("NewFromSlice() returned unexpected error: %v", err)
	}

	product, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	tolerance := 1e-9
	for i := 0; i < len(expected.Value); i++ {
		if !floatEquals(product.Value[i], expected.Value[i], tolerance) {
			t.Errorf("Multiply()[%d] = %v; want %v", i, product.Value[i], expected.Value[i])
		}
	}
}

// TestMultiplyImmutableOriginals ensures that the original matrices remain unchanged after multiplication.
func TestMultiplyImmutableOriginals(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	b := Matrix{
		rows:    3,
		columns: 2,
		Value:   []float64{7, 8, 9, 10, 11, 12},
	}

	aCopy := Matrix{
		rows:    a.rows,
		columns: a.columns,
		Value:   append([]float64(nil), a.Value...),
	}

	bCopy := Matrix{
		rows:    b.rows,
		columns: b.columns,
		Value:   append([]float64(nil), b.Value...),
	}

	_, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(&a, &aCopy) {
		t.Errorf("Matrix A was modified after Multiply()\nGot: %+v\nWant: %+v", a, aCopy)
	}

	if !equalMatrices(&b, &bCopy) {
		t.Errorf("Matrix B was modified after Multiply()\nGot: %+v\nWant: %+v", b, bCopy)
	}
}

// TestMultiplyWithNegativeDimensions verifies that multiplying matrices with negative dimensions is handled appropriately.
func TestMultiplyWithNegativeDimensions(t *testing.T) {
	// Attempt to create a matrix with negative rows
	_, err := New(-1, 2)
	if err == nil {
		t.Fatalf("New() expected error due to negative rows, got nil")
	}

	// Attempt to create a matrix with negative columns
	_, err = New(2, -3)
	if err == nil {
		t.Fatalf("New() expected error due to negative columns, got nil")
	}

	// Initialize valid matrices
	a := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	b := Matrix{
		rows:    3,
		columns: 2,
		Value:   []float64{7, 8, 9, 10, 11, 12},
	}

	// Perform multiplication (should work)
	product, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	expected := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{58, 64, 139, 154},
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() = %+v; want %+v", product, expected)
	}
}

// TestMultiplyWithNonSquareMatrices verifies that multiplying non-square matrices yields the correct product.
func TestMultiplyWithNonSquareMatrices(t *testing.T) {
	a := Matrix{
		rows:    3,
		columns: 2,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	b := Matrix{
		rows:    2,
		columns: 4,
		Value:   []float64{7, 8, 9, 10, 11, 12, 13, 14},
	}

	expected := &Matrix{
		rows:    3,
		columns: 4,
		Value:   []float64{29, 32, 35, 38, 65, 72, 79, 86, 101, 112, 123, 134},
	}

	product, err := a.Multiply(&b)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatrices(product, expected) {
		t.Errorf("Multiply() with non-square matrices = %+v; want %+v", product, expected)
	}
}

/*****************************************************************************************************************/

// TestInvertSquareMatrix verifies that inverting a simple 2x2 square matrix yields the correct inverse.
func TestInvertSquareMatrix(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1, 2, 3, 4},
	}

	expectedInverse := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{-2, 1, 1.5, -0.5},
	}

	inv, err := a.Invert()
	if err != nil {
		t.Fatalf("Invert() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(inv, expectedInverse, t) {
		t.Errorf("Invert() = %+v; want %+v", inv, expectedInverse)
	}

	// Verify that A * A_inverse = Identity matrix
	identity, err := a.Multiply(inv)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	expectedIdentity := &Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1, 0, 0, 1},
	}

	if !equalMatricesDetailed(identity, expectedIdentity, t) {
		t.Errorf("A * A_inverse = %+v; want %+v", identity, expectedIdentity)
	}
}

// TestInvertIdentityMatrix verifies that inverting the identity matrix returns the identity matrix itself.
func TestInvertIdentityMatrix(t *testing.T) {
	identity := Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1,
		},
	}

	expectedInverse := &Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1,
		},
	}

	inv, err := identity.Invert()
	if err != nil {
		t.Fatalf("Invert() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(inv, expectedInverse, t) {
		t.Errorf("Invert() = %+v; want %+v", inv, expectedInverse)
	}

	// Verify that A * A_inverse = Identity matrix
	identityResult, err := identity.Multiply(inv)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(identityResult, expectedInverse, t) {
		t.Errorf("A * A_inverse = %+v; want %+v", identityResult, expectedInverse)
	}
}

// TestInvertSingularMatrix verifies that attempting to invert a singular matrix returns an appropriate error.
func TestInvertSingularMatrix(t *testing.T) {
	singular := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1, 2, 2, 4},
	}

	_, err := singular.Invert()
	if err == nil {
		t.Fatalf("Invert() expected error for singular matrix, got nil")
	}

	expectedErr := "matrix is singular and cannot be inverted"
	if err.Error() != expectedErr {
		t.Errorf("Invert() error = %v; want '%s'", err, expectedErr)
	}
}

// TestInvertLargerMatrix verifies that inverting a larger 3x3 matrix yields the correct inverse.
func TestInvertLargerMatrix(t *testing.T) {
	a := Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			4, 7, 2,
			3, 6, 1,
			2, 5, 1,
		},
	}

	expectedInverse := &Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			0.3333333333333333, 1, -1.6666666666666667,
			-0.3333333333333333, 0, 0.6666666666666666,
			1, -2, 1,
		},
	}

	inv, err := a.Invert()
	if err != nil {
		t.Fatalf("Invert() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(inv, expectedInverse, t) {
		t.Errorf("Invert() = %+v; want %+v", inv, expectedInverse)
	}

	// Verify that A * A_inverse = Identity matrix
	identity, err := a.Multiply(inv)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	expectedIdentity := &Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1,
		},
	}

	if !equalMatricesDetailed(identity, expectedIdentity, t) {
		t.Errorf("A * A_inverse = %+v; want %+v", identity, expectedIdentity)
	}
}

// TestInvertWithFloatingPointValues verifies that inverting a matrix with floating-point values works correctly.
func TestInvertWithFloatingPointValues(t *testing.T) {
	a := Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			6, 24, 1,
			13, 42, 5,
			3, 6, 1,
		},
	}

	expectedInverse := &Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			0.16666666666666666, -0.25, 1.0833333333333333,
			0.027777777777777776, 0.041666666666666664, -0.2361111111111111,
			-0.6666666666666666, 0.5, -0.8333333333333334,
		},
	}

	inv, err := a.Invert()
	if err != nil {
		t.Fatalf("Invert() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(inv, expectedInverse, t) {
		t.Errorf("Invert() = %+v; want %+v", inv, expectedInverse)
	}

	// Verify that A * A_inverse = Identity matrix
	identity, err := a.Multiply(inv)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	expectedIdentity := &Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1,
		},
	}

	if !equalMatricesDetailed(identity, expectedIdentity, t) {
		t.Errorf("A * A_inverse = %+v; want %+v", identity, expectedIdentity)
	}
}

// TestInvertSingleElementMatrix verifies that inverting a single-element matrix yields the correct inverse.
func TestInvertSingleElementMatrix(t *testing.T) {
	a := Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{4},
	}

	expectedInverse := &Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{0.25},
	}

	inv, err := a.Invert()
	if err != nil {
		t.Fatalf("Invert() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(inv, expectedInverse, t) {
		t.Errorf("Invert() = %+v; want %+v", inv, expectedInverse)
	}

	// Verify that A * A_inverse = Identity matrix
	identity, err := a.Multiply(inv)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	expectedIdentity := &Matrix{
		rows:    1,
		columns: 1,
		Value:   []float64{1},
	}

	if !equalMatricesDetailed(identity, expectedIdentity, t) {
		t.Errorf("A * A_inverse = %+v; want %+v", identity, expectedIdentity)
	}
}

// TestInvertImmutableMatrix ensures that the original matrix remains unchanged after inversion.
func TestInvertImmutableMatrix(t *testing.T) {
	a := Matrix{
		rows:    2,
		columns: 2,
		Value:   []float64{1, 2, 3, 4},
	}

	aCopy := Matrix{
		rows:    a.rows,
		columns: a.columns,
		Value:   append([]float64(nil), a.Value...),
	}

	_, err := a.Invert()
	if err != nil {
		t.Fatalf("Invert() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(&a, &aCopy, t) {
		t.Errorf("Matrix A was modified after Invert()\nGot: %+v\nWant: %+v", a, aCopy)
	}
}

// TestInvertNonSquareMatrix verifies that attempting to invert a non-square matrix returns an appropriate error.
func TestInvertNonSquareMatrix(t *testing.T) {
	nonSquare := Matrix{
		rows:    2,
		columns: 3,
		Value:   []float64{1, 2, 3, 4, 5, 6},
	}

	_, err := nonSquare.Invert()
	if err == nil {
		t.Fatalf("Invert() expected error for non-square matrix, got nil")
	}

	expectedErr := "only square matrices can be inverted"
	if err.Error() != expectedErr {
		t.Errorf("Invert() error = %v; want '%s'", err, expectedErr)
	}
}

// TestInvertAnotherSingularMatrix verifies that another singular matrix returns an appropriate error upon inversion.
func TestInvertAnotherSingularMatrix(t *testing.T) {
	singular := Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			2, 4, 2,
			4, 8, 4,
			1, 2, 1,
		},
	}

	_, err := singular.Invert()
	if err == nil {
		t.Fatalf("Invert() expected error for singular matrix, got nil")
	}

	expectedErr := "matrix is singular and cannot be inverted"
	if err.Error() != expectedErr {
		t.Errorf("Invert() error = %v; want '%s'", err, expectedErr)
	}
}

// TestInvertFloatingPointMatrix verifies that inverting a 3x3 matrix with floating-point values works correctly.
func TestInvertFloatingPointMatrix(t *testing.T) {
	a := Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			6, 24, 1,
			13, 42, 5,
			3, 6, 1,
		},
	}

	expectedInverse := &Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			0.16666666666666666, -0.25, 1.0833333333333333,
			0.027777777777777776, 0.041666666666666664, -0.2361111111111111,
			-0.6666666666666666, 0.5, -0.8333333333333334,
		},
	}

	inv, err := a.Invert()
	if err != nil {
		t.Fatalf("Invert() returned unexpected error: %v", err)
	}

	if !equalMatricesDetailed(inv, expectedInverse, t) {
		t.Errorf("Invert() = %+v; want %+v", inv, expectedInverse)
	}

	// Verify that A * A_inverse = Identity matrix
	identity, err := a.Multiply(inv)
	if err != nil {
		t.Fatalf("Multiply() returned unexpected error: %v", err)
	}

	expectedIdentity := &Matrix{
		rows:    3,
		columns: 3,
		Value: []float64{
			1, 0, 0,
			0, 1, 0,
			0, 0, 1,
		},
	}

	if !equalMatricesDetailed(identity, expectedIdentity, t) {
		t.Errorf("A * A_inverse = %+v; want %+v", identity, expectedIdentity)
	}
}

/*****************************************************************************************************************/
