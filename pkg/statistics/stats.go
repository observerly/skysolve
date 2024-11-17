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

// Poisson generates a Poisson-distributed random number with mean lambda.
// lambda: the mean of the distribution.
func PoissonDistributedRandomNumber(lambda float64) float64 {
	// If lambda is less than 0, return 0:
	if lambda < 0 {
		return 0
	}

	// If lambda is 0, return 0:
	if lambda == 0 {
		return 0
	}

	// Calculate the exponential of -lambda:
	L := math.Exp(-lambda)

	// Initialize k, which is the number of random numbers generated:
	k := 0.0

	// Initialize p, which is the product of all random numbers generated:
	p := 1.0

	for p > L {
		k++
		u := rand.Float64()
		p *= u
	}

	return k - 1
}

/*****************************************************************************************************************/
