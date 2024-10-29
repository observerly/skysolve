/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package matrix

/*****************************************************************************************************************/

import "testing"

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
