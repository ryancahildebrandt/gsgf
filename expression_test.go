// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:12:27 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"reflect"
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
	for _, test := range table {
		res := test.e.ToTokens(lexer)
		if !reflect.DeepEqual(res, test.exp) {
			t.Errorf("%v.ToTokens(jsgflexer)\nGOT %v\nEXP %v", test.e, res, test.exp)
		}
	}
}

func TestExpressionParseWeight(t *testing.T) {
	dummy_error := errors.New("")
	table := []struct {
		e     Expression
		exp_e Expression
		exp_f float64
		err   error
	}{
		{Expression(""), Expression(""), 0.0, dummy_error},
		{Expression("//"), Expression("//"), 0.0, dummy_error},
		{Expression("abc"), Expression("abc"), 0.0, dummy_error},
		{Expression("abc//"), Expression("abc//"), 0.0, dummy_error},

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
	for _, test := range table {
		res_e, res_f, err := test.e.ParseWeight()
		if test.exp_e != res_e {
			t.Errorf("ParseWeight(%v)\nGOT %v\nEXP %v", test.e, res_e, test.exp_e)
		}
		if test.exp_f != res_f {
			t.Errorf("ParseWeight(%v)\nGOT %v\nEXP %v", test.e, res_f, test.exp_f)
		}
		if test.err != nil && err == nil {
			t.Errorf("ParseWeight(%v)\nGOT %v\nEXP %v", test.e, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("ParseWeight(%v)\nGOT %v\nEXP %v", test.e, err, test.err)
		}
	}
}
