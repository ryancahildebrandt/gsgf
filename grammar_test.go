// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:54 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"testing"
)

func TestGetCompositionOrder(t *testing.T) {
	table := []struct {
		g    Grammar
		want []string
	}{
		{
			g: Grammar{
				Rules:   map[string]Rule{},
				Imports: []string{},
			},
			want: []string{},
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("", false)},
				Imports: []string{},
			},
			want: []string{},
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("<b><c>", false)},
				Imports: []string{},
			},
			want: []string{},
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("", true)},
				Imports: []string{},
			},
			want: []string{"<a>"},
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("<b><c>", true)},
				Imports: []string{},
			},
			want: []string{"<a>", "<b>", "<c>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("", true),
					"<b>": NewRule("", true)},
				Imports: []string{},
			},
			want: []string{"<a>", "<b>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("", true),
					"<b>": NewRule("", false)},
				Imports: []string{},
			},
			want: []string{"<a>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("<b>", true),
					"<c>": NewRule("", true)},
				Imports: []string{},
			},
			want: []string{"<a>", "<b>", "<c>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("<b>", true),
					"<c>": NewRule("", false)},
				Imports: []string{},
			},
			want: []string{"<a>", "<b>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", true),
					"<d>": NewRule("", true),
				},
				Imports: []string{},
			},
			want: []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", true),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			want: []string{"<a>", "<b>", "<c>", "<c>", "<c>", "<d>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			want: []string{"<a>", "<b>", "<c>", "<c>", "<d>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", false),
					"<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			want: []string{"<a>", "<c>", "<d>"},
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", false),
					"<b>": NewRule("<c>", false),
					"<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			want: []string{},
		},
	}
	for i, test := range table {
		got := GetCompositionOrder(test.g)
		sort.Strings(got)
		sort.Strings(test.want)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.CompositionOrder()\nGOT %v\nWANT %v", i, test.g, got, test.want)
		}
	}
}

func TestGetAllProductions(t *testing.T) {
	lexer := NewJSGFLexer("\"")
	table := []struct {
		p       []string
		want    []string
		wantErr bool
	}{
		{
			p:       []string{""},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{";"},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{"", ""},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{";", ";", ";"},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{"abc;"},
			want:    []string{"abc"},
			wantErr: false},
		{
			p:       []string{"<b>;"},
			want:    []string{"1", "2", "3"},
			wantErr: false},
		{
			p:       []string{"<g>;"},
			want:    []string{"abc"},
			wantErr: false},
		{
			p:       []string{"<h>;"},
			want:    []string{"123bc"},
			wantErr: false},
		{
			p:       []string{"<i>;"},
			want:    []string{"a1c", "a2c", "a3c"},
			wantErr: false},
		{
			p: []string{"<j>;"},
			want: []string{
				"a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123",
				"a313123123", "a1123123123", "a2123123123", "a3123123123",
			},
			wantErr: false,
		},
		{
			p: []string{"<k>;"},
			want: []string{
				"a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221",
				"a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332",
				"a333",
			},
			wantErr: false,
		},
		{
			p:       []string{"a{}b//c/0.1/;"},
			want:    []string{"a{}b//c/0.1/"},
			wantErr: false},
		{
			p:       []string{"a{}/0.1/b/0.1/{}c;"},
			want:    []string{"a{}/0.1/b/0.1/{}c"},
			wantErr: false},
		{
			p:       []string{"abc<a>;"},
			want:    []string{"abc123"},
			wantErr: false},
		{
			p:       []string{"abc<l>;"},
			want:    []string{"abc<l>"},
			wantErr: true},
		{
			p:       []string{"abc<a><b><c>;"},
			want:    []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313"},
			wantErr: false,
		},
		{
			p:       []string{"abc(<g>|<h>|<i>);"},
			want:    []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c"},
			wantErr: false},
		{
			p:       []string{"abc<g><g><g>;"},
			want:    []string{"abcabcabcabc"},
			wantErr: false},
		{
			p:       []string{"abc;", "def;", "ghi;"},
			want:    []string{"abc", "def", "ghi"},
			wantErr: false},
		{
			p:       []string{"abc<a>;", "def<h>;", "<b>;"},
			want:    []string{"abc123", "def123bc", "1", "2", "3"},
			wantErr: false},
		{
			p: []string{"abc<a><b><c>;", "def[<h>];", "g((hi)|(jk));"},
			want: []string{
				"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313", "def", "def123bc", "ghi",
				"gjk",
			},
			wantErr: false,
		},
		{
			p: []string{"abc(<g>|<h>|<i>);", "def[<g>|<h>|<i>];", "ghi<b><b><b>;"},
			want: []string{
				"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c", "def", "defabc", "def123bc", "defa1c", "defa2c",
				"defa3c", "ghi111", "ghi112", "ghi113", "ghi121", "ghi122", "ghi123", "ghi131", "ghi132", "ghi133",
				"ghi211", "ghi212", "ghi213", "ghi221", "ghi222", "ghi223", "ghi231", "ghi232", "ghi233", "ghi311",
				"ghi312", "ghi313", "ghi321", "ghi322", "ghi323", "ghi331", "ghi332", "ghi333",
			},
			wantErr: false,
		},
		{
			p: []string{"abc<g><g><g>;", "<h><i>;", "<j>|<k>;"},
			want: []string{
				"abcabcabcabc", "123bca1c", "123bca2c", "123bca3c", "a11312312", "a21312312", "a31312312", "a112312312",
				"a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123",
				"a3123123123", "a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212",
				"a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323",
				"a331", "a332", "a333",
			},
			wantErr: false,
		},
	}
	for i, test := range table {
		g := NewGrammar()
		g.Rules = map[string]Rule{
			"<_>": {
				Exp:      "",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{},
					[]Expression{},
				),
			},
			"<a>": {
				Exp:      "123;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
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
			"<f>": {
				Exp:      "<l>;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<f>", ";", "<EOS>"},
				),
			},
			"<g>": {
				Exp:      "abc;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
					},
					[]Expression{"<SOS>", "abc", ";", "<EOS>"},
				),
			},
			"<h>": {
				Exp:      "<a>bc;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"},
				),
			},
			"<i>": {
				Exp:      "a<b>c;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
					},
					[]Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"},
				),
			},
			"<j>": {
				Exp:      "a<b><c><d><e>;",
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
					[]Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"},
				),
			},
			"<k>": {
				Exp:      "a<b><b><b>;",
				IsPublic: false,
				Graph: NewGraph(
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
		}
		for j, p := range test.p {
			rule := NewRule(p, true)
			rule.Tokens = ToTokens(rule.Exp, lexer)
			rule.Graph = NewGraph(ToEdgeList(rule.Tokens), rule.Tokens)
			g.Rules[fmt.Sprintf("<pub_%v>", j)] = rule
		}
		g, err := ResolveRules(g, lexer)
		got := GetAllProductions(g)
		sort.Strings(test.want)
		sort.Strings(got)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.Productions()\nGOT len %v %v\nWANT len %v %v", i, test.p, len(got), got, len(test.want), test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.Productions().err\nGOT %v\nWANT %v", i, test.p, err, test.wantErr)
		}
	}
}

func TestGetAllProductionsMinimized(t *testing.T) {
	lexer := NewJSGFLexer("\"")
	table := []struct {
		p       []string
		want    []string
		wantErr bool
	}{
		{
			p:       []string{""},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{";"},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{"", ""},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{";", ";", ";"},
			want:    []string{},
			wantErr: false},
		{
			p:       []string{"abc;"},
			want:    []string{"abc"},
			wantErr: false},
		{
			p:       []string{"<b>;"},
			want:    []string{"1", "2", "3"},
			wantErr: false},
		{
			p:       []string{"<g>;"},
			want:    []string{"abc"},
			wantErr: false},
		{
			p:       []string{"<h>;"},
			want:    []string{"123bc"},
			wantErr: false},
		{
			p:       []string{"<i>;"},
			want:    []string{"a1c", "a2c", "a3c"},
			wantErr: false},
		{
			p: []string{"<j>;"},
			want: []string{
				"a11312312", "a21312312", "a31312312", "a112312312", "a212312312", "a312312312", "a113123123", "a213123123",
				"a313123123", "a1123123123", "a2123123123", "a3123123123",
			},
			wantErr: false,
		},
		{
			p: []string{"<k>;"},
			want: []string{
				"a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212", "a213", "a221",
				"a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323", "a331", "a332",
				"a333",
			},
			wantErr: false,
		},
		{
			p:       []string{"a{}b//c/0.1/;"},
			want:    []string{"a{}b//c/0.1/"},
			wantErr: false},
		{
			p:       []string{"a{}/0.1/b/0.1/{}c;"},
			want:    []string{"a{}/0.1/b/0.1/{}c"},
			wantErr: false},
		{
			p:       []string{"abc<a>;"},
			want:    []string{"abc123"},
			wantErr: false},
		{
			p:       []string{"abc<l>;"},
			want:    []string{"abc<l>"},
			wantErr: true},
		{
			p:       []string{"abc<a><b><c>;"},
			want:    []string{"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313"},
			wantErr: false,
		},
		{
			p:       []string{"abc(<g>|<h>|<i>);"},
			want:    []string{"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c"},
			wantErr: false},
		{
			p:       []string{"abc<g><g><g>;"},
			want:    []string{"abcabcabcabc"},
			wantErr: false},
		{
			p:       []string{"abc;", "def;", "ghi;"},
			want:    []string{"abc", "def", "ghi"},
			wantErr: false},
		{
			p:       []string{"abc<a>;", "def<h>;", "<b>;"},
			want:    []string{"abc123", "def123bc", "1", "2", "3"},
			wantErr: false},
		{
			p: []string{"abc<a><b><c>;", "def[<h>];", "g((hi)|(jk));"},
			want: []string{
				"abc1231123", "abc123113", "abc1232123", "abc123213", "abc1233123", "abc123313", "def", "def123bc", "ghi",
				"gjk",
			},
			wantErr: false,
		},
		{
			p: []string{"abc(<g>|<h>|<i>);", "def[<g>|<h>|<i>];", "ghi<b><b><b>;"},
			want: []string{
				"abcabc", "abc123bc", "abca1c", "abca2c", "abca3c", "def", "defabc", "def123bc", "defa1c", "defa2c",
				"defa3c", "ghi111", "ghi112", "ghi113", "ghi121", "ghi122", "ghi123", "ghi131", "ghi132", "ghi133",
				"ghi211", "ghi212", "ghi213", "ghi221", "ghi222", "ghi223", "ghi231", "ghi232", "ghi233", "ghi311",
				"ghi312", "ghi313", "ghi321", "ghi322", "ghi323", "ghi331", "ghi332", "ghi333",
			},
			wantErr: false,
		},
		{
			p: []string{"abc<g><g><g>;", "<h><i>;", "<j>|<k>;"},
			want: []string{
				"abcabcabcabc", "123bca1c", "123bca2c", "123bca3c", "a11312312", "a21312312", "a31312312", "a112312312",
				"a212312312", "a312312312", "a113123123", "a213123123", "a313123123", "a1123123123", "a2123123123",
				"a3123123123", "a111", "a112", "a113", "a121", "a122", "a123", "a131", "a132", "a133", "a211", "a212",
				"a213", "a221", "a222", "a223", "a231", "a232", "a233", "a311", "a312", "a313", "a321", "a322", "a323",
				"a331", "a332", "a333",
			},
			wantErr: false,
		},
	}
	for i, test := range table {
		g := NewGrammar()
		g.Rules = map[string]Rule{
			"<_>": {
				Exp:      "",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{},
					[]Expression{},
				),
			},
			"<a>": {
				Exp:      "123;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
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
			"<f>": {
				Exp:      "<l>;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<f>", ";", "<EOS>"},
				),
			},
			"<g>": {
				Exp:      "abc;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
					[]Expression{"<SOS>", "abc", ";", "<EOS>"},
				),
			},
			"<h>": {
				Exp:      "<a>bc;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
					},
					[]Expression{"<SOS>", "<a>", "bc", ";", "<EOS>"},
				),
			},
			"<i>": {
				Exp:      "a<b>c;",
				IsPublic: false,
				Graph: NewGraph(
					EdgeList{
						{From: 0, To: 1, Weight: 1.0},
						{From: 1, To: 2, Weight: 1.0},
						{From: 2, To: 3, Weight: 1.0},
						{From: 3, To: 4, Weight: 1.0},
						{From: 4, To: 5, Weight: 1.0},
					},
					[]Expression{"<SOS>", "a", "<b>", "c", ";", "<EOS>"},
				),
			},
			"<j>": {
				Exp:      "a<b><c><d><e>;",
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
					[]Expression{"<SOS>", "a", "<b>", "<c>", "<d>", "<e>", ";", "<EOS>"},
				),
			},
			"<k>": {
				Exp:      "a<b><b><b>;",
				IsPublic: false,
				Graph: NewGraph(
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
		}
		for j, p := range test.p {
			rule := NewRule(p, true)
			rule.Tokens = ToTokens(rule.Exp, lexer)
			rule.Graph = NewGraph(ToEdgeList(rule.Tokens), rule.Tokens)
			rule.Graph = Minimize(rule.Graph)
			g.Rules[fmt.Sprintf("<pub_%v>", j)] = rule
		}
		g, err := ResolveRules(g, lexer)
		got := GetAllProductions(g)
		sort.Strings(test.want)
		sort.Strings(got)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.Productions()\nGOT len %v %v\nWANT len %v %v", i, test.p, len(got), got, len(test.want), test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.Productions().err\nGOT %v\nWANT %v", i, test.p, err, test.wantErr)
		}
	}
}

func TestGetAllProductionsE2E(t *testing.T) {
	lexer := NewJSGFLexer("\"")
	var productions []string
	f, err := os.Open("data/tests/productions.txt")
	if err != nil {
		t.Errorf("%s", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		productions = append(productions, scanner.Text())
	}
	table := []struct {
		p       string
		want    []string
		wantErr bool
	}{
		{
			p:       "data/tests/test0.jsgf",
			want:    productions,
			wantErr: false},
		{
			p:       "data/tests/test1.jsgf",
			want:    productions,
			wantErr: false},
		{
			p:       "data/tests/test2.jsgf",
			want:    productions,
			wantErr: true}, // returns error but is still resolvable, figure out how to validate completness before imports
		{
			p:       "data/tests/test3.jsgf",
			want:    productions,
			wantErr: false},
		{
			p:       "data/tests/test4.jsgf",
			want:    productions,
			wantErr: false},
		{
			p:       "data/tests/test5.jsgf",
			want:    productions,
			wantErr: false},
		{
			p:       "data/tests/a.jsgf",
			want:    []string{},
			wantErr: false},
		{
			p:       "data/tests/b.jsgf",
			want:    []string{},
			wantErr: false},
		{
			p:       "data/tests/dir0/c.jsgf",
			want:    []string{},
			wantErr: false},
		{
			p:       "data/tests/dir0/dir1/d.jsgf",
			want:    []string{},
			wantErr: true},
		{
			p:       "data/tests/dir0/dir1/dir2/e.jsgf",
			want:    productions,
			wantErr: false},
	}
	for i, test := range table {
		var err error
		grammar := NewGrammar()
		f, err1 := os.Open(test.p)
		scanner := bufio.NewScanner(f)
		grammar, err2 := FomJSGF(grammar, scanner, lexer)
		namespace, err3 := CreateNameSpace(test.p, ".jsgf")
		grammar = ImportNameSpace(grammar, namespace, lexer)
		grammar, err4 := ResolveRules(grammar, lexer)
		got := GetAllProductions(grammar)
		for _, e := range []error{err1, err2, err3, err4} {
			if e != nil {
				err = e
			}
		}
		sort.Strings(test.want)
		sort.Strings(got)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.Productions()\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.Productions().err\nGOT %v\nWANT %v", i, test.p, err, test.wantErr)
		}
	}
}

func TestValidateGrammarCompleteness(t *testing.T) {
	tests := []struct {
		g       Grammar
		wantErr bool
	}{
		{
			g: Grammar{
				Rules:   map[string]Rule{},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("", false)},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("<b><c>", false)},
				Imports: []string{},
			},
			wantErr: true,
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("", true)},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules:   map[string]Rule{"<a>": NewRule("<b><c>", true)},
				Imports: []string{},
			},
			wantErr: true,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("", true),
					"<b>": NewRule("", true)},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("", true),
					"<b>": NewRule("", false)},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("<b>", true),
					"<c>": NewRule("", true)},
				Imports: []string{},
			},
			wantErr: true,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{"<a>": NewRule("<b>", true),
					"<c>": NewRule("", false)},
				Imports: []string{},
			},
			wantErr: true,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", true),
					"<d>": NewRule("", true),
				},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", true),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", false),
					"<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			wantErr: false,
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", false),
					"<b>": NewRule("<c>", false),
					"<c>": NewRule("<d>", false),
					"<d>": NewRule("", false),
				},
				Imports: []string{},
			},
			wantErr: false,
		},
	}
	for i, test := range tests {
		if err := ValidateGrammarCompleteness(test.g); (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateGrammarCompleteness(%v).err\nGOT %v\nWANT %v", i, test.g, err, test.wantErr)
		}
	}
}
