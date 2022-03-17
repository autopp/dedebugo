package inspector

import (
	"bytes"
	"go/ast"
	"go/printer"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inspector", func() {
	Describe("Inspect()", func() {
		deniedList := DeniedList{
			"fmt": []string{"Print", "Printf", "Println"},
		}

		DescribeTable("returns FileSet and matched CallExpr nodes",
			func(filename string, expected [][]string) {
				wd, _ := os.Getwd()
				fullpath := filepath.Join(wd, filename)

				i := &Inspector{deniedList}
				fset, nodes, err := i.Inspect(filename)

				Expect(err).NotTo(HaveOccurred())

				for _, x := range expected {
					x[0] = fullpath + x[0]
				}
				Expect(nodes).To(WithTransform(
					func(nodes []*ast.CallExpr) [][]string {
						posAndExprs := make([][]string, 0, len(nodes))
						for _, node := range nodes {
							pos := fset.Position(node.Pos())
							expr := &bytes.Buffer{}
							printer.Fprint(expr, fset, node)
							posAndExprs = append(posAndExprs, []string{pos.String(), expr.String()})
						}
						return posAndExprs
					},
					Equal(expected),
				))
			},
			Entry(
				"with regular file (testdata/sample.go)",
				"testdata/sample.go",
				[][]string{
					{":23:2", `fmt.Print("hello ")`},
					{":24:2", `fmt.Println("world")`},
					{":25:2", `fmt.Printf("from %s\n", os.Args[0])`},
				},
			),
			Entry(
				"with test file (testdata/sample_test.go)",
				"testdata/sample_test.go",
				[][]string{
					{":9:2", `fmt.Print("hello ")`},
					{":10:2", `fmt.Println("world")`},
					{":11:2", `fmt.Printf("from %s\n", t.Name())`},
				},
			),
		)
	})
})
