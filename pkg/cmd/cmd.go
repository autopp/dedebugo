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
	"errors"
	"io"

	"github.com/autopp/dedebugo/pkg/finder"
	"github.com/autopp/dedebugo/pkg/inspector"
	"github.com/autopp/dedebugo/pkg/reporter"
	"github.com/spf13/cobra"
)

var ErrDeniedCallFound = errors.New("denied callls are found")
var defaultExcludedList = []string{".git", "vendor"}

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
			f, err := finder.New(finder.WithExcludedList(defaultExcludedList))
			if err != nil {
				return err
			}
			for _, a := range args {
				f, err := f.FindGoFiles(a)
				if err != nil {
					return err
				}
				gofiles = append(gofiles, f...)
			}

			found := false
			for _, filename := range gofiles {
				fset, nodes, err := i.Inspect(filename)
				if err != nil {
					return err
				}

				if len(nodes) != 0 {
					found = true
				}
				r.Report(cmd.OutOrStderr(), fset, nodes)
			}

			if found {
				return ErrDeniedCallFound
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
