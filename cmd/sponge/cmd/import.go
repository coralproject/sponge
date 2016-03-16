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
	importonlyfailedFlag bool
	localErrorsDBFlag    string

	// Query on the indicated field
	queryFlag string
)

const (
	defaultLimit            = 9999999999
	defaultOffset           = 0
	defaultOrderBy          = ""
	defaultImportonlyfailed = false
	defaultlocalErrorsDB    = "report.db"
	defaultQuery            = ""
)

func init() {
	RootCmd.PersistentFlags().IntVar(&limitFlag, "limit", defaultLimit, "number of rows that we are going to import (default is 9999999999)")
	RootCmd.PersistentFlags().IntVar(&offsetFlag, "offset", defaultOffset, "offset for rows to import (default is 0)")
	RootCmd.PersistentFlags().StringVar(&orderbyFlag, "orderby", defaultOrderBy, "order by field on the external source (default is not ordered)")
	RootCmd.PersistentFlags().StringVar(&queryFlag, "query", defaultQuery, "query on the external table (where condition on mysql, query document on mongodb). It only works with a specific --type. Example updated_date >'2003-12-31 01:02:03'")

	RootCmd.PersistentFlags().BoolVar(&importonlyfailedFlag, "onlyfails", defaultImportonlyfailed, "import only the the records that failed in the last import(default is import all)")

	RootCmd.PersistentFlags().StringVar(&localErrorsDBFlag, "errors", defaultlocalErrorsDB, "set the file path for the report on errors (default is report.db)")

	RootCmd.AddCommand(importCmd)
}

func addImport(cmd *cobra.Command, args []string) {

	sponge.Import(limitFlag, offsetFlag, orderbyFlag, queryFlag, typeFlag, importonlyfailedFlag, localErrorsDBFlag)
}
