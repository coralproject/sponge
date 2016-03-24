package cmd

import (
	"github.com/coralproject/sponge/pkg/report"
	"github.com/spf13/cobra"
)

// sponge index create --type

var (
	reportCmd = &cobra.Command{
		Use:   "show",
		Short: "Read the report on errors",
		Long:  `Read the report generated on errors.`,
		Run:   addReadReport,
	}

	dbnameFlag string
)

const (
	defaultdbnameFlag = "report.db"
)

func init() {
	RootCmd.PersistentFlags().StringVar(&dbnameFlag, "dbname", defaultdbnameFlag, "set the name for the db to read")

	RootCmd.AddCommand(reportCmd)
}

func addReadReport(cmd *cobra.Command, args []string) {

	report.Init("reading", dbnameFlag)
	report.Print()

}
