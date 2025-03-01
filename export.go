// -*- coding: utf-8 -*-

// Created on Sat Feb 22 10:19:30 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"slices"
	"strings"
)

type GrammarJson struct {
	Rules   map[string]RuleJson `json:"rules"`
	Imports []string            `json:"imports"`
}
type GraphJson struct {
	Tokens []string   `json:"tokens"`
	Edges  []EdgeJson `json:"edges"`
	Paths  [][]int    `json:"paths"`
}
type EdgeJson struct {
	From   int     `json:"source"`
	To     int     `json:"destination"`
	Weight float64 `json:"weight"`
}
type RuleJson struct {
	Expression string    `json:"expression"`
	IsPublic   bool      `json:"is_public"`
	Graph      GraphJson `json:"graph"`
}

func RuleToJson(r Rule) RuleJson {
	var tokens []string
	for _, i := range r.Tokens {
		tokens = append(tokens, i)
	}

	return RuleJson{Expression: r.Exp, IsPublic: r.IsPublic, Graph: GraphToJson(r.Graph)}
}

func EdgeToJson(e Edge) EdgeJson {
	return EdgeJson(e)
}

func GraphToJson(g Graph) GraphJson {
	var j GraphJson
	for _, i := range g.Tokens {
		j.Tokens = append(j.Tokens, i)
	}
	for _, i := range AllPaths(g) {
		j.Paths = append(j.Paths, i)
	}
	for _, i := range g.Edges {
		j.Edges = append(j.Edges, EdgeToJson(i))
	}

	return j
}

func GrammarToJson(g Grammar) GrammarJson {
	rules := make(map[string]RuleJson)
	for k, v := range g.Rules {
		rules[k] = RuleToJson(v)
	}

	return GrammarJson{Rules: rules, Imports: g.Imports}
}

func GraphToTxt(g Graph) (string, string) {
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

func GraphToDot(g Graph) string {
	var (
		b       strings.Builder
		entry   string
		visited []int
	)

	b.WriteString("digraph {\n\n")
	b.WriteString("\trankdir = \"LR\"\n\n")
	for _, e := range g.Edges {
		visited = append(visited, e.From)
		visited = append(visited, e.To)
	}
	for i, n := range g.Tokens {
		if slices.Contains(visited, i) {
			entry = fmt.Sprintf("\t_%v [label=\"%s\"];\n", i, n)
			b.WriteString(entry)
		}
	}
	b.WriteString("\n")
	for _, e := range g.Edges {
		entry = fmt.Sprintf("\t_%v -> _%v [weight=%v];\n", e.From, e.To, e.Weight)
		if e.Weight != 1.0 {
			entry = fmt.Sprintf("\t_%v -> _%v [label=\"%v\",weight=%v];\n", e.From, e.To, fmt.Sprint(e.Weight), e.Weight)
		}
		b.WriteString(entry)
	}
	b.WriteString("\n}")

	return b.String()
}

func ReferencesToDot(g Grammar) string {
	var b strings.Builder
	var entry string

	b.WriteString("digraph {\n\n")
	b.WriteString("\trankdir = \"LR\"\n\n")
	for k, v := range g.Rules {
		for _, ref := range References(v) {
			entry = fmt.Sprintf("\t%s -> %s;\n", ref, k)
			b.WriteString(entry)
		}
	}
	b.WriteString("\n}")

	return b.String()
}
