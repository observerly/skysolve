/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package star

/*****************************************************************************************************************/

type Star struct {
	Designation string  // e.g., some catalog ID or some colloquial name, e.g., "Sirius", or "HD 1" etc
	X           float64 // X pixel coordinate
	Y           float64 // Y pixel coordinate
	RA          float64 // Sky coordinates in the azimuthal plane (in degrees)
	Dec         float64 // Sky coordinates in the polar plane (in degrees)
	Intensity   float64 // Intensity of the star at the central pixel, X and Y
}

/*****************************************************************************************************************/
