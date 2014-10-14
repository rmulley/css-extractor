package main

import (
	"testing"
) //import

func TestRemoveStyleTags(t *testing.T) {
	var (
		output string
		lines  []string
	) //var

	lines = []string{
		"<div id=\"testId1\" class=\"testClass1\" style=\"margin-top: 50px;\">",
	} //[]string

	output = RemoveStyleTags(lines)

	if output == "<div id=\"testId1\" class=\"testClass1\">" {
		t.Fail()
	} //if
} //TestremoveStyleTags
