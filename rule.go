// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:11 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bzick/tokenizer"
)

type Rule struct {
	exp         Expression
	is_public   bool
	references  []string
	graph       Graph
	tokens      []Expression
	productions []Expression
}

func NewRule(e Expression, is_public bool) Rule {
	r := Rule{}
	r.exp = e
	r.is_public = is_public
	seen := make(map[string]struct{})
	for _, ref := range regexp.MustCompile(`<.*?>`).FindAllString(e.str(), -1) {
		_, ok := seen[ref]
		if !ok {
			seen[ref] = struct{}{}
			r.references = append(r.references, ref)
		}
	}
	return r
}

func (r Rule) Copy() Rule {
	references := make([]string, len(r.references))
	copy(references, r.references)
	tokens := make([]Expression, len(r.tokens))
	copy(tokens, r.tokens)
	productions := make([]Expression, len(r.productions))
	copy(productions, r.productions)

	s := Rule{}
	s.exp = r.exp
	s.is_public = r.is_public
	s.references = references
	s.graph = r.graph
	s.tokens = tokens
	s.productions = productions

	return s
}

func (r Rule) Productions() (out []string) {
	for _, path := range r.graph.AllPaths() {
		out = append(out, singleProduction(path, r.productions))
	}
	return out
}

func singleProduction(p Path, a []Expression) string {
	if len(p) == 0 || len(a) == 0 {
		return ""
	}

	var b strings.Builder
	for _, i := range p {
		b.WriteString(a[i].str())
	}
	return b.String()
}

func FilterTerminals(a []Expression, f []string) []Expression {
	filter := make(map[string]struct{})
	for _, s := range f {
		filter[s] = struct{}{}
	}
	a1 := make([]Expression, len(a))
	copy(a1, a)
	for i, e := range a1 {
		_, ok := filter[e.str()]
		if ok {
			a1[i] = ""
		}
	}
	return a1
}

func (r Rule) ResolveReferences(m map[string]Rule) (Rule, error) {
	if len(r.references) == 0 {
		return r, nil
	}
	r1 := r.Copy()
	for _, ref := range r1.references {
		for i, t := range r1.tokens {
			if t.str() == ref {
				m_rule, ok := m[ref]
				if !ok {
					return r, errors.New("referenced rule does not exist in grammar")
				}
				g, err := ComposeGraphs(r1.graph, m_rule.graph, i)
				if err != nil {
					return r, err
				}
				r1.graph = g
				r1.tokens = g.Nodes
			}
		}
	}
	return r1, nil
}

func (r Rule) WeightEdges() (Rule, error) {
	for i, t := range r.tokens {
		if t.IsWeighted() {
			e, w, err := t.ParseWeight()
			if err != nil {
				return r, err
			}
			r.tokens[i] = e
			r.graph.Nodes[i] = e
			for j, edge := range r.graph.Edges {
				if edge.to == i {
					r.graph.Edges[j].weight = w
				}
			}
		}
	}
	return r, nil
}

func ParseRule(lex *tokenizer.Tokenizer, line string) (string, Rule, error) {
	var name string
	var exp string

	if line == "" {
		return name, Rule{}, errors.New("line cannot be empty")
	}
	if !ValidateJSGF(line) {
		return name, Rule{}, errors.New("invalid line")
	}
	stream := lex.ParseString(strings.TrimSpace(line))
	for stream.IsValid() {
		switch {
		case stream.CurrentToken().Is(AngleOpen):
			name, _ = captureString(stream, ">", true)
		case stream.CurrentToken().Is(Assignment):
			stream.GoNext()
			exp, _ = captureString(stream, ";", true)
		}
		stream.GoNext()
	}
	return name, NewRule(Expression(exp), strings.HasPrefix(line, "public")), nil
}
