package cmd

import (
	"fmt"

	"github.com/coralproject/sponge/pkg/sponge"
	"github.com/spf13/cobra"
)

/*

sponge all

sponge import --all --offset --limit --orderby --type
sponge import --onlyfails

sponge index create --type

sponge version

*/

var (
	typeFlag string
)

// Config environmental variables.
const (
	defaultType = ""
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sponge",
	Short: "Sponge is an importer of data into the Coral database",
	Long: `An importer to get data from external sources into the coral database. It does the needed transofrmation
	to fit int the right schema and create needed indexes. For example:

Sponge is a CLI tool to import data from external sources.
It needs a strategy file that maps the transformation between external schema and the coral database.`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&typeFlag, "type", defaultType, "import or create indexes for only these types of data (default is everything)")

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(allCmd)
}

//* VERSION COMMAND *//

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Sponge",
	Long:  `This is Sponge's version.`,
	Run:   addVersion,
}

func addVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("sponge version: %v\n", sponge.VersionNumber)
}

//* RUNNING EVERYTHING *//

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Import and Create Indexes",
	Long:  `Import and create indexes at the same time`,
	Run:   addAll,
}

func addAll(cmd *cobra.Command, args []string) {

	addCreateIndexes(cmd, args)
	addImport(cmd, args)

}
