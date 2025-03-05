// -*- coding: utf-8 -*-

// Created on Sat Feb 22 10:19:30 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"slices"
	"strings"
)

type GrammarJSON struct {
	Rules   map[string]RuleJSON `json:"rules"`
	Imports []string            `json:"imports"`
}

type GraphJSON struct {
	Tokens []string   `json:"tokens"`
	Edges  []EdgeJSON `json:"edges"`
	Paths  [][]int    `json:"paths"`
}

type EdgeJSON struct {
	From   int     `json:"source"`
	To     int     `json:"destination"`
	Weight float64 `json:"weight"`
}

type RuleJSON struct {
	Expression string    `json:"expression"`
	IsPublic   bool      `json:"is_public"`
	Graph      GraphJSON `json:"graph"`
}

func RuleToJSON(r Rule) RuleJSON {
	return RuleJSON{Expression: r.Exp, IsPublic: r.IsPublic, Graph: GraphToJSON(r.Graph)}
}

func EdgeToJSON(e Edge) EdgeJSON {
	return EdgeJSON(e)
}

func GraphToJSON(g Graph) GraphJSON {
	var j GraphJSON

	j.Tokens = append(j.Tokens, g.Tokens...)
	j.Paths = append(j.Paths, GetAllPaths(g)...)
	for _, i := range g.Edges {
		j.Edges = append(j.Edges, EdgeToJSON(i))
	}

	return j
}

func GrammarToJSON(g Grammar) GrammarJSON {
	var rules map[string]RuleJSON = make(map[string]RuleJSON)

	for k, v := range g.Rules {
		rules[k] = RuleToJSON(v)
	}

	return GrammarJSON{Rules: rules, Imports: g.Imports}
}

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

func ReferencesToDOT(g Grammar) string {
	var (
		builder strings.Builder
		entry   string
		entries []string
	)

	builder.WriteString("digraph {\n\n")
	builder.WriteString("\trankdir = \"LR\"\n\n")
	for k, v := range g.Rules {
		for _, r := range GetReferences(v) {
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

func ReferencesToD2(g Grammar) string {
	var (
		builder strings.Builder
		entry   string
		entries []string
	)

	builder.WriteString("direction: right\n\n")
	for k, v := range g.Rules {
		for _, r := range GetReferences(v) {
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
