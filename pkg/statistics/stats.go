/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package stats

/*****************************************************************************************************************/

import (
	"math"
	"math/rand"
)

/*****************************************************************************************************************/

// NormalDistributedRandomNumber generates a normally distributed random number.
// mean: the mean of the distribution.
// stdDev: the standard deviation of the distribution.
func NormalDistributedRandomNumber(mean, stdDev float64) float64 {
	v := rand.Float64()
	return v*(stdDev*math.Sqrt(2*math.Pi)) + mean
}

/*****************************************************************************************************************/
