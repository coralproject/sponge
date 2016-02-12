package cmd

import (
	"github.com/coralproject/sponge/pkg/sponge"
	"github.com/spf13/cobra"
)

// sponge import all --offset --limit --orderby --type
// sponge import onlyfails

var (
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import data to the coral database",
		Long:  `Import the data from the external source configured in the strategy file.`,
		Run:   addImport,
	}

	// Limit on query
	limitFlag  int
	offsetFlag int

	// Order by field
	orderbyFlag string

	// Import from report on failed records (or not)
	importonlyfailedFlag string
	errorsfileFlag       string
)

const (
	defaultLimit            = 9999999999
	defaultOffset           = 0
	defaultOrderBy          = ""
	defaultImportonlyfailed = ""
	defaultErrorsfile       = "failed_import.csv"
)

func init() {
	RootCmd.Flags().IntVar(&limitFlag, "limit", defaultLimit, "number of rows that we are going to import (default is 9999999999)")
	RootCmd.Flags().IntVar(&offsetFlag, "offset", defaultOffset, "offset for rows to import (default is 0)")
	RootCmd.Flags().StringVar(&orderbyFlag, "orderby", defaultOrderBy, "order by field on the external source (default is not ordered)")

	RootCmd.Flags().StringVar(&importonlyfailedFlag, "onlyfails", defaultImportonlyfailed, "import only the the records that failed in the last import(default is import all)")

	RootCmd.Flags().StringVar(&errorsfileFlag, "errors", defaultErrorsfile, "set the file path for the report on errors (default is failed_import.csv)")

	RootCmd.AddCommand(importCmd)
}

func addImport(cmd *cobra.Command, args []string) {

	sponge.Import(limitFlag, offsetFlag, orderbyFlag, typeFlag, importonlyfailedFlag, errorsfileFlag)
}
