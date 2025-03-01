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
		{e: Expression(""), exp: []Expression{}},
		{e: Expression(";"), exp: []Expression{"<SOS>", ";", "<EOS>"}},
		{e: Expression(" "), exp: []Expression{"<SOS>", " ", "<EOS>"}},
		{e: Expression("test expression 123"), exp: []Expression{"<SOS>", "test expression 123", "<EOS>"}},
		{e: Expression("test expression 123;"), exp: []Expression{"<SOS>", "test expression 123", ";", "<EOS>"}},
		{
			e:   Expression("test expression 123 (abc);"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 [abc];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "[", "abc", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 [(abc)];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (abc) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (ab|c) | [de|f];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "ab", "|", "c", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 <rule> (abc) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 <rule> (abc) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 <rule> (abc) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123// [(abc)];"),
			exp: []Expression{"<SOS>", "test expression 123// ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 // [(abc)];"),
			exp: []Expression{"<SOS>", "test expression 123 // ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 /0.0/ [(abc)];"),
			exp: []Expression{"<SOS>", "test expression 123 /0.0/", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123{} [(abc)];"),
			exp: []Expression{"<SOS>", "test expression 123{}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 {} [(abc)];"),
			exp: []Expression{"<SOS>", "test expression 123 {}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 {tag} [(abc)];"),
			exp: []Expression{"<SOS>", "test expression 123 {tag}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (abc /1.0/) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc /1.0/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (abc /1000/) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc /1000/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (abc /-0.0001/) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc /-0.0001/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (abc { }) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc { }", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (abc {_}) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc {_}", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (abc {_t_a_g_}) [def];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "abc {_t_a_g_}", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (ab/1.0/|c/1.0/) | [de|f];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "ab/1.0/", "|", "c/1.0/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (ab/1000/|c/1000/) | [de|f];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "ab/1000/", "|", "c/1000/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (ab/-0.0001/|c/-0.0001/) | [de|f];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "ab/-0.0001/", "|", "c/-0.0001/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (ab{1}|c{1}) | [de|f];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1}", "|", "c{1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (ab{1.1}|c{1.1}) | [de|f];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1}", "|", "c{1.1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:   Expression("test expression 123 (ab{1.1/1}|c{1.1/1}) | [de|f];"),
			exp: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1/1}", "|", "c{1.1/1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
	}
	for i, test := range table {
		res := ToTokens(test.e, lexer)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.ToTokens(jsgflexer)\nGOT %v\nEXP %v", i, test.e, res, test.exp)
		}
	}
}

func TestExpressionParseWeight(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		e   Expression
		ee  Expression
		ef  float64
		err error
	}{
		{e: Expression(""), ee: Expression(""), ef: 0.0, err: dummyError},
		{e: Expression("/"), ee: Expression("/"), ef: 0.0, err: dummyError},
		{e: Expression("//"), ee: Expression("//"), ef: 0.0, err: dummyError},
		{e: Expression("/.//"), ee: Expression("/.//"), ef: 0.0, err: dummyError},
		{e: Expression("abc"), ee: Expression("abc"), ef: 0.0, err: dummyError},
		{e: Expression("abc//"), ee: Expression("abc//"), ef: 0.0, err: dummyError},
		{e: Expression("/aaa/"), ee: Expression("/aaa/"), ef: 0.0, err: dummyError},
		{e: Expression("/1.0/"), ee: Expression(""), ef: 1.0, err: nil},
		{e: Expression("abc/1.0/"), ee: Expression("abc"), ef: 1.0, err: nil},
		{e: Expression("abc/0.1/"), ee: Expression("abc"), ef: 0.1, err: nil},
		{e: Expression("abc/0.001/"), ee: Expression("abc"), ef: 0.001, err: nil},
		{e: Expression("abc/100000/"), ee: Expression("abc"), ef: 100000.0, err: nil},
		{e: Expression("abc/100000.0/"), ee: Expression("abc"), ef: 100000.0, err: nil},
		{e: Expression("abc/1/"), ee: Expression("abc"), ef: 1, err: nil},
		{e: Expression("abc/-1/"), ee: Expression("abc"), ef: -1.0, err: nil},
		{e: Expression("abc/-1.0/"), ee: Expression("abc"), ef: -1.0, err: nil},
		{e: Expression("abc/0/"), ee: Expression("abc"), ef: 0.0, err: nil},
	}
	for i, test := range table {
		e, f, err := ParseWeight(test.e)
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
