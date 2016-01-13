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

// Limit on query
var limitFlag int
var offsetFlag int

// Order by field
var orderbyFlag string

// Import from report on failed records (or not)
var importonlyfailedFlag bool

const (
	limitDefault            = 9999999999
	offsetDefault           = 0
	orderbyDefault          = ""
	importonlyfailedDefault = false
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

	flag.Parse()

}

func main() {

	log.Dev("main", "main", "Start")

	sponge.Import(limitFlag, offsetFlag, orderbyFlag, importonlyfailedFlag)

	log.Dev("shutdown", "main", "Complete")
}
