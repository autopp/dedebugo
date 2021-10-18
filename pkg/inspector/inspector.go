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
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type DeniedList map[string][]string

type Inspector struct {
	DeniedList DeniedList
}

// Inspect returns the ast.FileSet and all matched CallExpr nodes in given source.
// if src != nil, it should be []byte, string, or io.Reader and filename is used for position infomation only.
// if src == nil, Inspect read from filesystem.
func (i *Inspector) Inspect(filename string, src interface{}) (*token.FileSet, []*ast.CallExpr, error) {
	cfg := &packages.Config{
		Mode: packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, filename)
	if err != nil {
		return nil, nil, err
	}
	pkg := pkgs[0]

	converted := convertDeniedList(i.DeniedList, pkg)
	v := &visitor{pkg.TypesInfo, converted, []*ast.CallExpr{}}
	astF := pkg.Syntax[0]
	ast.Walk(v, astF)

	return pkg.Fset, v.mached, nil
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

func convertDeniedList(deniedList DeniedList, pkg *packages.Package) map[*types.Package]map[string]struct{} {
	converted := map[*types.Package]map[string]struct{}{}

	for _, p := range pkg.Types.Imports() {
		if methods, ok := deniedList[p.Path()]; ok {
			methodSet := map[string]struct{}{}
			for _, method := range methods {
				methodSet[method] = struct{}{}
			}
			converted[p] = methodSet
		}
	}

	return converted
}

type visitor struct {
	typeInfo   *types.Info
	deniedList map[*types.Package]map[string]struct{}
	mached     []*ast.CallExpr
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}

	switch n := node.(type) {
	case *ast.CallExpr:
		if s, ok := n.Fun.(*ast.SelectorExpr); ok && v.isDeniedSel(s) {
			v.mached = append(v.mached, n)
		}
	case *ast.GenDecl:
		if n.Tok != token.VAR {
			return nil
		}
	}

	return v
}

func (v *visitor) isDeniedSel(s *ast.SelectorExpr) bool {
	p := v.typeInfo.ObjectOf(s.Sel).Pkg()
	methods, ok := v.deniedList[p]
	if !ok {
		return false
	}

	_, ok = methods[s.Sel.Name]
	return ok
}
