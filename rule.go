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
	s := NewRule(r.exp, r.is_public)
	s.references = make([]string, len(r.references))
	copy(s.references, r.references)
	s.tokens = make([]Expression, len(r.tokens))
	copy(s.tokens, r.tokens)
	s.productions = make([]Expression, len(r.productions))
	copy(s.productions, r.productions)
	s.graph = r.graph.Copy()

	return s
}

func (r Rule) Productions() (out []string) {
	for _, path := range r.graph.AllPaths() {
		prod := singleProduction(path, r.productions)
		if prod != "" {
			out = append(out, prod)
		}
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

func (r Rule) ResolveReferences(m map[string]Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	var r1 Rule
	var err error
	if len(r.references) == 0 {
		return r, nil
	}
	r1 = r
	rules := make(map[string]Rule)
	for k, v := range m {
		rules[k] = v
	}
	for _, ref := range r.references {
		if ref == "" {
			continue
		}
		r2, ok := rules[ref]
		if !ok {
			return r, errors.New("referenced rule does not exist in grammar")
		}
		r1, err = r1.SingleResolveReference(ref, r2, lex)
		if err != nil {
			return r1, err
		}
	}
	return r1, nil
}

func (r Rule) SingleResolveReference(ref string, rule Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	r1 := r
	for i, t := range r1.exp.ToTokens(lex) {
		if t.str() == ref {
			g, err := ComposeGraphs(r1.graph, rule.graph, i)
			if err != nil {
				return r, err
			}
			r1.graph = g
			r1.tokens = g.Nodes
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
