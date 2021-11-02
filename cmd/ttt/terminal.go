// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urld/ttt"
)

// termApp provides means to interact with a passmgr store via terminal.
type termApp struct {
	filename string
	store    ttt.TimeTrackingDb
}

func (t *termApp) Init() {
	var err error
	t.store, err = ttt.LoadDb(t.filename)
	if err != nil {
		quitErr(err)
	}
}
func (t *termApp) InitEmpty() {
	var err error
	t.store, err = ttt.CreateDb(t.filename)
	if err != nil {
		quitErr(err)
	}
}
func (t *termApp) Start() {
	err := t.store.StartRecord(time.Now())
	if err != nil {
		quitErr(err)
	}
}
func (t *termApp) End() {
	err := t.store.EndRecord(time.Now())
	if err != nil {
		quitErr(err)
	}
}
func (t *termApp) Report() {
	records, err := t.store.GetRecords()
	if err != nil {
		quitErr(err)
	}
	for _, rec := range records {
		fmt.Println(rec)
	}
	err = t.store.GetWeekAggr()
	if err != nil {
		quitErr(err)
	}
	err = t.store.GetDayAggr()
	if err != nil {
		quitErr(err)
	}
}

func askConfirm(prompt string, a ...interface{}) bool {
	switch strings.ToLower(ask(prompt+" [Y/n] ", a...)) {
	case "y", "":
		return true
	case "n":
		return false
	default:
		return askConfirm(prompt)
	}
}

func ask(prompt string, a ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(prompt, a...)
	text, err := reader.ReadString('\n')
	if err != nil {
		quitErr(err)
	}
	return strings.TrimRight(text, "\r\n")
}
