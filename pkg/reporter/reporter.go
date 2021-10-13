// Copyright (C) 2021 Akira Tanimura (@autopp)
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reporter

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io"
)

type Reporter interface {
	Report(w io.Writer, fset *token.FileSet, nodes []*ast.CallExpr)
}

type reporter struct{}

func New() Reporter {
	return &reporter{}
}

func (r *reporter) Report(w io.Writer, fset *token.FileSet, nodes []*ast.CallExpr) {
	for _, node := range nodes {
		pos := fset.Position(node.Pos())
		w.Write([]byte(pos.String() + " "))
		printer.Fprint(w, fset, node)
		fmt.Fprintln(w, "")
	}
}
