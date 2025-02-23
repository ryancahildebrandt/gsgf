// -*- coding: utf-8 -*-

// Created on Thu Sep  5 07:38:44 PM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"

	"gonum.org/v1/gonum/stat/sampleuv"
)

type Graph struct {
	Nodes    []Expression
	Edges    EdgeList
	Children map[int][]int
	Weights  map[int]map[int]float64
}

func NewGraph(e EdgeList, n []Expression) Graph {
	g := Graph{}
	g.Children = make(map[int][]int)
	g.Weights = make(map[int]map[int]float64)
	for _, edge := range e {
		g = g.AddEdge(edge)
	}
	g.Nodes = n
	return g
}

func (g Graph) Copy() Graph {
	e := g.Edges.Copy()
	n := make([]Expression, len(g.Nodes))
	copy(n, g.Nodes)
	return NewGraph(e, n)
}

func (g Graph) From(i int) []int {
	v, ok := g.Children[i]
	if !ok {
		return []int{}
	}
	return v
}

func (g Graph) Weight(f int, t int) float64 {
	v, ok := g.Weights[f][t]
	if !ok {
		return 1.0
	}
	return v
}

func (g Graph) AddEdge(e Edge) Graph {
	if e.IsEmpty() {
		return g
	}
	g.Edges = append(g.Edges, e.Copy())
	g.Children[e.from] = append(g.Children[e.from], e.to)
	_, ok := g.Weights[e.from]
	if !ok {
		g.Weights[e.from] = make(map[int]float64)
	}
	g.Weights[e.from][e.to] = e.weight
	return g
}

func (g Graph) EndPoints() (i, f int) {
	var e1 = make(map[int]struct{})
	var e2 = make(map[int]struct{})
	edges := g.Edges.Copy().Sort()
	for _, edge := range edges {
		e1[edge.from] = struct{}{}
		e2[edge.to] = struct{}{}
	}
	for _, edge := range edges {
		_, ok := e2[edge.from]
		if !ok {
			i = edge.from
		}
		_, ok = e1[edge.to]
		if !ok {
			f = edge.to
		}
	}
	return i, f
}

type Path []int

func (g Graph) AllPaths() (res []Path) {
	var path Path
	var paths []Path
	var p Path
	var node int
	var f int
	var t int

	f, t = g.EndPoints()
	paths = []Path{{f}}
	for len(paths) > 0 {
		path, paths = paths[0], paths[1:]
		node = path[len(path)-1]
		if node == t {
			res = append(res, path)
			continue
		}
		for _, n := range g.From(node) {
			p = make(Path, len(path)+1)
			copy(p, path)
			p[len(path)] = n
			paths = append(paths, p)
		}
	}
	return res
}

func ComposeGraphs(g Graph, h Graph, i int) (Graph, error) {
	switch {
	case g.Edges.IsEmpty() || h.Edges.IsEmpty():
		return Graph{}, errors.New("one or more EdgeLists e and a are empty")
	case i < 0:
		return Graph{}, errors.New("cannot insert EdgeList a at negative index")
	case i > g.Edges.Max():
		return Graph{}, errors.New("cannot insert EdgeList a at index greater than EdgeList g.Max()")
	}

	h.Edges = h.Edges.Increment(g.Edges.Max() + 1)
	h_from, h_to := h.EndPoints()
	exp := append(g.Nodes, h.Nodes...)
	edg := h.Edges
	for _, edge := range g.Edges {
		e := edge.Copy()
		if edge.from == i {
			e.from = h_to
		}
		if edge.to == i {
			e.to = h_from
		}
		edg = append(edg, e)
	}
	return NewGraph(edg, exp), nil
}

func ChooseNext(c []int, w []float64) (int, error) {
	if len(c) == 0 || len(w) == 0 {
		return -1, errors.New("length of choices c and/or weights w is 0")
	}
	if len(c) != len(w) {
		return -1, errors.New("length of choices c and weights w do not match")
	}

	i, ok := sampleuv.NewWeighted(w, nil).Take()
	if !ok {
		return -1, errors.New("sampleuv.NewWeighted could not sample from choices c and weights w")
	}
	return c[i], nil
}

func (g Graph) RandomPath() (Path, error) {
	var res Path
	var choice int
	f, t := g.EndPoints()
	res = append(res, f)
	p := f
	for p != t {
		n := g.From(p)
		switch len(n) {
		case 0:
			return Path{}, errors.New("cannot proceed further down path")

		case 1:
			choice = n[0]
			res = append(res, choice)
			p = choice

		default:
			w := make([]float64, len(n))
			for i, dest := range n {
				w[i] = g.Weight(f, dest)
			}
			choice, err := ChooseNext(g.From(p), w)
			if err != nil {
				return Path{}, err
			}
			res = append(res, choice)
			p = choice
		}
	}
	return res, nil
}

func (g Graph) Minimize() Graph {
	var g1 Graph
	f := []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>", ""}
	g1 = g
	for i, t := range g1.Nodes {
		if slices.Contains(f, t.str()) {
			g1 = g1.DropNode(i)
		}
	}
	return g1
}

func (g Graph) DropNode(i int) Graph {
	var from []int
	var to []int
	var edg EdgeList
	var start int
	var end int

	start, end = g.EndPoints()
	for _, edge := range g.Edges {
		switch i {
		case start, end:
			edg = append(edg, edge)
		case edge.from:
			to = append(to, edge.to)
		case edge.to:
			from = append(from, edge.from)
		default:
			edg = append(edg, edge)
		}
	}

	for _, f := range from {
		for _, t := range to {
			edg = append(edg, Edge{f, t, 1.0})
		}
	}
	return NewGraph(edg.Unique(), g.Nodes)
}
