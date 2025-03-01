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
	Graph

	Exp      Expression
	IsPublic bool
}

func NewRule(e Expression, isPublic bool) Rule {
	r := Rule{}
	r.Exp = e
	r.IsPublic = isPublic

	return r
}

func Tokens(r Rule) []Expression {
	return r.Graph.Tokens
}

func References(r Rule) []string {
	var refs []string
	seen := make(map[string]struct{})

	for _, ref := range regexp.MustCompile(`<.*?>`).FindAllString(r.Exp, -1) {
		_, ok := seen[ref]
		if !ok {
			seen[ref] = struct{}{}
			refs = append(refs, ref)
		}
	}

	return refs
}

func ResolveReferences(r Rule, m map[string]Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	var r1 Rule
	var err error
	if len(References(r)) == 0 {
		return r, nil
	}

	r1 = r
	rules := make(map[string]Rule)
	for k, v := range m {
		rules[k] = v
	}
	for _, ref := range References(r) {
		if ref == "" {
			continue
		}
		r2, ok := rules[ref]
		if !ok {
			return r, errors.New("referenced rule does not exist in grammar")
		}
		r1, err = SingleResolveReference(r1, ref, r2, lex)
		if err != nil {
			return r1, err
		}
	}

	return r1, nil
}

func SingleResolveReference(r Rule, ref string, rule Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	var r1 Rule
	r1 = r
	for i, t := range ToTokens(r1.Exp, lex) {
		if t == ref {
			g, err := ComposeGraphs(r1.Graph, rule.Graph, i)
			if err != nil {
				return r, err
			}
			r1.Graph = g
			r1.Tokens = g.Tokens
		}
	}

	return r1, nil
}

func ParseRule(line string, lex *tokenizer.Tokenizer) (string, Rule, error) {
	var name string
	var exp string

	err := ValidateJSGFRule(line)
	if err != nil {
		return name, Rule{}, err
	}

	name, exp, found := strings.Cut(line, "=")
	if !found {
		return name, Rule{}, errors.New("jsgf line does not contain required assignment =")
	}
	name = strings.TrimPrefix(name, "public ")
	name = strings.TrimSpace(name)
	exp = strings.TrimSpace(exp)

	return name, NewRule(exp, strings.HasPrefix(line, "public")), nil
}

func ValidateRuleRecursion(r Rule, m map[string]Rule) error {
	if len(References(r)) == 0 {
		return nil
	}

	rules := make(map[string]Rule)
	for k, v := range m {
		rules[k] = v
	}
	for _, ref := range References(r) {
		if ref == "" {
			continue
		}
		_, ok := rules[ref]
		if !ok {
			return errors.New("referenced rule does not exist in grammar")
		}
	}

	return nil
}
