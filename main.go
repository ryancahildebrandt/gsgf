// -*- coding: utf-8 -*-

// Created on Sun Aug  4 11:54:10 AM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	mrand "math/rand/v2"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

/*
gsgf produces natural language expressions from context free grammars
The main executable can handle jsgf and jjsgf formats and outputs productions to stdout or plain text

Usage:

	gsgf [command] [options] example.jsgf

Options:

	--nProductions, -n (int)
		Number of productions to take from the top of the productions list

	--outFile, -o (string)
		Text file to write productions to.
		If blank, productions are returned to stdout

	--minimize, -m (bool)
		Minimze graph before calculating paths and productions.
		May boost performance on graphs with many flow control tokens ()[]|

	--shuffle, -s (bool)
		Shuffle production order before returning

	--singleQuote (bool)
		Changes lexer's default quote character from double quotes to single

	--wrapProductionsPrefix (string)
		Prefix applied to all productions

	--wrapProductionsSuffix (string)
		Suffix applied to all productions

	--collectTagsChar (string)
		Comment character to place between production and collected tags.
		If empty, tags are moved to the end of the production line

	--wrapTagsPrefix (string)
		Prefix applied to all tags

	--wrapTagsSuffix (string)
		Suffix applied to all tags

	--removeTags (bool)
		Remove all tags from productions

	--renderNewlines (bool)
		Replace \n with new line in productions

	--renderTabs (bool)
		Replace \t with tab character in productions

	--removeMultiSpaces (bool)
		Replace multiple consecutive spaces with single space in productions

	--removeEndSpaces (bool)
		Removing trailing and leading spaces in productions

	--exportDir, -e (string) (default: "./export")
		Directory to write export results to

	--help, -h
		Show help
*/
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := &cli.Command{
		Name:                  "GSGF",
		Usage:                 "Generate natural language expressions from context free grammars",
		UsageText:             "gsgf [COMMAND] [OPTIONS] example.jsgf",
		EnableShellCompletion: true,
		Suggest:               true,
		Commands: []*cli.Command{
			{
				Name:                  "generate",
				UsageText:             "gsgf generate [OPTIONS] example.jsgf",
				Usage:                 "Produce all expressions from a grammar file, disregarding token weights",
				EnableShellCompletion: true,
				Suggest:               true,
				Before:                prepareContext,
				Flags: []cli.Flag{
					&inFile,
					&ext,
					&quoteChar,
					&nProductions,
					&outFile,
					&minimize,
					&shuffle,
					&singleQuote,
					&wrapProductionsPrefix,
					&wrapProductionsSuffix,
					&collectTagsChar,
					&wrapTagsPrefix,
					&wrapTagsSuffix,
					&removeTags,
					&renderNewlines,
					&renderTabs,
					&removeMultiSpaces,
					&removeEndSpaces,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var (
						grammar     Grammar
						productions []string
						err         error
					)

					err = ValidateInFile(cmd.String("inFile"))
					if err != nil {
						log.Fatal(err)
					}
					err = ValidateOutFile(cmd.String("outFile"))
					if err != nil {
						log.Fatal(err)
					}

					grammar, err = buildGrammar(cmd)
					if err != nil {
						log.Fatal(err)
					}
					productions = GetAllProductions(grammar)
					productions = applyPostproc(productions, cmd)
					if cmd.String("outFile") == "" {
						for _, prod := range productions {
							fmt.Println(prod)
						}

						return nil
					}
					err = os.WriteFile(cmd.String("outFile"), []byte(strings.Join(productions, "\n")), 0644)
					if err != nil {
						log.Fatal(err)
					}

					return nil
				},
			},

			{
				Name:                  "sample",
				UsageText:             "gsgf sample [OPTIONS] example.jsgf",
				Usage:                 "Produce expressions from a grammar file, according to provided token weights",
				EnableShellCompletion: true,
				Suggest:               true,
				Before:                prepareContext,
				Flags: []cli.Flag{
					&inFile,
					&ext,
					&quoteChar,
					&nProductions,
					&outFile,
					&minimize,
					&shuffle,
					&singleQuote,
					&wrapProductionsPrefix,
					&wrapProductionsSuffix,
					&collectTagsChar,
					&wrapTagsPrefix,
					&wrapTagsSuffix,
					&removeTags,
					&renderNewlines,
					&renderTabs,
					&removeMultiSpaces,
					&removeEndSpaces,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var (
						grammar     Grammar
						productions []string
						keys        []string
						err         error
					)

					err = ValidateInFile(cmd.String("inFile"))
					if err != nil {
						log.Fatal(err)
					}
					err = ValidateOutFile(cmd.String("outFile"))
					if err != nil {
						log.Fatal(err)
					}

					grammar, err = buildGrammar(cmd)
					if err != nil {
						log.Fatal(err)
					}
					for k, v := range grammar.Rules {
						if v.IsPublic {
							keys = append(keys, k)
						}
					}
					for len(productions) < int(cmd.Int("nProductions")) {
						key := keys[mrand.IntN(len(keys))]
						graph := grammar.Rules[key].Graph
						path, err := getRandomPath(graph)
						if err != nil {
							log.Fatal(err)
						}
						prod := getSingleProduction(path, filterTokens(graph.Tokens, jsgfFilter))
						productions = append(productions, prod)
					}
					productions = applyPostproc(productions, cmd)
					if cmd.String("outFile") == "" {
						for _, prod := range productions {
							fmt.Println(prod)
						}

						return nil
					}
					err = os.WriteFile(cmd.String("outFile"), []byte(strings.Join(productions, "\n")), 0644)
					if err != nil {
						log.Fatal(err)
					}

					return nil
				},
			},
			{
				Name:                  "export",
				UsageText:             "gsgf export [OPTIONS] example.jsgf",
				Usage:                 "Save graph and grammar representations to disk",
				EnableShellCompletion: true,
				Suggest:               true,
				Before:                prepareContext,
				Flags: []cli.Flag{
					&inFile,
					&ext,
					&quoteChar,
					&exportDir,
					&minimize,
					&singleQuote,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var (
						grammar Grammar
						j       []byte
						err     error
					)
					err = ValidateInFile(cmd.String("inFile"))
					if err != nil {
						log.Fatal(err)
					}

					grammar, err = buildGrammar(cmd)
					if err != nil {
						log.Fatal(err)
					}
					j, err = json.Marshal(grammarToJSON(grammar))
					if err != nil {
						log.Fatal(err)
					}
					err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/grammar.json"), j, 0644)
					if err != nil {
						log.Fatal(err)
					}
					err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/references.d2"), []byte(ReferencesToD2(grammar)), 0644)
					if err != nil {
						log.Fatal(err)
					}
					err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/references.dot"), []byte(ReferencesToDOT(grammar)), 0644)
					if err != nil {
						log.Fatal(err)
					}
					for k, v := range grammar.Rules {
						if v.IsPublic {
							j, err := json.Marshal(graphToJSON(v.Graph))
							if err != nil {
								log.Fatal(err)
							}
							err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/", k, "_graph.json"), j, 0644)
							if err != nil {
								log.Fatal(err)
							}
							nodes, edges := GraphToTXT(v.Graph)
							err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/", k, "_edges.txt"), []byte(edges), 0644)
							if err != nil {
								log.Fatal(err)
							}
							err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/", k, "_nodes.txt"), []byte(nodes), 0644)
							if err != nil {
								log.Fatal(err)
							}
							err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/", k, "_graph.d2"), []byte(GraphToD2(v.Graph)), 0644)
							if err != nil {
								log.Fatal(err)
							}
							err = os.WriteFile(fmt.Sprint(cmd.String("exportDir"), "/", k, "_graph.dot"), []byte(GraphToDOT(v.Graph)), 0644)
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
