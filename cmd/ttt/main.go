// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os/user"
	"path/filepath"
	"time"
)

var usage string = `ttt - time tracking terminal

Usage:
	ttt <command> [options]

The commands are:

	start
	end
	status
	report

The options are:

`

func main() {
	cmd, app := parseCmd()

	// setup:u
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
	case statusCmd:
		app.Status()
	case reportCmd:
		app.Report(time.Now().Add(time.Hour*24*-365), time.Now())
	case defaultCmd:
		app.Report(time.Now().Add(time.Hour*24*-365), time.Now())
		app.Status()
	}
}

type command string

const (
	startCmd   command = "start"
	endCmd     command = "end"
	statusCmd  command = "status"
	reportCmd  command = "report"
	defaultCmd command = ""
)

func parseCmd() (command, appCtx) {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".ttt.sqlite")

	var app appCtx
	app.filename = *flag.String("file", defaultFilename, "specify the ttt database")
	flag.Func("duration", "the `format` that is used to print durations: clock|hours|decimal (default \"clock\")",
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
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	cmd := defaultCmd
	if flag.NArg() >= 1 {
		cmd = command(flag.Arg(0))
	}
	return cmd, app
}
