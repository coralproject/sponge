package cmd

import (
	"github.com/coralproject/sponge/pkg/sponge"
	"github.com/spf13/cobra"
)

// sponge index create --type

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Work with indexes in the coral database",
	Long:  `Create the indexes configured in the strategy file.`,
	Run:   addCreateIndexes,
}

func init() {
	RootCmd.AddCommand(indexCmd)
}

func addCreateIndexes(cmd *cobra.Command, args []string) {

	//log.Dev(u, "main", "Start")

	sponge.CreateIndex(typeFlag)

	//log.Dev(u, "main", "Complete")
}
