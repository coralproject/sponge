package fiddler_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFiddler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fiddler Suite")
}
