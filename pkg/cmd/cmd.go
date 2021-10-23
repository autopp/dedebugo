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

package cmd

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/autopp/dedebugo/pkg/inspector"
	"github.com/autopp/dedebugo/pkg/reporter"
	"github.com/spf13/cobra"
)

func Run(version string, stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	const versionFlag = "version"
	cmd := &cobra.Command{
		Use:           "dedebugo file",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if v, err := cmd.Flags().GetBool(versionFlag); err != nil {
				return err
			} else if v {
				cmd.Println(version)
				return nil
			}

			i := &inspector.Inspector{DeniedList: inspector.DefaultDeniedList()}
			r := reporter.New()

			gofiles := []string{}
			for _, a := range args {
				f, err := findGoFiles(a)
				if err != nil {
					return err
				}
				gofiles = append(gofiles, f...)
			}
			for _, filename := range gofiles {
				fset, nodes, err := i.Inspect(filename)
				if err != nil {
					return err
				}
				r.Report(cmd.OutOrStderr(), fset, nodes)
			}

			return nil
		},
	}

	cmd.Flags().Bool(versionFlag, false, "print version")

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	return cmd.Execute()
}

func findGoFiles(path string) ([]string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.Mode().IsRegular() {
		if strings.HasSuffix(path, ".go") {
			return []string{path}, nil
		}
		return nil, nil
	}

	if !stat.Mode().IsDir() {
		return nil, nil
	}

	gofiles := []string{}

	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() && strings.HasSuffix(path, ".go") {
			gofiles = append(gofiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return gofiles, nil
}
