// -*- coding: utf-8 -*-

// Created on Thu Nov  7 08:51.08 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"testing"
)

func TestGraphFrom(t *testing.T) {
	table := []struct {
		e   EdgeList
		n   int
		exp []int
	}{
		{e: EdgeList{}, n: -1, exp: []int{}},
		{e: EdgeList{}, n: 0, exp: []int{}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, n: 2, exp: []int{}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}}, n: 0, exp: []int{1}},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			}, n: 6, exp: []int{},
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
			}, n: 2, exp: []int{3, 4},
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
			}, n: 5, exp: []int{6},
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
			}, n: 2, exp: []int{3, 5, 7, 8},
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		res := g.GetFrom(test.n)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: Graph(%v).From(%v)\nGOT %v\nEXP %v", i, test.e, test.n, res, test.exp)
		}
	}
}

func TestGraphWeight(t *testing.T) {
	table := []struct {
		g   Graph
		f   int
		t   int
		exp float64
	}{
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: 0, t: 0, exp: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: -1, t: 0, exp: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: 0, t: -1, exp: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			}, f: -1, t: -1, exp: 1.0,
		},
		{
			g: Graph{
				Tokens: []Expression{}, Edges: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
				Children: map[int][]int{0: {1}, 1: {2}}, Weights: map[int]map[int]float64{0: {1: 1.0}, 1: {2: 1.0}},
			}, f: 0, t: 0, exp: 1.0,
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
			}, f: 0, t: 1, exp: 1.0,
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
			}, f: 0, t: 5, exp: 0.0,
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
			}, f: 1, t: 6, exp: 100,
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
			}, f: 5, t: 6, exp: 9,
		},
	}
	for i, test := range table {
		res := test.g.GetWeight(test.f, test.t)
		if test.exp != res {
			t.Errorf("test %v: %v.Weight(%v, %v)\nGOT %v\nEXP %v", i, test.g, test.f, test.t, res, test.exp)
		}
	}
}

func TestGraphAddEdge(t *testing.T) {
	table := []struct {
		e   EdgeList
		edg Edge
		exp Graph
	}{
		{
			e: EdgeList{}, exp: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			},
		},
		{
			e: EdgeList{}, edg: Edge{From: 0, To: 0, Weight: 1.0}, exp: Graph{
				Tokens: []Expression{}, Edges: EdgeList{}, Children: map[int][]int{}, Weights: map[int]map[int]float64{},
			},
		},
		{
			e: EdgeList{}, edg: Edge{From: 1, To: 10, Weight: 1.0}, exp: Graph{
				Tokens: []Expression{}, Edges: EdgeList{{From: 1, To: 10, Weight: 1.0}}, Children: map[int][]int{1: {10}},
				Weights: map[int]map[int]float64{1: {10: 1.0}},
			},
		},
		{
			e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, exp: Graph{
				Tokens: []Expression{}, Edges: EdgeList{{From: 0, To: 1, Weight: 1.0}}, Children: map[int][]int{0: {1}},
				Weights: map[int]map[int]float64{0: {1: 1.0}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 1.0}, exp: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, Children: map[int][]int{0: {1}, 1: {2, 2}}, Weights: map[int]map[int]float64{0: {1: 1.0}, 1: {2: 1.0}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 1.0}, exp: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0},
				}, Children: map[int][]int{0: {1}, 1: {2, 2}}, Weights: map[int]map[int]float64{0: {1: 0.99}, 1: {2: 1.0}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 0.99}, exp: Graph{
				Tokens: []Expression{}, Edges: EdgeList{
					{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 1, To: 2, Weight: 0.99},
				}, Children: map[int][]int{0: {1}, 1: {2, 2}}, Weights: map[int]map[int]float64{0: {1: 1.0}, 1: {2: 0.99}},
			},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 0.99}, {From: 1, To: 2, Weight: 1.0}},
			edg: Edge{From: 1, To: 2, Weight: 0.99}, exp: Graph{
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
			}, edg: Edge{From: 8, To: 9, Weight: 1.0}, exp: Graph{
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
			}, edg: Edge{From: 0, To: 2, Weight: 1.0}, exp: Graph{
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
		res := g.AddEdge(test.edg)
		if !slices.Equal(res.Tokens, test.exp.Tokens) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Nodes\nGOT %v\nEXP %v", i, test.e, test.edg, res.Tokens, test.exp.Tokens)
		}
		if !slices.Equal(Sort(res.Edges), Sort(test.exp.Edges)) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Edges\nGOT %v\nEXP %v", i, test.e, test.edg, res.Edges, test.exp.Edges)
		}
		if !maps.EqualFunc(res.Children, test.exp.Children, func(V1, V2 []int) bool { return slices.Equal(V1, V2) }) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT %v\nEXP %v", i, test.e, test.edg, res.Children, test.exp.Children)
		}
		if len(res.Weights) != len(test.exp.Weights) {
			t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT %v\nEXP %v", i, test.e, test.edg, res.Weights, test.exp.Weights)
		}
		for k1, v1 := range res.Weights {
			v2, ok := test.exp.Weights[k1]
			if !ok {
				t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT %v\nEXP %v", i, test.e, test.edg, v1, v2)
			}
			if !maps.Equal(v1, v2) {
				t.Errorf("test %v: Graph(%v).AddEdge(%v).Children\nGOT %v\nEXP %v", i, test.e, test.edg, v1, v2)
			}
		}
	}
}

func TestGraphEndPoints(t *testing.T) {
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
		if initial != test.i || final != test.f {
			t.Errorf("test %v: %v.EndPoints()\nGOT %v, %v\nEXP %v, %v", i, test.e, initial, final, test.i, test.f)
		}
	}
}

func TestGraphAllPaths(t *testing.T) {
	table := []struct {
		e   EdgeList
		exp []Path
	}{
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, exp: []Path{{0, 1}}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}}, exp: []Path{{0, 1, 2}}},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
			}, exp: []Path{{0, 1, 6}, {0, 3, 6}, {0, 5, 6}},
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
			}, exp: []Path{
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
			}, exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8}},
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
			}, exp: []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}, {0, 1, 2, 8, 9}},
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
			}, exp: []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}},
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
			}, exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
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
			}, exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
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
			}, exp: []Path{{0, 1, 2, 3, 4, 5, 14, 15}, {0, 1, 2, 7, 8, 9, 14, 15}, {0, 1, 2, 11, 12, 13, 14, 15}},
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
			}, exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 6, 7, 8, 9}, {0, 1, 2, 3, 7, 8, 9}},
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
			}, exp: []Path{
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
			}, exp: []Path{
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
			}, exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 8, 9}},
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
			}, exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, {0, 1, 2, 12, 13}},
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
			}, exp: []Path{
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
			}, exp: []Path{
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
			}, exp: []Path{
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
			}, exp: []Path{
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
			}, exp: []Path{
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
			}, exp: []Path{
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
			}, exp: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 12, 13, 14},
			},
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		res := GetAllPaths(g)
		sort.Slice(res, func(i, j int) bool { return fmt.Sprint(res[i]) < fmt.Sprint(res[j]) })
		sort.Slice(test.exp, func(i, j int) bool { return fmt.Sprint(test.exp[i]) < fmt.Sprint(test.exp[j]) })
		for n := range res {
			if !slices.Equal(res[n], test.exp[n]) {
				t.Errorf("test %v: %v.AllPaths()\nGOT %v\nEXP %v", i, test.e, res, test.exp)
			}
		}
	}
}

func TestGraphCompose(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		g   Graph
		h   Graph
		i   int
		exp Graph
		err error
	}{
		{
			g: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"a", "b"}),
			h: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"c", "d"}),
			i: 0,
			exp: NewGraph(EdgeList{{From: 2, To: 3, Weight: 1.0}, {From: 3, To: 1, Weight: 1.0}}, []Expression{
				"a", "b", "c", "d",
			}),
			err: nil,
		},
		{
			g: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"c", "d"}),
			h: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{"a", "b"}),
			i: 1,
			exp: NewGraph(EdgeList{{From: 0, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}}, []Expression{
				"c", "d", "a", "b",
			}),
			err: nil,
		},
		{
			g: NewGraph(EdgeList{
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
			h: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
			}, []Expression{"a", "b", "c", "d"}),
			i: 0,
			exp: NewGraph(EdgeList{
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
			err: nil,
		},
		{
			g: NewGraph(EdgeList{
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
			h: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
			}, []Expression{""}),
			i: 5,
			exp: NewGraph(EdgeList{
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
			err: nil,
		},
		{
			g: NewGraph(EdgeList{
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
			h: NewGraph(EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
			}, []Expression{}),
			i: 10,
			exp: NewGraph(EdgeList{
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
			err: nil,
		},
		{
			g: NewGraph(EdgeList{}, []Expression{}), h: NewGraph(EdgeList{}, []Expression{}), i: -1,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{}, []Expression{}), h: NewGraph(EdgeList{}, []Expression{}), i: 0,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{}, []Expression{}), h: NewGraph(EdgeList{}, []Expression{}), i: 2,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{}, []Expression{}),
			h: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: -1,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{}, []Expression{}),
			h: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: 0,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{}, []Expression{}),
			h: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: 2,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			h: NewGraph(EdgeList{}, []Expression{}), i: -1, exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			h: NewGraph(EdgeList{}, []Expression{}), i: 0, exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			h: NewGraph(EdgeList{}, []Expression{}), i: 2, exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			h: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: -1,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
		{
			g: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}),
			h: NewGraph(EdgeList{{From: 0, To: 1, Weight: 1.0}}, []Expression{}), i: 2,
			exp: NewGraph(EdgeList{}, []Expression{}), err: dummyError,
		},
	}
	for i, test := range table {
		res, err := ComposeGraphs(test.g, test.h, test.i)
		if !slices.Equal(res.Tokens, test.exp.Tokens) {
			t.Errorf("test %v: %v.Compose(%v, %v)\nGOT %v\nEXP %v", i, test.g, test.h, test.i, res.Tokens, test.exp.Tokens)
		}
		if !slices.Equal(Sort(res.Edges), Sort(test.exp.Edges)) {
			t.Errorf("test %v: %v.Compose(%v, %v)\nGOT %v\nEXP %v", i, test.g, test.h, test.i, res.Edges, test.exp.Edges)
		}
		if len(res.Children) != len(test.exp.Children) {
			t.Errorf("test %v: %v.Compose(%v, %v)\nGOT %v\nEXP %v", i, test.g, test.h, test.i, res.Children, test.exp.Children)
		}
		for k1, v1 := range res.Children {
			v2, ok := test.exp.Children[k1]
			if !ok {
				t.Errorf("test %v: %v.Compose(%v, %v)\nGOT %v\nEXP %v", i, test.g, test.h, test.i, res.Children, test.exp.Children)
			}
			sort.Ints(v1)
			sort.Ints(v2)
			if !slices.Equal(v1, v2) {
				t.Errorf("test %v: %v.Compose(%v, %v)\nGOT %v\nEXP %v", i, test.g, test.h, test.i, res.Children, test.exp.Children)
			}
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.Compose(%v, %v)\nGOT %v\nEXP %v", i, test.g, test.h, test.i, err, test.err)
		}
	}
}

func TestChooseNext(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		c   []int
		w   []float64
		p   bool
		exp int
		err error
	}{
		{c: []int{}, w: []float64{}, p: false, exp: -1, err: dummyError},
		{c: []int{0}, w: []float64{}, p: false, exp: -1, err: dummyError},
		{c: []int{}, w: []float64{0.0}, p: false, exp: -1, err: dummyError},
		{c: []int{0, 1, 2, 3}, w: []float64{0.0, 0.0, 0.0, 0.0}, p: false, exp: -1, err: dummyError},
		{c: []int{0, 1, 2, 3}, w: []float64{0.0, 1.0, 0.0, 0.0}, p: false, exp: 1, err: nil},
		{c: []int{0, 1, 2, 3}, w: []float64{0.0, 0.0, 0.0, 1.0}, p: false, exp: 3, err: nil},
		{c: []int{10, 11, 12, 13}, w: []float64{10.0, 1.0, 100.0, 0.0}, p: true, exp: 12, err: nil},
		{c: []int{10, 11, 12, 13}, w: []float64{0.001, 0.99, 0.01, 0.1}, p: true, exp: 11, err: nil},
		{c: []int{10, 11, 12, 13}, w: []float64{1.0, 11.0, 111.0, 1111.0}, p: true, exp: 13, err: nil},
	}
	for i, test := range table {
		var err error
		var res int
		if test.p {
			choices := make(map[int]float64)
			for range 1000 {
				c, _ := GetRandomChoice(test.c, test.w)
				choices[c]++
			}
			res = test.c[0]
			for k, v := range choices {
				if v > choices[res] {
					res = k
				}
			}
		} else {
			res, err = GetRandomChoice(test.c, test.w)
		}
		if res != test.exp {
			t.Errorf("test %v: ChooseNext(%v, %v)\nGOT %v\nEXP %v", i, test.c, test.w, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ChooseNext(%v, %v)\nGOT %v\nEXP %v", i, test.c, test.w, err, test.err)
		}
	}
}

func TestGraphRandomPath(t *testing.T) {
	table := []struct {
		e   EdgeList
		exp []Path
		err error
	}{
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}},
			exp: []Path{{0, 1}},
			err: nil,
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			exp: []Path{{0, 1, 2}},
			err: nil,
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
			exp: []Path{{0, 1, 6}, {0, 3, 6}, {0, 5, 6}},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8}, {0, 1, 2, 4, 5, 6, 7, 8}, {0, 1, 2, 4, 5, 7, 8}, {0, 1, 2, 3, 4, 5, 7, 8},
			},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8}},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}, {0, 1, 2, 8, 9}},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 8, 9}, {0, 1, 2, 5, 8, 9}, {0, 1, 2, 7, 8, 9}},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 4, 5, 14, 15}, {0, 1, 2, 7, 8, 9, 14, 15}, {0, 1, 2, 11, 12, 13, 14, 15}},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 6, 7, 8, 9}, {0, 1, 2, 3, 7, 8, 9}},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 7, 8, 9, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 10, 11, 12, 13},
				{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13},
			},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 3, 5, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 7, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 13, 14, 15},
			},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 8, 9}},
			err: nil,
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
			exp: []Path{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, {0, 1, 2, 12, 13}},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
			},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 6, 7, 8, 9}, {0, 1, 2, 3, 7, 8, 9}, {0, 1, 2, 8, 9},
			},
			err: nil,
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
			exp: []Path{
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
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 14, 15},
				{0, 1, 2, 3, 5, 14, 15},
				{0, 1, 2, 7, 8, 9, 14, 15},
				{0, 1, 2, 7, 9, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 13, 14, 15},
				{0, 1, 2, 14, 15},
			},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 12, 13, 14, 15},
			},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 11, 12, 13, 14, 15},
				{0, 1, 2, 3, 4, 5, 11, 12, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 8, 9, 10, 11, 12, 14, 15},
				{0, 1, 2, 7, 9, 10, 11, 12, 13, 14, 15},
				{0, 1, 2, 7, 9, 10, 11, 12, 14, 15},
				{0, 1, 2, 11, 12, 13, 14, 15},
				{0, 1, 2, 11, 12, 14, 15},
			},
			err: nil,
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
			exp: []Path{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 12, 13, 14},
			},
			err: nil,
		},
	}
	for i, test := range table {
		g := NewGraph(test.e, []Expression{})
		res, err := GetRandomPath(g)
		found := false
		for _, p := range test.exp {
			if slices.Equal(res, p) {
				found = true
			}
		}
		if !found {
			t.Errorf("test %v: %v.RandomPath()\nGOT %v\nEXP %v", i, test.e, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.RandomPath()\nGOT %v\nEXP %v", i, test.e, err, test.err)
		}
	}
}

func TestGraphDropNode(t *testing.T) {
	table := []struct {
		e   EdgeList
		i   int
		exp EdgeList
	}{
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}},
			i:   0,
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}},
			i:   1,
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}},
			i:   2,
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			i:   1,
			exp: EdgeList{{From: 0, To: 2, Weight: 1.0}},
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
		res := g.DropNode(test.i)
		if !slices.Equal(Sort(res.Edges), Sort(test.exp)) {
			t.Errorf("test %v: (%v).DropNode(%v)\nGOT %v\nEXP %v", i, test.e, test.i, Sort(res.Edges), Sort(test.exp))
		}
	}
}

func TestGraphMinimize(t *testing.T) {
	table := []struct {
		e   EdgeList
		n   []Expression
		exp EdgeList
	}{
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}},
			n:   []Expression{"a", "b"},
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			n:   []Expression{"a", "|", "b"},
			exp: EdgeList{{From: 0, To: 2, Weight: 1.0}},
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
			n:   []Expression{"<SOS>", "(", "|", ")", "[", "]", "<EOS>"},
			exp: EdgeList{{From: 0, To: 6, Weight: 1.0}},
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
			exp: EdgeList{
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
		res := Minimize(g)
		if !slices.Equal(Sort(res.Edges), Sort(test.exp)) {
			t.Errorf("test %v: (%v, %v).Minimize\nGOT %v\nEXP %v", i, test.e, test.n, Sort(res.Edges), Sort(test.exp))
		}
	}
}
