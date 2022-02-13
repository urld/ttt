// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/urld/ttt"
)

// appCtx provides means to interact with a passmgr store via terminal.
type appCtx struct {
	filename string
	*ttt.TimeTrackingDb

	durationFmt
}
type durationFmt int

const (
	clock durationFmt = iota
	hours
	decimal
)

func (app *appCtx) Init() {
	var err error
	app.TimeTrackingDb, err = ttt.LoadDb(app.filename)
	quitErr(err)
}
func (app *appCtx) InitEmpty() {
	var err error
	app.TimeTrackingDb, err = ttt.CreateDb(app.filename)
	quitErr(err)
}
func (app *appCtx) Start() {
	err := app.StartRecord(time.Now())
	quitErr(err)
}
func (app *appCtx) End() {
	err := app.EndRecord(time.Now())
	quitErr(err)
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
