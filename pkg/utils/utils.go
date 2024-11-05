/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package utils

/*****************************************************************************************************************/

import (
	"fmt"
	"math"
	"strconv"

	"github.com/observerly/skysolve/pkg/geometry"
)

/*****************************************************************************************************************/

func quantizeFeature(value float64, precision int) string {
	format := "%." + strconv.Itoa(precision) + "f"
	return fmt.Sprintf(format, value)
}

/*****************************************************************************************************************/

func QuantizeFeatures(features geometry.InvariantFeatures, precision int) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		quantizeFeature(features.RatioAB, precision),
		quantizeFeature(features.RatioAC, precision),
		quantizeFeature(features.AngleA, precision),
		quantizeFeature(features.AngleB, precision),
	)
}

/*****************************************************************************************************************/

// computePolynomialTerms generates the terms of a 2D polynomial up to the given order.
// For example, for order 3: [1, x, y, x^2, x*y, y^2, x^3, x^2*y, x*y^2, y^3]
func ComputePolynomialTerms(x, y float64, order int) []float64 {
	terms := []float64{}
	for i := 0; i <= order; i++ {
		for j := 0; j <= i; j++ {
			exponentX := i - j
			exponentY := j
			term := math.Pow(x, float64(exponentX)) * math.Pow(y, float64(exponentY))
			terms = append(terms, term)
		}
	}
	return terms
}

/*****************************************************************************************************************/

// generatePolynomialTermKeys generates the FITS-compatible keys for polynomial terms up to the given order.
// For example, for order 2: ["A_0_0", "A_0_1", "A_0_2", "A_1_0", "A_1_1", "A_2_0"]
func GeneratePolynomialTermKeys(prefix string, order int) []string {
	terms := []string{}
	for i := 0; i <= order; i++ {
		for j := 0; j <= i; j++ {
			exponentX := i - j
			exponentY := j
			term := fmt.Sprintf("%s_%d_%d", prefix, exponentX, exponentY)
			terms = append(terms, term)
		}
	}
	return terms
}

/*****************************************************************************************************************/
