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
	table := []struct {
		pubs []string
		exp  []string
		err  error
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
	for _, test := range table {
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
		for j, p := range test.pubs {
			rule := NewRule(Expression(p), true)
			rule.tokens = rule.exp.ToTokens(lexer)
			rule.graph = NewGraph(BuildEdgeList(rule.tokens), rule.tokens)
			rule.productions = FilterTerminals(rule.tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			g.Rules[fmt.Sprintf("<pub_%v>", j)] = rule
		}
		g, err := g.Resolve(lexer)
		res := g.Productions()
		sort.Strings(test.exp)
		sort.Strings(res)
		if fmt.Sprint(res) != fmt.Sprint(test.exp) {
			t.Errorf("%v.Productions()\nGOT len %v %v\nEXP len %v %v", test.pubs, len(res), res, len(test.exp), test.exp)
		}
		if test.err != nil && err == nil {
			t.Errorf("%v.Productions().err\nGOT %v\nEXP %v", test.pubs, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("%v.Productions().err\nGOT %v\nEXP %v", test.pubs, err, test.err)
		}
	}
}

func TestGrammarPeek(t *testing.T) {
	table := []struct {
		p           string
		name        string
		exp_imports []string
		exp_rules   map[string][]string
	}{
		{"data/tests/test0.jsgf", "test0", []string{"import <a.*>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/test1.jsgf", "test1", []string{"import <c.brew>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/test2.jsgf", "test2", []string{"import <a1.*>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/test3.jsgf", "test3", []string{"import <e.dne>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/test4.jsgf", "test4", []string{"import <d.*>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "quant": {}, "brew": {"quant"}}},
		{"data/tests/test5.jsgf", "test5", []string{"import <b.request>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
		{"data/tests/a.jsgf", "a", []string{}, map[string][]string{"request": {"brew"}, "order": {"quant"}}},
		{"data/tests/b.jsgf", "b", []string{"import <c.brew>"}, map[string][]string{"request": {"brew"}, "order": {"quant"}, "quant": {}}},
		{"data/tests/dir0/c.jsgf", "c", []string{}, map[string][]string{"teatype": {}, "brew": {"quant"}}},
		{"data/tests/dir0/dir1/d.jsgf", "d", []string{"import <c.teatype>", "import <a.order>"}, map[string][]string{}},
		{"data/tests/dir0/dir1/dir2/e.jsgf", "e", []string{}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
	}
	for _, test := range table {
		res_name, res_imports, res_rules, err := NewGrammar(test.p).Peek()
		if err != nil {
			t.Errorf("Grammar(%v).Peek()\nGOT error %v", test.p, err)
		}

		if res_name != test.name {
			t.Errorf("Grammar(%v).Peek().imports\nGOT %v\nEXP %v", test.p, res_name, test.name)
		}

		sort.Strings(res_imports)
		sort.Strings(test.exp_imports)
		if fmt.Sprint(res_imports) != fmt.Sprint(test.exp_imports) {
			t.Errorf("Grammar(%v).Peek().imports\nGOT %v\nEXP %v", test.p, res_imports, test.exp_imports)
		}

		for k, v_res := range res_rules {
			v_exp, ok := test.exp_rules[k]
			if !ok {
				t.Errorf("Grammar(%v).Peek().rules\nGOT %v\nEXP %v", test.p, res_rules, test.exp_rules)
			}
			sort.Strings(v_exp)
			sort.Strings(v_res)
			if fmt.Sprint((v_exp)) != fmt.Sprint((v_res)) {
				t.Errorf("Grammar(%v).Peek().rules\nGOT %v\nEXP %v", test.p, res_rules, test.exp_rules)
			}
		}
	}
}
