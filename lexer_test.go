// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:58 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestCaptureString(t *testing.T) {
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		s   string
		c   string
		inc bool
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
	for _, test := range table {
		stream := lexer.ParseString(test.s)
		res, _ := captureString(stream, test.c, test.inc)
		if res != test.exp {
			t.Errorf("captureString(%v, %v, %v)\nGOT %v\nEXP %v", test.s, test.c, test.inc, res, test.exp)
		}
	}
}

func TestParseRule(t *testing.T) {
	lexer := NewJSGFLexer()
	table := []struct {
		line string
		name string
		rule Rule
	}{
		{"", "", Rule{Expression(""), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{";", "", Rule{Expression(""), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{" ", "", Rule{Expression(""), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> =", "", Rule{Expression(""), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = ", "", Rule{Expression(""), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> =", "", Rule{Expression(""), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = ", "", Rule{Expression(""), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> =;", "<rule>", Rule{Expression(";"), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> =;", "<rule>", Rule{Expression(";"), true, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test expression 123;", "<rule>", Rule{Expression("test expression 123;"), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test \"expression\" 123;", "<rule>", Rule{Expression("test expression 123;"), false, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test expression 123;", "<rule>", Rule{Expression("test expression 123;"), true, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test \"expression\" 123;", "<rule>", Rule{Expression("test expression 123;"), true, []string{""}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test expression 123 <rule> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule> (abc) [def];"), false, []string{"<rule>"}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test expression 123 <rule> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule> (abc) [def];"), true, []string{"<rule>"}, Graph{}, []Expression{}, []Expression{}}},
		{"<rule> = test expression 123 <rule1> <rule2> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule1> <rule2> (abc) [def];"), false, []string{"<rule1>", "<rule2>"}, Graph{}, []Expression{}, []Expression{}}},
		{"public <rule> = test expression 123 <rule1> <rule2> (abc) [def];", "<rule>", Rule{Expression("test expression 123 <rule1> <rule2> (abc) [def];"), true, []string{"<rule1>", "<rule2>"}, Graph{}, []Expression{}, []Expression{}}},
	}
	for _, test := range table {
		n, r, _ := ParseRule(lexer, test.line)
		if n != test.name {
			t.Errorf("ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", test.line, n, test.name)
		}
		if r.Exp != test.rule.Exp {
			t.Errorf("ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", test.line, r, test.rule)
		}
		if r.Is_public != test.rule.Is_public {
			t.Errorf("ParseRule(jsgflexer, %v).Is_public\nGOT %v\nEXP %v", test.line, r, test.rule)
		}
		if fmt.Sprint(r.References) != fmt.Sprint(test.rule.References) {
			t.Errorf("ParseRule(jsgflexer, %v).References\nGOT %v\nEXP %v", test.line, r.References, test.rule.References)
		}
		if reflect.DeepEqual(r.Tokens, test.rule.Tokens) {
			t.Errorf("ParseRule(jsgflexer, %v).Tokens\nGOT %v\nEXP %v", test.line, r.Tokens, test.rule.Tokens)
		}
	}
}

func TestValidateJSGF(t *testing.T) {
	table := []struct {
		l   string
		exp bool
	}{
		{"", false},
		{";", false},
		{"=;", false},
		{"<>=;", false},
		{"public<>=;", false},
		{"public <>=;", false},
		{"< > = <>; ", false},
		{"< > = <>;", true},
		{"< >=;", true},
		{"public < >=;", true},
		{"public < > = ;", true},
		{"<abc> = def <ghi>;", true},
		{"<abc> = def = <ghi>;", true},
		{"<abc> = \"def\" = <ghi>;", true},
		{"<abc> = def <ghi>;;", true},
	}
	for _, test := range table {
		res := ValidateJSGF(test.l)
		if res != test.exp {
			t.Errorf("ValidateJSGF(%v)\nGOT %v\nEXP %v", test.l, res, test.exp)
		}
	}
}
