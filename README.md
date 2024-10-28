# skysolve

### Introduction

observerly's performant, zero-dependency Go plate solving library for astronomical images.

Some plate solving algorithms start with a zero sum assumption of blindness to where in the Sky an image is representing. However, most observatories have a 
rough idea of where they are pointing their telescopes. This library is designed to take advantage of that information to provide a faster, more accurate 
plate solving solution.

Hence, when the approximate equatorial coordinates of a image and the detecotors' field of view is known, skysolve can be used to compute a World Coordinate 
System (WCS) for the image. This is useful for astrometry, photometry, and other astronomical image processing tasks.

---

### Special Acknowledgements

As always, human knowledge is built on the shoulders of giants. This project would not be possible without the work of the following individuals and organizations:

- [astrometry.net](https://aa.usno.navy.mil/software/novas/novas_info.php) - The original and best, blind plate solving software that this project is based on.
- [Twirl](https://github.com/lgrcia/twirl) - Modern, astrometric plate solving package for Python.
- [The GAIA DR3 Catalog](https://www.cosmos.esa.int/web/gaia/dr3) - The most accurate and precise sky source catalog in existence.
- [The Sloan Digital Sky Survey](https://www.sdss.org/) - The most comprehensive digital sky survey in existence.
