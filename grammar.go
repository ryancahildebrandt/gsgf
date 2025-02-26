// -*- coding: utf-8 -*-

// Created on Tue Sep 17 11:55:23 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bzick/tokenizer"
)

type Grammar struct {
	Path    string
	Rules   map[string]Rule
	Imports []string
}

func NewGrammar(p string) Grammar {
	g := Grammar{}
	g.Path = p
	g.Rules = make(map[string]Rule)
	return g
}

func (g Grammar) Peek() (string, []string, map[string][]string, error) {
	var err error
	var name string
	var imports []string
	rules := make(map[string][]string)

	f, err := os.Open(g.Path)
	if err != nil {
		return name, imports, rules, errors.New(fmt.Sprint("unable to open grammar from import: ", g.Path))
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "grammar "):
			err = ValidateJSGFName(line)
			if err != nil {
				return name, imports, rules, err
			}
			name = CleanGrammarStatement(line)
		case strings.HasPrefix(line, "import <"):
			err = ValidateJSGFImport(line)
			if err != nil {
				return name, imports, rules, err
			}
			imports = append(imports, line)
		case strings.HasPrefix(line, "<") || strings.HasPrefix(line, "public <"):
			err = ValidateJSGFRule(line)
			if err != nil {
				return name, imports, rules, err
			}
			name, rule, _ := strings.Cut(line, "=")
			name = UnwrapRule(name)
			rules[name] = []string{}
			for _, ref := range regexp.MustCompile(`<.*?>`).FindAllString(rule, -1) {
				ref = UnwrapRule(ref)
				rules[name] = append(rules[name], ref)
			}
		default:
		}
	}
	return name, imports, rules, nil
}

func (g Grammar) CompositionOrder() []string {
	var rules []string
	var rule string
	var res []string

	for k, v := range g.Rules {
		if v.Is_public {
			rules = append(rules, k)
		}
	}
	for len(rules) > 0 {
		rule, rules = rules[0], rules[1:]
		rules = append(rules, g.Rules[rule].References...)
		res = append(res, rule)
	}
	return res
}

func (g Grammar) Productions() []string {
	out := []string{}
	for _, v := range g.Rules {
		if v.Is_public {
			out = append(out, v.Productions()...)
		}
	}
	return out
}

func (g Grammar) Resolve(lex *tokenizer.Tokenizer) (Grammar, error) {
	ord := g.CompositionOrder()
	seen := make(map[string]struct{})
	for i := len(ord) - 1; i >= 0; i-- {
		rname := ord[i]
		r1 := g.Rules[rname]
		_, ok := seen[rname]
		if !ok {
			seen[rname] = struct{}{}
			r2, err := r1.ResolveReferences(g.Rules, lex)
			if err != nil {
				return g, err
			}
			r2.productions = FilterTerminals(r2.Tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			g.Rules[rname] = r2
		}
	}
	return g, nil
}

func (g Grammar) ReadLines(s *bufio.Scanner, lex *tokenizer.Tokenizer) (Grammar, error) {
	var err error
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "import <"):
			err = ValidateJSGFImport(line)
			if err != nil {
				return NewGrammar(""), err
			}
			s := line
			s = strings.TrimPrefix(s, "import <")
			s = strings.TrimSuffix(s, ">")
			g.Imports = append(g.Imports, CleanImportStatement(s))
		case strings.HasPrefix(line, "public <"), strings.HasPrefix(line, "<"):
			err := ValidateJSGFRule(line)
			if err != nil {
				return NewGrammar(""), err
			}
			name, rule, err := ParseRule(lex, line)
			if err != nil {
				return NewGrammar(""), err
			}
			rule.Tokens = rule.Exp.ToTokens(lex)
			rule.Graph = NewGraph(BuildEdgeList(rule.Tokens), rule.Tokens)
			rule.productions = FilterTerminals(rule.Tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			g.Rules[name] = rule
		default:
			continue
		}
	}
	return g, nil
}

func (g Grammar) ReadNameSpace(r map[string]string, lex *tokenizer.Tokenizer) Grammar {
	for k, v := range r {
		rule := NewRule(Expression(v), false)
		rule.Tokens = rule.Exp.ToTokens(lex)
		rule.Graph = NewGraph(BuildEdgeList(rule.Tokens), rule.Tokens)
		rule.productions = FilterTerminals(rule.Tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
		_, ok := g.Rules[k]
		if !ok {
			g.Rules[k] = rule
		}
	}
	return g
}

func ValidateGrammarCompleteness(g Grammar) error {
	for _, v := range g.Rules {
		for _, ref := range v.References {
			_, ok := g.Rules[ref]
			if !ok {
				return errors.New("grammar references rule not present in namespace")
			}
		}
	}
	return nil
}

// basepath := "./data/tests/dir2/dir1/dir0/test.jsgf"
// 	// path := "./data/test.jsgf"
// 	f, err := os.Open(basepath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scanner := bufio.NewScanner(f)
// 	lexer := NewJSGFLexer()
// 	grammar := NewGrammar()
// 	grammar, err = grammar.ReadLines(scanner, lexer)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// grammar, err = grammar.Resolve(lexer)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
