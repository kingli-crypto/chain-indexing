package crossfire_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCrossfire(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crossfire DryRun Suite")
}
