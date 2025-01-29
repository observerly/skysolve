/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package solver

/*****************************************************************************************************************/

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/observerly/iris/pkg/fits"
	"github.com/observerly/skysolve/internal/utils"
	"github.com/observerly/skysolve/pkg/astrometry"
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/fov"
	"github.com/observerly/skysolve/pkg/solve"
	"github.com/spf13/cobra"
)

/*****************************************************************************************************************/

var (
	InputFileLocation          string
	RA                         float32
	Dec                        float32
	PixelScaleX                float64
	PixelScaleY                float64
	QuadTolerance              float64
	EuclidianDistanceTolerance float64
)

/*****************************************************************************************************************/

func getFilePathStem(file *os.File) string {
	path := file.Name()
	// Get the directory where the file is located (e.g. "./samples")
	directory := filepath.Dir(path)
	// Get the full filename (e.g. "astrometry.fits")
	base := filepath.Base(path)
	// Extract the extension (e.g. ".fits")
	extension := filepath.Ext(base)
	// Remove the extension from the filename (e.g. "astrometry"):
	name := strings.TrimSuffix(base, extension)
	// Return the filepath stem (e.g. "./samples/astrometry")
	return filepath.Join(directory, name)
}

/*****************************************************************************************************************/

var AstrometryCommand = &cobra.Command{
	Use:   "astrometry",
	Short: "astrometry",
	Long:  "astrometry",
	Run: func(cmd *cobra.Command, args []string) {
		// Attempt to open the file from the given filepath and validate it exists:
		inputFile, err := os.Open(InputFileLocation)
		if err != nil {
			fmt.Println("failed to open input file:", err)
			cmd.Usage()
			return
		}

		fmt.Println("Input File Location:", InputFileLocation)

		// Defer closing the input file:
		defer inputFile.Close()

		params := RunSolverParams{
			InputFile:                    inputFile,
			RA:                           RA,
			Dec:                          Dec,
			PixelScaleX:                  PixelScaleX,
			PixelScaleY:                  PixelScaleY,
			QuadTolerance:                QuadTolerance,
			EuclidianceDistanceTolerance: EuclidianDistanceTolerance,
		}

		// Attempt to run the solver with the given parameters:
		err = RunSolver(params)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

/*****************************************************************************************************************/

func init() {
	// Add the input flag to the astrometry command for reading the file from some input location:
	// example usage: --input ./astrometry.fits or -i ./astrometry.fits
	AstrometryCommand.Flags().StringVarP(
		&InputFileLocation,
		"input",
		"i",
		"",
		"The input file location on the filesystem",
	)
	AstrometryCommand.MarkFlagRequired("input")

	// Add the approximated point equatorial coordinate RA to the astrometry command for setting the approximate RA:
	// example usage: --ra 98.6
	AstrometryCommand.Flags().Float32VarP(
		&RA,
		"ra",
		"",
		float32(math.NaN()),
		"The approximate right ascension of the image",
	)

	// Add the approximated point equatorial coordinate dec to the astrometry command for setting the approximate dec:
	// example usage: --dec 2.5
	AstrometryCommand.Flags().Float32VarP(
		&Dec,
		"dec",
		"",
		float32(math.NaN()),
		"The approximate declination of the image",
	)

	// Add the pixel scale X flag to the astrometry command for setting the pixel scale in the x-axis:
	// example usage: --pixel-scale-x 0.000540 or -px 0.000540
	AstrometryCommand.Flags().Float64VarP(
		&PixelScaleX,
		"pixel-scale-x",
		"x",
		math.Inf(-1),
		"The pixel scale in the x-axis of the image",
	)

	// Add the pixel scale Y flag to the astrometry command for setting the pixel scale in the y-axis:
	// example usage: --pixel-scale-y 0.000540 or -py 0.000540
	AstrometryCommand.Flags().Float64VarP(
		&PixelScaleY,
		"pixel-scale-y",
		"y",
		math.Inf(-1),
		"The pixel scale in the y-axis of the image",
	)

	// Add the quad tolerance flag to the astrometry command for setting the quad tolerance:
	// example usage: --quad-tolerance 0.02
	AstrometryCommand.Flags().Float64VarP(
		&QuadTolerance,
		"quad-tolerance",
		"",
		0.02,
		"The quad tolerance for the solver",
	)

	// Add the euclidian distance tolerance flag to the astrometry command for setting the euclidian distance tolerance:
	// example usage: --euclidian-distance-tolerance 10
	AstrometryCommand.Flags().Float64VarP(
		&EuclidianDistanceTolerance,
		"euclidian-distance-tolerance",
		"",
		10.0,
		"The euclidian distance (in pixels) tolerance for the solver",
	)
}

/*****************************************************************************************************************/

type RunSolverParams struct {
	InputFile                    *os.File `json:"inputFile"`
	RA                           float32  `json:"ra"`
	Dec                          float32  `json:"dec"`
	PixelScaleX                  float64  `json:"pixelScaleX"`
	PixelScaleY                  float64  `json:"pixelScaleY"`
	QuadTolerance                float64  `json:"quadTolerance"`
	EuclidianceDistanceTolerance float64  `json:"euclidianDistanceTolerance"`
}

/*****************************************************************************************************************/

func RunSolver(params RunSolverParams) error {
	// Assume an image of 2x2 pixels with 16-bit depth, and no offset:
	fit := fits.NewFITSImage(2, 0, 0, 65535)

	// Read in our exposure data into the image:
	err := fit.Read(params.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Attempt to get the RA header from the FITS file, or resolve the user's input:
	ra, err := utils.ResolveOrExtractRAFromHeaders(params.RA, fit.Header)
	if err != nil {
		return fmt.Errorf("failed to resolve or extract RA from headers: %v", err)
	}

	fmt.Printf("Right Ascension: %v°\n", ra)

	// Attempt to get the Dec header from the FITS file, or resolve the user's input:
	dec, err := utils.ResolveOrExtractDecFromHeaders(params.Dec, fit.Header)
	if err != nil {
		return fmt.Errorf("failed to resolve or extract Dec from headers: %v", err)
	}

	fmt.Printf("Declination: %v°\n", dec)

	// Attempt to extract the height from the FITS file headers:
	height, err := utils.ExtractImageHeightFromHeaders(fit.Header)
	if err != nil {
		return fmt.Errorf("failed to extract height from headers: %v", err)
	}

	fmt.Printf("Height: %v pixels\n", height)

	// Attempt to extract the width from the FITS file header:                                                                                                                                     Attempt to extract the width from the FITS file headers:
	width, err := utils.ExtractImageWidthFromHeaders(fit.Header)
	if err != nil {
		return fmt.Errorf("failed to extract width from headers: %v", err)
	}

	fmt.Printf("Width: %v pixels\n", width)

	pixelScaleX := params.PixelScaleX

	// Validate the pixel scale in the x-axis:
	if pixelScaleX == math.Inf(-1) {
		return fmt.Errorf("pixel scale x is required")
	}

	if pixelScaleX == 0 {
		return fmt.Errorf("pixel scale x must be non-zero")
	}

	if pixelScaleX < 0 {
		pixelScaleX = math.Abs(pixelScaleX)
	}

	fmt.Printf("Pixel Scale X: %v\n", pixelScaleX)

	pixelScaleY := params.PixelScaleY

	// Validate the pixel scale in the y-axis:
	if pixelScaleY == math.Inf(-1) {
		return fmt.Errorf("pixel scale y is required")
	}

	if pixelScaleY == 0 {
		return fmt.Errorf("pixel scale y must be non-zero")
	}

	if pixelScaleY < 0 {
		pixelScaleY = math.Abs(pixelScaleY)
	}

	fmt.Printf("Pixel Scale Y: %v\n", pixelScaleY)

	// Get our approximate radial extent for the field of view of the image (in degrees):
	radius := fov.GetRadialExtent(
		float64(fit.Header.Naxis1),
		float64(fit.Header.Naxis2),
		fov.PixelScale{
			X: params.PixelScaleX,
			Y: params.PixelScaleY,
		})

	fmt.Printf("Search Radius: %v°\n", radius)

	// Create a new SIMBAD service client:
	service := catalog.NewCatalogService(catalog.GAIA, catalog.Params{
		Limit:     100, // Limit the number of records to 100
		Threshold: 16,  // Limiting Magntiude, filter out any stars that are magnitude 16 or above (fainter)
	})

	// Attempt to create a new PlateSolver:
	solver, err := solve.NewPlateSolver(solve.Params{
		Data:                fit.Data,           // The exposure data from the fits image
		Width:               int(width),         // The width of the image
		Height:              int(height),        // The height of the image
		PixelScaleX:         params.PixelScaleX, // The pixel scale in the x-axis
		PixelScaleY:         params.PixelScaleY, // The pixel scale in the y-axis
		ADU:                 fit.ADU,            // The analog-to-digital unit of the image
		ExtractionThreshold: 16,                 // Extract a minimum of 16 of the brightest stars
		Radius:              16,                 // 16 pixels radius for the star extraction
		Sigma:               2.5,                // 8 pixels sigma for the Gaussian kernel
	})
	if err != nil {
		fmt.Printf("there was an error while creating the plate solver: %v", err)
		return err
	}

	// Whilst we have no matches, and whilst we are within 1 degree of the initial { ra, dec } guess, keep solving:
	eq := astrometry.ICRSEquatorialCoordinate{
		RA:  float64(ra),
		Dec: float64(dec),
	}

	// Perform a radial search with the given center and radius, for all sources with a magnitude less than 10:
	sources, err := service.PerformRadialSearch(eq, radius)
	if err != nil {
		fmt.Printf("there was an error while performing the SIMBAD radial search: %v", err)
		return err
	}

	// Append the sources to the solver:
	solver.Sources = append(solver.Sources, sources...)

	// Define the tolerances for the solver, we can adjust these as needed:
	tolerance := solve.ToleranceParams{
		QuadTolerance:           params.QuadTolerance,
		EuclidianPixelTolerance: params.EuclidianceDistanceTolerance,
	}

	wcs, _, err := solver.Solve(tolerance, 3)

	if err != nil {
		fmt.Println("an error occured while plate solving:", err)
		return err
	}

	if wcs == nil {
		fmt.Println("no WCS solution found")
		return fmt.Errorf("no WCS solution found")
	}

	// Print fields from wcaxes through cd2_2
	fmt.Printf("WCAXES: %d\n", wcs.WCAXES)
	fmt.Printf("CRPIX1: %.6f\n", wcs.CRPIX1)
	fmt.Printf("CRPIX2: %.6f\n", wcs.CRPIX2)
	fmt.Printf("CRVAL1: %.6f\n", wcs.CRVAL1)
	fmt.Printf("CRVAL2: %.6f\n", wcs.CRVAL2)
	fmt.Printf("CTYPE1: %s\n", wcs.CTYPE1)
	fmt.Printf("CTYPE2: %s\n", wcs.CTYPE2)
	fmt.Printf("CDELT1: %.6f\n", wcs.CDELT1)
	fmt.Printf("CDELT2: %.6f\n", wcs.CDELT2)
	fmt.Printf("CUNIT1: %s\n", wcs.CUNIT1)
	fmt.Printf("CUNIT2: %s\n", wcs.CUNIT2)
	fmt.Printf("CD1_1:  %.6f\n", wcs.CD1_1)
	fmt.Printf("CD1_2:  %.6f\n", wcs.CD1_2)
	fmt.Printf("CD2_1:  %.6f\n", wcs.CD2_1)
	fmt.Printf("CD2_2:  %.6f\n", wcs.CD2_2)

	// Attempt to write the WCS solution to the output file:
	fit.Header.Set("WCSAXES", wcs.WCAXES, "Number of World Coordinate System axes")
	fit.Header.Set("CRPIX1", wcs.CRPIX1, "X Pixel coordinate of reference point")
	fit.Header.Set("CRPIX2", wcs.CRPIX2, "Y Pixel coordinate of reference point")
	fit.Header.Set("CDELT1", wcs.CDELT1, "Coordinate increment at reference point")
	fit.Header.Set("CDELT2", wcs.CDELT2, "Coordinate increment at reference point")
	fit.Header.Set("CUNIT1", wcs.CUNIT1, "Units of coordinate increment and value")
	fit.Header.Set("CUNIT2", wcs.CUNIT2, "Units of coordinate increment and value")
	fit.Header.Set("CTYPE1", wcs.CTYPE1, "Coordinate type code")
	fit.Header.Set("CTYPE2", wcs.CTYPE2, "Coordinate type code")
	fit.Header.Set("CRVAL1", wcs.CRVAL1, "Coordinate value at reference point")
	fit.Header.Set("CRVAL2", wcs.CRVAL2, "Coordinate value at reference point")
	fit.Header.Set("CD1_1", wcs.CD1_1, "Coordinate transformation matrix element")
	fit.Header.Set("CD1_2", wcs.CD1_2, "Coordinate transformation matrix element")
	fit.Header.Set("CD2_1", wcs.CD2_1, "Coordinate transformation matrix element")
	fit.Header.Set("CD2_2", wcs.CD2_2, "Coordinate transformation matrix element")

	// Attempt to write the WCS solution to the FITS file:
	buf, err := fit.WriteToBuffer()
	if err != nil {
		fmt.Println("failed to write to buffer:", err)
		return err
	}

	// Get the filepath stem where the file is located (e.g. "./samples/astrometry" from
	// "./samples/astrometry.fits"):
	wcsOutputFileStem := getFilePathStem(params.InputFile)

	// Join directory with the new filename and extension for the FITS output file:
	outputFile, err := os.Create(fmt.Sprintf("%s.wcs.fits", wcsOutputFileStem))

	if err != nil {
		fmt.Println("failed to create output file:", err)
		return err
	}

	// Attempt to write the buffer to the output file:
	_, err = buf.WriteTo(outputFile)
	if err != nil {
		fmt.Println("failed to write to output file:", err)
		return err
	}

	// Join directory with the new filename and extension for the JSON output file:
	wcsOutputFile, err := os.Create(fmt.Sprintf("%s.wcs.json", wcsOutputFileStem))
	if err != nil {
		fmt.Println("failed to create output file:", err)
		return err
	}

	// Defer closing the output file:
	defer wcsOutputFile.Close()

	// Attempt to write the WCS solution to the output file:
	encoder := json.NewEncoder(wcsOutputFile)
	// Set the indentation for the JSON encoder:
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(wcs); err != nil {
		// log.Fatalf("failed to encode WCS solution to JSON: %v", err)
		return err
	}

	fmt.Printf("Solution written to: %s\n", outputFile.Name())

	// Return nil if the solver ran successfully:
	return nil
}

/*****************************************************************************************************************/
