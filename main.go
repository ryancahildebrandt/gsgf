// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"encoding/json"
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

	ns, es := GraphToTxt(grammar.Rules["<main>"].Graph)
	WriteToFile([]byte(ns), "outputs/nodes.txt")
	WriteToFile([]byte(es), "outputs/edges.txt")

	var j []byte
	j, err = json.MarshalIndent(GraphToJson(grammar.Rules["<main>"].Graph), "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	WriteToFile(j, "outputs/graph.json")

	j, err = json.MarshalIndent(GrammarToJson(grammar), "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	WriteToFile(j, "outputs/grammar.json")

	j = []byte(GraphToDot(grammar.Rules["<main>"].Graph))
	WriteToFile(j, "outputs/full_graph.dot")
	j = []byte(GraphToDot(grammar.Rules["<main>"].Graph.Minimize()))
	WriteToFile(j, "outputs/minimized_graph.dot")
	j = []byte(ReferencesToDot(grammar))
	WriteToFile(j, "outputs/references.dot")

	fmt.Printf("Took %s", time.Since(start))
}
