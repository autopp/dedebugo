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
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type DeniedList map[string][]string

type Inspector struct {
	DeniedList DeniedList
}

// Inspect returns the ast.FileSet and all matched CallExpr nodes in given source.
func (i *Inspector) Inspect(filename string) (*token.FileSet, []*ast.CallExpr, error) {
	isTest := strings.HasSuffix(filename, "_test.go")
	cfg := &packages.Config{
		Mode:  packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: isTest,
	}
	pkgs, err := packages.Load(cfg, filename)
	if err != nil {
		return nil, nil, err
	}

	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return nil, nil, err
	}

	pkg, err := findPackageOfFile(pkgs, absFilename)
	if err != nil {
		return nil, nil, err
	}

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

func findPackageOfFile(pkgs []*packages.Package, filename string) (*packages.Package, error) {
	for _, pkg := range pkgs {
		for _, f := range pkg.CompiledGoFiles {
			if f == filename {
				return pkg, nil
			}
		}
	}

	return nil, fmt.Errorf("package which contains %s is not found", filename)
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
	o := v.typeInfo.ObjectOf(s.Sel)
	if o == nil {
		return false
	}
	p := o.Pkg()
	methods, ok := v.deniedList[p]
	if !ok {
		return false
	}

	_, ok = methods[s.Sel.Name]
	return ok
}
