// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	mrand "math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

func main() {
	start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
					var ext string = ".jsgf"
					var grammar Grammar = NewGrammar()
					var productions []string
					var err error
					var inFile string = cmd.Args().First()
					var quoteChar string = "\""

					err = ValidateInFile(inFile)
					if err != nil {
						log.Fatal(err)
					}
					err = ValidateOutFile(outFile)
					if err != nil {
						log.Fatal(err)
					}

					f, err := os.Open(inFile)
					if err != nil {
						log.Fatal(err)
					}
					scanner := bufio.NewScanner(f)

					if singleQuote {
						quoteChar = "'"
					}
					lex := NewJSGFLexer(quoteChar)
					grammar, err = ImportLines(grammar, scanner, lex)
					if err != nil {
						log.Fatal(err)
					}
					namespace, err := CreateNameSpace(inFile, ext)
					if err != nil {
						log.Fatal(err)
					}
					grammar = ImportNameSpace(grammar, namespace, lex)

					if minimize {
						for k, v := range grammar.Rules {
							v.Graph = Minimize(v.Graph)
							grammar.Rules[k] = v
						}
					}

					grammar, err = ResolveRules(grammar, lex)
					if err != nil {
						log.Fatal(err)
					}
					productions = GetAllProductions(grammar)

					if shuffle {
						mrand.Shuffle(len(productions), func(i, j int) { productions[i], productions[j] = productions[j], productions[i] })
					}

					if nProductions != -1 {
						productions = productions[0:nProductions]
					}

					if wrapProductionsPrefix != "" || wrapProductionsSuffix != "" {
						productions = WrapProductions(productions, wrapProductionsPrefix, wrapProductionsSuffix)
					}
					if wrapTagsPrefix != "" || wrapTagsSuffix != "" {
						productions = WrapProductions(productions, wrapTagsPrefix, wrapTagsSuffix)
					}
					if collectTagsChar != "" {
						productions = CollectTags(productions, collectTagsChar)
					}

					if removeTags {
						productions = RemoveTags(productions)
					}
					if removeMultiSpaces {
						productions = RemoveMultipleSpaces(productions)
					}
					if removeEndSpaces {
						productions = RemoveEndSpaces(productions)
					}
					if renderNewlines {
						productions = RenderNewLines(productions)
					}
					if renderTabs {
						productions = RenderTabs(productions)
					}

					if outFile == "" {
						for _, prod := range productions {
							fmt.Println(prod)
						}

						return nil
					}
					err = os.WriteFile(outFile, []byte(strings.Join(productions, "\n")), 0644)
					if err != nil {
						log.Fatal(err)
					}

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
					var ext string = ".jsgf"
					var grammar Grammar = NewGrammar()
					var inFile string = cmd.Args().First()
					var productions []string
					var err error
					var quoteChar string = "\""

					err = ValidateInFile(inFile)
					if err != nil {
						log.Fatal(err)
					}
					err = ValidateOutFile(outFile)
					if err != nil {
						log.Fatal(err)
					}

					f, err := os.Open(inFile)
					if err != nil {
						log.Fatal(err)
					}
					scanner := bufio.NewScanner(f)

					if singleQuote {
						quoteChar = "'"
					}
					lex := NewJSGFLexer(quoteChar)
					grammar, err = ImportLines(grammar, scanner, lex)
					if err != nil {
						log.Fatal(err)
					}
					namespace, err := CreateNameSpace(inFile, ext)
					if err != nil {
						log.Fatal(err)
					}
					grammar = ImportNameSpace(grammar, namespace, lex)

					if minimize {
						for k, v := range grammar.Rules {
							v.Graph = Minimize(v.Graph)
							grammar.Rules[k] = v
						}
					}

					grammar, err = ResolveRules(grammar, lex)
					if err != nil {
						log.Fatal(err)
					}

					var keys []string
					for k, v := range grammar.Rules {
						if v.IsPublic {
							keys = append(keys, k)
						}
					}

					if nProductions == -1 {
						nProductions = 1
					}
					for len(productions) < int(nProductions) {
						key := keys[mrand.IntN(len(keys))]
						graph := grammar.Rules[key].Graph
						path, err := GetRandomPath(graph)
						if err != nil {
							log.Fatal(err)
						}
						prod := GetSingleProduction(path, FilterTokens(graph.Tokens, []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>", ""}))
						productions = append(productions, prod)
					}

					if shuffle {
						mrand.Shuffle(len(productions), func(i, j int) { productions[i], productions[j] = productions[j], productions[i] })
					}

					if wrapProductionsPrefix != "" || wrapProductionsSuffix != "" {
						productions = WrapProductions(productions, wrapProductionsPrefix, wrapProductionsSuffix)
					}
					if wrapTagsPrefix != "" || wrapTagsSuffix != "" {
						productions = WrapProductions(productions, wrapTagsPrefix, wrapTagsSuffix)
					}
					if collectTagsChar != "" {
						productions = CollectTags(productions, collectTagsChar)
					}

					if removeTags {
						productions = RemoveTags(productions)
					}
					if removeMultiSpaces {
						productions = RemoveMultipleSpaces(productions)
					}
					if removeEndSpaces {
						productions = RemoveEndSpaces(productions)
					}
					if renderNewlines {
						productions = RenderNewLines(productions)
					}
					if renderTabs {
						productions = RenderTabs(productions)
					}

					if outFile == "" {
						for _, prod := range productions {
							fmt.Println(prod)
						}

						return nil
					}
					err = os.WriteFile(outFile, []byte(strings.Join(productions, "\n")), 0644)
					if err != nil {
						log.Fatal(err)
					}

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
					var inFile string = cmd.Args().First()
					var err error
					var quoteChar string = "\""

					err = ValidateInFile(inFile)
					if err != nil {
						log.Fatal(err)
					}
					err = ValidateExportDir(exportDir)
					if err != nil {
						log.Fatal(err)
					}

					f, err := os.Open(inFile)
					if err != nil {
						log.Fatal(err)
					}
					scanner := bufio.NewScanner(f)

					if singleQuote {
						quoteChar = "'"
					}
					lex := NewJSGFLexer(quoteChar)
					grammar, err = ImportLines(grammar, scanner, lex)
					if err != nil {
						log.Fatal(err)
					}
					namespace, err := CreateNameSpace(inFile, ext)
					if err != nil {
						log.Fatal(err)
					}
					grammar = ImportNameSpace(grammar, namespace, lex)

					if minimize {
						for k, v := range grammar.Rules {
							v.Graph = Minimize(v.Graph)
							grammar.Rules[k] = v
						}
					}

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
	fmt.Printf("Took %s", time.Since(start))
}

var (
	nProductions    int64
	ArgNProductions cli.IntFlag = cli.IntFlag{
		Name:        "nProductions",
		Aliases:     []string{"n"},
		Value:       -1,
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

func ValidateInFile(p string) error {
	// valid extension
	// exists
	_, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("in ValidateInFile(%v):\n%+w", p, err)
	}
	switch filepath.Ext(p) {
	case ".jsgf", ".jjsgf", ".ebnf":
		return nil
	default:
		return fmt.Errorf("in ValidateInFile(%v):\n%+w", p, errors.New("file extension is not one of .jsgf, .jjsgf, .ebnf"))
	}
}

func ValidateOutFile(p string) error {
	// dir exists if provided
	_, err := os.Stat(filepath.Dir(p))
	if err != nil {
		return fmt.Errorf("in ValidateOutFile(%v):\n%+w", p, err)
	}
	return nil
}

func ValidateExportDir(p string) error {
	// is dir
	// exists
	info, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("in ValidateExportDir(%v):\n%+w", p, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("in ValidateExportDir(%v):\n%+w", p, errors.New("provided path is not a directory"))
	}
	return nil
}
