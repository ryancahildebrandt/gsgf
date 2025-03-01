// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:58 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"
	"sort"
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

func TestParseRule(t *testing.T) {
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		l   string
		n   string
		r   Rule
		err error
	}{
		{
			l: "", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: ";", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: " ", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "<rule> =", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "<rule> = ", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "public <rule> =", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "public <rule> = ", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "<rule> =;", n: "<rule>", r: Rule{
				Exp: ";", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> =;", n: "<rule>", r: Rule{
				Exp: ";", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "<rule> = test expression 123;", n: "<rule>", r: Rule{
				Exp: "test expression 123;", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "<rule> = test \"expression\" 123;", n: "<rule>", r: Rule{
				Exp: "test \"expression\" 123;", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> = test expression 123;", n: "<rule>", r: Rule{
				Exp: "test expression 123;", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "public <rule> = test \"expression\" 123;", n: "<rule>", r: Rule{
				Exp: "test \"expression\" 123;", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "<rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule> (abc) [def];", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule> (abc) [def];", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "<rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: true,
			},
			err: nil,
		},
	}
	for i, test := range table {
		n, r, err := ParseRule(test.l, lexer)
		if n != test.n {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, n, test.n)
		}
		if r.Exp != test.r.Exp {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, r, test.r)
		}
		if r.IsPublic != test.r.IsPublic {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).Is_public\nGOT %v\nEXP %v", i, test.l, r, test.r)
		}
		if !slices.Equal(GetReferences(r), GetReferences(test.r)) {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).References\nGOT %v\nEXP %v", i, test.l, GetReferences(r), GetReferences(test.r))
		}
		sort.Slice(r.Tokens, func(i, j int) bool { return r.Tokens[i] < r.Tokens[j] })
		sort.Slice(test.r.Tokens, func(i, j int) bool { return r.Tokens[i] < r.Tokens[j] })
		if !slices.EqualFunc(r.Tokens, test.r.Tokens, func(E1, E2 Expression) bool { return E1 == E2 }) {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).Tokens\nGOT %v\nEXP %v", i, test.l, r.Tokens, test.r.Tokens)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ValidateJSGF(%v)\nGOT %v\nEXP %v", i, test.l, err, test.err)
		}
	}
}

func TestValidateJSGF(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		l   string
		err error
	}{
		{l: "", err: dummyError},
		{l: ";", err: dummyError},
		{l: "=;", err: dummyError},
		{l: "<>=;", err: dummyError},
		{l: "public<>=;", err: dummyError},
		{l: "public <>=;", err: dummyError},
		{l: "< > = <>; ", err: dummyError},
		{l: "< > = <>;", err: nil},
		{l: "< >=;", err: nil},
		{l: "public < >=;", err: nil},
		{l: "public < > = ;", err: nil},
		{l: "<abc> = def <ghi>;", err: nil},
		{l: "<abc> = def = <ghi>;", err: nil},
		{l: "<abc> = \"def\" = <ghi>;", err: nil},
		{l: "<abc> = def <ghi>;;", err: nil},
	}
	for i, test := range table {
		err := ValidateJSGFRule(test.l)
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ValidateJSGF(%v)\nGOT %v\nEXP %v", i, test.l, err, test.err)
		}
	}
}
