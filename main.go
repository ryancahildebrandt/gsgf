// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	basepath := "./data/tests/dir2/dir1/dir0/test0.jsgf"

	is, rs, err := CreateNameSpace(basepath, ".jsgf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(is)
	fmt.Println(rs)

	fmt.Println("----")
	grammar, imports, _, err := Peek(basepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(grammar, imports)
	for _, imp := range imports {
		imp = CleanImportStatement(imp)
		gram, rule, _ := strings.Cut(imp, ".")
		fmt.Println(gram, rule)
		v, ok := rs[gram]
		if !ok {
			fmt.Println("gram not ok")
		}
		fmt.Println(v)
		if rule == "*" {
			fmt.Println(v)
			continue
		}
		w, ok := v[rule]
		if !ok {
			fmt.Println("rule not ok")
		}
		fmt.Println(w)
	}

	fmt.Printf("Took %s", time.Since(start))
}
