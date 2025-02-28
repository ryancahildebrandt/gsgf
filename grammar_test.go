// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:54 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"testing"
)

func TestGrammarCompositionOrder(t *testing.T) {
	table := []struct {
		g   Grammar
		exp []string
	}{
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{}, Imports: []string{}},
			exp: []string{},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("", false)}, Imports: []string{}},
			exp: []string{},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("<b><c>", false)}, Imports: []string{}},
			exp: []string{},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("", true)}, Imports: []string{}},
			exp: []string{"<a>"},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("<b><c>", true)}, Imports: []string{}},
			exp: []string{"<a>", "<b>", "<c>"},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", true)}, Imports: []string{}},
			exp: []string{"<a>", "<b>"},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", false)}, Imports: []string{}},
			exp: []string{"<a>"},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", true)}, Imports: []string{}},
			exp: []string{"<a>", "<b>", "<c>"},
		},
		{
			g:   Grammar{Path: "", Rules: map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", false)}, Imports: []string{}},
			exp: []string{"<a>", "<b>"},
		},
		{
			g: Grammar{
				Path: "", Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true),
					"<d>": NewRule("", true),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Path: "", Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Path: "", Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<b>", "<c>", "<c>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Path: "", Rules: map[string]Rule{
					"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{"<a>", "<c>", "<d>"},
		},
		{
			g: Grammar{
				Path: "", Rules: map[string]Rule{
					"<a>": NewRule("<c>", false), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				}, Imports: []string{},
			}, exp: []string{},
		},
	}
	for i, test := range table {
		res := test.g.CompositionOrder()
		sort.Strings(res)
		sort.Strings(test.exp)

		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.CompositionOrder()\nGOT %v\nEXP %v", i, test.g, res, test.exp)
		}
	}
}

func TestGrammarProductions(t *testing.T) {
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
		g := NewGrammar("")
		g.Rules = map[string]Rule{
			"<_>": {
				Exp: "", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{}, []Expression{}),
				Tokens: []Expression{}, productions: []Expression{},
			},
			"<a>": {
				Exp: "123;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}), Tokens: []Expression{"<SOS>", "123", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<b>": {
				Exp: "1|2|3;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, productions: []Expression{},
			},
			"<c>": {
				Exp: "1[2]3;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, productions: []Expression{},
			},
			"<d>": {
				Exp: "1(2)3;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, productions: []Expression{},
			},
			"<e>": {
				Exp: "1(2[3]);", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
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
				Tokens:      []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<f>": {
				Exp: "<l>;", IsPublic: false, References: []string{"<f>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), Tokens: []Expression{"<SOS>", "<f>", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<g>": {
				Exp: "abc;", IsPublic: false, References: []string{""}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), Tokens: []Expression{"<SOS>", "abc", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<h>": {
				Exp: "<a>bc;", IsPublic: false, References: []string{"<a>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
				}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}, productions: []Expression{},
			},
			"<i>": {
				Exp: "a<b>c;", IsPublic: false, References: []string{"<b>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}, productions: []Expression{},
			},
			"<j>": {
				Exp: "a<b><c><d><e>;", IsPublic: false, References: []string{"<b>", "<c>", "<d>", "<e>"},
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}),
				Tokens:      []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<k>": {
				Exp: "a<b><b><b>;", IsPublic: false, References: []string{"<b>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}),
				Tokens:      []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"},
				productions: []Expression{},
			},
		}

		for j, p := range test.p {
			rule := NewRule(Expression(p), true)
			rule.Tokens = rule.Exp.ToTokens(lexer)
			rule.Graph = NewGraph(BuildEdgeList(rule.Tokens), rule.Tokens)
			rule.productions = FilterTerminals(rule.Tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			g.Rules[fmt.Sprintf("<pub_%v>", j)] = rule
		}

		g, err := g.Resolve(lexer)
		res := g.Productions()

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

func TestGrammarProductionsMinimized(t *testing.T) {
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
		g := NewGrammar("")
		g.Rules = map[string]Rule{
			"<_>": {
				Exp: "", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{}, []Expression{}),
				Tokens: []Expression{}, productions: []Expression{},
			},
			"<a>": {
				Exp: "123;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}), Tokens: []Expression{"<SOS>", "123", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<b>": {
				Exp: "1|2|3;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, productions: []Expression{},
			},
			"<c>": {
				Exp: "1[2]3;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, productions: []Expression{},
			},
			"<d>": {
				Exp: "1(2)3;", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, productions: []Expression{},
			},
			"<e>": {
				Exp: "1(2[3]);", IsPublic: false, References: []string{}, Graph: NewGraph(EdgeList{
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
				Tokens:      []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<f>": {
				Exp: "<l>;", IsPublic: false, References: []string{"<f>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), Tokens: []Expression{"<SOS>", "<f>", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<g>": {
				Exp: "abc;", IsPublic: false, References: []string{""}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), Tokens: []Expression{"<SOS>", "abc", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<h>": {
				Exp: "<a>bc;", IsPublic: false, References: []string{"<a>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
				}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}, productions: []Expression{},
			},
			"<i>": {
				Exp: "a<b>c;", IsPublic: false, References: []string{"<b>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}),
				Tokens: []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}, productions: []Expression{},
			},
			"<j>": {
				Exp: "a<b><c><d><e>;", IsPublic: false, References: []string{"<b>", "<c>", "<d>", "<e>"},
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}),
				Tokens:      []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"},
				productions: []Expression{},
			},
			"<k>": {
				Exp: "a<b><b><b>;", IsPublic: false, References: []string{"<b>"}, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}),
				Tokens:      []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"},
				productions: []Expression{},
			},
		}

		for j, p := range test.p {
			rule := NewRule(Expression(p), true)
			rule.Tokens = rule.Exp.ToTokens(lexer)
			rule.Graph = NewGraph(BuildEdgeList(rule.Tokens), rule.Tokens)
			rule.productions = FilterTerminals(rule.Tokens, []string{
				"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>", "",
			})
			g.Rules[fmt.Sprintf("<pub_%v>", j)] = rule
		}

		g, err := g.Resolve(lexer)
		res := g.Productions()

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

func TestGrammarPeek(t *testing.T) {
	table := []struct {
		p       string
		n       string
		imports []string
		rules   map[string][]string
	}{
		{
			p: "data/tests/test0.jsgf", n: "test0", imports: []string{"import <a.*>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "quant": {}, "teatype": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test1.jsgf", n: "test1", imports: []string{"import <c.brew>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {},
			},
		},
		{
			p: "data/tests/test2.jsgf", n: "test2", imports: []string{"import <a1.*>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test3.jsgf", n: "test3", imports: []string{"import <e.dne>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test4.jsgf", n: "test4", imports: []string{"import <d.*>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "quant": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test5.jsgf", n: "test5", imports: []string{"import <b.request>;"},
			rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "order": {"quant"}, "quant": {}, "teatype": {},
				"brew": {"quant"},
			},
		},
		{
			p: "data/tests/a.jsgf", n: "a", imports: []string{},
			rules: map[string][]string{"request": {"brew"}, "order": {"quant"}, "brew": {"quant"}, "quant": {}},
		},
		{
			p: "data/tests/b.jsgf", n: "b", imports: []string{"import <c.brew>;"},
			rules: map[string][]string{"request": {"brew"}, "order": {"quant"}, "quant": {}},
		},
		{
			p: "data/tests/dir0/c.jsgf", n: "c", imports: []string{},
			rules: map[string][]string{"teatype": {}, "brew": {"quant"}, "quant": {}},
		},
		{
			p: "data/tests/dir0/dir1/d.jsgf", n: "d", imports: []string{"import <c.teatype>;", "import <a.order>;"},
			rules: map[string][]string{},
		},
		{
			p: "data/tests/dir0/dir1/dir2/e.jsgf", n: "e", imports: []string{}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {}, "brew": {"quant"},
			},
		},
	}
	for i, test := range table {
		name, imports, rules, err := NewGrammar(test.p).Peek()
		if err != nil {
			t.Errorf("test %v: Grammar(%v).Peek()\nGOT error %v", i, test.p, err)
		}

		if name != test.n {
			t.Errorf("test %v: Grammar(%v).Peek().imports\nGOT %v\nEXP %v", i, test.p, name, test.n)
		}

		sort.Strings(imports)
		sort.Strings(test.imports)

		if !slices.Equal(imports, test.imports) {
			t.Errorf("test %v: Grammar(%v).Peek().imports\nGOT %v\nEXP %v", i, test.p, imports, test.imports)
		}

		if len(rules) != len(test.rules) {
			t.Errorf("test %v: Grammar(%v).Peek().rules\nGOT %v\nEXP %v", i, test.p, rules, test.rules)
		}

		for k1, v1 := range rules {
			v2, ok := test.rules[k1]
			if !ok {
				t.Errorf("test %v: Grammar(%v).Peek().rules\nGOT %v\nEXP %v", i, test.p, rules, test.rules)
			}

			sort.Strings(v1)
			sort.Strings(v2)

			if !slices.Equal(v1, v2) {
				t.Errorf("test %v: Grammar(%v).Peek().rules\nGOT %v\nEXP %v", i, test.p, rules, test.rules)
			}
		}
	}
}

func TestGrammarProductionsE2E(t *testing.T) {
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
		{
			p: "data/tests/test2.jsgf", exp: productions, err: dummyError,
		}, // returns error but is still resolvable, figure out how to validate completness before imports
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

		grammar := NewGrammar(test.p)
		f, err1 := os.Open(test.p)
		scanner := bufio.NewScanner(f)
		grammar, err2 := grammar.ReadLines(scanner, lexer)
		namespace, err3 := CreateNameSpace(grammar.Path, ".jsgf")
		grammar = grammar.ReadNameSpace(namespace, lexer)
		grammar, err4 := grammar.Resolve(lexer)
		res := grammar.Productions()

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
