// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:18 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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

func ValidateJSGFRule(s string) error {
	// optional public declaration
	// the name of the rule being defined, in <>
	// an equals sign `='
	// the expansion of the rule
	// a closing semi-colon `;'.
	if !regexp.MustCompile("^(public )?<.+?> ?= ?.*?;$").MatchString(s) {
		return errors.New("invalid jsgf line")
	}

	return nil
}

func ValidateJSGFName(s string) error {
	// grammar name;
	if !regexp.MustCompile("^grammar .+?;$").MatchString(s) {
		return errors.New("invalid jsgf name declaration")
	}

	return nil
}

func ValidateJSGFImport(s string) error {
	// import <gram.rule>;
	if !regexp.MustCompile("^import <.+?>;$").MatchString(s) {
		return errors.New("invalid jsgf import")
	}

	return nil
}

func CreateNameSpace(p string, e string) (map[string]string, error) {
	rs := make(map[string]string)
	imports, err := ImportOrder(p, e)
	if err != nil {
		return make(map[string]string), err
	}
	for _, t := range imports {
		gram, _, _ := strings.Cut(CleanImportStatement(t), ".")
		path, err := FindGrammar(p, gram, e)
		if err != nil {
			return make(map[string]string), err
		}
		_, _, rules, err := PeekGrammar(path)
		if err != nil {
			return make(map[string]string), err
		}
		for k, v := range rules {
			rs[k] = v
		}
	}

	return rs, nil
}

func PeekGrammar(p string) (string, []string, map[string]string, error) {
	var (
		err     error
		name    string
		imports []string
	)

	rules := make(map[string]string)
	f, err := os.Open(p)
	if err != nil {
		return name, imports, rules, errors.New(fmt.Sprint("unable to open grammar ", p))
	}
	info, err := f.Stat()
	if err != nil {
		return name, imports, rules, errors.New(fmt.Sprint("unable to open grammar ", p))
	}
	if info.IsDir() {
		return name, imports, rules, errors.New(fmt.Sprint("provided path is a directory", p))
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
			name = strings.TrimSpace(name)
			name = strings.TrimPrefix(name, "public ")
			rules[name] = strings.TrimSpace(rule)
		default:
			continue
		}
	}

	return name, imports, rules, nil
}

func FindGrammar(p string, t string, e string) (string, error) {
	var (
		target string
		found  bool
	)
	err := filepath.Walk(filepath.Dir(p), func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == e {
			name, _, _, err := PeekGrammar(path)
			if err != nil {
				return err
			}
			if name == t {
				found = true
				target = path

				return io.EOF
			}
		}

		return nil
	})
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return "", err
	}
	if !found {
		return "", errors.New(fmt.Sprint("grammar ", target, " not declared in available directories"))
	}

	return target, nil
}

func ImportOrder(p string, e string) ([]string, error) {
	var (
		imports []string
		imp     string
		res     []string
	)
	_, imports, _, err := PeekGrammar(p)
	if err != nil {
		return imports, err
	}
	for len(imports) > 0 {
		imp, imports = imports[0], imports[1:]
		gram, _, _ := strings.Cut(CleanImportStatement(imp), ".")
		path, err := FindGrammar(p, gram, e)
		if err != nil {
			return []string{}, err
		}
		_, imps, _, err := PeekGrammar(path)
		if err != nil {
			return []string{}, err
		}
		imports = append(imports, imps...)
		res = append(res, imp)
	}

	return res, nil
}
