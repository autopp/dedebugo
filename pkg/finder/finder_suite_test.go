package finder

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFinder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Finder Suite")
}
