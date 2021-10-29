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

package finder

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Finder interface {
	FindGoFiles(path string) ([]string, error)
}

type finder struct {
	excludedList []string
}

type Option func(f *finder) error

func New(opts ...Option) (Finder, error) {
	f := &finder{}

	for _, o := range opts {
		if err := o(f); err != nil {
			return nil, err
		}
	}
	return &finder{}, nil
}

func WithExcludedList(exlcludedList []string) Option {
	return func(f *finder) error {
		f.excludedList = exlcludedList
		return nil
	}
}

func (f *finder) FindGoFiles(path string) ([]string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.Mode().IsRegular() {
		if f.isTargetFile(path) {
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

		if d.Type().IsDir() {
			if f.isExcluded(path) {
				return filepath.SkipDir
			}
		} else if d.Type().IsRegular() && f.isTargetFile(path) {
			gofiles = append(gofiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return gofiles, nil
}

func (f *finder) isExcluded(path string) bool {
	for _, p := range f.excludedList {
		if m, ok := filepath.Match(p, path); ok != nil || m {
			return true
		}
	}

	return false
}

func (f *finder) isTargetFile(path string) bool {
	return !f.isExcluded(path) && strings.HasSuffix(path, ".go")
}
