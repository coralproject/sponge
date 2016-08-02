package sponge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSponge(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sponge Suite")
}
