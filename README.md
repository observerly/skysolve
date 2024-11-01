# skysolve

### Introduction

skysolve is a high-performance, zero-dependency Go library designed for plate solving astronomical images.

While many plate solving algorithms operate under the assumption of having no prior knowledge of an image’s location in the sky, most observatories typically have a rough estimate of where their telescopes are pointed. skysolve leverages this existing pointing information to deliver faster and more accurate plate solving solutions.

When provided with the approximate equatorial coordinates of an image and the detector’s field of view, skysolve can compute a World Coordinate System (WCS) for the image in under a second. This efficiency makes it an excellent choice for tasks in astrometry, photometry, and other astronomical image processing applications that demand high performance, including Space Situational Awareness (SSA) and Space Domain Awareness (SDA).

### Why SkySolve?

Here at observerly, we love to write type-safe, memory-safe and performant code. We regard Go lang as giving us the code writing-efficiency of Python, with the performance of C, all with the type-safe assuredness of ... Go. We aim to make it easy for feature contributions, bug fixes and optimisations whilst not compromising on code performance. When utilised with Go's coroutines, this package can blind plate solve in less than a minute.

### Key Features

- Zero Dependencies: Simplifies integration into your projects without the overhead of managing additional libraries.
- High Performance: Optimized for speed, enabling real-time processing and analysis.
- Accuracy: Utilizes prior pointing information to enhance the precision of plate solving.

### Algorithm & Methodology

#### Image Star Extraction

The algorithm first extracts stars from the image using a star detection algorithm. The star detection algorithm identifies bright spots in the image that are likely to be stars. These spots are then used to extract the stars' positions and intensities.

#### Catalog Matching

The extracted stars are then matched against a reference catalog of stars. The reference catalog contains the positions and magnitudes of stars in the sky. The algorithm uses the positions of the stars in the image and the reference catalog to find the best match between the two sets of stars.

#### Invariant Features

The algorithm used for plate solving leverages the invariant properties of triangles. Specifically, triangles retain consistent features: the ratio of side lengths, which normalizes scale differences, and the angles, which remain unchanged by rotation, scaling, and translation. This algorithm compares the invariant features of triangles formed by various stars in the image against a reference catalog to accurately solve the plate.

The accuracy of the solution improves with the number of stars detected in the image, and the probability of finding a solution increases with a larger catalog of reference sources.

#### Affine Transformations

The algorithm uses affine transformations to align the image with the reference catalog. Affine transformations are a class of linear transformations that preserve points, straight lines, and planes. The algorithm uses these transformations to rotate, scale, and translate the image to match the reference catalog.

---

### Special Acknowledgements

As always, human knowledge is built on the shoulders of giants. This project would not be possible without the work of the following individuals and organizations:

- [astrometry.net](https://aa.usno.navy.mil/software/novas/novas_info.php) - The original and best, blind plate solving software that this project is based on.
- [Twirl](https://github.com/lgrcia/twirl) - Modern, astrometric plate solving package for Python.
- [The GAIA DR3 Catalog](https://www.cosmos.esa.int/web/gaia/dr3) - The most accurate and precise sky source catalog in existence.
- [The Sloan Digital Sky Survey](https://www.sdss.org/) - The most comprehensive digital sky survey in existence.
