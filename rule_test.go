// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:12:02 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"
	"sort"
	"testing"
)

func TestRuleResolveReferences(t *testing.T) {
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	m := map[string]Rule{
		"<a>": {"123;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "123", ";", "<EOS>"}, []Expression{}},
		"<b>": {"1|2|3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 1.0}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
		"<c>": {"1[2]3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {2, 4, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, []Expression{}},
		"<d>": {"1(2)3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, []Expression{}},
		"<e>": {"1(2[3]);", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {4, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}, {7, 8, 1.0}, {8, 9, 1.0}}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{}},
	}
	table := []struct {
		r   Rule
		exp Rule
		err error
	}{
		{
			Rule{"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			Rule{"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			nil,
		},
		{
			Rule{"", true, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			Rule{"", true, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			nil,
		},
		{
			Rule{"<f>;", false, []string{"<f>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), []Expression{"<SOS>", "<f>", ";", "<EOS>"}, []Expression{}},
			Rule{"<f>;", false, []string{"<f>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), []Expression{"<SOS>", "<f>", ";", "<EOS>"}, []Expression{}},
			dummy_error,
		},
		{
			Rule{"<f>;", true, []string{"<f>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), []Expression{"<SOS>", "<f>", ";", "<EOS>"}, []Expression{}},
			Rule{"<f>;", true, []string{"<f>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}), []Expression{"<SOS>", "<f>", ";", "<EOS>"}, []Expression{}},
			dummy_error,
		},
		{
			Rule{"abc;", false, []string{""}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), []Expression{"<SOS>", "abc", ";", "<EOS>"}, []Expression{}},
			Rule{"abc;", false, []string{""}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), []Expression{"<SOS>", "abc", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"abc;", true, []string{""}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), []Expression{"<SOS>", "abc", ";", "<EOS>"}, []Expression{}},
			Rule{"abc;", true, []string{""}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "abc", ";", "<EOS>"}), []Expression{"<SOS>", "abc", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"<a>bc;", false, []string{"<a>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}), []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}, []Expression{}},
			Rule{"<a>bc;", false, []string{"<a>"}, NewGraph(EdgeList{{0, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 2, 1.0}}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"<a>bc;", true, []string{"<a>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}), []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}, []Expression{}},
			Rule{"<a>bc;", true, []string{"<a>"}, NewGraph(EdgeList{{0, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 2, 1.0}}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"a<b>c;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}, []Expression{}},
			Rule{"a<b>c;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 5, 1.0}, {10, 11, 1.0}, {11, 12, 1.0}, {3, 4, 1.0}, {5, 10, 1.0}, {5, 6, 1.0}, {5, 8, 1.0}, {6, 11, 1.0}, {8, 11, 1.0}, {12, 3, 1.0}}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"a<b>c;", true, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"}, []Expression{}},
			Rule{"a<b>c;", true, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 5, 1.0}, {10, 11, 1.0}, {11, 12, 1.0}, {3, 4, 1.0}, {5, 10, 1.0}, {5, 6, 1.0}, {5, 8, 1.0}, {6, 11, 1.0}, {8, 11, 1.0}, {12, 3, 1.0}}, []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"a<b><c><d><e>;", false, []string{"<b>", "<c>", "<d>", "<e>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}, []Expression{}},
			Rule{"a<b><c><d><e>;", false, []string{"<b>", "<c>", "<d>", "<e>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 8, 1.0}, {8, 9, 1.0}, {8, 11, 1.0}, {8, 13, 1.0}, {9, 14, 1.0}, {11, 14, 1.0}, {13, 14, 1.0}, {14, 15, 1.0}, {15, 16, 1.0}, {16, 17, 1.0}, {17, 18, 1.0}, {18, 19, 1.0}, {19, 20, 1.0}, {20, 21, 1.0}, {21, 22, 1.0}, {22, 23, 1.0}, {18, 20, 1.0}, {23, 24, 1.0}, {24, 25, 1.0}, {25, 26, 1.0}, {26, 27, 1.0}, {27, 28, 1.0}, {28, 29, 1.0}, {29, 30, 1.0}, {30, 31, 1.0}, {31, 32, 1.0}, {32, 33, 1.0}, {33, 34, 1.0}, {34, 35, 1.0}, {35, 36, 1.0}, {36, 37, 1.0}, {36, 38, 1.0}, {37, 38, 1.0}, {38, 39, 1.0}, {39, 40, 1.0}, {40, 41, 1.0}, {41, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"a<b><c><d><e>;", true, []string{"<b>", "<c>", "<d>", "<e>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"}, []Expression{}},
			Rule{"a<b><c><d><e>;", true, []string{"<b>", "<c>", "<d>", "<e>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 8, 1.0}, {8, 9, 1.0}, {8, 11, 1.0}, {8, 13, 1.0}, {9, 14, 1.0}, {11, 14, 1.0}, {13, 14, 1.0}, {14, 15, 1.0}, {15, 16, 1.0}, {16, 17, 1.0}, {17, 18, 1.0}, {18, 19, 1.0}, {19, 20, 1.0}, {20, 21, 1.0}, {21, 22, 1.0}, {22, 23, 1.0}, {18, 20, 1.0}, {23, 24, 1.0}, {24, 25, 1.0}, {25, 26, 1.0}, {26, 27, 1.0}, {27, 28, 1.0}, {28, 29, 1.0}, {29, 30, 1.0}, {30, 31, 1.0}, {31, 32, 1.0}, {32, 33, 1.0}, {33, 34, 1.0}, {34, 35, 1.0}, {35, 36, 1.0}, {36, 37, 1.0}, {36, 38, 1.0}, {37, 38, 1.0}, {38, 39, 1.0}, {39, 40, 1.0}, {40, 41, 1.0}, {41, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"a<b><b><b>;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}, []Expression{}},
			Rule{"a<b><b><b>;", false, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 7, 1.0}, {7, 8, 1.0}, {7, 10, 1.0}, {7, 12, 1.0}, {8, 13, 1.0}, {10, 13, 1.0}, {12, 13, 1.0}, {13, 14, 1.0}, {14, 15, 1.0}, {15, 16, 1.0}, {15, 18, 1.0}, {15, 20, 1.0}, {16, 21, 1.0}, {18, 21, 1.0}, {20, 21, 1.0}, {21, 22, 1.0}, {22, 23, 1.0}, {23, 24, 1.0}, {23, 26, 1.0}, {23, 28, 1.0}, {24, 29, 1.0}, {26, 29, 1.0}, {28, 29, 1.0}, {29, 30, 1.0}, {30, 5, 1.0}, {5, 6, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
			nil,
		},
		{
			Rule{"a<b><b><b>;", true, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}, []Expression{}},
			Rule{"a<b><b><b>;", true, []string{"<b>"}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 7, 1.0}, {7, 8, 1.0}, {7, 10, 1.0}, {7, 12, 1.0}, {8, 13, 1.0}, {10, 13, 1.0}, {12, 13, 1.0}, {13, 14, 1.0}, {14, 15, 1.0}, {15, 16, 1.0}, {15, 18, 1.0}, {15, 20, 1.0}, {16, 21, 1.0}, {18, 21, 1.0}, {20, 21, 1.0}, {21, 22, 1.0}, {22, 23, 1.0}, {23, 24, 1.0}, {23, 26, 1.0}, {23, 28, 1.0}, {24, 29, 1.0}, {26, 29, 1.0}, {28, 29, 1.0}, {29, 30, 1.0}, {30, 5, 1.0}, {5, 6, 1.0}}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{}},
			nil,
		},
	}
	for i, test := range table {
		res, err := test.r.ResolveReferences(m, lexer)
		if test.exp.Is_public != res.Is_public {
			t.Errorf("test %v: %v.ResolveReferences().Is_public\nGOT %v\nEXP %v", i, test.r, res.Is_public, test.exp.Is_public)
		}
		if !slices.Equal(test.exp.References, res.References) {
			t.Errorf("test %v: %v.ResolveReferences().References\nGOT %v\nEXP %v", i, test.r, res.References, test.exp.References)
		}
		if !slices.Equal(test.exp.Graph.Edges.Sort(), res.Graph.Edges.Sort()) {
			t.Errorf("test %v: %v.ResolveReferences().edges\nGOT %v\nEXP %v", i, test.r, res.Graph.Edges.Sort(), test.exp.Graph.Edges.Sort())
		}
		if !slices.Equal(test.exp.Graph.Nodes, res.Graph.Nodes) {
			t.Errorf("test %v: %v.ResolveReferences().nodes\nGOT %v\nEXP %v", i, test.r, res.Graph.Nodes, test.exp.Graph.Nodes)
		}
		if !slices.Equal(test.exp.Tokens, res.Tokens) {
			t.Errorf("test %v: %v.ResolveReferences().Tokens\nGOT %v\nEXP %v", i, test.r, res.Tokens, test.exp.Tokens)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.ResolveReferences().err\nGOT %v\nEXP %v", i, test.r, err, test.err)
		}
	}
}

func TestRuleProductions(t *testing.T) {
	table := []struct {
		r   Rule
		exp []string
	}{
		{
			Rule{"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			[]string{},
		},
		{
			Rule{";", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", ";", "<EOS>"}), []Expression{"<SOS>", ";", "<EOS>"}, []Expression{"", "", ""}},
			[]string{},
		},
		{
			Rule{"123;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "123", ";", "<EOS>"}, []Expression{"", "123", "", ""}},
			[]string{"123"},
		},
		{
			Rule{"1|2|3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 1.0}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			[]string{"1", "2", "3"},
		},
		{
			Rule{"1{}|2//|3/0.1/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 1.0}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1{}", "|", "2//", "|", "3/0.1/", ";", "<EOS>"}), []Expression{"<SOS>", "1{}", "|", "2//", "|", "3/0.1/", ";", "<EOS>"}, []Expression{"", "1{}", "", "2//", "", "3/0.1/", "", ""}},
			[]string{"1{}", "2//", "3/0.1/"},
		},
		{
			Rule{"1[2]3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {2, 4, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			[]string{"123", "13"},
		},
		{
			Rule{"1(2)3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			[]string{"123"},
		},
		{
			Rule{"1(2[3]);", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {4, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}, {7, 8, 1.0}, {8, 9, 1.0}}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", "", "", ""}},
			[]string{"12", "123"},
		},
	}
	for i, test := range table {
		res := test.r.Productions()
		sort.Strings(res)
		sort.Strings(test.exp)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.Productions()\nGOT %v\nEXP %v", i, test.r, res, test.exp)
		}
	}
}

func TestRuleWeightEdges(t *testing.T) {
	dummy_error := errors.New("")
	table := []struct {
		r   Rule
		exp Rule
		err error
	}{
		{
			Rule{"/.//", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{"/.//"}, []Expression{}},
			Rule{"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{"/.//"}, []Expression{}},
			dummy_error,
		},
		{
			Rule{"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			Rule{"", false, []string{}, NewGraph(EdgeList{}, []Expression{}), []Expression{}, []Expression{}},
			nil,
		},
		{
			Rule{"/.99/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}}, []Expression{"<SOS>", "/.99/", ";", "<EOS>"}), []Expression{"<SOS>", "/.99/", ";", "<EOS>"}, []Expression{"", "", ""}},
			Rule{"/.99/;", false, []string{}, NewGraph(EdgeList{{0, 1, 0.99}, {1, 2, 1.0}}, []Expression{"<SOS>", "", ";", "<EOS>"}), []Expression{"<SOS>", "", ";", "<EOS>"}, []Expression{"", "", ""}},
			nil,
		},
		{
			Rule{"123/.99/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "123/.99/", ";", "<EOS>"}), []Expression{"<SOS>", "123/.99/", ";", "<EOS>"}, []Expression{"", "123", "", ""}},
			Rule{"123/.99/;", false, []string{}, NewGraph(EdgeList{{0, 1, 0.99}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "123", ";", "<EOS>"}), []Expression{"<SOS>", "123", ";", "<EOS>"}, []Expression{"", "123", "", ""}},
			nil,
		},
		{
			Rule{"<123>/.99/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "<123>/.99/", ";", "<EOS>"}), []Expression{"<SOS>", "<123>/.99/", ";", "<EOS>"}, []Expression{"", "<123>", "", ""}},
			Rule{"<123>/.99/;", false, []string{}, NewGraph(EdgeList{{0, 1, 0.99}, {1, 2, 1.0}, {2, 3, 1.0}}, []Expression{"<SOS>", "<123>", ";", "<EOS>"}), []Expression{"<SOS>", "<123>", ";", "<EOS>"}, []Expression{"", "<123>", "", ""}},
			nil,
		},
		{
			Rule{"1|2|3/0.1/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 1.0}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "|", "2", "|", "3/0.1/", ";", "<EOS>"}), []Expression{"<SOS>", "1", "|", "2", "|", "3/0.1/", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			Rule{"1|2|3/0.1/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 0.1}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			nil,
		},
		{
			Rule{"1{}|2//|3/0.1/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 1.0}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1{}", "|", "2//", "|", "3/0.1/", ";", "<EOS>"}), []Expression{"<SOS>", "1{}", "|", "2//", "|", "3/0.1/", ";", "<EOS>"}, []Expression{"", "1{}", "", "2//", "", "3/0.1/", "", ""}},
			Rule{"1{}|2//|3/0.1/;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {0, 3, 1.0}, {0, 5, 0.1}, {1, 6, 1.0}, {3, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1{}", "|", "2//", "|", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1{}", "|", "2//", "|", "3", ";", "<EOS>"}, []Expression{"", "1{}", "", "2//", "", "3", "", ""}},
			nil,
		},
		{
			Rule{"1/0.1/[2]3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {2, 4, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1/1.01/", "[", "2", "]", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1/1.01/", "[", "2", "]", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			Rule{"1/0.1/[2]3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.01}, {1, 2, 1.0}, {2, 3, 1.0}, {2, 4, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			nil,
		},
		{
			Rule{"1(2/1.01/)3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "(", "2/1.01/", ")", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2/1.01/", ")", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			Rule{"1(2/1.01/)3;", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.01}, {3, 4, 1.0}, {4, 5, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", ""}},
			nil,
		},
		{
			Rule{"1/1.01/(2[3]);", false, []string{}, NewGraph(EdgeList{{0, 1, 1.0}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {4, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}, {7, 8, 1.0}, {8, 9, 1.0}}, []Expression{"<SOS>", "1/1.01/", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "1/1.01/", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", "", "", ""}},
			Rule{"1/1.01/(2[3]);", false, []string{}, NewGraph(EdgeList{{0, 1, 1.01}, {1, 2, 1.0}, {2, 3, 1.0}, {3, 4, 1.0}, {4, 5, 1.0}, {4, 6, 1.0}, {5, 6, 1.0}, {6, 7, 1.0}, {7, 8, 1.0}, {8, 9, 1.0}}, []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}), []Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}, []Expression{"", "1", "", "2", "", "3", "", "", "", ""}},
			nil,
		},
	}
	for i, test := range table {
		res, err := test.r.WeightEdges()
		if test.exp.Is_public != res.Is_public {
			t.Errorf("test %v: %v.WeightEdges().Is_public\nGOT %v\nEXP %v", i, test.r, res.Is_public, test.exp.Is_public)
		}
		if !slices.Equal(test.exp.References, res.References) {
			t.Errorf("test %v: %v.WeightEdges().References\nGOT %v\nEXP %v", i, test.r, res.References, test.exp.References)
		}
		if !slices.Equal(test.exp.Graph.Edges.Sort(), res.Graph.Edges.Sort()) {
			t.Errorf("test %v: %v.WeightEdges().edges\nGOT %v\nEXP %v", i, test.r, res.Graph.Edges.Sort(), test.exp.Graph.Edges.Sort())
		}
		if !slices.Equal(test.exp.Graph.Nodes, res.Graph.Nodes) {
			t.Errorf("test %v: %v.WeightEdges().nodes\nGOT %v\nEXP %v", i, test.r, res.Graph.Nodes, test.exp.Graph.Nodes)
		}
		if !slices.Equal(test.exp.Tokens, res.Tokens) {
			t.Errorf("test %v: %v.WeightEdges().Tokens\nGOT %v\nEXP %v", i, test.r, res.Tokens, test.exp.Tokens)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.WeightEdges().err\nGOT %v\nEXP %v", i, test.r, err, test.err)
		}
	}
}
