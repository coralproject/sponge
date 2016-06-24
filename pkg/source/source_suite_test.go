package source_test

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Source Suite")
}

const (
	cfgLoggingLevel = "LOGGING_LEVEL"
	DATABASE        = "coral_test"
)

var oStrategy string

var _ = BeforeSuite(func() {
	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel)

	// Save the strategy conf to set it back on tear down
	oStrategy = os.Getenv("STRATEGY_CONF")

	// SPIN UP TEST DATABASES
	createTestMongoDB()
})

var _ = AfterSuite(func() {
	if e := os.Setenv("STRATEGY_CONF", oStrategy); e != nil {
		Expect(e).NotTo(HaveOccurred())
	}

	// DISMANTLE TEST DATABASES
	removeMongoDB()
})
