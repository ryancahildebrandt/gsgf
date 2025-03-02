// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:12:02 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/bzick/tokenizer"
)

func TestResolveReferences(t *testing.T) {
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
		if !slices.Equal(GetReferences(res), GetReferences(test.exp)) {
			t.Errorf("test %v: %v.ResolveReferences().References\nGOT %v\nEXP %v", i, test.r, GetReferences(res), GetReferences(test.exp))
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

func TestParseRule(t *testing.T) {
	dummy_error := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		l   string
		n   string
		r   Rule
		err error
	}{
		{
			l: "", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: ";", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: " ", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "<rule> =", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "<rule> = ", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "public <rule> =", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "public <rule> = ", n: "", r: Rule{
				Exp: "", IsPublic: false,
			},
			err: dummy_error,
		},
		{
			l: "<rule> =;", n: "<rule>", r: Rule{
				Exp: ";", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> =;", n: "<rule>", r: Rule{
				Exp: ";", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "<rule> = test expression 123;", n: "<rule>", r: Rule{
				Exp: "test expression 123;", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "<rule> = test \"expression\" 123;", n: "<rule>", r: Rule{
				Exp: "test \"expression\" 123;", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> = test expression 123;", n: "<rule>", r: Rule{
				Exp: "test expression 123;", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "public <rule> = test \"expression\" 123;", n: "<rule>", r: Rule{
				Exp: "test \"expression\" 123;", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "<rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule> (abc) [def];", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule> (abc) [def];", IsPublic: true,
			},
			err: nil,
		},
		{
			l: "<rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: false,
			},
			err: nil,
		},
		{
			l: "public <rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{
				Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: true,
			},
			err: nil,
		},
	}
	for i, test := range table {
		n, r, err := ParseRule(test.l, lexer)
		if n != test.n {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, n, test.n)
		}
		if r.Exp != test.r.Exp {
			t.Errorf("test %v: ParseRule(jsgflexer, %v)\nGOT %v\nEXP %v", i, test.l, r, test.r)
		}
		if r.IsPublic != test.r.IsPublic {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).Is_public\nGOT %v\nEXP %v", i, test.l, r, test.r)
		}
		if !slices.Equal(GetReferences(r), GetReferences(test.r)) {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).References\nGOT %v\nEXP %v", i, test.l, GetReferences(r), GetReferences(test.r))
		}
		sort.Slice(r.Tokens, func(i, j int) bool { return r.Tokens[i] < r.Tokens[j] })
		sort.Slice(test.r.Tokens, func(i, j int) bool { return r.Tokens[i] < r.Tokens[j] })
		if !slices.EqualFunc(r.Tokens, test.r.Tokens, func(E1, E2 Expression) bool { return E1 == E2 }) {
			t.Errorf("test %v: ParseRule(jsgflexer, %v).Tokens\nGOT %v\nEXP %v", i, test.l, r.Tokens, test.r.Tokens)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ValidateJSGF(%v)\nGOT %v\nEXP %v", i, test.l, err, test.err)
		}
	}
}

func TestGetReferences(t *testing.T) {
	type args struct {
		r Rule
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetReferences(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleResolveReference(t *testing.T) {
	type args struct {
		r   Rule
		ref string
		r1  Rule
		lex *tokenizer.Tokenizer
	}
	tests := []struct {
		name    string
		args    args
		want    Rule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SingleResolveReference(tt.args.r, tt.args.ref, tt.args.r1, tt.args.lex)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleResolveReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleResolveReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateRuleRecursion(t *testing.T) {
	type args struct {
		r Rule
		m map[string]Rule
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
			if err := ValidateRuleRecursion(tt.args.r, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRuleRecursion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
