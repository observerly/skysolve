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
