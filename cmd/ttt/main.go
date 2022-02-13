// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os/user"
	"path/filepath"
	"time"
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
		app.Report(time.Now().Add(time.Hour*24*-365), time.Now())
	}
}

type command string

const (
	startCmd  command = "start"
	endCmd    command = "end"
	reportCmd command = "report"
)

func parseCmd() (command, appCtx) {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".ttt.sqlite")

	var app appCtx
	app.filename = *flag.String("file", defaultFilename, "specify the ttt database")
	flag.Func("duration", "the format that is used to print durations: clock|hours|decimal",
		func(arg string) error {
			switch arg {
			case "clock":
				app.durationFmt = clock
			case "hours":
				app.durationFmt = hours
			case "decimal":
				app.durationFmt = decimal
			}
			return nil
		})
	flag.Parse()

	cmd := reportCmd
	if flag.NArg() >= 1 {
		cmd = command(flag.Arg(0))
	}
	return cmd, app
}
