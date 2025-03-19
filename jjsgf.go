// -*- coding: utf-8 -*-

// Created on Fri Mar 14 09:33:35 PM EDT 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"slices"
	"strings"
)

// TODO: doc
type JJSGFGrammarJSON struct {
	Name    string            `json:"grammar"`
	Imports []string          `json:"imports"`
	Public  map[string]string `json:"public"`
	Rules   map[string]string `json:"rules"`
}

// TODO: doc
func JJSGFToJSGF(j JJSGFGrammarJSON) string {
	var (
		b       strings.Builder
		entry   string
		entries []string
	)

	b.WriteString("#JSGF V1.0 ISO8859-1 en;\n")
	b.WriteString(fmt.Sprintf("grammar %s;\n", j.Name))
	for _, i := range j.Imports {
		b.WriteString(fmt.Sprintf("import <%s>;\n", i))
	}
	for k, v := range j.Public {
		entry = fmt.Sprintf("public <%s> = %s;\n", k, v)
		entries = append(entries, entry)
	}
	slices.Sort(entries)
	for _, e := range entries {
		b.WriteString(e)
	}
	entries = []string{}
	for k, v := range j.Rules {
		entry = fmt.Sprintf("<%s> = %s;\n", k, v)
		entries = append(entries, entry)
	}
	slices.Sort(entries)
	for _, e := range entries {
		b.WriteString(e)
	}

	return b.String()
}
