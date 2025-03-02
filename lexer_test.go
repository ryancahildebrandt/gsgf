// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:58 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"testing"
)

func TestCaptureString(t *testing.T) {
	// dummyError := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		s   string
		e   string
		k   bool
		exp string
		err error
	}{
		{s: " ", e: "", k: false, exp: " ", err: nil},
		{s: " ", e: "", k: true, exp: " ", err: nil},
		{s: " ", e: " ", k: false, exp: "", err: nil},
		{s: " ", e: " ", k: true, exp: " ", err: nil},
		{s: "abc", e: "c", k: false, exp: "abc", err: nil},
		{s: "abc", e: "c", k: true, exp: "abc", err: nil},
		{s: "()()", e: ")", k: false, exp: "(", err: nil},
		{s: "()()", e: ")", k: true, exp: "()", err: nil},
	}
	for i, test := range table {
		stream := lexer.ParseString(test.s)
		res, err := CaptureString(stream, test.e, test.k)
		if res != test.exp {
			t.Errorf("test %v: captureString(%v, %v, %v)\nGOT %v\nEXP %v", i, test.s, test.e, test.k, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: captureString(%v, %v, %v)\nGOT %v\nEXP %v", i, test.s, test.e, test.k, err, test.err)
		}
	}
}
