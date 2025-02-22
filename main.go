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

	fmt.Println("----")
	grammar := NewGrammar(basepath)
	f, err := os.Open(basepath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	lex := NewJSGFLexer()
	grammar, err = grammar.ReadLines(scanner, lex)
	if err != nil {
		log.Fatal(err)
	}
	namespace, err := CreateNameSpace(grammar.Path, ext)
	if err != nil {
		log.Fatal(err)
	}
	if !grammar.IsComplete() {
		grammar = grammar.ReadNameSpace(namespace, lex)
	}
	grammar, err = grammar.Resolve(lex)
	if err != nil {
		log.Fatal(err)
	}
	// o, err := os.Create("./data/tests/productions.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// w := bufio.NewWriter(o)
	for _, p := range grammar.Productions() {
		fmt.Println(p)
		// w.WriteString(p)
		// w.WriteString("\n")
	}
	// w.Flush()

	fmt.Printf("Took %s", time.Since(start))
}
