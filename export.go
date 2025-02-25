// -*- coding: utf-8 -*-

// Created on Sat Feb 22 10:19:30 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type GrammarJson struct {
	Rules   map[string]RuleJson `json:rules`
	Imports []string            `json:imports`
}

type GraphJson struct {
	Nodes []string   `json:nodes`
	Edges []EdgeJson `json:edges`
	Paths [][]int    `json:paths`
}

type EdgeJson struct {
	From   int     `json:source`
	To     int     `json:destination`
	Weight float64 `json:weight`
}

type RuleJson struct {
	Expression string    `json:expression`
	Is_public  bool      `json:is_public`
	References []string  `json:references`
	Graph      GraphJson `json:graph`
	Tokens     []string  `json:tokens`
}

func RuleToJson(r Rule) RuleJson {
	tokens := []string{}
	for _, i := range r.Tokens {
		tokens = append(tokens, i.str())
	}
	return RuleJson{r.Exp.str(), r.Is_public, r.References, GraphToJson(r.Graph), tokens}
}

func EdgeToJson(e Edge) EdgeJson {
	return EdgeJson(e)
}

func GraphToJson(g Graph) GraphJson {
	gj := GraphJson{}
	for _, i := range g.Nodes {
		gj.Nodes = append(gj.Nodes, i.str())
	}
	for _, i := range g.AllPaths() {
		gj.Paths = append(gj.Paths, i)
	}
	for _, i := range g.Edges {
		gj.Edges = append(gj.Edges, EdgeToJson(i))
	}
	return gj
}

func GrammarToJson(g Grammar) GrammarJson {
	imports := g.Imports
	rules := make(map[string]RuleJson)
	for k, v := range g.Rules {
		rules[k] = RuleToJson(v)
	}
	return GrammarJson{rules, imports}

}

func GraphToTxt(g Graph) (string, string) {
	var nodes []string
	var edges []string

	for _, i := range g.Nodes {
		nodes = append(nodes, fmt.Sprintf("\"%s\"", i.str()))
	}
	for _, i := range g.Edges {
		edges = append(edges, fmt.Sprintf("%v,%v,%v", i.From, i.To, i.Weight))
	}

	return strings.Join(nodes, "\n"), strings.Join(edges, "\n")
}

func GraphToDot(g Graph) string {
	var b strings.Builder
	var entry string
	var visited []int

	b.WriteString("digraph {\n\n")
	b.WriteString("\trankdir = \"LR\"\n\n")
	for _, e := range g.Edges {
		visited = append(visited, e.From)
		visited = append(visited, e.To)
	}
	for i, n := range g.Nodes {
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
		for _, ref := range v.References {
			entry = fmt.Sprintf("\t%s -> %s;\n", ref, k)
			b.WriteString(entry)
		}
	}
	b.WriteString("\n}")
	return b.String()

}

func WriteToFile(b []byte, p string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	w.Write(b)
	w.Flush()

	return nil
}
