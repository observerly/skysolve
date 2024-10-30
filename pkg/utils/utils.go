/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package utils

/*****************************************************************************************************************/

import (
	"fmt"
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
