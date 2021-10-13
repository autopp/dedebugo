package reporter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var src = `package main

import (
	"fmt"
)

func main() {
	fmt.Print("hello ")
	fmt.Println("world")
	fmt.Printf("from %s\n", os.Args[0])
}
`

type visitor struct {
	nodes []*ast.CallExpr
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}

	if c, ok := node.(*ast.CallExpr); ok {
		v.nodes = append(v.nodes, c)
		fmt.Printf("%#v\n", v.nodes)
	}

	return v
}

var _ = g.Describe("reporter", func() {
	g.Describe("Report()", func() {
		filename := "sample.go"
		fset := token.NewFileSet()
		astf, _ := parser.ParseFile(fset, filename, src, parser.ParseComments)
		v := &visitor{[]*ast.CallExpr{}}
		ast.Walk(v, astf)
		nodes := v.nodes

		g.It("write position and code of each calls", func() {
			r := New()
			w := &bytes.Buffer{}
			r.Report(w, fset, nodes)

			expected := `sample.go:8:2 fmt.Print("hello ")` + "\n" +
				`sample.go:9:2 fmt.Println("world")` + "\n" +
				`sample.go:10:2 fmt.Printf("from %s\n", os.Args[0])` + "\n"
			Expect(w.String()).To(Equal(expected))
		})
	})
})
