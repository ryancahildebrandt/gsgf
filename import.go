// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:18 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Returns the gram.rule portion of a jsgf import statement
func cleanImportStatement(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "import")
	s = strings.TrimSuffix(s, ";")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")

	return s
}

// Returns the grammar name from a jsgf grammar statement
func cleanGrammarStatement(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "grammar")
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, ";")

	return s
}

// Checks that the string is a valid jsgf rule containing:
// - optional public declaration
// - the name of the rule being defined, in <>
// - an equals sign =
// - the expansion of the rule
// - a closing semi-colon ;
func ValidateJSGFRule(s string) error {
	if !regexp.MustCompile("^(public )?<.+?> ?= ?.*?;$").MatchString(s) {
		return fmt.Errorf("error when calling ValidateJSGFRule(%v):\n%+w", s, errors.New("invalid jsgf line"))
	}

	return nil
}

// Checks that the string is a valid jsgf grammar declaration, containing:
// - grammar
// - name
// - a closing semicolon ;
func ValidateJSGFName(s string) error {
	if !regexp.MustCompile("^grammar .+?;$").MatchString(s) {
		return fmt.Errorf("error when calling ValidateJSGFName(%v):\n%+w", s, errors.New("invalid jsgf name declaration"))
	}

	return nil
}

// Checks that the string is a valid jsgf import statement, containing:
// - import
// - opening and closing angle brackets <>
// - grammar name
// - optional rule name or *
// - a closing semicolon ;
func ValidateJSGFImport(s string) error {
	if !regexp.MustCompile("^import <.+?>;$").MatchString(s) {
		return fmt.Errorf("error when calling ValidateJSGFImport(%v):\n%+w", s, errors.New("invalid jsgf import"))
	}

	return nil
}

// Collects all required rules from grammar files in subdirectories of the provided path
// - Reads import order from the root grammar
// - For each imported grammar, reads each import statement and rule
// Returns an error if the required grammars cannot be found or opened
func CreateNameSpace(p string, e string) (map[string]string, error) {
	var res map[string]string = make(map[string]string)

	imports, err := getImportOrder(p, e)
	if err != nil {
		return make(map[string]string), fmt.Errorf("in CreateNameSpace(%v, %v):\n%+w", p, e, err)
	}
	for _, imp := range imports {
		gram, _, _ := strings.Cut(cleanImportStatement(imp), ".")
		path, err := findGrammar(p, gram, e)
		if err != nil {
			return make(map[string]string), fmt.Errorf("in CreateNameSpace(%v, %v):\n%+w", p, e, err)
		}
		_, _, rules, err := peekGrammar(path)
		if err != nil {
			return make(map[string]string), fmt.Errorf("in CreateNameSpace(%v, %v):\n%+w", p, e, err)
		}
		for k, v := range rules {
			res[k] = v
		}
	}

	return res, nil
}

// Checks the specified grammar file and returns the name, imports, and rules specified in the grammar
// Returns an error if the specified file cannot be opened or converted to grammar
func peekGrammar(p string) (string, []string, map[string]string, error) {
	var (
		err     error
		name    string
		imports []string
		rules   map[string]string = make(map[string]string)
		ext     string            = filepath.Ext(p)
		scanner *bufio.Scanner
	)

	f, err := os.Open(p)
	if err != nil {
		return "", []string{}, map[string]string{}, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, err)
	}
	info, err := f.Stat()
	if err != nil {
		return "", []string{}, map[string]string{}, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, err)
	}
	if info.IsDir() {
		return "", []string{}, map[string]string{}, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, errors.New("provided path is a directory"))
	}

	switch ext {
	case ".jsgf":
		scanner = bufio.NewScanner(f)
	case ".jjsgf":
		var jj JJSGFGrammarJSON
		err = json.NewDecoder(f).Decode(&jj)
		if err != nil {
			return "", []string{}, map[string]string{}, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, errors.New("error decoding json file"))
		}
		scanner = bufio.NewScanner(strings.NewReader(JJSGFToJSGF(jj)))
	default:
		return "", []string{}, map[string]string{}, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, errors.New("unsupported extension, not one of .jsgf, .jjsgf"))
	}

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "grammar "):
			err = ValidateJSGFName(line)
			if err != nil {
				return name, imports, rules, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, err)
			}
			name = cleanGrammarStatement(line)
		case strings.HasPrefix(line, "import <"):
			err = ValidateJSGFImport(line)
			if err != nil {
				return name, imports, rules, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, err)
			}
			imports = append(imports, line)
		case strings.HasPrefix(line, "<") || strings.HasPrefix(line, "public <"):
			err = ValidateJSGFRule(line)
			if err != nil {
				return name, imports, rules, fmt.Errorf("in PeekGrammar(%v):\n%+w", p, err)
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

// Returns the on-disk location of the specified grammar by checking each subdirectory of the specified path
// Returns an error if the target grammar is not found in files with given extension
func findGrammar(p string, t string, e string) (string, error) {
	var target string
	var found bool

	err := filepath.Walk(filepath.Dir(p), func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == e {
			name, _, _, err := peekGrammar(path)
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
		return "", fmt.Errorf("error when calling FindGrammar(%v, %v, %v):\n%+w", p, t, e, err)
	}
	if !found {
		return "", fmt.Errorf("error when calling FindGrammar(%v, %v, %v):\n%+w", p, t, e, errors.New("grammar not declared in available directories"))
	}

	return target, nil
}

// Returns the dependencies of the giben grammar file by traversing each listed grammar and its imports, in order
// Returns an error if the file cannot be opened or located
func getImportOrder(p string, e string) ([]string, error) {
	var (
		imports []string
		imp     string
		res     []string
	)
	_, imports, _, err := peekGrammar(p)
	if err != nil {
		return imports, fmt.Errorf("in GetImportOrder(%v, %v):\n%+w", p, e, err)
	}
	for len(imports) > 0 {
		imp, imports = imports[0], imports[1:]
		gram, _, _ := strings.Cut(cleanImportStatement(imp), ".")
		path, err := findGrammar(p, gram, e)
		if err != nil {
			return []string{}, fmt.Errorf("in GetImportOrder(%v, %v):\n%+w", p, e, err)
		}
		_, imps, _, err := peekGrammar(path)
		if err != nil {
			return []string{}, fmt.Errorf("in GetImportOrder(%v, %v):\n%+w", p, e, err)
		}
		imports = append(imports, imps...)
		res = append(res, imp)
	}

	return res, nil
}
