// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	// start := time.Now()
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// basepath := "./data/tests/test0.jsgf"
	// ext := ".jsgf"
	// fmt.Println(basepath)
	// grammar := NewGrammar()
	// f, err := os.Open(basepath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// scanner := bufio.NewScanner(f)
	// lex := NewJSGFLexer()
	// grammar, err = ImportLines(grammar, scanner, lex)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// namespace, err := CreateNameSpace(basepath, ext)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// grammar = ImportNameSpace(grammar, namespace, lex)
	// grammar, err = ResolveRules(grammar, lex)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, p := range GetAllProductions(grammar) {
	// 	fmt.Println(p)
	// }
	// for _, p := range WrapProductions([]string{"abc", "{}{}", "ab{cd}ef"}, "PRE: ", ": SUF") {
	// 	fmt.Println(p)
	// }
	// fmt.Printf("Took %s", time.Since(start))

	var stringFlag string
	var boolFlag bool
	var intFlag int64

	var flags []cli.Flag = []cli.Flag{
		&cli.IntFlag{
			Name:        "number",
			Aliases:     []string{"n"},
			Value:       0,
			Usage:       "",
			Destination: &intFlag,
		},
		&cli.StringFlag{
			Name:        "inPath",
			Aliases:     []string{"i"},
			Value:       "",
			Usage:       "",
			Destination: &stringFlag,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "outFile",
			Aliases:     []string{"o"},
			Value:       "./productions.txt",
			Usage:       "",
			Destination: &stringFlag,
		},
		&cli.StringFlag{
			Name:        "exportDir",
			Aliases:     []string{"e"},
			Value:       "",
			Usage:       "",
			Destination: &stringFlag,
		},
		&cli.BoolFlag{
			Name:        "minimize",
			Aliases:     []string{"m"},
			Usage:       "",
			Destination: &boolFlag,
		},
		&cli.BoolFlag{
			Name:        "shuffle",
			Aliases:     []string{"s"},
			Usage:       "",
			Destination: &boolFlag,
		},
		&cli.BoolFlag{
			Name:        "singleQuote",
			Aliases:     []string{"q"},
			Usage:       "",
			Destination: &boolFlag,
		},
	}
	app := &cli.Command{
		Name:                  "GSGF",
		Usage:                 "Generate natural language expressions from context free grammars",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "generate",
				Aliases:     []string{"gen"},
				Usage:       "gsgf generate",
				UsageText:   "",
				Description: "Produce all expressions from a grammar file, disregarding token weights",
				ArgsUsage:   "",
				Flags:       flags,
				Action:      func(ctx context.Context, cmd *cli.Command) error { return nil },
			},
			{
				Name:        "sample",
				Aliases:     []string{"sam"},
				Usage:       "gsgf sample",
				UsageText:   "",
				Description: "Produce expressions from a grammar file, according to provided token weights",
				ArgsUsage:   "",
				Flags:       flags,
				Action:      func(ctx context.Context, cmd *cli.Command) error { return nil },
			},
			{
				Name:        "export",
				Aliases:     []string{"exp"},
				Usage:       "gsgf export",
				UsageText:   "",
				Description: "Save graph and grammar representations to disk",
				ArgsUsage:   "",
				Flags:       flags,
				Action:      func(ctx context.Context, cmd *cli.Command) error { return nil },
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error { return nil },
	}
	app.Run(context.Background(), os.Args)
}
