// It provides teh coral project sponge CLI  platform.
package main

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/cmd/sponge/cmd"
	"github.com/coralproject/sponge/pkg/sponge"
	"github.com/pborman/uuid"
)

const (
	cfgLoggingLevel = "LOGGING_LEVEL"
	cfgStrategyFile = "STRATEGY_FILE"
	cfgPILLARURL    = "PILLAR_URL"
)

func main() {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	// Generate UUID to use with the logs
	uid := uuid.New()

	if err := sponge.Init(uid); err != nil {
		log.Error(uid, "main", err, "Unable to initialize configuration.")
		os.Exit(-1)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Error(uid, "main", err, "Unable to execute the command.")
		os.Exit(-1)
	}
}
