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
		{Grammar{"", map[string]Rule{}, []string{}}, []string{}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("", false)}, []string{}}, []string{}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<b><c>", false)}, []string{}}, []string{}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("", true)}, []string{}}, []string{"<a>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<b><c>", true)}, []string{}}, []string{"<a>", "<b>", "<c>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", true)}, []string{}}, []string{"<a>", "<b>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", false)}, []string{}}, []string{"<a>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", true)}, []string{}}, []string{"<a>", "<b>", "<c>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", false)}, []string{}}, []string{"<a>", "<b>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true), "<d>": NewRule("", true)}, []string{}}, []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>", "<d>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true), "<d>": NewRule("", false)}, []string{}}, []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", false), "<d>": NewRule("", false)}, []string{}}, []string{"<a>", "<b>", "<c>", "<c>", "<d>", "<d>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false), "<d>": NewRule("", false)}, []string{}}, []string{"<a>", "<c>", "<d>"}},
		{Grammar{"", map[string]Rule{"<a>": NewRule("<c>", false), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false), "<d>": NewRule("", false)}, []string{}}, []string{}},
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
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		p   []string
		exp []string
		err error
	}{
		{[]string{""}, []string{}, nil},
		{[]string{";"}, []string{}, nil},
		{[]string{"", ""}, []string{}, nil},
		{[]string{";", ";", ";"}, []string{}, nil},
		{[]string{"abc;"}, []string{"abc"}, nil},
		{[]string{"<b>;"}, []string{"1", "2", "3"}, nil},
		{[]string{"<g>;"}, []string{"abc"}, nil},
		{[]string{"<h>;"}, []string{"123bc"}, nil},
		{[]string{"<i>;"}, []string{"a1c", "a2c", "a3c"}, nil},
		{[]string{"<j>;"}, []string{"a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123", "a3123123123"}, nil},
		{[]string{"<k>;"}, []string{"a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332", "a333"}, nil},
		{[]string{"a{}b//c/0.1/;"}, []string{"a{}b//c/0.1/"}, nil},
		{[]string{"a{}/0.1/b/0.1/{}c;"}, []string{"a{}/0.1/b/0.1/{}c"}, nil},
		{[]string{"abc<a>;"}, []string{"abc123"}, nil},
		{[]string{"abc<l>;"}, []string{"abc<l>"}, dummy_error},
		{[]string{"abc<a><b><c>;"}, []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313"}, nil},
		{[]string{"abc(<g>|<h>|<i>);"}, []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c"}, nil},
		{[]string{"abc<g><g><g>;"}, []string{"abcabcabcabc"}, nil},
		{[]string{"abc;", "def;", "ghi;"}, []string{"abc", "def", "ghi"}, nil},
		{[]string{"abc<a>;", "def<h>;", "<b>;"}, []string{"abc123", "def123bc", "1", "2", "3"}, nil},
		{[]string{"abc<a><b><c>;", "def[<h>];", "g((hi)|(jk));"}, []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313", "def", "def123bc", "ghi", "gjk"}, nil},
		{[]string{"abc(<g>|<h>|<i>);", "def[<g>|<h>|<i>];", "ghi<b><b><b>;"}, []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c", "def", "defabc", "def123bc", "defa1c", "defa2c", "defa3c", "ghi111", "ghi112", "ghi113", "ghi121", "ghi122", "ghi123", "ghi131", "ghi132", "ghi133", "ghi211", "ghi212", "ghi213", "ghi221", "ghi222", "ghi223", "ghi231", "ghi232", "ghi233", "ghi311", "ghi312", "ghi313", "ghi321", "ghi322", "ghi323", "ghi331", "ghi332", "ghi333"}, nil},
		{[]string{"abc<g><g><g>;", "<h><i>;", "<j>|<k>;"}, []string{"abcabcabcabc", "123bca1c", "123bca2c", "123bca3c", "a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123", "a3123123123", "a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332", "a333"}, nil},
	}
	for i, test := range table {
		g := NewGrammar("")
		g.Rules = map[string]Rule{
			"<_>": {"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			"<a>": {"123;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "123", ";", "<EOS>"}, []Expression{}},
			"<b>": {"1|2|3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 1.0}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
			"<c>": {"1[2]3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {2, 4, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, []Expression{}},
			"<d>": {"1(2)3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, []Expression{}},
			"<e>": {"1(2[3]);", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {4, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}, {7, 8, 1.0}, {8, 9, 1.0}}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{}},
			"<f>": {"<l>;", false, []string{"<f>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), []Expression{"<SOS>", "<f>", ";", "<EOS>"}, []Expression{}},
			"<g>": {"abc;", false, []string{""}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), []Expression{"<SOS>", "abc", ";", "<EOS>"}, []Expression{}},
			"<h>": {"<a>bc;", false, []string{"<a>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}), []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}, []Expression{}},
			"<i>": {"a<b>c;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}, []Expression{}},
			"<j>": {"a<b><c><d><e>;", false, []string{"<b>", "<c>", "<d>", "<e>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}, []Expression{}},
			"<k>": {"a<b><b><b>;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}, []Expression{}},
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
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		p   []string
		exp []string
		err error
	}{
		{[]string{""}, []string{}, nil},
		{[]string{";"}, []string{}, nil},
		{[]string{"", ""}, []string{}, nil},
		{[]string{";", ";", ";"}, []string{}, nil},
		{[]string{"abc;"}, []string{"abc"}, nil},
		{[]string{"<b>;"}, []string{"1", "2", "3"}, nil},
		{[]string{"<g>;"}, []string{"abc"}, nil},
		{[]string{"<h>;"}, []string{"123bc"}, nil},
		{[]string{"<i>;"}, []string{"a1c", "a2c", "a3c"}, nil},
		{[]string{"<j>;"}, []string{"a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123", "a3123123123"}, nil},
		{[]string{"<k>;"}, []string{"a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332", "a333"}, nil},
		{[]string{"a{}b//c/0.1/;"}, []string{"a{}b//c/0.1/"}, nil},
		{[]string{"a{}/0.1/b/0.1/{}c;"}, []string{"a{}/0.1/b/0.1/{}c"}, nil},
		{[]string{"abc<a>;"}, []string{"abc123"}, nil},
		{[]string{"abc<l>;"}, []string{"abc<l>"}, dummy_error},
		{[]string{"abc<a><b><c>;"}, []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313"}, nil},
		{[]string{"abc(<g>|<h>|<i>);"}, []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c"}, nil},
		{[]string{"abc<g><g><g>;"}, []string{"abcabcabcabc"}, nil},
		{[]string{"abc;", "def;", "ghi;"}, []string{"abc", "def", "ghi"}, nil},
		{[]string{"abc<a>;", "def<h>;", "<b>;"}, []string{"abc123", "def123bc", "1", "2", "3"}, nil},
		{[]string{"abc<a><b><c>;", "def[<h>];", "g((hi)|(jk));"}, []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313", "def", "def123bc", "ghi", "gjk"}, nil},
		{[]string{"abc(<g>|<h>|<i>);", "def[<g>|<h>|<i>];", "ghi<b><b><b>;"}, []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c", "def", "defabc", "def123bc", "defa1c", "defa2c", "defa3c", "ghi111", "ghi112", "ghi113", "ghi121", "ghi122", "ghi123", "ghi131", "ghi132", "ghi133", "ghi211", "ghi212", "ghi213", "ghi221", "ghi222", "ghi223", "ghi231", "ghi232", "ghi233", "ghi311", "ghi312", "ghi313", "ghi321", "ghi322", "ghi323", "ghi331", "ghi332", "ghi333"}, nil},
		{[]string{"abc<g><g><g>;", "<h><i>;", "<j>|<k>;"}, []string{"abcabcabcabc", "123bca1c", "123bca2c", "123bca3c", "a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123", "a3123123123", "a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332", "a333"}, nil},
	}
	for i, test := range table {
		g := NewGrammar("")
		g.Rules = map[string]Rule{
			"<_>": {"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			"<a>": {"123;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "123", ";", "<EOS>"}, []Expression{}},
			"<b>": {"1|2|3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 1.0}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
			"<c>": {"1[2]3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {2, 4, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, []Expression{}},
			"<d>": {"1(2)3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, []Expression{}},
			"<e>": {"1(2[3]);", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {4, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}, {7, 8, 1.0}, {8, 9, 1.0}}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{}},
			"<f>": {"<l>;", false, []string{"<f>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), []Expression{"<SOS>", "<f>", ";", "<EOS>"}, []Expression{}},
			"<g>": {"abc;", false, []string{""}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), []Expression{"<SOS>", "abc", ";", "<EOS>"}, []Expression{}},
			"<h>": {"<a>bc;", false, []string{"<a>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}), []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}, []Expression{}},
			"<i>": {"a<b>c;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}, []Expression{}},
			"<j>": {"a<b><c><d><e>;", false, []string{"<b>", "<c>", "<d>", "<e>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}, []Expression{}},
			"<k>": {"a<b><b><b>;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}, []Expression{}},
		}
		for j, p := range test.p {
			rule := NewRule(Expression(p), true)
			rule.Tokens = rule.Exp.ToTokens(lexer)
			rule.Graph = NewGraph(BuildEdgeList(rule.Tokens), rule.Tokens)
			rule.productions = FilterTerminals(rule.Tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>", ""})
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
		{"data/tests/test0.jsgf", "test0", []string{"import <a.*>;"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/test1.jsgf", "test1", []string{"import <c.brew>;"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}}},
		{"data/tests/test2.jsgf", "test2", []string{"import <a1.*>;"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/test3.jsgf", "test3", []string{"import <e.dne>;"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/test4.jsgf", "test4", []string{"import <d.*>;"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "quant": {}, "brew": {"quant"}}},
		{"data/tests/test5.jsgf", "test5", []string{"import <b.request>;"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/a.jsgf", "a", []string{}, map[string][]string{"request": {"brew"}, "order": {"quant"}, "brew": {"quant"}, "quant": {}}},
		{"data/tests/b.jsgf", "b", []string{"import <c.brew>;"}, map[string][]string{"request": {"brew"}, "order": {"quant"}, "quant": {}}},
		{"data/tests/dir0/c.jsgf", "c", []string{}, map[string][]string{"teatype": {}, "brew": {"quant"}, "quant": {}}},
		{"data/tests/dir0/dir1/d.jsgf", "d", []string{"import <c.teatype>;", "import <a.order>;"}, map[string][]string{}},
		{"data/tests/dir0/dir1/dir2/e.jsgf", "e", []string{}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
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
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	productions := []string{}
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
		{"data/tests/test0.jsgf", productions, nil},
		{"data/tests/test1.jsgf", productions, nil},
		{"data/tests/test2.jsgf", productions, dummy_error}, // returns error but is still resolvable, figure out how to validate completness before imports
		{"data/tests/test3.jsgf", productions, nil},
		{"data/tests/test4.jsgf", productions, nil},
		{"data/tests/test5.jsgf", productions, nil},
		{"data/tests/a.jsgf", []string{}, nil},
		{"data/tests/b.jsgf", []string{}, nil},
		{"data/tests/dir0/c.jsgf", []string{}, nil},
		{"data/tests/dir0/dir1/d.jsgf", []string{}, dummy_error},
		{"data/tests/dir0/dir1/dir2/e.jsgf", productions, nil},
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
