package finder

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Finder", func() {
	DescribeTable("FindGoFiles",
		func(excludedList []string, expected []string) {
			f, _ := New(WithExcludedList(excludedList))
			actual, err := f.FindGoFiles("testdata")
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(expected))
		},
		Entry("with no exclude", []string{}, []string{"testdata/a/foo.go", "testdata/b/c/bar.go", "testdata/main.go"}),
		Entry("with no exclude", []string{"testdata/main.go"}, []string{"testdata/a/foo.go", "testdata/b/c/bar.go"}),
		Entry("with no exclude", []string{"testdata/a"}, []string{"testdata/b/c/bar.go", "testdata/main.go"}),
		Entry("with no exclude", []string{"testdata/*/c"}, []string{"testdata/a/foo.go", "testdata/main.go"}),
	)
})
