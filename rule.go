// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:11 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"regexp"
	"slices"
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

func GetTokens(r Rule) []Expression {
	return r.Graph.Tokens
}

func GetReferences(r Rule) []string {
	var refs []string
	var seen map[string]struct{} = make(map[string]struct{})

	for _, r := range regexp.MustCompile(`<.*?>`).FindAllString(r.Exp, -1) {
		_, ok := seen[r]
		if !ok {
			seen[r] = struct{}{}
			refs = append(refs, r)
		}
	}

	return refs
}

func ResolveReferences(r Rule, m map[string]Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	if len(GetReferences(r)) == 0 {
		return r, nil
	}

	var (
		r1    Rule = r
		err   error
		rules map[string]Rule = make(map[string]Rule)
	)

	for k, v := range m {
		rules[k] = v
	}
	for _, ref := range GetReferences(r) {
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

func SingleResolveReference(r Rule, ref string, r1 Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	var r2 Rule = r

	for i, t := range ToTokens(r2.Exp, lex) {
		if t == ref {
			g, err := ComposeGraphs(r2.Graph, r1.Graph, i)
			if err != nil {
				return r, err
			}
			r2.Graph = g
			r2.Tokens = g.Tokens
		}
	}

	return r2, nil
}

func ParseRule(line string, lex *tokenizer.Tokenizer) (string, Rule, error) {
	err := ValidateJSGFRule(line)
	if err != nil {
		return "", Rule{}, err
	}

	var name string
	var exp string

	name, exp, found := strings.Cut(line, "=")
	if !found {
		return "", Rule{}, errors.New("jsgf line does not contain required assignment =")
	}
	name = strings.TrimPrefix(name, "public ")
	name = strings.TrimSpace(name)
	exp = strings.TrimSpace(exp)

	return name, NewRule(exp, strings.HasPrefix(line, "public")), nil
}

func ValidateRuleRecursion(n string, r Rule, m map[string]Rule) error {
	refs := GetReferences(r)

	if len(refs) == 0 {
		return nil
	}
	if slices.Contains(refs, n) {
		return errors.New("rule references self")
	}
	for _, v := range m {
		if slices.Contains(GetReferences(v), n) {
			return errors.New("rule references self")
		}
	}

	return nil
}
