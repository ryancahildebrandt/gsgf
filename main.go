// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.Open("./data/tea_base.jsgf")
	// f, err := os.Open("./data/test.jsgf")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	lexer := NewJSGFLexer()
	grammar := NewGrammar()
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "import"):
			fmt.Println(line)
		case strings.HasPrefix(line, "public"), strings.HasPrefix(line, "<"):
			name, rule, err := ParseRule(lexer, line)
			if err != nil {
				log.Fatal(err)
			}
			rule.tokens = rule.exp.ToTokens(lexer)
			rule.graph = NewGraph(BuildEdgeList(rule.tokens), rule.tokens)
			rule.productions = FilterTerminals(rule.tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			grammar.Rules[name] = rule
			for _, p := range rule.Productions() {
				fmt.Println(p)
			}
		default:
			continue
		}
	}
	r1 := grammar.Rules["<order>"]
	fmt.Println(r1.tokens)
	fmt.Println(r1.graph.Nodes)
	fmt.Println(r1.graph.Edges)
	fmt.Println(r1.graph.AllPaths())
	for _, p := range r1.Productions() {
		fmt.Println(p)
	}
	fmt.Println("")
	
	
	// grammar, err = grammar.Resolve()
	// fmt.Println(err)
	// // grammar.Productions()
	// for _, p := range grammar.Productions() {
	// 	fmt.Println(p)
	// }

	// for i, path := range []string{
	// 	"dir0/dir1/dir2/gram.rule",
	// 	"dir0/dir1/dir2/gram.*",
	// 	"dir1/dir2/gram.rule",
	// 	"dir1/dir2/gram.*",
	// 	"dir2/gram.rule",
	// 	"dir2/gram.*",
	// 	"gram.rule",
	// 	"gram.*",
	// 	"../dir2/gram.rule",
	// 	"../dir2/gram.*",
	// 	"../../dir1/dir2/gram.rule",
	// 	"../../dir1/dir2/gram.*",
	// 	"../../../dir0/dir1/dir2/gram.rule",
	// 	"../../../dir0/dir1/dir2/gram.*",
	// } {
	// 	root, ru, err := ParseImportPath(path)
	// 	fmt.Println(i, root, ru, err)
	// 	// f, err = os.Open(fmt.Sprint("data/tests/", root, ".jsgf"))
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// 	// fmt.Println(f)
	// }

	// //----
	// f, err = os.Open("./data/tea.jsgf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// scanner = bufio.NewScanner(f)
	// fmt.Println(ReadRule(scanner, "main"))
	fmt.Printf("Took %s", time.Since(start))
}
