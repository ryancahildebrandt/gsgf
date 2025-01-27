// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:54 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	"sort"
	"testing"
)

func TestGrammarCompositionOrder(t *testing.T) {
	table := []struct {
		g   Grammar
		exp []string
	}{
		{Grammar{map[string]Rule{}, []string{}}, []string{}},
		{Grammar{map[string]Rule{"<a>": NewRule("", false)}, []string{}}, []string{}},
		{Grammar{map[string]Rule{"<a>": NewRule("<b><c>", false)}, []string{}}, []string{}},
		{Grammar{map[string]Rule{"<a>": NewRule("", true)}, []string{}}, []string{"<a>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<b><c>", true)}, []string{}}, []string{"<a>", "<b>", "<c>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", true)}, []string{}}, []string{"<a>", "<b>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("", true), "<b>": NewRule("", false)}, []string{}}, []string{"<a>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", true)}, []string{}}, []string{"<a>", "<b>", "<c>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<b>", true), "<c>": NewRule("", false)}, []string{}}, []string{"<a>", "<b>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true), "<d>": NewRule("", true)}, []string{}}, []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>", "<d>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", true), "<d>": NewRule("", false)}, []string{}}, []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", true), "<c>": NewRule("<d>", false), "<d>": NewRule("", false)}, []string{}}, []string{"<a>", "<b>", "<c>", "<c>", "<d>", "<d>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<c>", true), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false), "<d>": NewRule("", false)}, []string{}}, []string{"<a>", "<c>", "<d>"}},
		{Grammar{map[string]Rule{"<a>": NewRule("<c>", false), "<b>": NewRule("<c>", false), "<c>": NewRule("<d>", false), "<d>": NewRule("", false)}, []string{}}, []string{}},
	}
	for _, test := range table {
		res := test.g.CompositionOrder()
		sort.Strings(res)
		sort.Strings(test.exp)
		if fmt.Sprint(res) != fmt.Sprint(test.exp) {
			t.Errorf("%v.CompositionOrder()\nGOT %v\nEXP %v", test.g, res, test.exp)
		}
	}
}

func TestGrammarProductions(t *testing.T) {
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	m := map[string]Rule{
		"<_>": {"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
		"<a>": {"123;", false, []string{}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}}, []Expression{"<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "123", ";", "<EOS>"}, []Expression{}},
		"<b>": {"1|2|3;", false, []string{}, NewGraph(EdgeList{{0, 1, 0.0}, {0, 3, 0.0}, {0, 5, 0.0}, {1, 6, 0.0}, {3, 6, 0.0}, {5, 6, 0.0}, {6, 7, 0.0}}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
		"<c>": {"1[2]3;", false, []string{}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}, {2, 3, 0.0}, {2, 4, 0.0}, {3, 4, 0.0}, {4, 5, 0.0}, {5, 6, 0.0}, {6, 7, 0.0}}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, []Expression{}},
		"<d>": {"1(2)3;", false, []string{}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}, {2, 3, 0.0}, {3, 4, 0.0}, {4, 5, 0.0}, {5, 6, 0.0}, {6, 7, 0.0}}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, []Expression{}},
		"<e>": {"1(2[3]);", false, []string{}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}, {2, 3, 0.0}, {3, 4, 0.0}, {4, 5, 0.0}, {4, 6, 0.0}, {5, 6, 0.0}, {6, 7, 0.0}, {7, 8, 0.0}, {8, 9, 0.0}}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{}},
		"<f>": {"<l>;", false, []string{"<f>"}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), []Expression{"<SOS>", "<f>", ";", "<EOS>"}, []Expression{}},
		"<g>": {"abc;", false, []string{""}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), []Expression{"<SOS>", "abc", ";", "<EOS>"}, []Expression{}},
		"<h>": {"<a>bc;", false, []string{"<a>"}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}), []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}, []Expression{}},
		"<i>": {"a<b>c;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}, {2, 3, 0.0}, {3, 4, 0.0}}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}, []Expression{}},
		"<j>": {"a<b><c><d><e>;", false, []string{"<b>", "<c>", "<d>", "<e>"}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}, {2, 3, 0.0}, {3, 4, 0.0}, {4, 5, 0.0}, {5, 6, 0.0}, {6, 7, 0.0}}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}, []Expression{}},
		"<k>": {"a<b><b><b>;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 0.0}, {1, 2, 0.0}, {2, 3, 0.0}, {3, 4, 0.0}, {4, 5, 0.0}, {5, 6, 0.0}}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}, []Expression{}},
	}
	table := []struct {
		pubs []string
		exp  []string
		err  error
	}{
		{[]string{""}, []string{""}, nil},
		{[]string{";"}, []string{""}, nil},
		{[]string{"abc;"}, []string{"abc"}, nil},
		{[]string{"a{}b//c/0.1/;"}, []string{"a{}b//c/0.1/"}, nil},
		{[]string{"a{}/0.1/b/0.1/{}c;"}, []string{"a{}/0.1/b/0.1/{}c"}, nil},
		{[]string{"abc<a>;"}, []string{"abc123"}, nil},
		{[]string{"abc<l>;"}, []string{"abc<l>"}, dummy_error},
		{[]string{"abc<a><b><c>;"}, []string{"abc12313123", "abc1231323", "abc12323123", "abc1232323", "abc1233123", "abc123323"}, nil},
		{[]string{"abc(<g>|<h>|<i>);"}, []string{"abca1", "abcc1", "abc1"}, nil},
		{[]string{"abc<g><g><g>;"}, []string{"abcabcabc"}, nil},
		{[]string{"", ""}, []string{"", ""}, nil},
		{[]string{";", ";", ";"}, []string{"", "", ""}, nil},
		{[]string{"abc;", "def;", "ghi;"}, []string{"abc", "def", "ghi"}, nil},
		{[]string{"abc<a>;", "def<h>;", "<b>;"}, []string{"3", "3", "3", "123", "abc123"}, nil},
		{[]string{"abc<a><b><c>;", "def[<h>];", "g(hi)|(jk);"}, []string{"abc12313123", "abc1231323", "abc12323123", "abc1232323", "abc1233123", "abc123323", "ghi", "jk"}, nil},
		{[]string{"abc(<g>|<h>|<i>);", "def[<g>|<h>|<i>];", "ghi<b><b><b>;"}, []string{"ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333", "ghi333"}, nil},
		{[]string{"abc<g><g><g>;", "<h><i>;", "<j>|<k>;"}, []string{"3", "3", "3", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "3123123123", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "333", "abcabcabc"}, nil},
	}
	for _, test := range table {
		g := NewGrammar()
		g.Rules = m
		for j, p := range test.pubs {
			rule := NewRule(Expression(p), true)
			rule.tokens = rule.exp.ToTokens(lexer)
			rule.graph = NewGraph(BuildEdgeList(rule.tokens), rule.tokens)
			rule.productions = FilterTerminals(rule.tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			g.Rules[fmt.Sprintf("<%v>", j)] = rule
		}
		g, err := g.Resolve()
		res := g.Productions()
		sort.Strings(test.exp)
		sort.Strings(res)
		if fmt.Sprint(res) != fmt.Sprint(test.exp) {
			t.Errorf("%v.Productions()\nGOT %v\nEXP %v", test.pubs, res, test.exp)
		}
		if test.err != nil && err == nil {
			t.Errorf("%v.ResolveReferences().err\nGOT %v\nEXP %v", test.pubs, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("%v.ResolveReferences().err\nGOT %v\nEXP %v", test.pubs, err, test.err)
		}
	}
}
