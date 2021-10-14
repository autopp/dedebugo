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

package inspector

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/loader"
)

type DeniedList map[string][]string

type Inspector struct {
	DeniedList DeniedList
}

// Inspect returns the ast.FileSet and all matched CallExpr nodes in given source.
// if src != nil, it should be []byte, string, or io.Reader and filename is used for position infomation only.
// if src == nil, Inspect read from filesystem.
func (i *Inspector) Inspect(filename string, src interface{}) (*token.FileSet, []*ast.CallExpr, error) {
	l := &loader.Config{ParserMode: parser.ParseComments}
	astf, err := l.ParseFile(filename, src)
	if err != nil {
		return nil, nil, err
	}
	l.CreateFromFiles("", astf)
	prog, err := l.Load()
	if err != nil {
		return nil, nil, err
	}

	converted := convertDeniedList(i.DeniedList, prog)
	pkg := prog.Package(astf.Name.Name)
	v := &visitor{pkg, converted, []*ast.CallExpr{}}
	ast.Walk(v, astf)

	return l.Fset, v.mached, nil
}

func DefaultDeniedList() DeniedList {
	return DeniedList{
		"fmt": {
			"Print",
			"Printf",
			"Println",
		},
	}
}

func convertDeniedList(deniedList DeniedList, prog *loader.Program) map[*types.Package]map[string]struct{} {
	converted := map[*types.Package]map[string]struct{}{}
	for pkg, methods := range deniedList {
		methodSet := map[string]struct{}{}
		for _, method := range methods {
			methodSet[method] = struct{}{}
		}
		converted[prog.Package(pkg).Pkg] = methodSet
	}

	return converted
}

type visitor struct {
	pkg        *loader.PackageInfo
	deniedList map[*types.Package]map[string]struct{}
	mached     []*ast.CallExpr
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	if c, ok := node.(*ast.CallExpr); ok {
		if s, ok := c.Fun.(*ast.SelectorExpr); ok && v.isDeniedSel(s) {
			v.mached = append(v.mached, c)
		}
	}
	return v
}

func (v *visitor) isDeniedSel(s *ast.SelectorExpr) bool {
	p := v.pkg.Info.ObjectOf(s.Sel).Pkg()
	methods, ok := v.deniedList[p]
	if !ok {
		return false
	}

	_, ok = methods[s.Sel.Name]
	return ok
}
