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
	dummyError := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		s   string
		e   string
		k   bool
		exp string
		err error
	}{
		{s: " ", e: "", k: false, exp: " ", err: dummyError},
		{s: " ", e: "", k: true, exp: " ", err: dummyError},
		{s: " ", e: " ", k: false, exp: "", err: dummyError},
		{s: " ", e: " ", k: true, exp: " ", err: dummyError},
		{s: "abc", e: "c", k: false, exp: "abc", err: dummyError},
		{s: "abc", e: "c", k: true, exp: "abc", err: dummyError},
		{s: "()()", e: ")", k: false, exp: "(", err: dummyError},
		{s: "()()", e: ")", k: true, exp: "()", err: dummyError},
	}

	for i, test := range table {
		stream := lexer.ParseString(test.s)

		res, _ := captureString(stream, test.e, test.k)
		if res != test.exp {
			t.Errorf("test %v: captureString(%v, %v, %v)\nGOT %v\nEXP %v", i, test.s, test.e, test.k, res, test.exp)
		}
	}
}

func TestParseRule(t *testing.T) {
	lexer := NewJSGFLexer()
	table := []struct {
		l string
		n string
		r Rule
	}{
		{
			l: "", n: "", r: Rule{
				Exp: Expression(""), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: ";", n: "", r: Rule{
				Exp: Expression(""), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: " ", n: "", r: Rule{
				Exp: Expression(""), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "<rule> =", n: "", r: Rule{
				Exp: Expression(""), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "<rule> = ", n: "", r: Rule{
				Exp: Expression(""), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "public <rule> =", n: "", r: Rule{
				Exp: Expression(""), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "public <rule> = ", n: "", r: Rule{
				Exp: Expression(""), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "<rule> =;", n: "<rule>", r: Rule{
				Exp: Expression(";"), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "public <rule> =;", n: "<rule>", r: Rule{
				Exp: Expression(";"), IsPublic: true, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "<rule> = test expression 123;", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123;"), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "<rule> = test \"expression\" 123;", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123;"), IsPublic: false, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "public <rule> = test expression 123;", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123;"), IsPublic: true, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "public <rule> = test \"expression\" 123;", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123;"), IsPublic: true, References: []string{}, Tokens: []Expression{},
				productions: []Expression{},
			},
		},
		{
			l: "<rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123 <rule> (abc) [def];"), IsPublic: false,
				References: []string{"<rule>"}, Tokens: []Expression{}, productions: []Expression{},
			},
		},
		{
			l: "public <rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123 <rule> (abc) [def];"), IsPublic: true, References: []string{"<rule>"},
				Tokens: []Expression{}, productions: []Expression{},
			},
		},
		{
			l: "<rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123 <rule1> <rule2> (abc) [def];"), IsPublic: false,
				References: []string{"<rule1>", "<rule2>"}, Tokens: []Expression{}, productions: []Expression{},
			},
		},
		{
			l: "public <rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{
				Exp: Expression("test expression 123 <rule1> <rule2> (abc) [def];"), IsPublic: true,
				References: []string{"<rule1>", "<rule2>"}, Tokens: []Expression{}, productions: []Expression{},
			},
		},
	}

	for i, test := range table {
		n, r, _ := ParseRule(lexer, test.l)
		if n != test.n {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, n, test.n)
		}

		if r.Exp != test.r.Exp {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, r, test.r)
		}

		if r.IsPublic != test.r.IsPublic {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).Is_public\nGOT %v\nEXP %v", i, test.l, r, test.r)
		}

		if !slices.Equal(r.References, test.r.References) {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).References\nGOT %v\nEXP %v", i, test.l, r.References, test.r.References)
		}

		sort.Slice(r.Tokens, func(i, j int) bool { return r.Tokens[i].str() < r.Tokens[j].str() })
		sort.Slice(test.r.Tokens, func(i, j int) bool { return r.Tokens[i].str() < r.Tokens[j].str() })

		if !slices.EqualFunc(r.Tokens, test.r.Tokens, func(E1, E2 Expression) bool { return E1.str() == E2.str() }) {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).Tokens\nGOT %v\nEXP %v", i, test.l, r.Tokens, test.r.Tokens)
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
