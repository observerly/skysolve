/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package matrix

/*****************************************************************************************************************/

import (
	"errors"
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
