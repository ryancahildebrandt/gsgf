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
	dummyError := errors.New("")
	lexer := NewJSGFLexer()
	m := map[string]Rule{
		"<a>": {
			Exp:      "123;",
			IsPublic: false,
			Graph: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
			}, []Expression{"<SOS>", "123", ";", "<EOS>"}),
		},
		"<b>": {
			Exp:      "1|2|3;",
			IsPublic: false,
			Graph: NewGraph(EdgeList{
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
			Exp:      "1[2]3;",
			IsPublic: false,
			Graph: NewGraph(EdgeList{
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
			Exp:      "1(2)3;",
			IsPublic: false,
			Graph: NewGraph(EdgeList{
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
			Exp:      "1(2[3]);",
			IsPublic: false,
			Graph: NewGraph(EdgeList{
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
	}
	table := []struct {
		r   Rule
		exp Rule
		err error
	}{
		{
			r: Rule{
				Exp:      "",
				IsPublic: false,
				Graph:    NewGraph(EdgeList{}, []Expression{}),
			},
			exp: Rule{
				Exp:      "",
				IsPublic: false,
				Graph:    NewGraph(EdgeList{}, []Expression{}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "", IsPublic: true, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			exp: Rule{
				Exp: "", IsPublic: true, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "<f>;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "<f>;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}),
			},
			err: dummyError,
		},
		{
			r: Rule{
				Exp: "<f>;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "<f>;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<f>", ";", "<EOS>"}),
			},
			err: dummyError,
		},
		{
			r: Rule{
				Exp: "abc;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "abc;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "abc;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "abc;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "abc", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "<a>bc;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "<a>bc;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 2, Weight: 1.0},
				}, []Expression{
					"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "<a>bc;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "<a>bc;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 2, Weight: 1.0},
				}, []Expression{
					"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "a<b>c;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "c", ";", "<EOS>",
				}),
			},
			exp: Rule{
				Exp: "a<b>c;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 5, Weight: 1.0},
					{From: 10, To: 11, Weight: 1.0},
					{From: 11, To: 12, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 5, To: 10, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 5, To: 8, Weight: 1.0},
					{From: 6, To: 11, Weight: 1.0},
					{From: 8, To: 11, Weight: 1.0},
					{From: 12, To: 3, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "a<b>c;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "c", ";", "<EOS>",
				}),
			},
			exp: Rule{
				Exp: "a<b>c;", IsPublic: true, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 5, Weight: 1.0},
					{From: 10, To: 11, Weight: 1.0},
					{From: 11, To: 12, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 5, To: 10, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 5, To: 8, Weight: 1.0},
					{From: 6, To: 11, Weight: 1.0},
					{From: 8, To: 11, Weight: 1.0},
					{From: 12, To: 3, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
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
			exp: Rule{
				Exp: "a<b><c><d><e>;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 8, Weight: 1.0},
					{From: 8, To: 9, Weight: 1.0},
					{From: 8, To: 11, Weight: 1.0},
					{From: 8, To: 13, Weight: 1.0},
					{From: 9, To: 14, Weight: 1.0},
					{From: 11, To: 14, Weight: 1.0},
					{From: 13, To: 14, Weight: 1.0},
					{From: 14, To: 15, Weight: 1.0},
					{From: 15, To: 16, Weight: 1.0},
					{From: 16, To: 17, Weight: 1.0},
					{From: 17, To: 18, Weight: 1.0},
					{From: 18, To: 19, Weight: 1.0},
					{From: 19, To: 20, Weight: 1.0},
					{From: 20, To: 21, Weight: 1.0},
					{From: 21, To: 22, Weight: 1.0},
					{From: 22, To: 23, Weight: 1.0},
					{From: 18, To: 20, Weight: 1.0},
					{From: 23, To: 24, Weight: 1.0},
					{From: 24, To: 25, Weight: 1.0},
					{From: 25, To: 26, Weight: 1.0},
					{From: 26, To: 27, Weight: 1.0},
					{From: 27, To: 28, Weight: 1.0},
					{From: 28, To: 29, Weight: 1.0},
					{From: 29, To: 30, Weight: 1.0},
					{From: 30, To: 31, Weight: 1.0},
					{From: 31, To: 32, Weight: 1.0},
					{From: 32, To: 33, Weight: 1.0},
					{From: 33, To: 34, Weight: 1.0},
					{From: 34, To: 35, Weight: 1.0},
					{From: 35, To: 36, Weight: 1.0},
					{From: 36, To: 37, Weight: 1.0},
					{From: 36, To: 38, Weight: 1.0},
					{From: 37, To: 38, Weight: 1.0},
					{From: 38, To: 39, Weight: 1.0},
					{From: 39, To: 40, Weight: 1.0},
					{From: 40, To: 41, Weight: 1.0},
					{From: 41, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";",
					"<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";",
					"<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "a<b><c><d><e>;", IsPublic: true,
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
			exp: Rule{
				Exp: "a<b><c><d><e>;", IsPublic: true,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 8, Weight: 1.0},
					{From: 8, To: 9, Weight: 1.0},
					{From: 8, To: 11, Weight: 1.0},
					{From: 8, To: 13, Weight: 1.0},
					{From: 9, To: 14, Weight: 1.0},
					{From: 11, To: 14, Weight: 1.0},
					{From: 13, To: 14, Weight: 1.0},
					{From: 14, To: 15, Weight: 1.0},
					{From: 15, To: 16, Weight: 1.0},
					{From: 16, To: 17, Weight: 1.0},
					{From: 17, To: 18, Weight: 1.0},
					{From: 18, To: 19, Weight: 1.0},
					{From: 19, To: 20, Weight: 1.0},
					{From: 20, To: 21, Weight: 1.0},
					{From: 21, To: 22, Weight: 1.0},
					{From: 22, To: 23, Weight: 1.0},
					{From: 18, To: 20, Weight: 1.0},
					{From: 23, To: 24, Weight: 1.0},
					{From: 24, To: 25, Weight: 1.0},
					{From: 25, To: 26, Weight: 1.0},
					{From: 26, To: 27, Weight: 1.0},
					{From: 27, To: 28, Weight: 1.0},
					{From: 28, To: 29, Weight: 1.0},
					{From: 29, To: 30, Weight: 1.0},
					{From: 30, To: 31, Weight: 1.0},
					{From: 31, To: 32, Weight: 1.0},
					{From: 32, To: 33, Weight: 1.0},
					{From: 33, To: 34, Weight: 1.0},
					{From: 34, To: 35, Weight: 1.0},
					{From: 35, To: 36, Weight: 1.0},
					{From: 36, To: 37, Weight: 1.0},
					{From: 36, To: 38, Weight: 1.0},
					{From: 37, To: 38, Weight: 1.0},
					{From: 38, To: 39, Weight: 1.0},
					{From: 39, To: 40, Weight: 1.0},
					{From: 40, To: 41, Weight: 1.0},
					{From: 41, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";",
					"<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";",
					"<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "a<b><b><b>;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "a<b><b><b>;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 7, Weight: 1.0},
					{From: 7, To: 8, Weight: 1.0},
					{From: 7, To: 10, Weight: 1.0},
					{From: 7, To: 12, Weight: 1.0},
					{From: 8, To: 13, Weight: 1.0},
					{From: 10, To: 13, Weight: 1.0},
					{From: 12, To: 13, Weight: 1.0},
					{From: 13, To: 14, Weight: 1.0},
					{From: 14, To: 15, Weight: 1.0},
					{From: 15, To: 16, Weight: 1.0},
					{From: 15, To: 18, Weight: 1.0},
					{From: 15, To: 20, Weight: 1.0},
					{From: 16, To: 21, Weight: 1.0},
					{From: 18, To: 21, Weight: 1.0},
					{From: 20, To: 21, Weight: 1.0},
					{From: 21, To: 22, Weight: 1.0},
					{From: 22, To: 23, Weight: 1.0},
					{From: 23, To: 24, Weight: 1.0},
					{From: 23, To: 26, Weight: 1.0},
					{From: 23, To: 28, Weight: 1.0},
					{From: 24, To: 29, Weight: 1.0},
					{From: 26, To: 29, Weight: 1.0},
					{From: 28, To: 29, Weight: 1.0},
					{From: 29, To: 30, Weight: 1.0},
					{From: 30, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
					"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "a<b><b><b>;", IsPublic: true, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "a<b><b><b>;", IsPublic: true, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 7, Weight: 1.0},
					{From: 7, To: 8, Weight: 1.0},
					{From: 7, To: 10, Weight: 1.0},
					{From: 7, To: 12, Weight: 1.0},
					{From: 8, To: 13, Weight: 1.0},
					{From: 10, To: 13, Weight: 1.0},
					{From: 12, To: 13, Weight: 1.0},
					{From: 13, To: 14, Weight: 1.0},
					{From: 14, To: 15, Weight: 1.0},
					{From: 15, To: 16, Weight: 1.0},
					{From: 15, To: 18, Weight: 1.0},
					{From: 15, To: 20, Weight: 1.0},
					{From: 16, To: 21, Weight: 1.0},
					{From: 18, To: 21, Weight: 1.0},
					{From: 20, To: 21, Weight: 1.0},
					{From: 21, To: 22, Weight: 1.0},
					{From: 22, To: 23, Weight: 1.0},
					{From: 23, To: 24, Weight: 1.0},
					{From: 23, To: 26, Weight: 1.0},
					{From: 23, To: 28, Weight: 1.0},
					{From: 24, To: 29, Weight: 1.0},
					{From: 26, To: 29, Weight: 1.0},
					{From: 28, To: 29, Weight: 1.0},
					{From: 29, To: 30, Weight: 1.0},
					{From: 30, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
				}, []Expression{
					"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
					"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
				}),
			},
			err: nil,
		},
	}
	for i, test := range table {
		res, err := ResolveReferences(test.r, m, lexer)
		if test.exp.IsPublic != res.IsPublic {
			t.Errorf("test %v: %v.ResolveReferences().Is_public\nGOT %v\nEXP %v", i, test.r, res.IsPublic, test.exp.IsPublic)
		}
		if !slices.Equal(References(res), References(test.exp)) {
			t.Errorf("test %v: %v.ResolveReferences().References\nGOT %v\nEXP %v", i, test.r, References(res), References(test.exp))
		}
		if !slices.Equal(Sort(test.exp.Graph.Edges), Sort(res.Graph.Edges)) {
			t.Errorf("test %v: %v.ResolveReferences().edges\nGOT %v\nEXP %v", i, test.r, Sort(res.Graph.Edges), Sort(test.exp.Graph.Edges))
		}
		if !slices.Equal(test.exp.Graph.Tokens, res.Graph.Tokens) {
			t.Errorf("test %v: %v.ResolveReferences().nodes\nGOT %v\nEXP %v", i, test.r, res.Graph.Tokens, test.exp.Graph.Tokens)
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
			r: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			exp: []string{},
		},
		{
			r: Rule{
				Exp: ";", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", ";", "<EOS>"}),
			},
			exp: []string{},
		},
		{
			r: Rule{
				Exp: "123;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}),
			},
			exp: []string{"123"},
		},
		{
			r: Rule{
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
			exp: []string{"1", "2", "3"},
		},
		{
			r: Rule{
				Exp: "1{}|2//|3/0.1/;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1{}", "|", "2//", "|", "3/0.1/", ";", "<EOS>"}),
			},
			exp: []string{"1{}", "2//", "3/0.1/"},
		},
		{
			r: Rule{
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
			exp: []string{"123", "13"},
		},
		{
			r: Rule{
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
			exp: []string{"123"},
		},
		{
			r: Rule{
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
			exp: []string{"12", "123"},
		},
	}
	for i, test := range table {
		res := Productions(test.r)
		sort.Strings(res)
		sort.Strings(test.exp)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.Productions()\nGOT %v\nEXP %v", i, test.r, res, test.exp)
		}
	}
}

func TestRuleWeightEdges(t *testing.T) {
	// dummyError := errors.New("")
	table := []struct {
		r   Rule
		exp Rule
		err error
	}{
		{
			r: Rule{
				Exp: "/.//", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			exp: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			exp: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "/.99/", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "123/.99/;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{
					"<SOS>", "123/.99/", ";", "<EOS>",
				}),
			},
			exp: Rule{
				Exp: "123/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "<123>/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{
					"<SOS>", "<123>/.99/", ";", "<EOS>",
				}),
			},
			exp: Rule{
				Exp: "<123>/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{
					"<SOS>", "<123>", ";", "<EOS>",
				}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "1|2|3/0.1/;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "|", "2", "|", "3/0.1/", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "1|2|3/0.1/;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 0.1},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "1{}|2//|3/0.1/;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1{}", "|", "2//", "|", "3/0.1/", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "1{}|2//|3/0.1/;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 0.1},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1{}", "|", "2//", "|", "3", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "1/0.1/[2]3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1/1.01/", "[", "2", "]", "3", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "1/0.1/[2]3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.01},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "1(2/1.01/)3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2/1.01/", ")", "3", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "1(2/1.01/)3;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.01},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				}, []Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"}),
			},
			err: nil,
		},
		{
			r: Rule{
				Exp: "1/1.01/(2[3]);", IsPublic: false, Graph: NewGraph(EdgeList{
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
				}, []Expression{"<SOS>", "1/1.01/", "(", "2", "[", "3", "]", ")", ";", "<EOS>"}),
			},
			exp: Rule{
				Exp: "1/1.01/(2[3]);", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.01},
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
			err: nil,
		},
	}
	for i, test := range table {
		res, err := WeightEdges(test.r)
		if test.exp.IsPublic != res.IsPublic {
			t.Errorf("test %v: %v.WeightEdges().Is_public\nGOT %v\nEXP %v", i, test.r, res.IsPublic, test.exp.IsPublic)
		}
		if !slices.Equal(References(res), References(test.exp)) {
			t.Errorf("test %v: %v.WeightEdges().References\nGOT %v\nEXP %v", i, test.r, References(res), References(test.exp))
		}
		if !slices.Equal(Sort(test.exp.Graph.Edges), Sort(res.Graph.Edges)) {
			t.Errorf("test %v: %v.WeightEdges().edges\nGOT %v\nEXP %v", i, test.r, Sort(res.Graph.Edges), Sort(test.exp.Graph.Edges))
		}
		if !slices.Equal(test.exp.Graph.Tokens, res.Graph.Tokens) {
			t.Errorf("test %v: %v.WeightEdges().nodes\nGOT %v\nEXP %v", i, test.r, res.Graph.Tokens, test.exp.Graph.Tokens)
		}
		if !slices.Equal(test.exp.Tokens, res.Tokens) {
			t.Errorf("test %v: %v.WeightEdges().Tokens\nGOT %v\nEXP %v", i, test.r, res.Tokens, test.exp.Tokens)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.WeightEdges().err\nGOT %v\nEXP %v", i, test.r, err, test.err)
		}
	}
}
