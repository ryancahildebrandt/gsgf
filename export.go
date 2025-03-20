// -*- coding: utf-8 -*-

// Created on Sat Feb 22 10:19:30 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"slices"
	"strings"
)

// JSON wrapper for grammar struct
type grammarJSON struct {
	Rules   map[string]ruleJSON `json:"rules"`
	Imports []string            `json:"imports"`
}

// JSON wrapper for graph struct
type graphJSON struct {
	Tokens []string   `json:"tokens"`
	Edges  []edgeJSON `json:"edges"`
	Paths  [][]int    `json:"paths"`
}

// JSON wrapper for edge struct
type edgeJSON struct {
	From   int     `json:"source"`
	To     int     `json:"destination"`
	Weight float64 `json:"weight"`
}

// JSON wrapper for rule struct
type ruleJSON struct {
	Expression string    `json:"expression"`
	IsPublic   bool      `json:"is_public"`
	Graph      graphJSON `json:"graph"`
}

// Construct ruleJSON from rule
func ruleToJSON(r Rule) ruleJSON {
	return ruleJSON{Expression: r.exp, IsPublic: r.IsPublic, Graph: graphToJSON(r.Graph)}
}

// Construct edgeJSON from edge
func edgeToJSON(e Edge) edgeJSON {
	return edgeJSON(e)
}

// Construct graphJSON from graph
func graphToJSON(g Graph) graphJSON {
	var j graphJSON

	j.Tokens = append(j.Tokens, g.Tokens...)
	j.Paths = append(j.Paths, getAllPaths(g)...)
	for _, i := range g.Edges {
		j.Edges = append(j.Edges, edgeToJSON(i))
	}

	return j
}

// Construct grammarJSON from grammar
func grammarToJSON(g Grammar) grammarJSON {
	var rules map[string]ruleJSON = make(map[string]ruleJSON)

	for k, v := range g.Rules {
		rules[k] = ruleToJSON(v)
	}

	return grammarJSON{Rules: rules, Imports: g.Imports}
}

// Export graph to a slices of edges and a slice of nodes, each separated by newlines
func GraphToTXT(g Graph) (string, string) {
	var nodes []string
	var edges []string

	for _, i := range g.Tokens {
		nodes = append(nodes, fmt.Sprintf("\"%s\"", i))
	}
	for _, i := range g.Edges {
		edges = append(edges, fmt.Sprintf("%v,%v,%v", i.From, i.To, i.Weight))
	}

	return strings.Join(nodes, "\n"), strings.Join(edges, "\n")
}

// Export graph to graphviz DOT format
func GraphToDOT(g Graph) string {
	var (
		builder strings.Builder
		entry   string
		visited []int
	)

	builder.WriteString("digraph {\n\n")
	builder.WriteString("\trankdir = \"LR\"\n\n")
	for _, e := range g.Edges {
		visited = append(visited, e.From)
		visited = append(visited, e.To)
	}
	for i, t := range g.Tokens {
		if slices.Contains(visited, i) {
			entry = fmt.Sprintf("\t_%v [label=\"%s\"];\n", i, t)
			builder.WriteString(entry)
		}
	}
	builder.WriteString("\n")
	for _, e := range g.Edges {
		entry = fmt.Sprintf("\t_%v -> _%v [weight=%v];\n", e.From, e.To, e.Weight)
		if e.Weight != 1.0 {
			entry = fmt.Sprintf("\t_%v -> _%v [label=\"%v\",weight=%v];\n", e.From, e.To, fmt.Sprint(e.Weight), e.Weight)
		}
		builder.WriteString(entry)
	}
	builder.WriteString("\n}")

	return builder.String()
}

// Export rule reference tree to graphviz DOT format
func ReferencesToDOT(g Grammar) string {
	var (
		builder strings.Builder
		entry   string
		entries []string
	)

	builder.WriteString("digraph {\n\n")
	builder.WriteString("\trankdir = \"LR\"\n\n")
	for k, v := range g.Rules {
		for _, r := range getReferences(v) {
			entry = fmt.Sprintf("\t%s -> %s;\n", r, k)
			entries = append(entries, entry)
		}
	}
	slices.Sort(entries)
	for _, e := range entries {
		builder.WriteString(e)
	}
	builder.WriteString("\n}")

	return builder.String()
}

// Export graph to D2 diagram format
func GraphToD2(g Graph) string {
	var (
		builder strings.Builder
		entry   string
		visited []int
	)

	builder.WriteString("direction: right\n\n")
	for _, e := range g.Edges {
		visited = append(visited, e.From)
		visited = append(visited, e.To)
	}
	for i, t := range g.Tokens {
		if slices.Contains(visited, i) {
			entry = fmt.Sprintf("_%v: \"%s\"\n", i, t)
			builder.WriteString(entry)
		}
	}
	builder.WriteString("\n")
	for _, e := range g.Edges {
		entry = fmt.Sprintf("_%v -> _%v: \"%v\"\n", e.From, e.To, e.Weight)
		builder.WriteString(entry)
	}

	return builder.String()
}

// Export rule reference tree to D2 diagram format
func ReferencesToD2(g Grammar) string {
	var (
		builder strings.Builder
		entry   string
		entries []string
	)

	builder.WriteString("direction: right\n\n")
	for k, v := range g.Rules {
		for _, r := range getReferences(v) {
			entry = fmt.Sprintf("\"%s\" -> \"%s\"\n", r, k)
			entries = append(entries, entry)
		}
	}
	slices.Sort(entries)
	for _, e := range entries {
		builder.WriteString(e)
	}

	return builder.String()
}
