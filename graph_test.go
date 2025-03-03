// -*- coding: utf-8 -*-

// Created on Thu Nov  7 08:51.08 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"testing"
)

func TestGetFrom(t *testing.T) {
	table := []struct {
		e    EdgeList
		n    int
		want []int
	}{
		{e: EdgeList{}, n: -1, want: []int{}},
		{e: EdgeList{}, n: 0, want: []int{}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, n: 2, want: []int{}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}}, n: 0, want: []int{1}},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			}, n: 6, want: []int{},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			}, n: 2, want: []int{3, 4},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			}, n: 5, want: []int{6},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			}, n: 2, want: []int{3, 5, 7, 8},
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		got := g.GetFrom(test.n)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: Graph(%v).GetFrom(%v)\nGOT  %v\nWANT %v", i, test.e, test.n, got, test.want)
		}
	}
}

func TestGetWeight(t *testing.T) {
	table := []struct {
		g    Graph
		f    int
		t    int
		want float64
	}{
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: 0, t: 0, want: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: -1, t: 0, want: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: 0, t: -1, want: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: -1, t: -1, want: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
				Children: map[int][]int{0: {1}, 1: {2}}, Weights: map[int]map[int]float64{0: {1: 1.0}, 1: {2: 1.0}},
			}, f: 0, t: 0, want: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				}, Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights: map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			}, f: 0, t: 1, want: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				}, Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights: map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			}, f: 0, t: 5, want: 0.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				}, Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights: map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			}, f: 1, t: 6, want: 100,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 0.99},
					{From: 0, To: 5, Weight: 0.0},
					{From: 1, To: 6, Weight: 100},
					{From: 3, To: 6, Weight: 99},
					{From: 5, To: 6, Weight: 9},
				}, Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}},
				Weights: map[int]map[int]float64{0: {1: 1.0, 3: 0.99, 5: 0.0}, 1: {6: 100}, 3: {6: 99}, 5: {6: 9}},
			}, f: 5, t: 6, want: 9,
		},
	}
	for i, test := range table {
		got := test.g.GetWeight(test.f, test.t)
		if test.want != got {
			t.Errorf("test %v: %v.GetWeight(%v, %v)\nGOT  %v\nWANT %v", i, test.g, test.f, test.t, got, test.want)
		}
	}
}

func TestAddEdge(t *testing.T) {
	table := []struct {
		e    EdgeList
		edg  Edge
		want Graph
	}{
		{
			e: EdgeList{}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			},
		},
		{
			e: EdgeList{}, edg: Edge{From: 0, To: 0, Weight: 1.0}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			},
		},
		{
			e: EdgeList{}, edg: Edge{From: 1, To: 10, Weight: 1.0}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{{From: 1, To: 10, Weight: 1.0}}, Children: map[int][]int{1: {10}},
				Weights: map[int]map[int]float64{1: {10: 1.0}},
			},
		},
		{
			e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{{From: 0, To: 1, Weight: 1.0}}, Children: map[int][]int{0: {1}},
				Weights: map[int]map[int]float64{0: {1: 1.0}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 1.0}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, Children: map[int][]int{0: {1}, 1: {2, 2}}, Weights: map[int]map[int]float64{0: {1: 1.0}, 1: {2: 1.0}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 1.0}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, Children: map[int][]int{0: {1}, 1: {2, 2}}, Weights: map[int]map[int]float64{0: {1: 0.99}, 1: {2: 1.0}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 0.99}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 1, To: 2, Weight: 0.99},
				}, Children: map[int][]int{0: {1}, 1: {2, 2}}, Weights: map[int]map[int]float64{0: {1: 1.0}, 1: {2: 0.99}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 0.99}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}, {From: 1, To: 2, Weight: 0.99},
				}, Children: map[int][]int{0: {1}, 1: {2, 2}}, Weights: map[int]map[int]float64{0: {1: 0.99}, 1: {2: 0.99}},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			}, edg: Edge{From: 8, To: 9, Weight: 1.0}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 0, To: 3, Weight: 1.0},
					{From: 0, To: 5, Weight: 1.0},
					{From: 1, To: 6, Weight: 1.0},
					{From: 3, To: 6, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 8, To: 9, Weight: 1.0},
				}, Children: map[int][]int{0: {1, 3, 5}, 1: {6}, 3: {6}, 5: {6}, 8: {9}}, Weights: map[int]map[int]float64{
					0: {1: 1.0, 3: 1.0, 5: 1.0}, 1: {6: 1.0}, 3: {6: 1.0}, 5: {6: 1.0}, 8: {9: 1.0},
				},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			}, edg: Edge{From: 0, To: 2, Weight: 1.0}, want: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0},
					{From: 1, To: 2, Weight: 1.0},
					{From: 2, To: 3, Weight: 1.0},
					{From: 2, To: 4, Weight: 1.0},
					{From: 3, To: 4, Weight: 1.0},
					{From: 4, To: 5, Weight: 1.0},
					{From: 5, To: 6, Weight: 1.0},
					{From: 5, To: 7, Weight: 1.0},
					{From: 6, To: 7, Weight: 1.0},
					{From: 7, To: 8, Weight: 1.0},
					{From: 0, To: 2, Weight: 1.0},
				}, Children: map[int][]int{0: {1, 2}, 1: {2}, 2: {3, 4}, 3: {4}, 4: {5}, 5: {6, 7}, 6: {7}, 7: {8}},
				Weights: map[int]map[int]float64{
					0: {1: 1.0, 2: 1.0}, 1: {2: 1.0}, 2: {3: 1.0, 4: 1.0}, 3: {4: 1.0}, 4: {5: 1.0}, 5: {6: 1.0, 7: 1.0},
					6: {7: 1.0}, 7: {8: 1.0},
				},
			},
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		got := g.AddEdge(test.edg)
		if !slices.Equal(got.Tokens, test.want.Tokens) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Nodes\nGOT  %v\nWANT %v", i, test.e, test.edg, got.Tokens, test.want.Tokens)
		}
		if !slices.Equal(Sort(got.Edges), Sort(test.want.Edges)) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Edges\nGOT  %v\nWANT %v", i, test.e, test.edg, got.Edges, test.want.Edges)
		}
		if !maps.EqualFunc(got.Children, test.want.Children, func(V1, V2 []int) bool { return slices.Equal(V1, V2) }) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT  %v\nWANT %v", i, test.e, test.edg, got.Children, test.want.Children)
		}
		if len(got.Weights) != len(test.want.Weights) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT  %v\nWANT %v", i, test.e, test.edg, got.Weights, test.want.Weights)
		}
		for k1, v1 := range got.Weights {
			v2, ok := test.want.Weights[k1]
			if !ok {
				t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT  %v\nWANT %v", i, test.e, test.edg, v1, v2)
			}
			if !maps.Equal(v1, v2) {
				t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT  %v\nWANT %v", i, test.e, test.edg, v1, v2)
			}
		}
	}
}

func TestGetEndPoints(t *testing.T) {
	table := []struct {
		e EdgeList
		i int
		f int
	}{
		{e: EdgeList{}, i: 0, f: 0},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, i: 0, f: 1},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}}, i: 0, f: 2},
		{
			e: EdgeList{
				{From: 10, To: 0, Weight: 1.0},
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 20, Weight: 1.0},
			}, i: 10, f: 20,
		},
		{
			e: EdgeList{
				{From: 1, To: 11, Weight: 1.0},
				{From: 11, To: 62, Weight: 1.0},
				{From: 62, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 2, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 4, To: 8, Weight: 1.0},
				{From: 8, To: 44, Weight: 1.0},
			}, i: 1, f: 44,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 2, To: 15, Weight: 1.0},
				{From: 2, To: 24, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 8, Weight: 1.0},
				{From: 15, To: 20, Weight: 1.0},
				{From: 17, To: 18, Weight: 1.0},
				{From: 18, To: 19, Weight: 1.0},
				{From: 19, To: 8, Weight: 1.0},
				{From: 20, To: 21, Weight: 1.0},
				{From: 21, To: 22, Weight: 1.0},
				{From: 22, To: 23, Weight: 1.0},
				{From: 23, To: 17, Weight: 1.0},
				{From: 24, To: 25, Weight: 1.0},
				{From: 25, To: 30, Weight: 1.0},
				{From: 27, To: 28, Weight: 1.0},
				{From: 28, To: 29, Weight: 1.0},
				{From: 29, To: 8, Weight: 1.0},
				{From: 30, To: 31, Weight: 1.0},
				{From: 30, To: 33, Weight: 1.0},
				{From: 30, To: 35, Weight: 1.0},
				{From: 31, To: 36, Weight: 1.0},
				{From: 33, To: 36, Weight: 1.0},
				{From: 35, To: 36, Weight: 1.0},
				{From: 36, To: 37, Weight: 1.0},
				{From: 37, To: 27, Weight: 1.0},
			}, i: 0, f: 10,
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		initial, final := GetEndPoints(g)
		if initial != test.i {
			t.Errorf("test %v: GetEndPoints(%v).initial\nGOT  %v\nWANT %v", i, test.e, initial, test.i)
		}
		if final != test.f {
			t.Errorf("test %v: GetEndPoints(%v).final\nGOT  %v\nWANT %v", i, test.e, final, test.f)
		}
	}
}

func TestGetAllPaths(t *testing.T) {
	table := []struct {
		e    EdgeList
		want []Path
	}{
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, want: []Path{{0, 1}}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}}, want: []Path{{0, 1, 2}}},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			}, want: []Path{{0, 1, 6}, {0, 3, 6}, {0, 5, 6}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8}, {0, 1, 2, 4, 5, 6, 7, 8}, {0, 1, 2, 4, 5, 7, 8}, {0, 1, 2, 3, 4, 5, 7, 8},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}, {0, 1, 2, 8, 9}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 4, 5, 14, 15}, {0, 1, 2, 7, 8, 9, 14, 15}, {0, 1, 2, 11, 12, 13, 14, 15}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 6, 7, 8, 9}, {0, 1, 2, 3, 7, 8, 9}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 3, 5, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 7, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 13, 14, 15},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 8, 9}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
			}, want: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, {0, 1, 2, 12, 13}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15}, {0, 1, 2, 14, 15}, {0, 1, 2, 7, 8, 9, 14, 15}, {0, 1, 2, 11, 12, 13, 14, 15},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 6, 7, 8, 9}, {0, 1, 2, 3, 7, 8, 9}, {0, 1, 2, 8, 9},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13},
				{0, 1, 2, 12, 13},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 3, 5, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 7, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 13, 14, 15},
				{0, 1, 2, 14, 15},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 13, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 10, To: 12, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 12, 13, 14, 15},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 11, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 12, To: 14, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 11, 12, 13, 14, 15},
				{0, 1, 2, 3, 4, 5, 11, 12, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 14, 15},
				{0, 1, 2, 7, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 9, 10, 11, 12, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 12, 14, 15},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 1, To: 12, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
			}, want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 12, 13, 14},
			},
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		got := GetAllPaths(g)
		sort.Slice(got, func(i, j int) bool { return fmt.Sprint(got[i]) < fmt.Sprint(got[j]) })
		sort.Slice(test.want, func(i, j int) bool { return fmt.Sprint(test.want[i]) < fmt.Sprint(test.want[j]) })
		for n := range got {
			if !slices.Equal(got[n], test.want[n]) {
				t.Errorf("test %v: GetAllPaths(%v)\nGOT  %v\nWANT %v", i, test.e, got[n], test.want[n])
			}
		}
	}
}

func TestComposeGraphs(t *testing.T) {
	table := []struct {
		g1      Graph
		g2      Graph
		i       int
		want    Graph
		wantErr bool
	}{
		{
			g1: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"a", "b"}),
			g2: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"c", "d"}),
			i:  0,
			want: NewGraph(EdgeList{{From: 2, To: 3, Weight: 1.0}, {From: 3, To: 1, Weight: 1.0}}, []Expression{
				"a", "b", "c", "d",
			}),
			wantErr: false,
		},
		{
			g1: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"c", "d"}),
			g2: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"a", "b"}),
			i:  1,
			want: NewGraph(EdgeList{{From: 0, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}}, []Expression{
				"c", "d", "a", "b",
			}),
			wantErr: false,
		},
		{
			g1: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 10, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
			}, []Expression{"a", "b", "c", "d"}),
			g2: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
			}, []Expression{"a", "b", "c", "d"}),
			i: 0,
			want: NewGraph(EdgeList{
				{From: 15, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 10, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			}, []Expression{"a", "b", "c", "d", "a", "b", "c", "d"}),
			wantErr: false,
		},
		{
			g1: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 10, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
			}, []Expression{""}),
			g2: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
			}, []Expression{""}),
			i: 5,
			want: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 11, Weight: 1.0},
				{From: 15, To: 10, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			}, []Expression{"", ""}),
			wantErr: false,
		},
		{
			g1: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 10, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
			}, []Expression{"a", "b", "c"}),
			g2: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
			}, []Expression{}),
			i: 10,
			want: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 11, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			}, []Expression{"a", "b", "c"}),
			wantErr: false,
		},
		{
			g1: NewGraph(EdgeList{}, []Expression{}), g2: NewGraph(EdgeList{}, []Expression{}), i: -1,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{}, []Expression{}), g2: NewGraph(EdgeList{}, []Expression{}), i: 0,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{}, []Expression{}), g2: NewGraph(EdgeList{}, []Expression{}), i: 2,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{}, []Expression{}),
			g2: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: -1,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{}, []Expression{}),
			g2: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: 0,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{}, []Expression{}),
			g2: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: 2,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			g2: NewGraph(EdgeList{}, []Expression{}), i: -1, want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			g2: NewGraph(EdgeList{}, []Expression{}), i: 0, want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			g2: NewGraph(EdgeList{}, []Expression{}), i: 2, want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			g2: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: -1,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
		{
			g1: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			g2: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: 2,
			want: NewGraph(EdgeList{}, []Expression{}), wantErr: true,
		},
	}
	for i, test := range table {
		got, err := ComposeGraphs(test.g1, test.g2, test.i)
		if !slices.Equal(got.Tokens, test.want.Tokens) {
			t.Errorf("test %v: ComposeGraphs(%v, %v, %v)\nGOT  %v\nWANT %v", i, test.g1, test.g2, test.i, got.Tokens, test.want.Tokens)
		}
		if !slices.Equal(Sort(got.Edges), Sort(test.want.Edges)) {
			t.Errorf("test %v: ComposeGraphs(%v, %v, %v)\nGOT  %v\nWANT %v", i, test.g1, test.g2, test.i, got.Edges, test.want.Edges)
		}
		if len(got.Children) != len(test.want.Children) {
			t.Errorf("test %v: ComposeGraphs(%v, %v, %v)\nGOT  %v\nWANT %v", i, test.g1, test.g2, test.i, got.Children, test.want.Children)
		}
		for k1, v1 := range got.Children {
			v2, ok := test.want.Children[k1]
			if !ok {
				t.Errorf("test %v: ComposeGraphs(%v, %v, %v)\nGOT  %v\nWANT %v", i, test.g1, test.g2, test.i, got.Children, test.want.Children)
			}
			sort.Ints(v1)
			sort.Ints(v2)
			if !slices.Equal(v1, v2) {
				t.Errorf("test %v: ComposeGraphs(%v, %v, %v)\nGOT  %v\nWANT %v", i, test.g1, test.g2, test.i, got.Children, test.want.Children)
			}
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ComposeGraphs(%v, %v, %v)\nGOT  %v\nWANT %v", i, test.g1, test.g2, test.i, err, test.wantErr)
		}
	}
}

func TestGetRandomChoice(t *testing.T) {
	table := []struct {
		c       []int
		w       []float64
		p       bool
		want    int
		wantErr bool
	}{
		{c: []int{}, w: []float64{}, p: false, want: -1, wantErr: true},
		{c: []int{0}, w: []float64{}, p: false, want: -1, wantErr: true},
		{c: []int{}, w: []float64{0.0}, p: false, want: -1, wantErr: true},
		{c: []int{0, 1, 2, 3}, w: []float64{0.0, 0.0, 0.0, 0.0}, p: false, want: -1, wantErr: true},
		{c: []int{0, 1, 2, 3}, w: []float64{0.0, 1.0, 0.0, 0.0}, p: false, want: 1, wantErr: false},
		{c: []int{0, 1, 2, 3}, w: []float64{0.0, 0.0, 0.0, 1.0}, p: false, want: 3, wantErr: false},
		{c: []int{10, 11, 12, 13}, w: []float64{10.0, 1.0, 100.0, 0.0}, p: true, want: 12, wantErr: false},
		{c: []int{10, 11, 12, 13}, w: []float64{0.001, 0.99, 0.01, 0.1}, p: true, want: 11, wantErr: false},
		{c: []int{10, 11, 12, 13}, w: []float64{1.0, 11.0, 111.0, 1111.0}, p: true, want: 13, wantErr: false},
	}
	for i, test := range table {
		var (
			c   int
			got int
			err error
		)
		if test.p {
			choices := make(map[int]float64)
			for range 1000 {
				c, err = GetRandomChoice(test.c, test.w)
				choices[c]++
			}
			got = test.c[0]
			for k, v := range choices {
				if v > choices[got] {
					got = k
				}
			}
		} else {
			got, err = GetRandomChoice(test.c, test.w)
		}
		if got != test.want {
			t.Errorf("test %v: ChooseNext(%v, %v)\nGOT  %v\nWANT %v", i, test.c, test.w, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ChooseNext(%v, %v)\nGOT  %v\nWANT %v", i, test.c, test.w, err, test.wantErr)
		}
	}
}

func TestGetRandomPath(t *testing.T) {
	table := []struct {
		e       EdgeList
		want    []Path
		wantErr bool
	}{
		{
			e:       EdgeList{{From: 0, To: 1, Weight: 1.0}},
			want:    []Path{{0, 1}},
			wantErr: false,
		},
		{
			e:       EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			want:    []Path{{0, 1, 2}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			},
			want:    []Path{{0, 1, 6}, {0, 3, 6}, {0, 5, 6}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8}, {0, 1, 2, 4, 5, 6, 7, 8}, {0, 1, 2, 4, 5, 7, 8}, {0, 1, 2, 3, 4, 5, 7, 8},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}, {0, 1, 2, 8, 9}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 4, 5, 14, 15}, {0, 1, 2, 7, 8, 9, 14, 15}, {0, 1, 2, 11, 12, 13, 14, 15}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 6, 7, 8, 9}, {0, 1, 2, 3, 7, 8, 9}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 3, 5, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 7, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 13, 14, 15},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 8, 9}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
			},
			want:    []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, {0, 1, 2, 12, 13}},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 6, 7, 8, 9}, {0, 1, 2, 3, 7, 8, 9}, {0, 1, 2, 8, 9},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13},
				{0, 1, 2, 12, 13},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 3, 5, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 7, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 13, 14, 15},
				{0, 1, 2, 14, 15},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 13, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 10, To: 12, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 12, 13, 14, 15},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 11, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 12, To: 14, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 11, 12, 13, 14, 15},
				{0, 1, 2, 3, 4, 5, 11, 12, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 14, 15},
				{0, 1, 2, 7, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 9, 10, 11, 12, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 12, 14, 15},
			},
			wantErr: false,
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 1, To: 12, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
			},
			want: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 12, 13, 14},
			},
			wantErr: false,
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		got, err := GetRandomPath(g)
		found := false
		for _, p := range test.want {
			if slices.Equal(got, p) {
				found = true
			}
		}
		if !found {
			t.Errorf("test %v: %v.RandomPath()\nGOT  %v\nWANT %v", i, test.e, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.RandomPath()\nGOT  %v\nWANT %v", i, test.e, err, test.wantErr)
		}
	}
}

func TestGraphDropNode(t *testing.T) {
	table := []struct {
		e    EdgeList
		i    int
		want EdgeList
	}{
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}},
			i:    0,
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}},
			i:    1,
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}},
			i:    2,
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			i:    1,
			want: EdgeList{{From: 0, To: 2, Weight: 1.0}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			},
			i: 5,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 0, To: 6, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
			i: 2,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 3, Weight: 1.0},
				{From: 1, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
			i: 3,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			i: 2,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 3, Weight: 1.0},
				{From: 1, To: 5, Weight: 1.0},
				{From: 1, To: 7, Weight: 1.0},
				{From: 1, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			i: 8,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 9, Weight: 1.0},
				{From: 5, To: 9, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			i: 4,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
			},
			i: 6,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 0, To: 7, Weight: 1.0},
				{From: 0, To: 9, Weight: 1.0},
				{From: 1, To: 10, Weight: 1.0},
				{From: 3, To: 10, Weight: 1.0},
				{From: 5, To: 10, Weight: 1.0},
				{From: 7, To: 10, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
			},
			i: 0,
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 0, To: 7, Weight: 1.0},
				{From: 0, To: 9, Weight: 1.0},
				{From: 1, To: 10, Weight: 1.0},
				{From: 3, To: 10, Weight: 1.0},
				{From: 5, To: 10, Weight: 1.0},
				{From: 7, To: 10, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
			},
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		got := g.DropNode(test.i)
		if !slices.Equal(Sort(got.Edges), Sort(test.want)) {
			t.Errorf("test %v: (%v).DropNode(%v)\nGOT  %v\nWANT %v", i, test.e, test.i, Sort(got.Edges), Sort(test.want))
		}
	}
}

func TestGraphMinimize(t *testing.T) {
	table := []struct {
		e    EdgeList
		n    []Expression
		want EdgeList
	}{
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}},
			n:    []Expression{"a", "b"},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			n:    []Expression{"a", "|", "b"},
			want: EdgeList{{From: 0, To: 2, Weight: 1.0}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			},
			n:    []Expression{"<SOS>", "(", "|", ")", "[", "]", "<EOS>"},
			want: EdgeList{{From: 0, To: 6, Weight: 1.0}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
			n: []Expression{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
			n: []Expression{"a", "b", "c", "d", "e", "f", "g", "h", ""},
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			n: []Expression{"", "", "", "d", "e", "f", "<EOS>", "h", "i", "j"},
			want: EdgeList{
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 0, To: 7, Weight: 1.0},
				{From: 0, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			n: []Expression{"a", "b", "c", "d", "|", "f", "|", "h", "i", ";"},
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			n: []Expression{"<SOS>", "<SOS>", "c", "d", "|", "f", "|", "h", ";", ";"},
			want: EdgeList{
				{From: 0, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
			},
			n: []Expression{"a", "(", "b", ")", "[", "c", "]", "d", "e", "|", "f", "g", "h", ";"},
			want: EdgeList{
				{From: 0, To: 2, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
			},
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, test.n)
		got := Minimize(g)
		if !slices.Equal(Sort(got.Edges), Sort(test.want)) {
			t.Errorf("test %v: (%v, %v).Minimize\nGOT  %v\nWANT %v", i, test.e, test.n, Sort(got.Edges), Sort(test.want))
		}
	}
}

func TestRuleWeightEdges(t *testing.T) {
	table := []struct {
		r       Rule
		want    Rule
		wantErr bool
	}{
		{
			r: Rule{
				Exp: "/.//", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			want: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			want: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "/.99/", ";", "<EOS>"}),
			},
			want: Rule{
				Exp: "/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", "", ";", "<EOS>"}),
			},
			wantErr: false,
		},
		{
			r: Rule{
				Exp: "123/.99/;", IsPublic: false, Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{
					"<SOS>", "123/.99/", ";", "<EOS>",
				}),
			},
			want: Rule{
				Exp: "123/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}),
			},
			wantErr: false,
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
			want: Rule{
				Exp: "<123>/.99/;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{
					"<SOS>", "<123>", ";", "<EOS>",
				}),
			},
			wantErr: false,
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
			want: Rule{
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
			wantErr: false,
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
			want: Rule{
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
			wantErr: false,
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
			want: Rule{
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
			wantErr: false,
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
			want: Rule{
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
			wantErr: false,
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
			want: Rule{
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
			wantErr: false,
		},
	}
	for i, test := range table {
		got, err := WeightEdges(test.r)
		if test.want.IsPublic != got.IsPublic {
			t.Errorf("test %v: %v.WeightEdges().Is_public\nGOT  %v\nWANT %v", i, test.r, got.IsPublic, test.want.IsPublic)
		}
		if !slices.Equal(GetReferences(got), GetReferences(test.want)) {
			t.Errorf("test %v: %v.WeightEdges().References\nGOT  %v\nWANT %v", i, test.r, GetReferences(got), GetReferences(test.want))
		}
		if !slices.Equal(Sort(test.want.Graph.Edges), Sort(got.Graph.Edges)) {
			t.Errorf("test %v: %v.WeightEdges().edges\nGOT  %v\nWANT %v", i, test.r, Sort(got.Graph.Edges), Sort(test.want.Graph.Edges))
		}
		if !slices.Equal(test.want.Graph.Tokens, got.Graph.Tokens) {
			t.Errorf("test %v: %v.WeightEdges().nodes\nGOT  %v\nWANT %v", i, test.r, got.Graph.Tokens, test.want.Graph.Tokens)
		}
		if !slices.Equal(test.want.Tokens, got.Tokens) {
			t.Errorf("test %v: %v.WeightEdges().Tokens\nGOT  %v\nWANT %v", i, test.r, got.Tokens, test.want.Tokens)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.WeightEdges().err\nGOT  %v\nWANT %v", i, test.r, err, test.wantErr)
		}
	}
}

func TestGetProductions(t *testing.T) {
	table := []struct {
		r    Rule
		want []string
	}{
		{
			r: Rule{
				Exp: "", IsPublic: false, Graph: NewGraph(EdgeList{}, []Expression{}),
			},
			want: []string{},
		},
		{
			r: Rule{
				Exp: ";", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, []Expression{"<SOS>", ";", "<EOS>"}),
			},
			want: []string{},
		},
		{
			r: Rule{
				Exp: "123;", IsPublic: false,
				Graph: NewGraph(EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0},
				}, []Expression{"<SOS>", "123", ";", "<EOS>"}),
			},
			want: []string{"123"},
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
			want: []string{"1", "2", "3"},
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
			want: []string{"1{}", "2//", "3/0.1/"},
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
			want: []string{"123", "13"},
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
			want: []string{"123"},
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
			want: []string{"12", "123"},
		},
	}
	for i, test := range table {
		got := GetProductions(test.r)
		sort.Strings(got)
		sort.Strings(test.want)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.Productions()\nGOT  %v\nWANT %v", i, test.r, got, test.want)
		}
	}
}

func TestFilterTokens(t *testing.T) {
	tests := []struct {
		e    []Expression
		f    []string
		want []Expression
	}{
		{
			e:    []Expression{},
			f:    []string{},
			want: []Expression{},
		},
		{
			e:    []Expression{"<SOS>", ";", "<EOS>"},
			f:    []string{""},
			want: []Expression{"<SOS>", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", " ", "<EOS>"},
			f:    []string{""},
			want: []Expression{"<SOS>", " ", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123", "<EOS>"},
			f:    []string{"<SOS>", " ", "<EOS>"},
			want: []Expression{"", "test expression 123", ""},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123", ";", "<EOS>"},
			f:    []string{"<SOS>", ";"},
			want: []Expression{"", "test expression 123", "", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", ";", "<EOS>"},
			f:    []string{"()", "abc", "<SOS>"},
			want: []Expression{"", "test expression 123 ", "(", "", ")", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "[", "abc", "]", ";", "<EOS>"},
			f:    []string{"<SOS>"},
			want: []Expression{"", "test expression 123 ", "[", "abc", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
			f:    []string{"a", "bc"},
			want: []Expression{"<SOS>", "test expression 123 ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"123"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "ab", "|", "c", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
			f:    []string{" ", "", "|"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab", "", "c", ")", "", "", "", "[", "de", "", "f", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"<rule>", "abc"},
			want: []Expression{"<SOS>", "test expression 123 ", "", " ", "(", "", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"<SOS>", "<rule>", "("},
			want: []Expression{"", "test expression 123 ", "", " ", "", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"<SOS>", "<rule>", "("},
			want: []Expression{"", "test expression 123 ", "", " ", "", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123// ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
			f:    []string{"<SOS>", "test expression 123//"},
			want: []Expression{"", "test expression 123// ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 // ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
			f:    []string{"<SOS>", "test expression 123//"},
			want: []Expression{"", "test expression 123 // ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 /0.0/", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
			f:    []string{"test expression 123 /0.0/"},
			want: []Expression{"<SOS>", "", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123{}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
			f:    []string{"test expression 123{}"},
			want: []Expression{"<SOS>", "", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 {}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
			f:    []string{"test expression 123{}", "test expression 123{}"},
			want: []Expression{"<SOS>", "test expression 123 {}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 {tag}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
			f:    []string{"test expression 123{}"},
			want: []Expression{"<SOS>", "test expression 123 {tag}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc /1.0/", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"test expression 123{}"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc /1.0/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc /1000/", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"/1000/", "abc"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc /1000/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc /-0.0001/", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"abc"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc /-0.0001/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc { }", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"abc {}"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc { }", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc {_}", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{"abc {_}"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "abc {_t_a_g_}", ")", " ", "[", "def", "]", ";", "<EOS>"},
			f:    []string{";"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc {_t_a_g_}", ")", " ", "[", "def", "]", "", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "ab/1.0/", "|", "c/1.0/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
			f:    []string{"ab/1.0/", "c"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "", "|", "c/1.0/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "ab/1000/", "|", "c/1000/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
			f:    []string{"ab/1.0/", "c"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab/1000/", "|", "c/1000/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "ab/-0.0001/", "|", "c/-0.0001/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
			f:    []string{"ab", "c"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab/-0.0001/", "|", "c/-0.0001/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "ab{1}", "|", "c{1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
			f:    []string{"|"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1}", "", "c{1}", ")", " ", "", " ", "[", "de", "", "f", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1}", "|", "c{1.1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
			f:    []string{"ab{1.1}"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "", "|", "c{1.1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1/1}", "|", "c{1.1/1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
			f:    []string{"ab{1.1/}"},
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1/1}", "|", "c{1.1/1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
	}
	for i, test := range tests {
		got := FilterTokens(test.e, test.f)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: FilterTerminals(%v, %v)\nGOT  %v\nWANT %v", i, test.e, test.f, got, test.want)
		}
	}
}
