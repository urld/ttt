// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os/user"
	"path/filepath"
)

type command int

const (
	startCmd command = iota
	endCmd
	reportCmd
	noCmd
)

func main() {
	cmd, app := parseCmd()

	// setup:
	if isFile(app.filename) {
		app.Init()
	} else {
		app.InitEmpty()
	}

	// exec:
	switch cmd {
	case startCmd:
		app.Start()
	case endCmd:
		app.End()
	case reportCmd:
		app.Report()
	}
}

func parseCmd() (command, termApp) {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".ttt.db")

	// cmd parsing:
	filename := flag.String("file", defaultFilename, "specify the ttt database")
	flag.Parse()

	cmd := noCmd
	if flag.NArg() >= 1 {
		switch flag.Arg(0) {
		case "start":
			cmd = startCmd
		case "end":
			cmd = endCmd
		case "report":
			cmd = reportCmd
		}
	}

	return cmd, termApp{filename: calcFilename(*filename)}
}
