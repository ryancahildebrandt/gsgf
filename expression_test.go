// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:12:27 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"
	"testing"
)

func TestExpressionToTokens(t *testing.T) {
	lexer := NewJSGFLexer()
	table := []struct {
		e   Expression
		exp []Expression
	}{
		{Expression(""), []Expression{}},
		{Expression(";"), []Expression{"<SOS>", ";", "<EOS>"}},
		{Expression(" "), []Expression{"<SOS>", " ", "<EOS>"}},
		{Expression("test expression 123"), []Expression{"<SOS>", "test expression 123", "<EOS>"}},
		{Expression("test expression 123;"), []Expression{"<SOS>", "test expression 123", ";", "<EOS>"}},
		{Expression("test expression 123 (abc);"), []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", ";", "<EOS>"}},
		{Expression("test expression 123 [abc];"), []Expression{"<SOS>", "test expression 123 ", "[", "abc", "]", ";", "<EOS>"}},
		{Expression("test expression 123 [(abc)];"), []Expression{"<SOS>", "test expression 123 ", "[", "(", "abc", ")", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (abc) [def];"), []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (ab|c) | [de|f];"), []Expression{"<SOS>", "test expression 123 ", "(", "ab", "|", "c", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"}},
		{Expression("test expression 123 <rule> (abc) [def];"), []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 <rule> (abc) [def];"), []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 <rule> (abc) [def];"), []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123// [(abc)];"), []Expression{"<SOS>", "test expression 123// ", "[", "(", "abc", ")", "]", ";", "<EOS>"}},
		{Expression("test expression 123 // [(abc)];"), []Expression{"<SOS>", "test expression 123 // ", "[", "(", "abc", ")", "]", ";", "<EOS>"}},
		{Expression("test expression 123 /0.0/ [(abc)];"), []Expression{"<SOS>", "test expression 123 /0.0/", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"}},
		{Expression("test expression 123{} [(abc)];"), []Expression{"<SOS>", "test expression 123{}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"}},
		{Expression("test expression 123 {} [(abc)];"), []Expression{"<SOS>", "test expression 123 {}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"}},
		{Expression("test expression 123 {tag} [(abc)];"), []Expression{"<SOS>", "test expression 123 {tag}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (abc /1.0/) [def];"), []Expression{"<SOS>", "test expression 123 ", "(", "abc /1.0/", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (abc /1000/) [def];"), []Expression{"<SOS>", "test expression 123 ", "(", "abc /1000/", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (abc /-0.0001/) [def];"), []Expression{"<SOS>", "test expression 123 ", "(", "abc /-0.0001/", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (abc { }) [def];"), []Expression{"<SOS>", "test expression 123 ", "(", "abc { }", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (abc {_}) [def];"), []Expression{"<SOS>", "test expression 123 ", "(", "abc {_}", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (abc {_t_a_g_}) [def];"), []Expression{"<SOS>", "test expression 123 ", "(", "abc {_t_a_g_}", ")", " ", "[", "def", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (ab/1.0/|c/1.0/) | [de|f];"), []Expression{"<SOS>", "test expression 123 ", "(", "ab/1.0/", "|", "c/1.0/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (ab/1000/|c/1000/) | [de|f];"), []Expression{"<SOS>", "test expression 123 ", "(", "ab/1000/", "|", "c/1000/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (ab/-0.0001/|c/-0.0001/) | [de|f];"), []Expression{"<SOS>", "test expression 123 ", "(", "ab/-0.0001/", "|", "c/-0.0001/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (ab{1}|c{1}) | [de|f];"), []Expression{"<SOS>", "test expression 123 ", "(", "ab{1}", "|", "c{1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (ab{1.1}|c{1.1}) | [de|f];"), []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1}", "|", "c{1.1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"}},
		{Expression("test expression 123 (ab{1.1/1}|c{1.1/1}) | [de|f];"), []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1/1}", "|", "c{1.1/1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"}},
	}
	for i, test := range table {
		res := test.e.ToTokens(lexer)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.ToTokens(jsgflexer)\nGOT %v\nEXP %v", i, test.e, res, test.exp)
		}
	}
}

func TestExpressionParseWeight(t *testing.T) {
	dummy_error := errors.New("")
	table := []struct {
		e   Expression
		ee  Expression
		ef  float64
		err error
	}{
		{Expression(""), Expression(""), 0.0, dummy_error},
		{Expression("/"), Expression("/"), 0.0, dummy_error},
		{Expression("//"), Expression("//"), 0.0, dummy_error},
		{Expression("/.//"), Expression("/.//"), 0.0, dummy_error},
		{Expression("abc"), Expression("abc"), 0.0, dummy_error},
		{Expression("abc//"), Expression("abc//"), 0.0, dummy_error},
		{Expression("/aaa/"), Expression("/aaa/"), 0.0, dummy_error},

		{Expression("/1.0/"), Expression(""), 1.0, nil},
		{Expression("abc/1.0/"), Expression("abc"), 1.0, nil},
		{Expression("abc/0.1/"), Expression("abc"), 0.1, nil},
		{Expression("abc/0.001/"), Expression("abc"), 0.001, nil},
		{Expression("abc/100000/"), Expression("abc"), 100000.0, nil},
		{Expression("abc/100000.0/"), Expression("abc"), 100000.0, nil},
		{Expression("abc/1/"), Expression("abc"), 1, nil},
		{Expression("abc/-1/"), Expression("abc"), -1.0, nil},
		{Expression("abc/-1.0/"), Expression("abc"), -1.0, nil},
		{Expression("abc/0/"), Expression("abc"), 0.0, nil},
	}
	for i, test := range table {
		e, f, err := test.e.ParseWeight()
		if test.ee != e {
			t.Errorf("test %v: ParseWeight(%v)\nGOT %v\nEXP %v", i, test.e, e, test.ee)
		}
		if test.ef != f {
			t.Errorf("test %v: ParseWeight(%v)\nGOT %v\nEXP %v", i, test.e, f, test.ef)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ParseWeight(%v)\nGOT %v\nEXP %v", i, test.e, err, test.err)
		}
	}
}
