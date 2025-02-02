// -*- coding: utf-8 -*-

// Created on Tue Sep 17 11:55:23 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

type Grammar struct {
	Rules map[string]Rule
	// productions []string
}

func NewGrammar() Grammar {
	g := Grammar{}
	g.Rules = make(map[string]Rule)
	return g
}

func (g Grammar) CompositionOrder() []string {
	var rules []string
	var rule string
	var res []string

	for k, v := range g.Rules {
		if v.is_public {
			rules = append(rules, k)
		}
	}
	for len(rules) > 0 {
		rule, rules = rules[0], rules[1:]
		rules = append(rules, g.Rules[rule].references...)
		res = append(res, rule)
	}
	return res
}

func (g Grammar) Productions() []string {
	out := []string{}
	for _, v := range g.Rules {
		if v.is_public {
			out = append(out, v.Productions()...)
		}
	}
	return out
}

func (g Grammar) Resolve() (Grammar, error) {
	rs := g.CompositionOrder()
	rules := make(map[string]Rule)
	for k, v := range g.Rules {
		rules[k] = v
	}
	seen := make(map[string]struct{})
	for i := len(rs) - 1; i >= 0; i-- {
		rname := rs[i]
		r1 := g.Rules[rname]
		_, ok := seen[rname]
		if !ok {
			seen[rname] = struct{}{}
			r2, err := r1.ResolveReferences(g.Rules)
			if err != nil {
				return g, err
			}
			r2.productions = FilterTerminals(r2.tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			g.Rules[rname] = r2
		}
	}
	return g, nil
}
