// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func quitErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func quitParamErr(err string) {
	fmt.Println(err)
	flag.Usage()
	os.Exit(2)
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
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
