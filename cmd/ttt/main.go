// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
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
	day
	report

The options are:

`

func main() {
	app := new(appCtx)
	app.InitDefaults()
	app.ParseCmd()

	// setup
	if isFile(app.filename) {
		app.InitDb()
	} else {
		app.InitEmptyDb()
	}

	// exec:
	switch app.cmd {
	case startCmd:
		app.Start()
	case endCmd:
		app.End()
	case statusCmd:
		app.Status()
	case recordDayCmd:
		app.RecordDay()
	case reportCmd:
		app.Report(time.Now().Add(time.Hour*24*-365), time.Now())
	case defaultCmd:
		app.Report(time.Now().Add(time.Hour*24*-365), time.Now())
		app.Status()
	}
}

type command string

const (
	startCmd     command = "start"
	endCmd       command = "end"
	statusCmd    command = "status"
	recordDayCmd command = "day"
	reportCmd    command = "report"
	defaultCmd   command = ""
)

func (app *appCtx) ParseCmd() {
	user, err := user.Current()
	if err != nil {
		quitErr(err)
	}
	defaultFilename := filepath.Join(user.HomeDir, ".ttt.sqlite")

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
			default:
				return errors.New("unrecognized duration format")
			}
			return nil
		})
	flag.Func("time", "specify record start or end time with 1m resolution. format=now|HH:mm|YYYY-MM-DD HH:mm (default \"now\" with specified resolution)",
		func(arg string) error {
			if arg == "now" {
				app.opTime = time.Now().Round(time.Minute)
			} else if len(arg) > 5 {
				app.opTime, err = time.Parse("2006-01-02 15:04", arg)
			} else {
				app.opTime, err = time.Parse("15:04", arg)
			}
			return err
		})
	flag.DurationVar(&app.resolution, "resolution", 15*time.Minute, "specify the resolution the input timestamps should be rounded to.")
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 1 {
		app.cmd = command(flag.Arg(0))
	} else if flag.NArg() > 1 {
		quitParamErr("Unrecognized arguments after command: " + fmt.Sprint(flag.Args()[1:]))
	}
}
