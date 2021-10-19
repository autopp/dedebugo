package inspector

import (
	"bytes"
	"go/ast"
	"go/printer"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inspector", func() {
	Describe("Inspect()", func() {
		deniedList := DeniedList{
			"fmt": []string{"Print", "Printf", "Println"},
		}
		filename := "testdata/sample.go"
		wd, _ := os.Getwd()
		fullpath := filepath.Join(wd, filename)

		It("returns FileSet and matched CallExpr nodes", func() {
			i := &Inspector{deniedList}
			fset, nodes, err := i.Inspect(filename)

			Expect(err).NotTo(HaveOccurred())
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
				Equal([][]string{
					{fullpath + ":23:2", `fmt.Print("hello ")`},
					{fullpath + ":24:2", `fmt.Println("world")`},
					{fullpath + ":25:2", `fmt.Printf("from %s\n", os.Args[0])`},
				}),
			))
		})
	})
})
