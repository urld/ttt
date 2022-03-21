// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/urld/ttt"
)

type appCtx struct {
	filename string
	*ttt.TimeTrackingDb

	durationFmt
	resolution time.Duration

	cmd    command
	opTime time.Time
}
type durationFmt int

const (
	clock durationFmt = iota
	hours
	decimal
)

func (app *appCtx) InitDefaults() {
	app.cmd = defaultCmd
	app.opTime = app.cleanDts(time.Now())
}

func (app *appCtx) InitDb() {
	var err error
	app.TimeTrackingDb, err = ttt.LoadDb(app.filename)
	quitErr(err)
}
func (app *appCtx) InitEmptyDb() {
	var err error
	app.TimeTrackingDb, err = ttt.CreateDb(app.filename)
	quitErr(err)
}
func (app *appCtx) Start() {
	err := app.StartRecord(app.opTime)
	quitErr(err)
}
func (app *appCtx) End() {
	err := app.EndRecord(app.opTime)
	quitErr(err)
}

func (app *appCtx) Status() {
	rec, err := app.GetCurrentRecord()
	quitErr(err)
	if rec == (ttt.Record{}) {
		fmt.Println("No records have been tracked yet. Use `ttt start` to begin a new record.")
	} else if rec.Active() {
		duration := time.Now().Sub(*rec.Start)
		fmt.Printf("Current record is active since %s, (%s)\n", fmtTime(*rec.Start), fmtDurationTrim(duration, app.durationFmt))
		fmt.Println("Use `ttt end` to end the record.")
	} else {
		duration := rec.End.Sub(*rec.Start)
		fmt.Println("No record is active at the moment. Use `ttt start` to begin a new record.")
		fmt.Printf("Last record was active from %s to %s, (%s)\n", fmtTime(*rec.Start), fmtTime(*rec.End), fmtDurationTrim(duration, app.durationFmt))
	}
}

func (app *appCtx) Report(reportStart, reportEnd time.Time) {
	days, err := app.GetBusinessDays(reportStart, reportEnd)
	quitErr(err)

	var saldo time.Duration
	var totalWorked, totalWork time.Duration
	var aggrWorked, aggrWork time.Duration
	var previous ttt.BusinessDay

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.SetStyle(table.StyleColoredBright)
	tw.Style().Format.Header = text.FormatTitle
	tw.Style().Format.Footer = text.FormatTitle
	tw.SetTitle("Report: %s to %s", fmtDate(reportStart), fmtDate(reportEnd))
	tw.AppendHeader(table.Row{"week", "date", "worked", "delta", "saldo"})

	for _, d := range days {
		if d.ISOWeek != previous.ISOWeek && previous.ISOWeek != 0 {
			//print summary
			tw.AppendRow(app.weekRow(previous.ISOWeek, aggrWorked, aggrWorked-aggrWork))
			tw.AppendSeparator()
			//reset
			aggrWorked = 0
			aggrWork = 0
		}
		// aggregate
		delta := d.WorkedHours - d.WorkHours
		saldo += delta
		totalWorked += d.WorkedHours
		totalWork += d.WorkHours
		aggrWorked += d.WorkedHours
		aggrWork += d.WorkHours
		// print day
		tw.AppendRow(app.dayRow(d.Date, d.WorkedHours, delta, saldo))
		previous = d
	}
	tw.AppendRow(app.weekRow(previous.ISOWeek, aggrWorked, aggrWorked-aggrWork))
	tw.AppendFooter(app.totalRow(totalWorked, totalWorked-totalWork))
	tw.Render()
}

func (app *appCtx) dayRow(date time.Time, worked, delta, saldo time.Duration) table.Row {
	return table.Row{
		"",
		fmtDate(date),
		fmtDuration(worked, app.durationFmt),
		fmtDuration(delta, app.durationFmt),
		fmtDuration(saldo, app.durationFmt),
	}
}

func (app *appCtx) weekRow(week int, worked, delta time.Duration) table.Row {
	return table.Row{
		week,
		"",
		fmtDuration(worked, app.durationFmt),
		fmtDuration(delta, app.durationFmt),
		"",
	}
}

func (app *appCtx) totalRow(worked, saldo time.Duration) table.Row {
	return table.Row{
		"Total",
		"",
		fmtDuration(worked, app.durationFmt),
		"",
		fmtDuration(saldo, app.durationFmt),
	}
}

func (app *appCtx) cleanDts(dts time.Time) time.Time {
	return dts.Round(app.resolution)
}
