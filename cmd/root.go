/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package cmd

/*****************************************************************************************************************/

import (
	"github.com/observerly/skysolve/internal/solver"
	"github.com/spf13/cobra"
)

/*****************************************************************************************************************/

var rootCommand = &cobra.Command{
	Use:   "solve",
	Short: "SkySolve CLI is a command-line tool for performing an astrometric plate solve on your astronomical images.",
	Long:  "SkySolve CLI is a command-line tool for performing an astrometric plate solve on your astronomical images.",
}

/*****************************************************************************************************************/

func init() {
	rootCommand.AddCommand(solver.AstrometryCommand)
}

/*****************************************************************************************************************/

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}

/*****************************************************************************************************************/
