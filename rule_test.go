// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:12:02 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"slices"
	"sort"
	"testing"
)

func TestResolveReferences(t *testing.T) {
	lexer := NewJSGFLexer()
	m := map[string]Rule{
		"<a>": {
			Exp:      "123;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
				},
				[]Expression{"<SOS>", "123", ";", "<EOS>"},
			),
		},
		"<b>": {
			Exp:      "1|2|3;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				},
				[]Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"},
			),
		},
		"<c>": {
			Exp:      "1[2]3;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				},
				[]Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"},
			),
		},
		"<d>": {
			Exp:      "1(2)3;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				},
				[]Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"},
			),
		},
		"<e>": {
			Exp:      "1(2[3]);",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
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
				},
				[]Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"},
			),
		},
	}
	table := []struct {
		r       Rule
		want    Rule
		wantErr bool
	}{
		{
			r: Rule{
				Exp:      "",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{},
					[]Expression{},
				),
			},
			want: Rule{
				Exp:      "",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{},
					[]Expression{},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "", IsPublic: true, Graph: NewGraph(
					EdgeList{},
					[]Expression{},
				),
			},
			want: Rule{
				Exp: "", IsPublic: true, Graph: NewGraph(
					EdgeList{},
					[]Expression{},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "<f>;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<f>", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "<f>;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<f>", ";", "<EOS>"},
				),
			},
			wantErr: true,
		},
		{
			r: Rule{
				Exp: "<f>;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<f>", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "<f>;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<f>", ";", "<EOS>"},
				),
			},
			wantErr: true,
		},
		{
			r: Rule{
				Exp: "abc;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "abc", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "abc;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "abc", ";", "<EOS>"},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "abc;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "abc", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "abc;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "abc", ";", "<EOS>"},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "<a>bc;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "<a>bc;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
						{From: 5, To: 2, Weight: 1.0},
					},
					[]Expression{
						"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "<a>bc;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "<a>bc;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
						{From: 5, To: 2, Weight: 1.0},
					},
					[]Expression{
						"<SOS>", "<a>", "bc", ";", "<EOS>", "<SOS>", "123", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "a<b>c;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
					},
					[]Expression{
						"<SOS>", "a", "<b>", "c", ";", "<EOS>",
					},
				),
			},
			want: Rule{
				Exp: "a<b>c;", IsPublic: false, Graph: NewGraph(
					EdgeList{
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
					},
					[]Expression{
						"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "a<b>c;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
					},
					[]Expression{
						"<SOS>", "a", "<b>", "c", ";", "<EOS>",
					},
				),
			},
			want: Rule{
				Exp: "a<b>c;", IsPublic: true, Graph: NewGraph(
					EdgeList{
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
					},
					[]Expression{
						"<SOS>", "a", "<b>", "c", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "a<b><c><d><e>;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
						{From: 5, To: 6, Weight: 1.0},
						{From: 6, To: 7, Weight: 1.0},
					},
					[]Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "a<b><c><d><e>;", IsPublic: false,
				Graph: NewGraph(
					EdgeList{
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
					},
					[]Expression{
						"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";",
						"<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";",
						"<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "a<b><c><d><e>;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
						{From: 5, To: 6, Weight: 1.0},
						{From: 6, To: 7, Weight: 1.0},
					},
					[]Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "a<b><c><d><e>;", IsPublic: true,
				Graph: NewGraph(
					EdgeList{
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
					},
					[]Expression{
						"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";",
						"<EOS>", "<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>", "<SOS>", "1", "(", "2", ")", "3", ";",
						"<EOS>", "<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "a<b><b><b>;", IsPublic: false, Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
						{From: 5, To: 6, Weight: 1.0},
					},
					[]Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "a<b><b><b>;", IsPublic: false, Graph: NewGraph(
					EdgeList{
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
					},
					[]Expression{
						"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
						"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "a<b><b><b>;", IsPublic: true, Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
						{From: 5, To: 6, Weight: 1.0},
					},
					[]Expression{"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>"},
				),
			},
			want: Rule{
				Exp: "a<b><b><b>;", IsPublic: true, Graph: NewGraph(
					EdgeList{
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
					},
					[]Expression{
						"<SOS>", "a", "<b>", "<b>", "<b>", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
						"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>", "<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>",
					},
				),
			},
			wantErr: false,
		},
	}
	for i, test := range table {
		got, err := ResolveReferences(test.r, m, lexer)
		if test.want.IsPublic != got.IsPublic {
			t.Errorf("test %v: ResolveReferences(%v).Is_public\nGOT %v\nWANT %v", i, test.r, got.IsPublic, test.want.IsPublic)
		}
		if !slices.Equal(GetReferences(got), GetReferences(test.want)) {
			t.Errorf("test %v: ResolveReferences(%v).References\nGOT %v\nWANT %v", i, test.r, GetReferences(got), GetReferences(test.want))
		}
		if !slices.Equal(Sort(test.want.Graph.Edges), Sort(got.Graph.Edges)) {
			t.Errorf("test %v: ResolveReferences(%v).edges\nGOT %v\nWANT %v", i, test.r, Sort(got.Graph.Edges), Sort(test.want.Graph.Edges))
		}
		if !slices.Equal(test.want.Graph.Tokens, got.Graph.Tokens) {
			t.Errorf("test %v: ResolveReferences(%v).nodes\nGOT %v\nWANT %v", i, test.r, got.Graph.Tokens, test.want.Graph.Tokens)
		}
		if !slices.Equal(test.want.Tokens, got.Tokens) {
			t.Errorf("test %v: ResolveReferences(%v).Tokens\nGOT %v\nWANT %v", i, test.r, got.Tokens, test.want.Tokens)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ResolveReferences(%v).err\nGOT %v\nWANT %v", i, test.r, err, test.wantErr)
		}
	}
}

func TestParseRule(t *testing.T) {
	lexer := NewJSGFLexer()
	table := []struct {
		l       string
		n       string
		r       Rule
		wantErr bool
	}{
		{l: "", n: "", r: Rule{Exp: "", IsPublic: false}, wantErr: true},
		{l: ";", n: "", r: Rule{Exp: "", IsPublic: false}, wantErr: true},
		{l: " ", n: "", r: Rule{Exp: "", IsPublic: false}, wantErr: true},
		{l: "<rule> =", n: "", r: Rule{Exp: "", IsPublic: false}, wantErr: true},
		{l: "<rule> = ", n: "", r: Rule{Exp: "", IsPublic: false}, wantErr: true},
		{l: "public <rule> =", n: "", r: Rule{Exp: "", IsPublic: false}, wantErr: true},
		{l: "public <rule> = ", n: "", r: Rule{Exp: "", IsPublic: false}, wantErr: true},
		{l: "<rule> =;", n: "<rule>", r: Rule{Exp: ";", IsPublic: false}, wantErr: false},
		{l: "public <rule> =;", n: "<rule>", r: Rule{Exp: ";", IsPublic: true}, wantErr: false},
		{l: "<rule> = test expression 123;", n: "<rule>", r: Rule{Exp: "test expression 123;", IsPublic: false}, wantErr: false},
		{l: "<rule> = test \"expression\" 123;", n: "<rule>", r: Rule{Exp: "test \"expression\" 123;", IsPublic: false}, wantErr: false},
		{l: "public <rule> = test expression 123;", n: "<rule>", r: Rule{Exp: "test expression 123;", IsPublic: true}, wantErr: false},
		{l: "public <rule> = test \"expression\" 123;", n: "<rule>", r: Rule{Exp: "test \"expression\" 123;", IsPublic: true}, wantErr: false},
		{l: "<rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{Exp: "test expression 123 <rule> (abc) [def];", IsPublic: false}, wantErr: false},
		{l: "public <rule> = test expression 123 <rule> (abc) [def];", n: "<rule>", r: Rule{Exp: "test expression 123 <rule> (abc) [def];", IsPublic: true}, wantErr: false},
		{l: "<rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: false}, wantErr: false},
		{l: "public <rule> = test expression 123 <rule1> <rule2> (abc) [def];", n: "<rule>", r: Rule{Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: true}, wantErr: false},
	}
	for i, test := range table {
		n, r, err := ParseRule(test.l, lexer)
		if n != test.n {
			t.Errorf("test %v: ParseRule(%v, jsgfLexer)\nGOT %v\nWANT %v", i, test.l, n, test.n)
		}
		if r.Exp != test.r.Exp {
			t.Errorf("test %v: ParseRule(%v, jsgfLexer)\nGOT %v\nWANT %v", i, test.l, r, test.r)
		}
		if r.IsPublic != test.r.IsPublic {
			t.Errorf("test %v: ParseRule(%v, jsgfLexer).Is_public\nGOT %v\nWANT %v", i, test.l, r, test.r)
		}
		if !slices.Equal(GetReferences(r), GetReferences(test.r)) {
			t.Errorf("test %v: ParseRule(%v, jsgfLexer).References\nGOT %v\nWANT %v", i, test.l, GetReferences(r), GetReferences(test.r))
		}
		sort.Slice(r.Tokens, func(i, j int) bool { return r.Tokens[i] < r.Tokens[j] })
		sort.Slice(test.r.Tokens, func(i, j int) bool { return r.Tokens[i] < r.Tokens[j] })
		if !slices.EqualFunc(r.Tokens, test.r.Tokens, func(E1, E2 Expression) bool { return E1 == E2 }) {
			t.Errorf("test %v: ParseRule(%v, jsgfLexer).Tokens\nGOT %v\nWANT %v", i, test.l, r.Tokens, test.r.Tokens)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ParseRule(%v, jsgfLexer).err\nGOT %v\nWANT %v", i, test.l, err, test.wantErr)
		}
	}
}

func TestGetReferences(t *testing.T) {
	tests := []struct {
		r    Rule
		want []string
	}{
		{r: Rule{Exp: "", IsPublic: false}, want: []string{}},
		{r: Rule{Exp: "", IsPublic: true}, want: []string{}},
		{r: Rule{Exp: ";", IsPublic: false}, want: []string{}},
		{r: Rule{Exp: ";", IsPublic: true}, want: []string{}},
		{r: Rule{Exp: "test expression 123;", IsPublic: false}, want: []string{}},
		{r: Rule{Exp: "test \"expression\" 123;", IsPublic: true}, want: []string{}},
		{r: Rule{Exp: "test expression 123 <rule> (abc) [def];", IsPublic: false}, want: []string{"<rule>"}},
		{r: Rule{Exp: "test expression 123 <rule> (abc) [def];", IsPublic: true}, want: []string{"<rule>"}},
		{r: Rule{Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: false}, want: []string{"<rule1>", "<rule2>"}},
		{r: Rule{Exp: "test expression 123 <rule1> <rule2> (abc) [def];", IsPublic: true}, want: []string{"<rule1>", "<rule2>"}},
	}
	for i, test := range tests {
		got := GetReferences(test.r)
		slices.Sort(got)
		slices.Sort(test.want)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: GetReferences(%v).Peek()\nGOT %v\nWANT %v", i, test.r, got, test.want)
		}
	}
}

func TestValidateRuleRecursion(t *testing.T) {
	m := map[string]Rule{
		"<a>": {
			Exp:      "123<b>;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
				},
				[]Expression{"<SOS>", "123", ";", "<EOS>"},
			),
		},
		"<b>": {
			Exp:      "1|2|3;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				},
				[]Expression{"<SOS>", "1", "|", "2", "|", "3", ";", "<EOS>"},
			),
		},
		"<c>": {
			Exp:      "1[2]3<d>;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				},
				[]Expression{"<SOS>", "1", "[", "2", "]", "3", ";", "<EOS>"},
			),
		},
		"<d>": {
			Exp:      "1(2)3<e>;",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
				},
				[]Expression{"<SOS>", "1", "(", "2", ")", "3", ";", "<EOS>"},
			),
		},
		"<e>": {
			Exp:      "1(2[3<a>]);",
			IsPublic: false,
			Graph: NewGraph(
				EdgeList{
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
				},
				[]Expression{"<SOS>", "1", "(", "2", "[", "3", "]", ")", ";", "<EOS>"},
			),
		},
	}
	tests := []struct {
		n       string
		e       string
		wantErr bool
	}{
		{n: "", e: "", wantErr: false},
		{n: "<rule>", e: ";", wantErr: false},
		{n: "<rule>", e: "abc;", wantErr: false},
		{n: "<rule>", e: "<rule1>;", wantErr: false},
		{n: "<rule>", e: "<a> <b> <c>;", wantErr: false},
		{n: "<rule>", e: "<rule>;", wantErr: true},
		{n: "<rule>", e: "<rule><rule1>;", wantErr: true},
		{n: "<d>", e: "abc<c>;", wantErr: true},
		{n: "<b>", e: "123<e>;", wantErr: true},
		{n: "<e>", e: "<a><a>;", wantErr: true},
	}
	for i, test := range tests {
		err := ValidateRuleRecursion(test.n, NewRule(test.e, false), m)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateRuleRecursion(%v, %v, %v)\nGOT %v\nWANT %v", i, test.n, test.e, m, err, test.wantErr)
		}
	}
}
