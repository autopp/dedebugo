package reporter

import (
	"testing"

	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestReporter(t *testing.T) {
	RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Reporter Suite")
}
