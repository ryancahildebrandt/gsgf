// -*- coding: utf-8 -*-

// Created on Sun Mar 16 05:13:48 PM EDT 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"testing"
)

func TestValidateExportDir(t *testing.T) {
	table := []struct {
		p       string
		wantErr bool
	}{
		{p: "", wantErr: true},
		{p: ".", wantErr: false},
		{p: " ", wantErr: true},
		{p: ".txt/", wantErr: true},
		{p: "../a", wantErr: true},
		{p: "out/a", wantErr: true},
		{p: "export/out.txt", wantErr: true},
		{p: "data/data/b", wantErr: true},
		{p: "data/tests/dir0/", wantErr: false},
	}
	for i, test := range table {
		err := ValidateExportDir(test.p)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateExportDir(%v)\nGOT %v\nWANT %v", i, test.p, err, test.wantErr)
		}
	}
}

func TestValidateOutFile(t *testing.T) {
	table := []struct {
		p       string
		wantErr bool
	}{
		{p: "", wantErr: false},
		{p: ".", wantErr: false},
		{p: " ", wantErr: false},
		{p: ".txt", wantErr: false},
		{p: "../a.txt", wantErr: false},
		{p: "out/a", wantErr: true},
		{p: "export/out.txt", wantErr: false},
		{p: "data/data/b", wantErr: true},
		{p: "data/tests/dir0/", wantErr: false},
	}
	for i, test := range table {
		err := ValidateOutFile(test.p)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateOutFile(%v)\nGOT %v\nWANT %v", i, test.p, err, test.wantErr)
		}
	}
}

func TestValidateInFile(t *testing.T) {
	table := []struct {
		p       string
		wantErr bool
	}{
		{p: "", wantErr: true},
		{p: ".", wantErr: true},
		{p: " ", wantErr: true},
		{p: ".jsgf", wantErr: true},
		{p: "../a.jsgf", wantErr: true},
		{p: "data/tests/a.jsgf", wantErr: false},
		{p: "data/tests/dir0/c.jsgf", wantErr: false},
		{p: "data/tests/test0.jsgf", wantErr: false},
	}
	for i, test := range table {
		err := ValidateInFile(test.p)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateInFile(%v)\nGOT %v\nWANT %v", i, test.p, err, test.wantErr)
		}
	}
}
