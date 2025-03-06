// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	basepath := "./data/tests/test0.jsgf"
	ext := ".jsgf"
	fmt.Println(basepath)
	grammar := NewGrammar()
	f, err := os.Open(basepath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	lex := NewJSGFLexer()
	grammar, err = ImportLines(grammar, scanner, lex)
	if err != nil {
		log.Fatal(err)
	}
	namespace, err := CreateNameSpace(basepath, ext)
	if err != nil {
		log.Fatal(err)
	}
	grammar = ImportNameSpace(grammar, namespace, lex)
	grammar, err = ResolveRules(grammar, lex)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range GetAllProductions(grammar) {
		fmt.Println(p)
	}
	for _, p := range WrapProductions([]string{"abc", "{}{}", "ab{cd}ef"}, "PRE: ", ": SUF") {
		fmt.Println(p)
	}
	fmt.Printf("Took %s", time.Since(start))
}
