/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package indexer

/*****************************************************************************************************************/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/healpix"
	"github.com/observerly/skysolve/pkg/index"
	"github.com/spf13/cobra"
)

/*****************************************************************************************************************/

var (
	NSide  int
	Scheme string
)

/*****************************************************************************************************************/

var IndexCommand = &cobra.Command{
	Use:   "indexer",
	Short: "indexer",
	Long:  "indexer",
	Run: func(cmd *cobra.Command, args []string) {
		var scheme healpix.Scheme

		// Extract the scheme from the user input:
		switch strings.ToUpper(Scheme) {
		case "NESTED":
			// Set the scheme to the nested scheme:
			scheme = healpix.NESTED
		case "RING":
			// Set the scheme to the ring scheme:
			scheme = healpix.RING
		default:
			// Set the scheme to the nested scheme by default:
			scheme = healpix.NESTED
		}

		params := RunIndexerParams{
			NSide:  NSide,
			Scheme: scheme,
		}

		// Attempt to run the indexer with the given parameters:
		err := RunIndexer(params)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

/*****************************************************************************************************************/

func init() {
	// Add the nside flag tp the indexer command for the number of sides for the HealPIX scheme:
	// example usage: --nside 256 or -n 256
	IndexCommand.Flags().IntVarP(
		&NSide,
		"nside",
		"n",
		2,
		"The number of sides for the HealPIX scheme",
	)
	IndexCommand.MarkFlagRequired("nside")

	// Add the scheme flag to the indexer command for the HealPIX scheme:
	// example usage: --scheme NESTED or -s NESTED
	IndexCommand.Flags().StringVarP(
		&Scheme,
		"scheme",
		"s",
		"NESTED",
		"The HealPIX scheme to use",
	)
}

/*****************************************************************************************************************/

// Track created files for rollback
var createdFilePaths []string

/*****************************************************************************************************************/

type RunIndexerParams struct {
	NSide  int
	Scheme healpix.Scheme
}

/*****************************************************************************************************************/

func RunIndexer(params RunIndexerParams) error {
	// Create a new SIMBAD service client:
	service := catalog.NewCatalogService(catalog.GAIA, catalog.Params{
		Limit:     8,  // Limit the number of records to 8
		Threshold: 16, // Limiting Magntiude, filter out any stars that are magnitude 16 or above (fainter)
	})

	// Number of sides for the HealPIX:
	sides := params.NSide

	// Create a new HealPIX instance with the given sides:
	healPix := healpix.NewHealPIX(sides, params.Scheme)

	// Get the number of pixels for the given sides:
	pixels := healPix.GetNumberOfPixels()

	// Setup signal handling for cleanup on interrupt:
	signalChannel := make(chan os.Signal, 1)

	// Listen for interrupt signals on the SIGTERM system call, or os.Interrupt:
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// Listen for interrupt signal in a concurrent goroutine, and rollback if received:
	go func() {
		<-signalChannel
		fmt.Println("\nInterrupt received. Rolling back...")
		rollback(createdFilePaths)
		os.Exit(1)
	}()

	// Create a new indexer instance:
	indexer := index.NewIndexer(*healPix, *service)

	// Generate quads for each pixel in the HealPIX grid for the given sides parameter:
	for pixel := 0; pixel <= pixels; pixel++ {
		// Generate quads for the given pixel:
		quads, err := indexer.GenerateQuadsForPixel(pixel)

		if err != nil {
			fmt.Printf("failed to generate quads: %v", err)
			return err
		}

		quadsMarshalled, err := json.Marshal(quads)
		if err != nil {
			log.Fatal(err)
		}

		// Create the directory structure: indexes/<SIDES>
		directoryPath := filepath.Join("indexes", fmt.Sprint(sides))
		if err := os.MkdirAll(directoryPath, 0755); err != nil {
			log.Fatal("Failed to create directory:", err)
		}

		// File path: indexes/<SIDES>/<PIXEL_INDEX>.json
		filePath := filepath.Join(directoryPath, fmt.Sprintf("%d.json", pixel))

		// Write the JSON data to the file
		if err := os.WriteFile(filePath, quadsMarshalled, 0644); err != nil {
			log.Fatal("Failed to write file:", err)
		}

		createdFilePaths = append(createdFilePaths, filePath)

		fmt.Printf("Quads saved to %s\n", filePath)
	}

	// Return nil if the indexer ran successfully:
	return nil
}

/*****************************************************************************************************************/

// rollbackFiles deletes created files in case of failure or interruption
func rollback(filepaths []string) {
	for _, file := range filepaths {
		if err := os.Remove(file); err != nil {
			fmt.Printf("Warning: Failed to remove file %s: %v\n", file, err)
		} else {
			fmt.Printf("Rolled back: %s\n", file)
		}
	}
}

/*****************************************************************************************************************/
