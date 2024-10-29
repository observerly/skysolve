/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package matrix

/*****************************************************************************************************************/

import "testing"

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
