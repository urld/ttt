// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func quitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func calcFilename(filename string) string {
	if isDir(filename) {
		return filepath.Join(filename, ".passmgr_store")
	}
	return filename
}

func isDir(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
