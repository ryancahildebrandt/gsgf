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

	"github.com/urfave/cli"
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

// cli
// generate
// sample
// export
// -n
// -shuffle
// -infile
// -outfile
// -outdir
// -minimize
// -quote

  configPath := ""
	lastrun := true
	intflag := 0
  app := &cli.App{
		Name:  "GSGF",
		Usage: "Generate natural language expressions from context free grammars",
		// EnableShellCompletion: true,
		Commands: []cli.Command{
			cli.Command{
				Name:        "generate",
				Usage:       "",
				UsageText:   "",
        Description: "",
				ArgsUsage:   "",
				Flags:       []cli.Flag{},
			  Action: func() bool {return true},
      },
      {
				Name:        "sample",
				Usage:       "",
				UsageText:   "",
        Description: "",
				ArgsUsage:   "",
				Flags:       []cli.Flag{},
			  Action: func() bool {return true},
      },
      {
				Name:        "export",
				Usage:       "",
				UsageText:   "",
        Description: "",
				ArgsUsage:   "",
				Flags:       []cli.Flag{},
			  Action: func() bool {return true},
      },
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "c",
				Value:       0,
				Usage:       "",
				Destination: &intflag,
			},
      &cli.StringFlag{
				Name:        "inFile",
				Value:       "",
				Usage:       "",
				Destination: &configPath,
			},
      &cli.StringFlag{
				Name:        "outFile",
				Value:       "",
				Usage:       "",
				Destination: &configPath,
			},
      &cli.StringFlag{
				Name:        "exportDir",
				Value:       "",
				Usage:       "",
				Destination: &configPath,
			},
			&cli.BoolFlag{
				Name:        "minimize",
				Usage:       "",
				Destination: &lastrun,
			},
      &cli.BoolFlag{
				Name:        "shuffle",
				Usage:       "",
				Destination: &lastrun,
			},
      &cli.BoolFlag{
				Name:        "singleQuote",
				Usage:       "",
				Destination: &lastrun,
			},
		},
		Action: func(*cli.Context) {}}
	app.Run(os.Args)
}
