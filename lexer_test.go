// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:58 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"testing"
)

func TestCaptureString(t *testing.T) {
	lexer := NewJSGFLexer("\"")
	table := []struct {
		s       string
		e       string
		i       bool
		want    string
		wantErr bool
	}{
		{s: " ", e: "a", i: true, want: "", wantErr: true},
		{s: " ", e: "a", i: false, want: "", wantErr: true},
		{s: " ", e: "", i: false, want: " ", wantErr: false},
		{s: " ", e: "", i: false, want: " ", wantErr: false},
		{s: " ", e: "", i: true, want: " ", wantErr: false},
		{s: " ", e: " ", i: false, want: "", wantErr: false},
		{s: " ", e: " ", i: true, want: " ", wantErr: false},
		{s: "abc", e: "c", i: false, want: "abc", wantErr: false},
		{s: "abc", e: "c", i: true, want: "abc", wantErr: false},
		{s: "()()", e: ")", i: false, want: "(", wantErr: false},
		{s: "()()", e: ")", i: true, want: "()", wantErr: false},
	}
	for i, test := range table {
		stream := lexer.ParseString(test.s)
		got, err := CaptureString(stream, test.e, test.i)
		if got != test.want {
			t.Errorf("test %v: captureString(%v, %v, %v)\nGOT %v\nWANT %v", i, test.s, test.e, test.i, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: captureString(%v, %v, %v)\nGOT %v\nWANT %v", i, test.s, test.e, test.i, err, test.wantErr)
		}
	}
}

func TestValidateLexerString(t *testing.T) {
	tests := []struct {
		s       string
		wantErr bool
	}{
		{s: "", wantErr: true},
		{s: "\x00", wantErr: true},
		{s: "\x00\x00", wantErr: true},
		{s: " ", wantErr: false},
		{s: "abc", wantErr: false},
		{s: "()", wantErr: false},
		{s: "\x01", wantErr: false},
	}
	for i, test := range tests {
		err := ValidateLexerString(test.s)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateLexerString(%v)\nGOT %v\nWANT %v", i, test.s, err, test.wantErr)
		}
	}
}
