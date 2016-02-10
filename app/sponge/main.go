/*
Package main
Import external source database into local source and transform it
*/
package main

import (
	"flag"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/pkg/sponge"
	"github.com/pborman/uuid"
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

	createindexFlag bool
)

const (
	limitDefault            = 9999999999
	offsetDefault           = 0
	orderbyDefault          = ""
	importonlyfailedDefault = ""
	errorsfileDefault       = "failed_import.csv"
	tableDefault            = ""
)

// Init initialize log, get flag variables, initialize report
func Init() {

	// Logs

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	// Flags

	flag.IntVar(&limitFlag, "limit", limitDefault, "-limit= Number of rows that we are going to import at a time")
	flag.IntVar(&offsetFlag, "offset", offsetDefault, "-offset= Offset for the sql query")
	flag.StringVar(&orderbyFlag, "orderby", orderbyDefault, "-orderby= Order by field of the query on external source")

	flag.StringVar(&importonlyfailedFlag, "onlyfails", importonlyfailedDefault, "-onlyfails Import only the failed documents recorded in report")
	flag.StringVar(&errorsfileFlag, "errors", errorsfileDefault, "-errors Set the path to the file path where to record errors to")

	flag.StringVar(&tableFlag, "type", tableDefault, "-type Import only that type  or types of data. Separate types by ','.")

	flag.Parse()

}

func main() {

	// Generate UUID to use with the logs
	u := uuid.New()

	Init()

	log.Dev(u, "main", "Start")

	sponge.Init(u)

	sponge.CreateIndex(tableFlag)

	sponge.Import(limitFlag, offsetFlag, orderbyFlag, tableFlag, importonlyfailedFlag, errorsfileFlag)

	log.Dev(u, "main", "Complete")
}
