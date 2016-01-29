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
	importonlyfailedFlag string
	errorsfileFlag       string

	tableFlag string
)

const (
	limitDefault            = 9999999999
	offsetDefault           = 0
	orderbyDefault          = ""
	importonlyfailedDefault = ""
	errorsfileDefault       = "failed_import.csv"
	tableDefault            = ""
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
	flag.StringVar(&importonlyfailedFlag, "onlyfails", importonlyfailedDefault, "-onlyfails Import only the failed documents recorded in report")
	flag.StringVar(&errorsfileFlag, "errors", errorsfileDefault, "-errors Set the path to the file path where to record errors to")
	flag.StringVar(&tableFlag, "type", tableDefault, "-type Import only that type of data.")

	flag.Parse()

}

func main() {
	log.Dev("cmd", "main", "Start")

	sponge.Import(limitFlag, offsetFlag, orderbyFlag, tableFlag, importonlyfailedFlag, errorsfileFlag)

	log.Dev("cmd", "main", "Complete")
}
