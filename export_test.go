// -*- coding: utf-8 -*-// Created on Sat Feb 22 10:19:34 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt
package main

import (
	"testing"
)

func TestGraphToTXT(t *testing.T) {
	tests := []struct {
		g     Graph
		want1 string
		want2 string
	}{
		{
			g: Graph{
				Tokens:   []Expression{},
				Edges:    EdgeList{},
				Children: map[int][]int{},
				Weights:  map[int]map[int]float64{},
			},
			want1: "",
			want2: "",
		},
		{
			g: Graph{
				Tokens:   []Expression{"a", "b", "c"},
				Edges:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
				Children: map[int][]int{0: {1}, 1: {2}},
				Weights:  map[int]map[int]float64{0: {1: 1.0}, 1: {2: 1.0}},
			},
			want1: "\"a\"\n\"b\"\n\"c\"",
			want2: "0,1,1\n1,2,1",
		},
		{
			g: Graph{
				Tokens: []Expression{"a", "b", "c", "", "e"},
				Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				},
				Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights:  map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			},
			want1: "\"a\"\n\"b\"\n\"c\"\n\"\"\n\"e\"",
			want2: "0,1,1\n0,3,0.99\n0,5,0\n1,6,100\n3,6,99\n5,6,9",
		},
		{
			g: Graph{
				Tokens: []Expression{"a", "b", "c", "d", "e"},
				Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				},
				Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights:  map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			},
			want1: "\"a\"\n\"b\"\n\"c\"\n\"d\"\n\"e\"",
			want2: "0,1,1\n0,3,0.99\n0,5,0\n1,6,100\n3,6,99\n5,6,9",
		},
	}
	for i, test := range tests {
		got, got1 := GraphToTXT(test.g)
		if got != test.want1 {
			t.Errorf("test %v: GraphToTXT(%v)\nGOT %v\nWANT %v", i, test.g, got, test.want1)
		}
		if got1 != test.want2 {
			t.Errorf("test %v: GraphToTXT(%v)\nGOT %v\nWANT %v", i, test.g, got1, test.want2)
		}
	}
}

func TestGraphToDOT(t *testing.T) {
	tests := []struct {
		g    Graph
		want string
	}{
		{
			g: Graph{
				Tokens:   []Expression{},
				Edges:    EdgeList{},
				Children: map[int][]int{},
				Weights:  map[int]map[int]float64{},
			},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\n\n}",
		},
		{
			g: Graph{
				Tokens:   []Expression{"a", "b", "c"},
				Edges:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
				Children: map[int][]int{0: {1}, 1: {2}},
				Weights:  map[int]map[int]float64{0: {1: 1.0}, 1: {2: 1.0}},
			},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\t_0 [label=\"a\"];\n\t_1 [label=\"b\"];\n\t_2 [label=\"c\"];\n\n\t_0 -> _1 [weight=1];\n\t_1 -> _2 [weight=1];\n\n}",
		},
		{
			g: Graph{
				Tokens: []Expression{"a", "b", "c", "", "e"},
				Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				},
				Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights:  map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\t_0 [label=\"a\"];\n\t_1 [label=\"b\"];\n\t_3 [label=\"\"];\n\n\t_0 -> _1 [weight=1];\n\t_0 -> _3 [label=\"0.99\",weight=0.99];\n\t_0 -> _5 [label=\"0\",weight=0];\n\t_1 -> _6 [label=\"100\",weight=100];\n\t_3 -> _6 [label=\"99\",weight=99];\n\t_5 -> _6 [label=\"9\",weight=9];\n\n}",
		},
		{
			g: Graph{
				Tokens: []Expression{"a", "b", "c", "d", "e"},
				Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				},
				Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights:  map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\t_0 [label=\"a\"];\n\t_1 [label=\"b\"];\n\t_3 [label=\"d\"];\n\n\t_0 -> _1 [weight=1];\n\t_0 -> _3 [label=\"0.99\",weight=0.99];\n\t_0 -> _5 [label=\"0\",weight=0];\n\t_1 -> _6 [label=\"100\",weight=100];\n\t_3 -> _6 [label=\"99\",weight=99];\n\t_5 -> _6 [label=\"9\",weight=9];\n\n}",
		},
	}
	for i, test := range tests {
		got := GraphToDOT(test.g)
		if got != test.want {
			t.Errorf("test %v: GraphToDOT(%v)\nGOT %v\nWANT %v", i, test.g, got, test.want)
		}
	}
}

func TestReferencesToDOT(t *testing.T) {
	tests := []struct {
		g    Grammar
		want string
	}{
		{
			g:    Grammar{Rules: map[string]Rule{}, Imports: []string{}},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\n}",
		},
		{
			g:    Grammar{Rules: map[string]Rule{"<a>": NewRule("", false)}, Imports: []string{}},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\n}",
		},
		{
			g:    Grammar{Rules: map[string]Rule{"<a>": NewRule("<b><c>", false)}, Imports: []string{}},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\t<b> -> <a>;\n\t<c> -> <a>;\n\n}",
		},
		{
			g: Grammar{Rules: map[string]Rule{"<a>": NewRule("<b>", true),
				"<c>": NewRule("", false)}, Imports: []string{}},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\t<b> -> <a>;\n\n}",
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", true),
					"<d>": NewRule("", true),
				}, Imports: []string{}},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\t<c> -> <a>;\n\t<c> -> <b>;\n\t<d> -> <c>;\n\n}",
		},
		{
			g: Grammar{
				Rules: map[string]Rule{
					"<a>": NewRule("<c>", true),
					"<b>": NewRule("<c>", true),
					"<c>": NewRule("<d>", true),
					"<d>": NewRule("", false),
				}, Imports: []string{}},
			want: "digraph {\n\n\trankdir = \"LR\"\n\n\t<c> -> <a>;\n\t<c> -> <b>;\n\t<d> -> <c>;\n\n}",
		},
	}
	for i, test := range tests {
		got := ReferencesToDOT(test.g)
		if got != test.want {
			t.Errorf("test %v: ReferencesToDOT(%v)\nGOT %v\nWANT %v", i, test.g, got, test.want)
		}
	}
}
