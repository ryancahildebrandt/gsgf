// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:54 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/bzick/tokenizer"
)

func TestGetCompositionOrder(t *testing.T) {
	table := []struct {
		g   Grammar
		exp []string
	}{
		{
			g:   Grammar{Rules: map[string]Rule{}, Imports: []string{}},
			exp: []string{},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("", false)}, Imports: []string{}},
			exp: []string{},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("<b><c>", false)}, Imports: []string{}},
			exp: []string{},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("", true)}, Imports: []string{}},
			exp: []string{"<a>"},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("<b><c>", true)}, Imports: []string{}},
			exp: []string{"<a>", "<b>", "<c>"},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", true)}, Imports: []string{}},
			exp: []string{"<a>", "<b>"},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", false)}, Imports: []string{}},
			exp: []string{"<a>"},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", true)}, Imports: []string{}},
			exp: []string{"<a>", "<b>", "<c>"},
		},
		{
			g:   Grammar{Rules: map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", false)}, Imports: []string{}},
			exp: []string{"<a>", "<b>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true),
					"<d>": NewRule("", true),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<b>", "<c>", "<c>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<c>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", false), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{},
		},
	}
	for i, test := range table {
		res := GetCompositionOrder(test.g)
		sort.Strings(res)
		sort.Strings(test.exp)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.CompositionOrder()\nGOT %v\nEXP %v", i, test.g, res, test.exp)
		}
	}
}

func TestGetAllProductions(t *testing.T) {
	dummyError := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		p   []string
		exp []string
		err error
	}{
		{p: []string{""}, exp: []string{}, err: nil},
		{p: []string{";"}, exp: []string{}, err: nil},
		{p: []string{"", ""}, exp: []string{}, err: nil},
		{p: []string{";", ";", ";"}, exp: []string{}, err: nil},
		{p: []string{"abc;"}, exp: []string{"abc"}, err: nil},
		{p: []string{"<b>;"}, exp: []string{"1", "2", "3"}, err: nil},
		{p: []string{"<g>;"}, exp: []string{"abc"}, err: nil},
		{p: []string{"<h>;"}, exp: []string{"123bc"}, err: nil},
		{p: []string{"<i>;"}, exp: []string{"a1c", "a2c", "a3c"}, err: nil},
		{
			p: []string{"<j>;"}, exp: []string{
				"a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123",
				"a313123123", "a1123123123", "a2123123123", "a3123123123",
			}, err: nil,
		},
		{
			p: []string{"<k>;"}, exp: []string{
				"a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221",
				"a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332",
				"a333",
			}, err: nil,
		},
		{p: []string{"a{}b//c/0.1/;"}, exp: []string{"a{}b//c/0.1/"}, err: nil},
		{p: []string{"a{}/0.1/b/0.1/{}c;"}, exp: []string{"a{}/0.1/b/0.1/{}c"}, err: nil},
		{p: []string{"abc<a>;"}, exp: []string{"abc123"}, err: nil},
		{p: []string{"abc<l>;"}, exp: []string{"abc<l>"}, err: dummyError},
		{
			p:   []string{"abc<a><b><c>;"},
			exp: []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313"}, err: nil,
		},
		{p: []string{"abc(<g>|<h>|<i>);"}, exp: []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c"}, err: nil},
		{p: []string{"abc<g><g><g>;"}, exp: []string{"abcabcabcabc"}, err: nil},
		{p: []string{"abc;", "def;", "ghi;"}, exp: []string{"abc", "def", "ghi"}, err: nil},
		{p: []string{"abc<a>;", "def<h>;", "<b>;"}, exp: []string{"abc123", "def123bc", "1", "2", "3"}, err: nil},
		{
			p: []string{"abc<a><b><c>;", "def[<h>];", "g((hi)|(jk));"}, exp: []string{
				"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313", "def", "def123bc", "ghi",
				"gjk",
			}, err: nil,
		},
		{
			p: []string{"abc(<g>|<h>|<i>);", "def[<g>|<h>|<i>];", "ghi<b><b><b>;"}, exp: []string{
				"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c", "def", "defabc", "def123bc", "defa1c", "defa2c",
				"defa3c", "ghi111", "ghi112", "ghi113", "ghi121", "ghi122", "ghi123", "ghi131", "ghi132", "ghi133",
				"ghi211", "ghi212", "ghi213", "ghi221", "ghi222", "ghi223", "ghi231", "ghi232", "ghi233", "ghi311",
				"ghi312", "ghi313", "ghi321", "ghi322", "ghi323", "ghi331", "ghi332", "ghi333",
			}, err: nil,
		},
		{
			p: []string{"abc<g><g><g>;", "<h><i>;", "<j>|<k>;"}, exp: []string{
				"abcabcabcabc", "123bca1c", "123bca2c", "123bca3c", "a11312312", "a21312312", "a31312312", "a112312312",
				"a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123",
				"a3123123123", "a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212",
				"a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323",
				"a331", "a332", "a333",
			}, err: nil,
		},
	}
	for i, test := range table {
		g := NewGrammar()
		g.Rules = map[string]Rule{
			"<_>": {
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			"<a>": {
				Exp: "123;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}),
			},
			"<b>": {
				Exp: "1|2|3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}),
			},
			"<c>": {
				Exp: "1[2]3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}),
			},
			"<d>": {
				Exp: "1(2)3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}),
			},
			"<e>": {
				Exp: "1(2[3]);", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 4, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
					{From: 7, To: 8, Weight: 1.0},
					{From: 8, To: 9, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}),
			},
			"<f>": {
				Exp: "<l>;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}),
			},
			"<g>": {
				Exp: "abc;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}),
			},
			"<h>": {
				Exp: "<a>bc;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
				}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}),
			},
			"<i>": {
				Exp: "a<b>c;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}),
			},
			"<j>": {
				Exp: "a<b><c><d><e>;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}),
			},
			"<k>": {
				Exp: "a<b><b><b>;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}),
			},
		}
		for j, p := range test.p {
			rule := NewRule(p, true)
			rule.Tokens = ToTokens(rule.Exp, lexer)
			rule.Graph = NewGraph(ToEdgeList(rule.Tokens), rule.Tokens)
			// rule.productions = FilterTerminals(rule.Tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			g.Rules[fmt.Sprintf("<pub_%v>", j)] = rule
		}
		g, err := ResolveRules(g, lexer)
		res := GetAllProductions(g)
		sort.Strings(test.exp)
		sort.Strings(res)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.Productions()\nGOT len %v %v\nEXP len %v %v", i, test.p, len(res), res, len(test.exp), test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.Productions().err\nGOT %v\nEXP %v", i, test.p, err, test.err)
		}
	}
}

func TestGetAllProductionsMinimized(t *testing.T) {
	dummyError := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		p   []string
		exp []string
		err error
	}{
		{p: []string{""}, exp: []string{}, err: nil},
		{p: []string{";"}, exp: []string{}, err: nil},
		{p: []string{"", ""}, exp: []string{}, err: nil},
		{p: []string{";", ";", ";"}, exp: []string{}, err: nil},
		{p: []string{"abc;"}, exp: []string{"abc"}, err: nil},
		{p: []string{"<b>;"}, exp: []string{"1", "2", "3"}, err: nil},
		{p: []string{"<g>;"}, exp: []string{"abc"}, err: nil},
		{p: []string{"<h>;"}, exp: []string{"123bc"}, err: nil},
		{p: []string{"<i>;"}, exp: []string{"a1c", "a2c", "a3c"}, err: nil},
		{
			p: []string{"<j>;"}, exp: []string{
				"a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123",
				"a313123123", "a1123123123", "a2123123123", "a3123123123",
			}, err: nil,
		},
		{
			p: []string{"<k>;"}, exp: []string{
				"a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221",
				"a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332",
				"a333",
			}, err: nil,
		},
		{p: []string{"a{}b//c/0.1/;"}, exp: []string{"a{}b//c/0.1/"}, err: nil},
		{p: []string{"a{}/0.1/b/0.1/{}c;"}, exp: []string{"a{}/0.1/b/0.1/{}c"}, err: nil},
		{p: []string{"abc<a>;"}, exp: []string{"abc123"}, err: nil},
		{p: []string{"abc<l>;"}, exp: []string{"abc<l>"}, err: dummyError},
		{
			p:   []string{"abc<a><b><c>;"},
			exp: []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313"}, err: nil,
		},
		{p: []string{"abc(<g>|<h>|<i>);"}, exp: []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c"}, err: nil},
		{p: []string{"abc<g><g><g>;"}, exp: []string{"abcabcabcabc"}, err: nil},
		{p: []string{"abc;", "def;", "ghi;"}, exp: []string{"abc", "def", "ghi"}, err: nil},
		{p: []string{"abc<a>;", "def<h>;", "<b>;"}, exp: []string{"abc123", "def123bc", "1", "2", "3"}, err: nil},
		{
			p: []string{"abc<a><b><c>;", "def[<h>];", "g((hi)|(jk));"}, exp: []string{
				"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313", "def", "def123bc", "ghi",
				"gjk",
			}, err: nil,
		},
		{
			p: []string{"abc(<g>|<h>|<i>);", "def[<g>|<h>|<i>];", "ghi<b><b><b>;"}, exp: []string{
				"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c", "def", "defabc", "def123bc", "defa1c", "defa2c",
				"defa3c", "ghi111", "ghi112", "ghi113", "ghi121", "ghi122", "ghi123", "ghi131", "ghi132", "ghi133",
				"ghi211", "ghi212", "ghi213", "ghi221", "ghi222", "ghi223", "ghi231", "ghi232", "ghi233", "ghi311",
				"ghi312", "ghi313", "ghi321", "ghi322", "ghi323", "ghi331", "ghi332", "ghi333",
			}, err: nil,
		},
		{
			p: []string{"abc<g><g><g>;", "<h><i>;", "<j>|<k>;"}, exp: []string{
				"abcabcabcabc", "123bca1c", "123bca2c", "123bca3c", "a11312312", "a21312312", "a31312312", "a112312312",
				"a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123",
				"a3123123123", "a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212",
				"a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323",
				"a331", "a332", "a333",
			}, err: nil,
		},
	}
	for i, test := range table {
		g := NewGrammar()
		g.Rules = map[string]Rule{
			"<_>": {
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			"<a>": {
				Exp: "123;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}),
			},
			"<b>": {
				Exp: "1|2|3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}),
			},
			"<c>": {
				Exp: "1[2]3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}),
			},
			"<d>": {
				Exp: "1(2)3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}),
			},
			"<e>": {
				Exp: "1(2[3]);", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 4, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
					{From: 7, To: 8, Weight: 1.0},
					{From: 8, To: 9, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}),
			},
			"<f>": {
				Exp: "<l>;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}),
			},
			"<g>": {
				Exp: "abc;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}),
			},
			"<h>": {
				Exp: "<a>bc;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
				}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}),
			},
			"<i>": {
				Exp: "a<b>c;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}),
			},
			"<j>": {
				Exp: "a<b><c><d><e>;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}),
			},
			"<k>": {
				Exp: "a<b><b><b>;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}),
			},
		}
		for j, p := range test.p {
			rule := NewRule(p, true)
			rule.Tokens = ToTokens(rule.Exp, lexer)
			rule.Graph = NewGraph(ToEdgeList(rule.Tokens), rule.Tokens)
			g.Rules[fmt.Sprintf("<pub_%v>", j)] = rule
		}
		g, err := ResolveRules(g, lexer)
		res := GetAllProductions(g)
		sort.Strings(test.exp)
		sort.Strings(res)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.Productions()\nGOT len %v %v\nEXP len %v %v", i, test.p, len(res), res, len(test.exp), test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.Productions().err\nGOT %v\nEXP %v", i, test.p, err, test.err)
		}
	}
}

func TestGetAllProductionsE2E(t *testing.T) {
	dummyError := errors.New("")
	lexer := NewJSGFLexer()
	var productions []string
	f, _ := os.Open("data/tests/productions.txt")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		productions = append(productions, scanner.Text())
	}
	table := []struct {
		p   string
		exp []string
		err error
	}{
		{p: "data/tests/test0.jsgf", exp: productions, err: nil},
		{p: "data/tests/test1.jsgf", exp: productions, err: nil},
		{p: "data/tests/test2.jsgf", exp: productions, err: dummyError}, // returns error but is still resolvable, figure out how to validate completness before imports
		{p: "data/tests/test3.jsgf", exp: productions, err: nil},
		{p: "data/tests/test4.jsgf", exp: productions, err: nil},
		{p: "data/tests/test5.jsgf", exp: productions, err: nil},
		{p: "data/tests/a.jsgf", exp: []string{}, err: nil},
		{p: "data/tests/b.jsgf", exp: []string{}, err: nil},
		{p: "data/tests/dir0/c.jsgf", exp: []string{}, err: nil},
		{p: "data/tests/dir0/dir1/d.jsgf", exp: []string{}, err: dummyError},
		{p: "data/tests/dir0/dir1/dir2/e.jsgf", exp: productions, err: nil},
	}
	for i, test := range table {
		var err error
		grammar := NewGrammar()
		f, err1 := os.Open(test.p)
		scanner := bufio.NewScanner(f)
		grammar, err2 := ImportLines(grammar, scanner, lexer)
		namespace, err3 := CreateNameSpace(test.p, ".jsgf")
		grammar = ImportNameSpace(grammar, namespace, lexer)
		grammar, err4 := ResolveRules(grammar, lexer)
		res := GetAllProductions(grammar)
		for _, e := range []error{err1, err2, err3, err4} {
			if e != nil {
				err = e
			}
		}
		sort.Strings(test.exp)
		sort.Strings(res)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.Productions()\nGOT %v\nEXP %v", i, test.p, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.Productions().err\nGOT %v\nEXP %v", i, test.p, err, test.err)
		}
	}
}

func TestValidateGrammarCompleteness(t *testing.T) {
	type args struct {
		g Grammar
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateGrammarCompleteness(tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("ValidateGrammarCompleteness() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImportLines(t *testing.T) {
	type args struct {
		g   Grammar
		s   *bufio.Scanner
		lex *tokenizer.Tokenizer
	}
	tests := []struct {
		name    string
		args    args
		want    Grammar
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ImportLines(tt.args.g, tt.args.s, tt.args.lex)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImportNameSpace(t *testing.T) {
	type args struct {
		g   Grammar
		r   map[string]string
		lex *tokenizer.Tokenizer
	}
	tests := []struct {
		name string
		args args
		want Grammar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ImportNameSpace(tt.args.g, tt.args.r, tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportNameSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveRules(t *testing.T) {
	type args struct {
		g   Grammar
		lex *tokenizer.Tokenizer
	}
	tests := []struct {
		name    string
		args    args
		want    Grammar
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveRules(tt.args.g, tt.args.lex)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
