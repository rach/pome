package main

// https://husobee.github.io/golang/testing/unit-test/2015/06/08/golang-unit-testing.html
// check defer

import (
	"testing"
)

func TestCmdLine(t *testing.T) {
	exit := false
	app.Terminate(func(int) { exit = true })
	//without args
	parseCmdLine([]string{})
	if exit {
		t.Error("Program should have terminated")
	}
	//wrong flag
	exit = false
	parseCmdLine([]string{"-x"})
	if exit {
		t.Error("Program should have terminated")
	}

	//minimum flags
	exit = false
	parseCmdLine([]string{"-x"})
	if exit {
		t.Error("Program should have terminated")
	}
}
