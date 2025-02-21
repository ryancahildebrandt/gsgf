// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	basepath := "./data/tests/test0.jsgf"
	// ext := ".jsgf"
	fmt.Println(basepath)

	fmt.Println("----")
	// grammar, imports, rules, err := NewGrammar(basepath).Peek()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(grammar, imports)
	// for _, v := range rules {
	// 	for _, ref := range v {
	// 		_, ok := rules[ref]
	// 		if !ok {
	// 			fmt.Println(false, ref)
	// 		}
	// 	}
	// }

	// for _, imp := range imports {
	// 	imp = CleanImportStatement(imp)
	// 	gram, rule, _ := strings.Cut(imp, ".")
	// 	fmt.Println("GRAM", gram)
	// 	fmt.Println("RULE", rule)
	// 	v, ok := rs[gram]
	// 	if !ok {
	// 		fmt.Println("gram not ok")
	// 	}
	// 	fmt.Println(v)
	// 	if rule == "*" {
	// 		fmt.Println(v)
	// 		continue
	// 	}
	// 	w, ok := v[rule]
	// 	if !ok {
	// 		fmt.Println("rule not ok")
	// 	}
	// 	fmt.Println(w)
	// }

	fmt.Printf("Took %s", time.Since(start))
}
