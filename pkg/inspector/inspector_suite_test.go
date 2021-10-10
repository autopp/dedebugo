package inspector

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInspector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inspector Suite")
}
