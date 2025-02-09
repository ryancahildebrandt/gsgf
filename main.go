// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	basepath := "./data/tests/dir2/dir1/dir0/test.jsgf"
	// path := "./data/test.jsgf"
	f, err := os.Open(basepath)
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
			s := line
			s = strings.TrimPrefix(s, "import <")
			s = strings.TrimSuffix(s, ">")
			grammar.Imports = append(grammar.Imports, s)
		case strings.HasPrefix(line, "public"), strings.HasPrefix(line, "<"):
			name, rule, err := ParseRule(lexer, line)
			if err != nil {
				log.Fatal(err)
			}
			rule.tokens = rule.exp.ToTokens(lexer)
			rule.graph = NewGraph(BuildEdgeList(rule.tokens), rule.tokens)
			rule.productions = FilterTerminals(rule.tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>"})
			grammar.Rules[name] = rule
			// s, e := rule.graph.EndPoints()
			// fmt.Println(name, s, e)
			// for _, p := range rule.Productions() {
			// 	fmt.Println(p)
			// }
		default:
			continue
		}
	}

	// grammar, err = grammar.Resolve(lexer)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, p := range grammar.Productions() {
	// 	fmt.Println(p)
	// }
	fmt.Println(basepath)             // ./dir/file.ext
	fmt.Println(path.Base(basepath))  // file.ext
	fmt.Println(path.Clean(basepath)) // dir/file.ext
	// Replace multiple slashes with a single slash.
	// Eliminate each . path name element (the current directory).
	// Eliminate each inner .. path name element (the parent directory) along with the non-.. element that precedes it.
	// Eliminate .. elements that begin a rooted path: that is, replace "/.." by "/" at the beginning of a path.
	fmt.Println(path.Dir(basepath))   // dir
	fmt.Println(path.Ext(basepath))   // .ext
	fmt.Println(path.Split(basepath)) // [./dir, file.ext]
	p := strings.LastIndex(basepath, "/")
	pp := basepath[0 : p+1]
	for _, i := range grammar.Imports {
		if i == "" {
			continue
		}
		// if i == filename {continue}
		//fmt.Println(i)
		ppp := fmt.Sprint(pp, i)
		_, gram, _, err := SplitImportPath(ppp)
		if err != nil {
			log.Fatal(err)
		}
		pppp := fmt.Sprint(gram, ".jsgf")
		f, err = os.Open(pppp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("--------", pppp)
		scanner = bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}
	//----
	// f, err = os.Open(path)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// scanner = bufio.NewScanner(f)
	// fmt.Println(ReadAllRules(scanner))
	fmt.Printf("Took %s", time.Since(start))
}
