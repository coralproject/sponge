package sponge

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/pkg/coral"
	"github.com/coralproject/sponge/pkg/fiddler"
	uuidimported "github.com/pborman/uuid"
)

var (
	oStrategy        string
	oPillarURL       string
	oPollingInterval string
)

func setup() {

	// Save original enviroment variables
	oStrategy = os.Getenv("STRATEGY_CONF")
	oPillarURL = os.Getenv("PILLAR_URL")
	oPollingInterval = os.Getenv("POLLING_INTERVAL")

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	// MOCK STRATEGY CONF
	strategyConf := "../../tests/strategy_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

}

func teardown() {

	// recover the environment variables

	e := os.Setenv("STRATEGY_CONF", oStrategy)
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	e = os.Setenv("PILLAR_URL", oPillarURL)
	if e != nil {
		fmt.Println("It could not setup the mock pillar url variable")
	}

	e = os.Setenv("POLLING_INTERVAL", oPollingInterval)
	if e != nil {
		fmt.Println("It could not setup the mock polling interval variable")
	}
}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

// Signature: process(modelName string, data []map[string]interface{})
func TestProcess(t *testing.T) {

	u := uuidimported.New()

	fiddler.Init(u)
	coral.Init(u)

	modelName := "comment"
	var data []map[string]interface{}
	reportOnFailedRecords := false

	AddOptions(999, 0, "", "", "", false, reportOnFailedRecords, "", 5)

	process(modelName, data)

	// check data is sent to pillar with the right transformations
}
