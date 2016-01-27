/*
Package main
Import external source database into local source and transform it
*/
package main

import (
	"flag"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/coralproject/sponge/pkg/log"
	"github.com/coralproject/sponge/pkg/sponge"
)

var (
	// Limit on query
	limitFlag  int
	offsetFlag int

	// Order by field
	orderbyFlag string

	// Import from report on failed records (or not)
	importonlyfailedFlag bool

	tableFlag string

	createindexFlag bool
)

const (
	limitDefault            = 9999999999
	offsetDefault           = 0
	orderbyDefault          = ""
	importonlyfailedDefault = false
	tableDefault            = ""
	createindexDefault      = false
)

// Initialize log, get flag variables, initialize report
func init() {

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	flag.IntVar(&limitFlag, "limit", limitDefault, "-limit= Number of rows that we are going to import at a time")
	flag.IntVar(&offsetFlag, "offset", offsetDefault, "-offset= Offset for the sql query")
	flag.StringVar(&orderbyFlag, "orderby", orderbyDefault, "-orderby= Order by field of the query on external source")
	flag.BoolVar(&importonlyfailedFlag, "onlyfails", importonlyfailedDefault, "-onlyfails Import only the failed documents recorded in report")
	flag.StringVar(&tableFlag, "type", tableDefault, "Import only this type")

	flag.BoolVar(&createindexFlag, "index", createindexDefault, "-index Create index for all the tables or only for the one specifies")

	flag.Parse()

}

func main() {
	log.Dev("cmd", "main", "Start")
	sponge.Import(limitFlag, offsetFlag, orderbyFlag, tableFlag, importonlyfailedFlag)

	if createindexFlag {
		sponge.CreateIndex(tableFlag)
	}

	log.Dev("cmd", "main", "Complete")
}
