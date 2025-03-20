// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:11 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bzick/tokenizer"
)

// Contains information for a jsgf rule, including graph, public/private, and base expression
type Rule struct {
	Graph

	IsPublic bool
	exp      Expression
}

func NewRule(e Expression, isPublic bool) Rule {
	r := Rule{}
	r.exp = e
	r.IsPublic = isPublic

	return r
}

// Helper function to pull rule tokens from the nodes of the rule's graph
func getTokens(r Rule) []Expression {
	return r.Graph.Tokens
}

// Returns any rules referenced in a rule definition
func getReferences(r Rule) []string {
	var refs []string
	var seen map[string]struct{} = make(map[string]struct{})

	for _, r := range regexp.MustCompile(`<.*?>`).FindAllString(r.exp, -1) {
		_, ok := seen[r]
		if !ok {
			seen[r] = struct{}{}
			refs = append(refs, r)
		}
	}

	return refs
}

// Composes graphs of rules referenced by main rule
// Returns error if referenced rule is not availabe
func ResolveReferences(r Rule, m map[string]Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	if len(getReferences(r)) == 0 {
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
	for _, ref := range getReferences(r) {
		if ref == "" {
			continue
		}
		r2, ok := rules[ref]
		if !ok {
			return r, fmt.Errorf("error when calling ResolveReferences(%v, %v, %v), rule %v, ref %v:\n%+w", r, m, lex, rules, ref, errors.New("referenced rule does not exist in grammar"))
		}
		r1, err = singleResolveReference(r1, ref, r2, lex)
		if err != nil {
			return r1, err
		}
	}

	return r1, nil
}

// Composes a single referenced rule into its parent rule
// Returns an error if graphs cannot be composed
func singleResolveReference(r Rule, ref string, r1 Rule, lex *tokenizer.Tokenizer) (Rule, error) {
	var r2 Rule = r

	for i, t := range ToTokens(r2.exp, lex) {
		if t == ref {
			g, err := composeGraphs(r2.Graph, r1.Graph, i)
			if err != nil {
				return r, err
			}
			r2.Graph = g
			r2.Tokens = g.Tokens
		}
	}

	return r2, nil
}

// Splits a jsgf rule statement into its constituent name and rule expression
// Returns an error if the line is not a valid rule statement
func ParseRule(line string, lex *tokenizer.Tokenizer) (string, Rule, error) {
	var (
		name string
		exp  string
		err  error
	)

	err = ValidateJSGFRule(line)
	if err != nil {
		return "", Rule{}, err
	}
	name, exp, _ = strings.Cut(line, "=")
	name = strings.TrimPrefix(name, "public ")
	name = strings.TrimSpace(name)
	exp = strings.TrimSpace(exp)

	return name, NewRule(exp, strings.HasPrefix(line, "public")), nil
}

// Checks that a rule does not reference itself either directly or indirectly
func ValidateRuleRecursion(n string, r Rule, m map[string]Rule) error {
	refs := getReferences(r)

	if len(refs) == 0 {
		return nil
	}
	if slices.Contains(refs, n) {
		return fmt.Errorf("error when calling ValidateRuleRecursion(%v, %v, %v), references %v:\n%+w", n, r, m, refs, errors.New("rule references self directly"))
	}
	for _, v := range m {
		if slices.Contains(getReferences(v), n) {
			return fmt.Errorf("error when calling ValidateRuleRecursion(%v, %v, %v), references %v, rule %v:\n%+w", n, r, m, refs, v, errors.New("rule references self indirectly"))
		}
	}

	return nil
}
