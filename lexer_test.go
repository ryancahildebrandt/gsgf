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
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		s   string
		c   string
		k   bool
		exp string
		err error
	}{
		{" ", "", false, " ", dummy_error},
		{" ", "", true, " ", dummy_error},
		{" ", " ", false, "", dummy_error},
		{" ", " ", true, " ", dummy_error},
		{"abc", "c", false, "abc", dummy_error},
		{"abc", "c", true, "abc", dummy_error},
		{"()()", ")", false, "(", dummy_error},
		{"()()", ")", true, "()", dummy_error},
	}
	for i, test := range table {
		stream := lexer.ParseString(test.s)
		res, _ := captureString(stream, test.c, test.k)
		if res != test.exp {
			t.Errorf("test %v: captureString(%v, %v, %v)\nGOT %v\nEXP %v", i, test.s, test.c, test.k, res, test.exp)
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
		{"", "", Rule{Expression(""), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{";", "", Rule{Expression(""), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{" ", "", Rule{Expression(""), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> =", "", Rule{Expression(""), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = ", "", Rule{Expression(""), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> =", "", Rule{Expression(""), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = ", "", Rule{Expression(""), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> =;", "<rule>", Rule{Expression(";"), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> =;", "<rule>", Rule{Expression(";"), true, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test expression 123;", "<rule>", Rule{Expression("test expression 123;"), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test \"expression\" 123;", "<rule>", Rule{Expression("test expression 123;"), false, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test expression 123;", "<rule>", Rule{Expression("test expression 123;"), true, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test \"expression\" 123;", "<rule>", Rule{Expression("test expression 123;"), true, []string{}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test expression 123 <rule> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule> (abc) [def];"), false, []string{"<rule>"}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test expression 123 <rule> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule> (abc) [def];"), true, []string{"<rule>"}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test expression 123 <rule1> <rule2> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule1> <rule2> (abc) [def];"), false, []string{"<rule1>", "<rule2>"}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test expression 123 <rule1> <rule2> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule1> <rule2> (abc) [def];"), true, []string{"<rule1>", "<rule2>"}, Graph{}, []Expression{}, []Expression{}}},
	}
	for i, test := range table {
		n, r, _ := ParseRule(lexer, test.l)
		if n != test.n {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, n, test.n)
		}
		if r.Exp != test.r.Exp {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, r, test.r)
		}
		if r.Is_public != test.r.Is_public {
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
	dummy_error := errors.New("")
	table := []struct {
		l   string
		err error
	}{
		{"", dummy_error},
		{";", dummy_error},
		{"=;", dummy_error},
		{"<>=;", dummy_error},
		{"public<>=;", dummy_error},
		{"public <>=;", dummy_error},
		{"< > = <>; ", dummy_error},
		{"< > = <>;", nil},
		{"< >=;", nil},
		{"public < >=;", nil},
		{"public < > = ;", nil},
		{"<abc> = def <ghi>;", nil},
		{"<abc> = def = <ghi>;", nil},
		{"<abc> = \"def\" = <ghi>;", nil},
		{"<abc> = def <ghi>;;", nil},
	}
	for i, test := range table {
		err := ValidateJSGFRule(test.l)
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ValidateJSGF(%v)\nGOT %v\nEXP %v", i, test.l, err, test.err)
		}
	}
}
