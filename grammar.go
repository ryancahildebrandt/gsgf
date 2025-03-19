// -*- coding: utf-8 -*-

// Created on Tue Sep 17 11:55:23 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/bzick/tokenizer"
)

// TODO: doc
type Grammar struct {
	Rules   map[string]Rule
	Imports []string
}

func NewGrammar() Grammar {
	var grammar Grammar

	grammar.Rules = make(map[string]Rule)

	return grammar
}

// TODO: doc
func getCompositionOrder(g Grammar) []string {
	var (
		rules []string
		rule  string
		res   []string
	)

	for k, v := range g.Rules {
		if v.IsPublic {
			rules = append(rules, k)
		}
	}
	for len(rules) > 0 {
		rule, rules = rules[0], rules[1:]
		rules = append(rules, getReferences(g.Rules[rule])...)
		res = append(res, rule)
	}

	return res
}

// TODO: doc
func GetAllProductions(g Grammar) []string {
	var productions []string

	for _, v := range g.Rules {
		if v.IsPublic {
			productions = append(productions, getProductions(v)...)
		}
	}

	return productions
}

// TODO: doc
func ResolveRules(g Grammar, lex *tokenizer.Tokenizer) (Grammar, error) {
	var order []string = getCompositionOrder(g)
	var seen map[string]struct{} = make(map[string]struct{})

	for i := len(order) - 1; i >= 0; i-- {
		rname := order[i]
		r1 := g.Rules[rname]
		_, ok := seen[rname]
		if !ok {
			seen[rname] = struct{}{}

			r2, err := ResolveReferences(r1, g.Rules, lex)
			if err != nil {
				return g, err
			}
			g.Rules[rname] = r2
		}
	}

	return g, nil
}

// TODO: doc
func FomJSGF(g Grammar, s *bufio.Scanner, lex *tokenizer.Tokenizer) (Grammar, error) {
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "import <"):
			err := ValidateJSGFImport(line)
			if err != nil {
				return NewGrammar(), err
			}
			g.Imports = append(g.Imports, cleanImportStatement(line))
		case strings.HasPrefix(line, "public <"), strings.HasPrefix(line, "<"):
			err := ValidateJSGFRule(line)
			if err != nil {
				return NewGrammar(), err
			}
			name, rule, err := ParseRule(line, lex)
			if err != nil {
				return NewGrammar(), err
			}
			rule.Tokens = ToTokens(rule.exp, lex)
			rule.Graph = NewGraph(ToEdgeList(rule.Tokens), rule.Tokens)
			g.Rules[name] = rule
		default:
			continue
		}
	}

	return g, nil
}

// TODO: doc
func ImportNameSpace(g Grammar, r map[string]string, lex *tokenizer.Tokenizer) Grammar {
	for k, v := range r {
		rule := NewRule(v, false)
		rule.Tokens = ToTokens(rule.exp, lex)
		rule.Graph = NewGraph(ToEdgeList(rule.Tokens), rule.Tokens)
		_, ok := g.Rules[k]
		if !ok {
			g.Rules[k] = rule
		}
	}

	return g
}

// TODO: doc
func ValidateGrammarCompleteness(g Grammar) error {
	for _, v := range g.Rules {
		for _, r := range getReferences(v) {
			_, ok := g.Rules[r]
			if !ok {
				return fmt.Errorf("error when calling ValidateGrammarCompleteness(%v), on rule %v, reference %v:\n%+w", g, v, r, errors.New("grammar references rule not present in namespace"))
			}
		}
	}

	return nil
}
