// -*- coding: utf-8 -*-

// Created on Sun Mar  9 06:52:40 PM EDT 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"github.com/urfave/cli/v3"
)

var (
	stringFlag string
	boolFlag   bool
	intFlag    int64

	inPath cli.StringArg = cli.StringArg{
		Name:        "inPath",
		Value:       "Grammar file to read. If directory, all grammar files matching inFormat in the top level of the directory will be read",
		Destination: &stringFlag,
	}

	nProductions cli.IntFlag = cli.IntFlag{
		Name:        "nProductions",
		Aliases:     []string{"n"},
		Value:       0,
		Usage:       "Number of productions to take from the top of the productions list",
		Destination: &intFlag,
	}

	minimize cli.BoolFlag = cli.BoolFlag{
		Name:        "minimize",
		Aliases:     []string{"m"},
		Value:       false,
		Usage:       "Minimze graph before calculating paths and productions. May boost performance on graphs with many flow control tokens ()[]|",
		Destination: &boolFlag,
	}
	shuffle cli.BoolFlag = cli.BoolFlag{
		Name:        "shuffle",
		Aliases:     []string{"s"},
		Value:       false,
		Usage:       "Shuffle production order before returning",
		Destination: &boolFlag,
	}
	singleQuote cli.BoolFlag = cli.BoolFlag{
		Name:        "singleQuote",
		Value:       false,
		Usage:       "Changes lexer's default quote character from double quotes to single",
		Destination: &boolFlag,
	}
	removeTags cli.BoolFlag = cli.BoolFlag{
		Name:        "removeTags",
		Value:       false,
		Usage:       "Remove all tags from productions",
		Destination: &boolFlag,
	}
	renderNewlines cli.BoolFlag = cli.BoolFlag{
		Name:        "renderNewlines",
		Value:       false,
		Usage:       "Replace \n with new line in productions",
		Destination: &boolFlag,
	}
	renderTabs cli.BoolFlag = cli.BoolFlag{
		Name:        "renderTabs",
		Value:       false,
		Usage:       "Replace \t with tab character in productions",
		Destination: &boolFlag,
	}
	removeMultiSpaces cli.BoolFlag = cli.BoolFlag{
		Name:        "removeMultiSpaces",
		Value:       false,
		Usage:       "Replace multiple consecutive spaces with single space in productions",
		Destination: &boolFlag,
	}
	removeEndSpaces cli.BoolFlag = cli.BoolFlag{
		Name:        "removeEndSpaces",
		Value:       false,
		Usage:       "Removing trailing and leading spaces in productions",
		Destination: &boolFlag,
	}

	inFormat cli.StringFlag = cli.StringFlag{
		Name:        "inFormat",
		Aliases:     []string{"f"},
		Value:       "",
		Usage:       "File format to read from input directory, one of .jsgf, .jjsgf, .ebnf",
		Destination: &stringFlag,
	}
	outFile cli.StringFlag = cli.StringFlag{
		Name:        "outFile",
		Aliases:     []string{"o"},
		Value:       "",
		Usage:       "Text file to write productions to. If blank, productions are returned to stdout",
		Destination: &stringFlag,
	}
	exportDir cli.StringFlag = cli.StringFlag{
		Name:        "exportDir",
		Aliases:     []string{"e"},
		Value:       "./export",
		Usage:       "Directory to write export results to",
		Destination: &stringFlag,
	}
	wrapProductionsPrefix cli.StringFlag = cli.StringFlag{
		Name:        "wrapProductionsPrefix",
		Value:       "",
		Usage:       "Prefix applied to all productions",
		Destination: &stringFlag,
	}
	wrapProductionsSuffix cli.StringFlag = cli.StringFlag{
		Name:        "wrapProductionsSuffix",
		Value:       "",
		Usage:       "Suffix applied to all productions",
		Destination: &stringFlag,
	}
	collectTagsChar cli.StringFlag = cli.StringFlag{
		Name:        "collectTagsChar",
		Value:       "",
		Usage:       "Comment character to place between production and collected tags. If empty, tags are moved to the end of the production line",
		Destination: &stringFlag,
	}
	wrapTagsPrefix cli.StringFlag = cli.StringFlag{
		Name:        "wrapTagsPrefix",
		Value:       "",
		Usage:       "Prefix applied to all tags",
		Destination: &stringFlag,
	}
	wrapTagsSuffix cli.StringFlag = cli.StringFlag{
		Name:        "wrapTagsSuffix",
		Value:       "",
		Usage:       "Suffix applied to all tags",
		Destination: &stringFlag,
	}
)
