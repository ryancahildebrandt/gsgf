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
	graph := Graph{}
	graph.Tokens = n
	graph.Children = make(map[int][]int)
	graph.Weights = make(map[int]map[int]float64)
	for _, edge := range e {
		graph = graph.AddEdge(edge)
	}

	return graph
}

func (g Graph) GetFrom(i int) []int {
	children, ok := g.Children[i]
	if !ok {
		return []int{}
	}

	return children
}

func (g Graph) GetWeight(f int, t int) float64 {
	weight, ok := g.Weights[f][t]
	if !ok {
		return 1.0
	}

	return weight
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
		from       []int
		to         []int
		edges      EdgeList
		start, end int = GetEndPoints(g)
	)

	for _, edge := range g.Edges {
		switch i {
		case start, end:
			edges = append(edges, edge)
		case edge.From:
			to = append(to, edge.To)
		case edge.To:
			from = append(from, edge.From)
		default:
			edges = append(edges, edge)
		}
	}

	for _, f := range from {
		for _, t := range to {
			edges = append(edges, Edge{From: f, To: t, Weight: 1.0})
		}
	}

	return NewGraph(Unique(edges), g.Tokens)
}

func GetEndPoints(g Graph) (int, int) {
	var (
		i         int
		f         int
		fromNodes map[int]struct{} = make(map[int]struct{})
		toNodes   map[int]struct{} = make(map[int]struct{})
		edges     EdgeList         = Sort(g.Edges)
	)

	for _, edge := range edges {
		fromNodes[edge.From] = struct{}{}
		toNodes[edge.To] = struct{}{}
	}
	for _, e := range edges {
		_, ok := toNodes[e.From]
		if !ok {
			i = e.From
		}
		_, ok = fromNodes[e.To]
		if !ok {
			f = e.To
		}
	}

	return i, f
}

type Path = []int

func GetAllPaths(g Graph) []Path {
	var (
		from, to int    = GetEndPoints(g)
		paths    []Path = []Path{{from}}
		path     Path
		res      []Path
		tmp      Path
		node     int
	)

	for len(paths) > 0 {
		path, paths = paths[0], paths[1:]
		node = path[len(path)-1]
		if node == to {
			res = append(res, path)

			continue
		}
		for _, n := range g.GetFrom(node) {
			tmp = make(Path, len(path)+1)
			copy(tmp, path)
			tmp[len(path)] = n
			paths = append(paths, tmp)
		}
	}

	return res
}

func ComposeGraphs(g Graph, g1 Graph, i int) (Graph, error) {
	switch {
	case g.Edges.IsEmpty() || g1.Edges.IsEmpty():
		return Graph{}, errors.New("one or more EdgeLists e and a are empty")
	case i < 0:
		return Graph{}, errors.New("cannot insert EdgeList a at negative index")
	case i > g.Edges.Max():
		return Graph{}, errors.New("cannot insert EdgeList a at index greater than EdgeList g.Max()")
	}

	g1.Edges = Increment(g1.Edges, g.Edges.Max()+1)
	g1From, g1To := GetEndPoints(g1)
	exp := append(g.Tokens, g1.Tokens...)
	edg := g1.Edges

	for _, edge := range g.Edges {
		e := edge
		if edge.From == i {
			e.From = g1To
		}
		if edge.To == i {
			e.To = g1From
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
	choice, ok := sampleuv.NewWeighted(w, nil).Take()
	if !ok {
		return -1, errors.New("sampleuv.NewWeighted could not sample from choices c and weights w")
	}

	return c[choice], nil
}

func GetRandomPath(g Graph) (Path, error) {
	var (
		from, to int  = GetEndPoints(g)
		res      Path = Path{from}
		node     int  = from
		choice   int
	)

	for node != to {
		n := g.GetFrom(node)
		switch len(n) {
		case 0:
			return Path{}, errors.New("cannot proceed further down path")
		case 1:
			choice = n[0]
			res = append(res, choice)
			node = choice
		default:
			w := make([]float64, len(n))
			for i, dest := range n {
				w[i] = g.GetWeight(from, dest)
			}
			choice, err := GetRandomChoice(g.GetFrom(node), w)
			if err != nil {
				return Path{}, err
			}
			res = append(res, choice)
			node = choice
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
			exp, weight, err := ParseWeight(t)
			if err != nil {
				return r, err
			}
			r.Tokens[i] = exp
			r.Graph.Tokens[i] = exp
			for j, edge := range r.Graph.Edges {
				if edge.To == i {
					r.Graph.Edges[j].Weight = weight
				}
			}
		}
	}

	return r, nil
}

func GetProductions(r Rule) []string {
	var productions []string
	for _, path := range GetAllPaths(r.Graph) {
		prod := GetSingleProduction(path, FilterTokens(GetTokens(r), []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"}))
		if prod != "" {
			productions = append(productions, prod)
		}
	}

	return productions
}

func GetSingleProduction(p Path, a []Expression) string {
	if len(p) == 0 || len(a) == 0 {
		return ""
	}

	var builder strings.Builder

	for _, i := range p {
		builder.WriteString(a[i])
	}

	return builder.String()
}

func FilterTokens(e []Expression, f []string) []Expression {
	var filter map[string]struct{} = make(map[string]struct{})
	var e1 []Expression = make([]Expression, len(e))
	copy(e1, e)

	for _, s := range f {
		filter[s] = struct{}{}
	}
	for i, s := range e1 {
		_, ok := filter[s]
		if ok {
			e1[i] = ""
		}
	}

	return e1
}
