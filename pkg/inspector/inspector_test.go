package inspector

import (
	"bytes"
	"go/ast"
	"go/printer"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inspector", func() {
	Describe("Inspect()", func() {
		deniedList := DeniedList{
			"fmt": []string{"Print", "Printf", "Println"},
		}
		filename := "testdata/sample.go"
		contents, _ := os.ReadFile(filename)
		wd, _ := os.Getwd()
		fullpath := filepath.Join(wd, filename)

		DescribeTable("success cases",
			func(src interface{}) {
				i := &Inspector{deniedList}
				fset, nodes, err := i.Inspect(filename, src)

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
			},
			Entry("with filename only", nil),
			Entry("with filename and []byte", contents),
			Entry("with filename and string", string(contents)),
			Entry("with filename and io.Reader", bytes.NewReader(contents)),
		)
	})
})
