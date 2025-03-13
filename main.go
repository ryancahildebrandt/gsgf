// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	// start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// fmt.Printf("Took %s", time.Since(start))

	app := &cli.Command{
		Name:                  "GSGF",
		Usage:                 "Generate natural language expressions from context free grammars",
		UsageText:             "gsgf [COMMAND] [OPTIONS...] example.jsgf",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:      "generate",
				Aliases:   []string{"gen"},
				UsageText: "gsgf generate [OPTIONS...] example.jsgf",
				Usage:     "Produce all expressions from a grammar file, disregarding token weights",
				Flags: []cli.Flag{
					&ArgNProductions,
					&ArgOutFile,
					&ArgExportDir,
					&ArgMinimize,
					&ArgShuffle,
					&ArgSingleQuote,
					&ArgWrapProductionsPrefix,
					&ArgWrapProductionsSuffix,
					&ArgCollectTagsChar,
					&ArgWrapTagsPrefix,
					&ArgWrapTagsSuffix,
					&ArgRemoveTags,
					&ArgRenderNewlines,
					&ArgRenderTabs,
					&ArgRemoveMultiSpaces,
					&ArgRemoveEndSpaces,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					return nil
				},
			},
			{
				Name:      "sample",
				Aliases:   []string{"sam"},
				UsageText: "gsgf sample [OPTIONS...] example.jsgf",
				Usage:     "Produce expressions from a grammar file, according to provided token weights",
				Flags: []cli.Flag{
					&ArgNProductions,
					&ArgOutFile,
					&ArgExportDir,
					&ArgMinimize,
					&ArgShuffle,
					&ArgSingleQuote,
					&ArgWrapProductionsPrefix,
					&ArgWrapProductionsSuffix,
					&ArgCollectTagsChar,
					&ArgWrapTagsPrefix,
					&ArgWrapTagsSuffix,
					&ArgRemoveTags,
					&ArgRenderNewlines,
					&ArgRenderTabs,
					&ArgRemoveMultiSpaces,
					&ArgRemoveEndSpaces,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
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
					return nil
				},
			},
			{
				Name:      "export",
				Aliases:   []string{"exp"},
				UsageText: "gsgf export [OPTIONS...] example.jsgf",
				Usage:     "Save graph and grammar representations to disk",
				Flags: []cli.Flag{
					&ArgExportDir,
					&ArgMinimize,
					&ArgSingleQuote,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var ext string = ".jsgf"
					var grammar Grammar = NewGrammar()

					inPath := cmd.Args().First()
					f, err := os.Open(inPath)
					if err != nil {
						log.Fatal(err)
					}
					scanner := bufio.NewScanner(f)
					lex := NewJSGFLexer()
					grammar, err = ImportLines(grammar, scanner, lex)
					if err != nil {
						log.Fatal(err)
					}
					namespace, err := CreateNameSpace(inPath, ext)
					if err != nil {
						log.Fatal(err)
					}
					grammar = ImportNameSpace(grammar, namespace, lex)
					grammar, err = ResolveRules(grammar, lex)
					if err != nil {
						log.Fatal(err)
					}

					j, err := json.MarshalIndent(GrammarToJSON(grammar), "", "\t")
					if err != nil {
						log.Fatal(err)
					}
					err = os.WriteFile(fmt.Sprint(exportDir, "/grammar.json"), j, 0644)
					if err != nil {
						log.Fatal(err)
					}

					err = os.WriteFile(fmt.Sprint(exportDir, "/references.d2"), []byte(ReferencesToD2(grammar)), 0644)
					if err != nil {
						log.Fatal(err)
					}

					err = os.WriteFile(fmt.Sprint(exportDir, "/references.dot"), []byte(ReferencesToDOT(grammar)), 0644)
					if err != nil {
						log.Fatal(err)
					}

					for k, v := range grammar.Rules {
						if v.IsPublic {
							j, err := json.MarshalIndent(GraphToJSON(v.Graph), "", "\t")
							if err != nil {
								log.Fatal(err)
							}
							err = os.WriteFile(fmt.Sprint(exportDir, "/", k, "_graph.json"), j, 0644)
							if err != nil {
								log.Fatal(err)
							}

							edges, nodes := GraphToTXT(v.Graph)
							err = os.WriteFile(fmt.Sprint(exportDir, "/", k, "_edges.txt"), []byte(edges), 0644)
							if err != nil {
								log.Fatal(err)
							}

							err = os.WriteFile(fmt.Sprint(exportDir, "/", k, "_nodes.txt"), []byte(nodes), 0644)
							if err != nil {
								log.Fatal(err)
							}
							err = os.WriteFile(fmt.Sprint(exportDir, "/", k, "_graph.d2"), []byte(GraphToD2(v.Graph)), 0644)
							if err != nil {
								log.Fatal(err)
							}
							err = os.WriteFile(fmt.Sprint(exportDir, "/", k, "_graph.dot"), []byte(GraphToDOT(v.Graph)), 0644)
							if err != nil {
								log.Fatal(err)
							}
						}
					}
					return nil
				},
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	// inPath    string
	// ArgInPath cliArg = cliArg{
	// 	Name:        "inPath",
	// 	Value:       "Grammar file to read",
	// 	Destination: &inPath,
	// }

	nProductions    int64
	ArgNProductions cli.IntFlag = cli.IntFlag{
		Name:        "nProductions",
		Aliases:     []string{"n"},
		Value:       0,
		Usage:       "Number of productions to take from the top of the productions list",
		Destination: &nProductions,
	}

	minimize    bool
	ArgMinimize cli.BoolFlag = cli.BoolFlag{
		Name:        "minimize",
		Aliases:     []string{"m"},
		Value:       false,
		Usage:       "Minimze graph before calculating paths and productions. May boost performance on graphs with many flow control tokens ()[]|",
		Destination: &minimize,
	}
	shuffle    bool
	ArgShuffle cli.BoolFlag = cli.BoolFlag{
		Name:        "shuffle",
		Aliases:     []string{"s"},
		Value:       false,
		Usage:       "Shuffle production order before returning",
		Destination: &shuffle,
	}
	singleQuote    bool
	ArgSingleQuote cli.BoolFlag = cli.BoolFlag{
		Name:        "singleQuote",
		Value:       false,
		Usage:       "Changes lexer's default quote character from double quotes to single",
		Destination: &singleQuote,
	}
	removeTags    bool
	ArgRemoveTags cli.BoolFlag = cli.BoolFlag{
		Name:        "removeTags",
		Value:       false,
		Usage:       "Remove all tags from productions",
		Destination: &removeTags,
	}
	renderNewlines    bool
	ArgRenderNewlines cli.BoolFlag = cli.BoolFlag{
		Name:        "renderNewlines",
		Value:       false,
		Usage:       "Replace \n with new line in productions",
		Destination: &renderNewlines,
	}
	renderTabs    bool
	ArgRenderTabs cli.BoolFlag = cli.BoolFlag{
		Name:        "renderTabs",
		Value:       false,
		Usage:       "Replace \t with tab character in productions",
		Destination: &renderTabs,
	}
	removeMultiSpaces    bool
	ArgRemoveMultiSpaces cli.BoolFlag = cli.BoolFlag{
		Name:        "removeMultiSpaces",
		Value:       false,
		Usage:       "Replace multiple consecutive spaces with single space in productions",
		Destination: &removeMultiSpaces,
	}
	removeEndSpaces    bool
	ArgRemoveEndSpaces cli.BoolFlag = cli.BoolFlag{
		Name:        "removeEndSpaces",
		Value:       false,
		Usage:       "Removing trailing and leading spaces in productions",
		Destination: &removeEndSpaces,
	}

	outFile    string
	ArgOutFile cli.StringFlag = cli.StringFlag{
		Name:        "outFile",
		Aliases:     []string{"o"},
		Value:       "",
		Usage:       "Text file to write productions to. If blank, productions are returned to stdout",
		Destination: &outFile,
	}
	exportDir    string
	ArgExportDir cli.StringFlag = cli.StringFlag{
		Name:        "exportDir",
		Aliases:     []string{"e"},
		Value:       "./export",
		Usage:       "Directory to write export results to",
		Destination: &exportDir,
	}
	wrapProductionsPrefix    string
	ArgWrapProductionsPrefix cli.StringFlag = cli.StringFlag{
		Name:        "wrapProductionsPrefix",
		Value:       "",
		Usage:       "Prefix applied to all productions",
		Destination: &wrapProductionsPrefix,
	}
	wrapProductionsSuffix    string
	ArgWrapProductionsSuffix cli.StringFlag = cli.StringFlag{
		Name:        "wrapProductionsSuffix",
		Value:       "",
		Usage:       "Suffix applied to all productions",
		Destination: &wrapProductionsSuffix,
	}
	collectTagsChar    string
	ArgCollectTagsChar cli.StringFlag = cli.StringFlag{
		Name:        "collectTagsChar",
		Value:       "",
		Usage:       "Comment character to place between production and collected tags. If empty, tags are moved to the end of the production line",
		Destination: &collectTagsChar,
	}
	wrapTagsPrefix    string
	ArgWrapTagsPrefix cli.StringFlag = cli.StringFlag{
		Name:        "wrapTagsPrefix",
		Value:       "",
		Usage:       "Prefix applied to all tags",
		Destination: &wrapTagsPrefix,
	}
	wrapTagsSuffix    string
	ArgWrapTagsSuffix cli.StringFlag = cli.StringFlag{
		Name:        "wrapTagsSuffix",
		Value:       "",
		Usage:       "Suffix applied to all tags",
		Destination: &wrapTagsSuffix,
	}
)
