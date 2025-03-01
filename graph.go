// -*- coding: utf-8 -*-

// Created on Thu Sep  5 07:38:44 PM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"
	"strings"

	"gonum.org/v1/gonum/stat/sampleuv"
)

type Graph struct {
	Tokens   []Expression
	Edges    EdgeList
	Children map[int][]int
	Weights  map[int]map[int]float64
}

func NewGraph(e EdgeList, n []Expression) Graph {
	g := Graph{}
	g.Tokens = n
	g.Children = make(map[int][]int)
	g.Weights = make(map[int]map[int]float64)
	for _, edge := range e {
		g = g.AddEdge(edge)
	}

	return g
}

func (g Graph) GetFrom(i int) []int {
	v, ok := g.Children[i]
	if !ok {
		return []int{}
	}

	return v
}

func (g Graph) GetWeight(f int, t int) float64 {
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
	g.Edges = append(g.Edges, e)
	g.Children[e.From] = append(g.Children[e.From], e.To)
	_, ok := g.Weights[e.From]
	if !ok {
		g.Weights[e.From] = make(map[int]float64)
	}
	g.Weights[e.From][e.To] = e.Weight

	return g
}

func (g Graph) DropNode(i int) Graph {
	var (
		from  []int
		to    []int
		edg   EdgeList
		start int
		end   int
	)

	start, end = GetEndPoints(g)

	for _, edge := range g.Edges {
		switch i {
		case start, end:
			edg = append(edg, edge)
		case edge.From:
			to = append(to, edge.To)
		case edge.To:
			from = append(from, edge.From)
		default:
			edg = append(edg, edge)
		}
	}

	for _, f := range from {
		for _, t := range to {
			edg = append(edg, Edge{From: f, To: t, Weight: 1.0})
		}
	}

	return NewGraph(Unique(edg), g.Tokens)
}

func GetEndPoints(g Graph) (i, f int) {
	e1 := make(map[int]struct{})
	e2 := make(map[int]struct{})
	edges := Sort(g.Edges)
	for _, edge := range edges {
		e1[edge.From] = struct{}{}
		e2[edge.To] = struct{}{}
	}
	for _, edge := range edges {
		_, ok := e2[edge.From]
		if !ok {
			i = edge.From
		}
		_, ok = e1[edge.To]
		if !ok {
			f = edge.To
		}
	}

	return i, f
}

type Path = []int

func GetAllPaths(g Graph) []Path {
	var (
		f, t  int    = GetEndPoints(g)
		paths []Path = []Path{{f}}
		path  Path
		res   []Path
		p     Path
		node  int
	)

	for len(paths) > 0 {
		path, paths = paths[0], paths[1:]
		node = path[len(path)-1]
		if node == t {
			res = append(res, path)

			continue
		}
		for _, n := range g.GetFrom(node) {
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

	h.Edges = Increment(h.Edges, g.Edges.Max()+1)
	hFrom, hTo := GetEndPoints(h)
	exp := append(g.Tokens, h.Tokens...)
	edg := h.Edges

	for _, edge := range g.Edges {
		e := edge
		if edge.From == i {
			e.From = hTo
		}
		if edge.To == i {
			e.To = hFrom
		}
		edg = append(edg, e)
	}

	return NewGraph(edg, exp), nil
}

func GetRandomChoice(c []int, w []float64) (int, error) {
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

func GetRandomPath(g Graph) (Path, error) {
	var (
		res    Path
		choice int
		f, t   int = GetEndPoints(g)
	)

	res = append(res, f)
	p := f
	for p != t {
		n := g.GetFrom(p)
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
				w[i] = g.GetWeight(f, dest)
			}
			choice, err := GetRandomChoice(g.GetFrom(p), w)
			if err != nil {
				return Path{}, err
			}
			res = append(res, choice)
			p = choice
		}
	}

	return res, nil
}

func Minimize(g Graph) Graph {
	var g1 Graph = g
	var f []string = []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>", ""}

	for i, t := range g1.Tokens {
		if slices.Contains(f, t) {
			g1 = g1.DropNode(i)
		}
	}

	return g1
}

func WeightEdges(r Rule) (Rule, error) {
	for i, t := range r.Tokens {
		if IsWeighted(t) {
			e, w, err := ParseWeight(t)
			if err != nil {
				return r, err
			}
			r.Tokens[i] = e
			r.Graph.Tokens[i] = e
			for j, edge := range r.Graph.Edges {
				if edge.To == i {
					r.Graph.Edges[j].Weight = w
				}
			}
		}
	}

	return r, nil
}

func GetProductions(r Rule) []string {
	var out []string
	for _, path := range GetAllPaths(r.Graph) {
		prod := GetSingleProduction(path, FilterTerminals(GetTokens(r), []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"}))
		if prod != "" {
			out = append(out, prod)
		}
	}

	return out
}

func GetSingleProduction(p Path, a []Expression) string {
	if len(p) == 0 || len(a) == 0 {
		return ""
	}

	var b strings.Builder

	for _, i := range p {
		b.WriteString(a[i])
	}

	return b.String()
}

func FilterTerminals(a []Expression, f []string) []Expression {
	var filter map[string]struct{} = make(map[string]struct{})
	var a1 []Expression = make([]Expression, len(a))

	copy(a1, a)
	for _, s := range f {
		filter[s] = struct{}{}
	}
	for i, e := range a1 {
		_, ok := filter[e]
		if ok {
			a1[i] = ""
		}
	}

	return a1
}
