// -*- coding: utf-8 -*-

// Created on Thu Sep  5 07:38:44 PM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	mrand "math/rand/v2"
	"slices"
	"strings"

	xrand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/sampleuv"
)

// TODO: doc
type Graph struct {
	Tokens   []Expression
	Edges    EdgeList
	children map[int][]int
	weights  map[int]map[int]float64
}

// TODO: doc
func NewGraph(e EdgeList, n []Expression) Graph {
	graph := Graph{}
	graph.Tokens = n
	graph.children = make(map[int][]int)
	graph.weights = make(map[int]map[int]float64)
	for _, edge := range e {
		graph = graph.addEdge(edge)
	}

	return graph
}

// TODO: doc
func (g Graph) getFrom(i int) []int {
	children, ok := g.children[i]
	if !ok {
		return []int{}
	}

	return children
}

func (g Graph) getWeight(f int, t int) float64 {
	weight, ok := g.weights[f][t]
	if !ok {
		return 1.0
	}

	return weight
}

// TODO: doc
func (g Graph) addEdge(e Edge) Graph {
	if e.isEmpty() {
		return g
	}

	g.Edges = append(g.Edges, e)
	g.children[e.From] = append(g.children[e.From], e.To)
	_, ok := g.weights[e.From]
	if !ok {
		g.weights[e.From] = make(map[int]float64)
	}
	g.weights[e.From][e.To] = e.Weight

	return g
}

// TODO: doc
func (g Graph) dropNode(i int) Graph {
	var (
		from       []int
		to         []int
		edges      EdgeList
		start, end int = getEndPoints(g)
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

// TODO: doc
func getEndPoints(g Graph) (int, int) {
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

// TODO: doc
type Path = []int

// TODO: doc
func getAllPaths(g Graph) []Path {
	var (
		from, to int    = getEndPoints(g)
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
		for _, n := range g.getFrom(node) {
			tmp = make(Path, len(path)+1)
			copy(tmp, path)
			tmp[len(path)] = n
			paths = append(paths, tmp)
		}
	}

	return res
}

// TODO: doc
func composeGraphs(g Graph, g1 Graph, i int) (Graph, error) {
	switch {
	case g.Edges.isEmpty() || g1.Edges.isEmpty():
		return Graph{}, fmt.Errorf("error when calling ComposeGraphs(%v, %v, %v):\n%+w", g, g1, i, errors.New("one or more EdgeLists e and a are empty"))
	case i < 0:
		return Graph{}, fmt.Errorf("error when calling ComposeGraphs(%v, %v, %v):\n%+w", g, g1, i, errors.New("cannot insert EdgeList a at negative index"))
	case i > g.Edges.max():
		return Graph{}, fmt.Errorf("error when calling ComposeGraphs(%v, %v, %v):\n%+w", g, g1, i, errors.New("cannot insert EdgeList g1 at index greater than EdgeList g.Max()"))
	}

	g1.Edges = increment(g1.Edges, g.Edges.max()+1)
	g1From, g1To := getEndPoints(g1)
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

// TODO: doc
func getRandomChoice(c []int, w []float64, s xrand.Source) (int, error) {
	if len(c) == 0 || len(w) == 0 {
		return -1, fmt.Errorf("error when calling GetRandomChoice(%v, %v):\n%+w", c, w, errors.New("length of choices c and/or weights w is 0"))
	}
	if len(c) != len(w) {
		return -1, fmt.Errorf("error when calling GetRandomChoice(%v, %v):\n%+w", c, w, errors.New("length of choices c and weights w do not match"))
	}

	choice, ok := sampleuv.NewWeighted(w, s).Take()
	if !ok {
		return -1, fmt.Errorf("error when calling GetRandomChoice(%v, %v):\n%+w", c, w, errors.New("sampleuv.NewWeighted could not sample from choices c and weights w"))
	}

	return c[choice], nil
}

// TODO: doc
func getRandomPath(g Graph) (Path, error) {
	var (
		source   xrand.Source = xrand.NewSource(mrand.Uint64())
		from, to int          = getEndPoints(g)
		res      Path         = Path{from}
		node     int          = from
		choice   int
	)

	for node != to {
		n := g.getFrom(node)
		switch len(n) {
		case 0:
			return Path{}, fmt.Errorf("error when calling GetRandomPath(%v), GetFrom(%v):\n%+w", g, n, errors.New("cannot proceed further down path, no nodes are reachable from n"))
		case 1:
			choice = n[0]
			res = append(res, choice)
			node = choice
		default:
			w := make([]float64, len(n))
			for i, dest := range n {
				w[i] = g.getWeight(from, dest)
			}

			choice, err := getRandomChoice(g.getFrom(node), w, source)
			if err != nil {
				return Path{}, fmt.Errorf("in GetRandomPath(%v):\n%+w", g, err)
			}
			res = append(res, choice)
			node = choice
		}
	}

	return res, nil
}

// TODO: doc
func Minimize(g Graph, f []string) Graph {
	var g1 Graph = g

	for i, t := range g1.Tokens {
		if slices.Contains(f, t) {
			g1 = g1.dropNode(i)
		}
	}

	return g1
}

// TODO: doc
func weightEdges(r Rule) (Rule, error) {
	for i, t := range r.Tokens {
		if isWeighted(t) {
			exp, weight, err := ParseWeight(t)
			if err != nil {
				return r, fmt.Errorf("in WeightEdges(%v):\n%+w", r, err)
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

// TODO: doc
func getProductions(r Rule) []string {
	var productions []string
	for _, path := range getAllPaths(r.Graph) {
		prod := getSingleProduction(path, filterTokens(getTokens(r), jsgfFilter))
		if prod != "" {
			productions = append(productions, prod)
		}
	}

	return productions
}

// TODO: doc
func getSingleProduction(p Path, a []Expression) string {
	if len(p) == 0 || len(a) == 0 {
		return ""
	}

	var builder strings.Builder

	for _, i := range p {
		builder.WriteString(a[i])
	}

	return builder.String()
}

// TODO: doc
func filterTokens(e []Expression, f []string) []Expression {
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
