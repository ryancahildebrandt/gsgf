// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:11 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"maps"
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
	seen := make(map[string]struct{})
	var refs []string

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
	maps.Copy(m, rules)

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
	r1 := r
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

func ParseRule(lex *tokenizer.Tokenizer, line string) (string, Rule, error) {
	var (
		name string
		exp  string
	)

	err := ValidateJSGFRule(line)
	if err != nil {
		return name, Rule{}, err
	}

	stream := lex.ParseString(strings.TrimSpace(line))
	for stream.IsValid() {
		switch {
		case stream.CurrentToken().Is(AngleOpen):
			name, err = captureString(stream, ">", true)
			if err != nil {
				return name, Rule{}, err
			}
		case stream.CurrentToken().Is(Assignment):
			stream.GoNext()
			exp, err = captureString(stream, ";", true)
			if err != nil {
				return name, Rule{}, err
			}
		}
		stream.GoNext()
	}

	return name, NewRule(Expression(exp), strings.HasPrefix(line, "public")), nil
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
