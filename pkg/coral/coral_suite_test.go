package coral_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCoral(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Coral Suite")
}
