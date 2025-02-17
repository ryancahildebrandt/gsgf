// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:18 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// type Import struct {
// 	path string
// ext    string
// file   string
// target string
// gram   string
// rule   string
// dir    string
// }

// func NewImport(s string) Import {
// if strings.HasPrefix(s, "import <") {
// 	s = CleanImportStatement(s)
// }
// if s == "" {
// 	return Import{}
// }
// i := Import{}
// i.path = s
// i.ext = e
// i.dir = filepath.Dir(s)
// i.target = filepath.Base(s)
// i.rule = strings.TrimPrefix(filepath.Ext(i.target), ".")
// i.gram = strings.TrimSuffix(i.target, fmt.Sprint(".", i.rule))
// i.file = fmt.Sprint(i.gram, i.ext)
// return i
// }

func Peek(p string) (string, []string, map[string][]string, error) {
	grammar := ""
	imports := []string{}
	rules := make(map[string][]string)

	f, err := os.Open(p)
	if err != nil {
		return grammar, imports, rules, errors.New(fmt.Sprint("unable to open grammar from import: ", p))
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "grammar "):
			grammar = CleanGrammarStatement(line)
		case strings.HasPrefix(line, "import <"):
			imports = append(imports, line)
		case strings.HasPrefix(line, "<") || strings.HasPrefix(line, "public <"):
			name, rule, _ := strings.Cut(line, "=")
			for _, ref := range regexp.MustCompile(`<.*?>`).FindAllString(rule, -1) {
				name = UnwrapRule(name)
				ref = UnwrapRule(ref)
				rules[name] = append(rules[name], ref)
			}
		default:
		}
	}
	return grammar, imports, rules, nil
}

func WrapRule(s string) string {
	return fmt.Sprint("<", s, ">")
}

func UnwrapRule(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "public ")
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")
	return s
}

func CleanImportStatement(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "import ")
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")
	return s
}

func CleanGrammarStatement(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "grammar ")
	s = strings.TrimSuffix(s, ";")
	return s
}

func CreateNameSpace(p string, e string) (map[string][]string, map[string]map[string][]string, error) {
	var rs = make(map[string]map[string][]string)
	var is = make(map[string][]string)

	err := filepath.Walk(filepath.Dir(p), func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == e {
			grammar, imports, rules, err := Peek(path)
			if err != nil {
				return err
			}
			rs[grammar] = rules
			is[grammar] = imports
		}
		return nil
	})
	if err != nil {
		return is, rs, err
	}
	return is, rs, nil
}
