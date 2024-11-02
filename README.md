# skysolve

### Introduction

skysolve is a high-performance, zero-dependency Go library designed for plate solving astronomical images.

While many plate solving algorithms operate under the assumption of having no prior knowledge of an image’s location in the sky, most observatories typically have a rough estimate of where their telescopes are pointed. skysolve leverages this existing pointing information to deliver faster and more accurate plate solving solutions.

When provided with the approximate equatorial coordinates of an image and the detector’s field of view, skysolve can compute a World Coordinate System (WCS) for the image in under a second. This efficiency makes it an excellent choice for tasks in astrometry, photometry, and other astronomical image processing applications that demand high performance, including Space Situational Awareness (SSA) and Space Domain Awareness (SDA).

### Why SkySolve?

### Exceptional Precision & Performance

Enhanced Precision: Leverages prior pointing information to improve the accuracy of plate solving. When combined locally with the GAIA DR3 catalog, it can achieve sub-arcsecond precision within ~<100ms on a modern CPU.

#### Zero Dependencies

Seamless Integration: Easily incorporate into your projects without the need to manage additional libraries, interoperating with C or C++ modules, simplifying deployment and maintenance.
  
#### High Performance

Optimised Speed: Facilitates real-time processing and analysis for efficient workflows, for example in Space Situational Awareness (SSA) and Space Domain Awareness (SDA).
  
#### Star Identification Protocol (SIP)

Reliable Star Matching: Utilises SIP to accurately identify and match stars within the image, enhancing overall solving reliability.

#### Adheres To FITS Standards

WCS Compliant: Generates World Coordinate System (WCS) solutions that adhere to the FITS standard, ensuring compatibility with existing astronomical software and tools.

### Algorithm & Methodology

#### 1. Image Star Extraction

The process begins with extracting stars from the input image using a sophisticated star detection algorithm. This algorithm identifies bright spots that are potential stars, accurately determining their positions and intensities. These extracted data points form the foundation for subsequent matching and analysis.

#### 2. Catalog Matching

With the stars extracted, the algorithm matches them against a comprehensive reference catalog containing star positions and magnitudes. By comparing the spatial arrangement and brightness of the detected stars with those in the catalog, the algorithm identifies the best possible match, facilitating accurate plate solving.

#### 3. Invariant Features

To ensure robust plate solving, the algorithm utilizes the invariant properties of triangles formed by stars. Specifically, the ratio of side lengths remains constant, normalizing any scale differences, while the angles between sides stay the same regardless of rotation, scaling, or translation. By comparing these invariant features from the image with those in the reference catalog, the algorithm can accurately determine the plate’s orientation and position. This approach enhances solution precision as more stars are detected and increases the likelihood of finding a matching solution with a larger reference catalog.

#### 4. Affine Transformations

To align the captured image with the reference catalog, the algorithm employs affine transformations. These transformations include rotation, which adjusts the image orientation to match the catalog; scaling, which normalizes the size differences between the image and reference data; and translation, which shifts the image position to align with the catalog coordinates. Affine transformations preserve points, straight lines, and planes, ensuring that the alignment process maintains the geometric integrity of the image.

#### 5. Simple Imaging Polynomial (SIP) Integration

**Simple Imaging Polynomial (SIP)** is employed to account for optical distortions and projection effects inherent in astronomical imaging systems. SIP uses polynomial coefficients to map pixel coordinates to world coordinates, enabling precise correction of image distortions. By modeling the relationship between pixel positions and celestial coordinates, SIP corrects for lens aberrations and other distortions. Additionally, applying SIP transformations ensures that the star positions are accurately represented in a standardized projection system, facilitating reliable catalog matching.

---

### Special Acknowledgements

As always, human knowledge is built on the shoulders of giants. This project would not be possible without the work of the following individuals and organizations:

- [astrometry.net](https://aa.usno.navy.mil/software/novas/novas_info.php) - The original and best, blind plate solving software that this project is based on.
- [Twirl](https://github.com/lgrcia/twirl) - Modern, astrometric plate solving package for Python.
- [The GAIA DR3 Catalog](https://www.cosmos.esa.int/web/gaia/dr3) - The most accurate and precise sky source catalog in existence.
- [The Sloan Digital Sky Survey](https://www.sdss.org/) - The most comprehensive digital sky survey in existence.
