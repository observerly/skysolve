![@observerly:skysolve](./.github/assets/banner.png)

---

![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/observerly/skysolve/main?filename=go.mod&label=Go)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/observerly/skysolve)](https://pkg.go.dev/github.com/observerly/skysolve)
[![Go Report Card](https://goreportcard.com/badge/github.com/observerly/skysolve)](https://goreportcard.com/report/github.com/observerly/skysolve)
[![SkySolve Actions Status](https://github.com/observerly/skysolve/actions/workflows/ci.yml/badge.svg)](https://github.com/observerly/skysolve/actions/workflows/ci.yml)

### Introduction

skysolve is a high-performance, zero-dependency Go library designed for plate solving astronomical images.

While many plate solving algorithms operate under the assumption of having no prior knowledge of an image’s location in the sky, most observatories typically have a rough estimate of where their telescopes are pointed. skysolve leverages this existing pointing information to deliver faster and more accurate plate solving solutions.

When provided with the approximate equatorial coordinates of an image and the detector’s field of view, skysolve can compute a World Coordinate System (WCS) for the image in under a second. This efficiency makes it an excellent choice for tasks in astrometry, photometry, and other astronomical image processing applications that demand high performance, including Space Situational Awareness (SSA) and Space Domain Awareness (SDA).

### Prerequisites

- [go](https://go.dev/) (>= 1.21.*)

### Why SkySolve?

#### Exceptional Precision & Performance

Enhanced Precision: Leverages prior pointing information to improve the accuracy of plate solving. When combined locally with the GAIA DR3 catalog, it can achieve sub-arcsecond precision within ~<100ms on a modern CPU.

#### Zero Dependencies

Seamless Integration: Easily incorporate into your projects *without* the need to manage additional libraries or interoperating with C or C++ modules, third-party dependencies or Linux packages, simplifying usage, deployment and maintenance.
  
#### High Performance

Optimised Speed: Facilitates real-time processing and analysis for efficient workflows, for example in Space Situational Awareness (SSA) and Space Domain Awareness (SDA).
  
#### Simple Imaging Polynomial (SIP)

Corrects for optical distortions and projection effects, ensuring accurate representation of star positions in a standardized projection system.

#### Adheres To FITS Standards

WCS Compliant: Generates World Coordinate System (WCS) solutions that adhere to the FITS standard, ensuring compatibility with existing astronomical software and tools.


### Usage

SkySolve is designed to interoperate between IRIS, a FITs image processing library [observerly/iris](), and the GAIA DR3 catalog. The following example demonstrates how to use SkySolve to plate solve an astronomical image:

```go
package main

import (
	"fmt"
	"os"

	"github.com/observerly/iris/pkg/fits"
	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/geometry"
	"github.com/observerly/skysolve/pkg/solve"
)

func main() {
	// Attempt to open the file from the given filepath:
	file, err := os.Open("<PATH_TO_YOUR_IMAGE>.fits")
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
		return
	}
	// Defer closing the file:
	defer file.Close()

	// Assume an image of 2x2 pixels with 16-bit depth, and no offset:
	fit := fits.NewFITSImage(2, 0, 0, 65535)

	// Read in our exposure data into the image:
	err = fit.Read(file)

	if err != nil {
		fmt.Printf("failed to read fits file: %v", err)
	}

	// Attempt to get the RA header from the FITS file:
	ra, exists := fit.Header.Floats["RA"]
	if !exists {
		fmt.Printf("ra header not found")
		return
	}

	// Attempt to get the Dec header from the FITS file:
	dec, exists := fit.Header.Floats["DEC"]
	if !exists {
		fmt.Println("dec header not found")
		return
	}

	// Get the radius based on the known but approximate pixel scale of the image:
	radius := fov.GetRadialExtent(float64(fit.Header.Naxis1), float64(fit.Header.Naxis2), fov.PixelScale{X: pixelScaleX, Y: pixelScaleY})

	// Create a new GAIA service client:
	service := catalog.NewCatalogService(catalog.GAIA, catalog.Params{
		Limit:     int(math.Ceil(42 / 0.5 * radius)),
		Threshold: 16, // Limiting Magntiude, filter out any stars that are magnitude 16 or above (fainter)
	})
	
	eq := astrometry.ICRSEquatorialCoordinate{
		RA:  float64(ra.Value),
		Dec: float64(dec.Value),
	}

	// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
	sources, err := service.PerformRadialSearch(eq, radius)
	if err != nil {
		fmt.Printf("there was an error while performing the radial search: %v", err)
		return
	}

	// Attempt to create a new PlateSolver:
	solver, err := solve.NewPlateSolver(solve.Params{
		Data:                fit.Data,               // The exposure data from the fits image
		Width:               int(fit.Header.Naxis1), // The width of the image
		Height:              int(fit.Header.Naxis2), // The height of the image
		PixelScaleX:         pixelScaleX,            // The pixel scale in the x-axis
		PixelScaleY:         pixelScaleY,            // The pixel scale in the y-axis
		ADU:                 fit.ADU,                // The analog-to-digital unit of the image
		ExtractionThreshold: 16,                     // Extract a minimum of 20 of the brightest stars
		Radius:              16,                     // 16 pixels radius for the star extraction
		Sigma:               2.5,                    // 8 pixels sigma for the Gaussian kernel
	})
	if err != nil {
		fmt.Printf("there was an error while creating the plate solver: %v", err)
		return
	}

	solver.Sources = append(solver.Sources, sources...)

	// Define the tolerances for the solver, we can adjust these as needed:
	tolerance := solve.ToleranceParams{
		QuadTolerance:           0.02,
		EuclidianPixelTolerance: 10,
	}

	// Extract the WCS solution, as well as the matches, from the solver:
	wcs, matches, err := solver.Solve(tolerance, 3)
	if err != nil {
		fmt.Printf("an error occured while plate solving: %v", err)
		return
	}
}
```

### Examples

Examples can be found in the [./examples](./examples) directory:

- [Astrometric Plate Solver](./examples/solve/main.go) - Plate solves a FITs image using the SkySolve library.
- [Generated GAIA Sky](./examples/sky/main.go) - Generates a sky map of the GAIA DR2 & DR3 catalog, and returns a ZScale normalized image of the sky for your chosen parameters.

### Algorithm & Methodology

#### 1. Image Star Extraction

The process begins with extracting stars from the input image using a sophisticated star detection algorithm. This algorithm identifies bright spots that are potential stars, accurately determining their positions and intensities. These extracted data points form the foundation for subsequent matching and analysis.

#### 2. Catalog Matching

With the stars extracted, the algorithm matches them against a comprehensive reference catalog containing star positions and magnitudes. By comparing the spatial arrangement and brightness of the detected stars with those in the catalog, the algorithm identifies the best possible match, facilitating accurate plate solving.

#### 3. Invariant Features

Invariant geometric features are crucial for robust image registration and plate solving, particularly when astronomical images inhabit a vast hypothesis space in which exhaustive searches are impractical. Our approach employs a continuous geometric hashing technique, selecting four stars (a “quad”), defining a local coordinate system using the two most widely separated stars, and mapping the remaining two stars’ positions into a four-dimensional hash code. This code is invariant under translation, rotation, and scaling, and it is resilient to small positional noise. Furthermore, because the codes are well-distributed for uniformly scattered stars, storing them in a vp-tree enables efficient retrieval of matching quads. By identifying these invariant features in both the image and the reference catalogue, the algorithm can rapidly and accurately align the datasets, ensuring a precise and consistent plate-solving solution.

#### 4. Affine Transformations

To align the captured image with the reference catalog, the algorithm employs affine transformations. These transformations include rotation, which adjusts the image orientation to match the catalog; scaling, which normalizes the size differences between the image and reference data; and translation, which shifts the image position to align with the catalog coordinates. Affine transformations preserve points, straight lines, and planes, ensuring that the alignment process maintains the geometric integrity of the image.

#### 5. Simple Imaging Polynomial (SIP) Integration

**Simple Imaging Polynomial (SIP)** is employed to account for optical distortions and projection effects inherent in astronomical imaging systems. SIP uses polynomial coefficients to map pixel coordinates to world coordinates, enabling precise correction of image distortions. By modeling the relationship between pixel positions and celestial coordinates, SIP corrects for lens aberrations and other distortions. Additionally, applying SIP transformations ensures that the star positions are accurately represented in a standardized projection system, facilitating reliable catalog matching.

---

### To Do

- [x] For an approximate `{ ra, dec }` and a known field of view, compute the corresponding HEALPix cell for N side 128 and 256 to ensure fast and efficient index lookups for a given equatorial coordinate space region.
- [ ] Implement an offline-first index of HEALPix cells for the GAIA DR2 & DR3 catalog.
- [ ] Implement an online way to retrieve all notable objects from SIMBAD in an image region.
- [ ] Improve robustness, accuracy and precision of the plate solving algorithm down to a precision of +/- 2 degrees in RA and Dec from the approximate known coordinates.
- [ ] Implement a way to retrieve the approximate coordinates of an image from the FITS header.

---

### License

Mozilla Public License 2.0

Permissions of this weak copyleft license are conditioned on making available source code of licensed files and modifications of those files under the same license (or in certain cases, one of the GNU licenses). Copyright and license notices must be preserved. Contributors provide an express grant of patent rights. However, a larger work using the licensed work may be distributed under different terms and without source code for files added in the larger work.

### Special Acknowledgements

As always, human knowledge is built on the shoulders of giants. This project would not be possible without the work of the following individuals and organizations:

- [astrometry.net](https://aa.usno.navy.mil/software/novas/novas_info.php) - The original and best, blind plate solving software that this project is based on.
- [Twirl](https://github.com/lgrcia/twirl) - Modern, astrometric plate solving package for Python.
- [The GAIA DR3 Catalog](https://www.cosmos.esa.int/web/gaia/dr3) - The most accurate and precise sky source catalog in existence.
- [The Sloan Digital Sky Survey](https://www.sdss.org/) - The most comprehensive digital sky survey in existence.
- [HEALPix](https://healpix.sourceforge.io/) - The Hierarchical Equal Area isoLatitude Pixelization of the sphere, without which, indexing of the night sky would not be possible.

### Standing On The Shoulders Of Giants

To the Harvard Computers, Annie Jump Cannon, Antonia Maury, Henrietta Leavitt, Evelyn Leland, Anna Winlock, Muriel Mussells Seyfert, Cecilia Payne-Gaposchkin, Williamina Fleming and Florence Cushman.


