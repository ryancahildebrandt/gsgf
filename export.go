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
	var tokens []string
	for _, t := range r.Tokens {
		tokens = append(tokens, t)
	}

	return RuleJSON{Expression: r.Exp, IsPublic: r.IsPublic, Graph: GraphToJSON(r.Graph)}
}

func EdgeToJSON(e Edge) EdgeJSON {
	return EdgeJSON(e)
}

func GraphToJSON(g Graph) GraphJSON {
	var j GraphJSON

	for _, i := range g.Tokens {
		j.Tokens = append(j.Tokens, i)
	}
	for _, i := range GetAllPaths(g) {
		j.Paths = append(j.Paths, i)
	}
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
	var builder strings.Builder
	var entry string

	builder.WriteString("digraph {\n\n")
	builder.WriteString("\trankdir = \"LR\"\n\n")
	for k, v := range g.Rules {
		for _, r := range GetReferences(v) {
			entry = fmt.Sprintf("\t%s -> %s;\n", r, k)
			builder.WriteString(entry)
		}
	}
	builder.WriteString("\n}")

	return builder.String()
}
