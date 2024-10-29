# skysolve

### Introduction

skysolve is a high-performance, zero-dependency Go library designed for plate solving astronomical images.

While many plate solving algorithms operate under the assumption of having no prior knowledge of an image’s location in the sky, most observatories typically have a rough estimate of where their telescopes are pointed. skysolve leverages this existing pointing information to deliver faster and more accurate plate solving solutions.

When provided with the approximate equatorial coordinates of an image and the detector’s field of view, skysolve can compute a World Coordinate System (WCS) for the image in under a second. This efficiency makes it an excellent choice for tasks in astrometry, photometry, and other astronomical image processing applications that demand high performance, including Space Situational Awareness (SSA) and Space Domain Awareness (SDA).

### Key Features

- Zero Dependencies: Simplifies integration into your projects without the overhead of managing additional libraries.
- High Performance: Optimized for speed, enabling real-time processing and analysis.
- Accuracy: Utilizes prior pointing information to enhance the precision of plate solving.

---

### Special Acknowledgements

As always, human knowledge is built on the shoulders of giants. This project would not be possible without the work of the following individuals and organizations:

- [astrometry.net](https://aa.usno.navy.mil/software/novas/novas_info.php) - The original and best, blind plate solving software that this project is based on.
- [Twirl](https://github.com/lgrcia/twirl) - Modern, astrometric plate solving package for Python.
- [The GAIA DR3 Catalog](https://www.cosmos.esa.int/web/gaia/dr3) - The most accurate and precise sky source catalog in existence.
- [The Sloan Digital Sky Survey](https://www.sdss.org/) - The most comprehensive digital sky survey in existence.
