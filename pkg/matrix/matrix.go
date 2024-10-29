/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package matrix

/*****************************************************************************************************************/

import (
	"errors"
	"fmt"
)

/*****************************************************************************************************************/

// Matrix represents a 2D matrix in row-major order.
type Matrix struct {
	rows    int
	columns int
	Value   []float64
}

/*****************************************************************************************************************/

// New creates a new matrix with the specified number of rows and columns.
// All elements are initialized to zero.
func New(rows, columns int) (*Matrix, error) {
	if rows <= 0 || columns <= 0 {
		return nil, errors.New("matrix dimensions must be positive")
	}

	value := make([]float64, rows*columns)

	return &Matrix{
		rows:    rows,
		columns: columns,
		Value:   value,
	}, nil
}

/*****************************************************************************************************************/

// NewFromSlice creates a new matrix from a given slice.
// The slice should have exactly rows*columns elements.
func NewFromSlice(value []float64, rows, columns int) (*Matrix, error) {
	// Check if the matrix dimensions are valid
	if rows <= 0 || columns <= 0 {
		return nil, errors.New("matrix dimensions must be positive")
	}

	length := len(value)

	// Check if the data length matches the matrix dimensions
	if length != rows*columns {
		return nil, fmt.Errorf("length %d does not match matrix dimensions %dx%d", length, rows, columns)
	}

	// Create a copy to prevent external modifications
	v := make([]float64, length)

	// Copy the values from the given slice to the new matrix
	copy(v, value)

	return &Matrix{
		rows:    rows,
		columns: columns,
		Value:   v,
	}, nil
}

/*****************************************************************************************************************/

// Rows returns the number of rows in the matrix.
func (m *Matrix) Rows() int {
	return m.rows
}

/*****************************************************************************************************************/

// Columns returns the number of columns in the matrix.
func (m *Matrix) Columns() int {
	return m.columns
}

/*****************************************************************************************************************/

// At returns the element at the specified row and column.
// Rows and columns are zero-indexed.
func (m *Matrix) At(row, col int) (float64, error) {
	if row < 0 || row >= m.rows || col < 0 || col >= m.columns {
		return 0, fmt.Errorf("index out of bounds: row=%d, col=%d", row, col)
	}

	return m.Value[row*m.columns+col], nil
}

/*****************************************************************************************************************/

// Set sets the element at the specified row and column to the given value.
// Rows and columns are zero-indexed.
func (m *Matrix) Set(row, col int, value float64) error {
	if row < 0 || row >= m.rows || col < 0 || col >= m.columns {
		return fmt.Errorf("index out of bounds: row=%d, col=%d", row, col)
	}

	m.Value[row*m.columns+col] = value

	return nil
}

/*****************************************************************************************************************/
