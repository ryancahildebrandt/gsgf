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

	"github.com/urfave/cli/v3"
)

var (
	inFile       cli.StringFlag = cli.StringFlag{Name: "inFile", Hidden: true}
	ext          cli.StringFlag = cli.StringFlag{Name: "ext", Hidden: true}
	quoteChar    cli.StringFlag = cli.StringFlag{Name: "quoteChar", Value: "\"", Hidden: true}
	nProductions cli.IntFlag    = cli.IntFlag{
		Name:    "nProductions",
		Aliases: []string{"n"},
		Value:   -1,
		Usage:   "Number of productions to take from the top of the productions list",
	}
	minimize cli.BoolFlag = cli.BoolFlag{
		Name:    "minimize",
		Aliases: []string{"m"},
		Usage:   "Minimze graph before calculating paths and productions. May boost performance on graphs with many flow control tokens ()[]|",
	}
	shuffle cli.BoolFlag = cli.BoolFlag{
		Name:    "shuffle",
		Aliases: []string{"s"},
		Usage:   "Shuffle production order before returning",
	}
	singleQuote cli.BoolFlag = cli.BoolFlag{
		Name:  "singleQuote",
		Usage: "Changes lexer's default quote character from double quotes to single",
	}
	removeTags cli.BoolFlag = cli.BoolFlag{
		Name:  "removeTags",
		Usage: "Remove all tags from productions",
	}
	renderNewlines cli.BoolFlag = cli.BoolFlag{
		Name:  "renderNewlines",
		Usage: "Replace \\n with new line in productions",
	}
	renderTabs cli.BoolFlag = cli.BoolFlag{
		Name:  "renderTabs",
		Usage: "Replace \\t with tab character in productions",
	}
	removeMultiSpaces cli.BoolFlag = cli.BoolFlag{
		Name:  "removeMultiSpaces",
		Usage: "Replace multiple consecutive spaces with single space in productions",
	}
	removeEndSpaces cli.BoolFlag = cli.BoolFlag{
		Name:  "removeEndSpaces",
		Usage: "Removing trailing and leading spaces in productions",
	}
	outFile cli.StringFlag = cli.StringFlag{
		Name:    "outFile",
		Aliases: []string{"o"},
		Usage:   "Text file to write productions to. If blank, productions are returned to stdout",
	}
	exportDir cli.StringFlag = cli.StringFlag{
		Name:    "exportDir",
		Aliases: []string{"e"},
		Value:   "./export",
		Usage:   "Directory to write export results to",
	}
	wrapProductionsPrefix cli.StringFlag = cli.StringFlag{
		Name:  "wrapProductionsPrefix",
		Usage: "Prefix applied to all productions",
	}
	wrapProductionsSuffix cli.StringFlag = cli.StringFlag{
		Name:  "wrapProductionsSuffix",
		Usage: "Suffix applied to all productions",
	}
	collectTagsChar cli.StringFlag = cli.StringFlag{
		Name:  "collectTagsChar",
		Usage: "Comment character to place between production and collected tags. If empty, tags are moved to the end of the production line",
	}
	wrapTagsPrefix cli.StringFlag = cli.StringFlag{
		Name:  "wrapTagsPrefix",
		Usage: "Prefix applied to all tags",
	}
	wrapTagsSuffix cli.StringFlag = cli.StringFlag{
		Name:  "wrapTagsSuffix",
		Usage: "Suffix applied to all tags",
	}
)

func ValidateInFile(p string) error {
	// exists, valid ext
	_, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("in ValidateInFile(%v):\n%+w", p, err)
	}
	switch filepath.Ext(p) {
	case ".jsgf", ".jjsgf":
		return nil
	default:
		return fmt.Errorf("in ValidateInFile(%v):\n%+w", p, errors.New("file extension is not one of .jsgf, .jjsgf"))
	}
}

func ValidateOutFile(p string) error {
	// provided dir exists
	_, err := os.Stat(filepath.Dir(p))
	if err != nil {
		return fmt.Errorf("in ValidateOutFile(%v):\n%+w", p, err)
	}
	return nil
}

func ValidateExportDir(p string) error {
	// exists, is dir
	info, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("in ValidateExportDir(%v):\n%+w", p, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("in ValidateExportDir(%v):\n%+w", p, errors.New("provided path is not a directory"))
	}
	return nil
}

func FileScanner(p string) (*bufio.Scanner, error) {
	switch filepath.Ext(p) {
	case ".jsgf":
		f, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}
		return bufio.NewScanner(f), nil
	case ".jjsgf":
		f, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}
		var jj JJSGFGrammarJSON
		err = json.NewDecoder(f).Decode(&jj)
		if err != nil {
			log.Fatal(err)
		}
		return bufio.NewScanner(strings.NewReader(JJSGFToJSGF(jj))), nil
	}

	return &bufio.Scanner{}, fmt.Errorf("error when calling gsgf generate --inFile=%v:\n%+w", p, errors.New("invalid file extension, not one of .jsgf, .jjsgf"))
}

func ApplyPostproc(p []string, cmd *cli.Command) []string {
	if cmd.Bool("shuffle") {
		mrand.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
	}
	if cmd.String("wrapProductionsPrefix") != "" || cmd.String("wrapProductionsSuffix") != "" {
		p = WrapProductions(p, cmd.String("wrapProductionsPrefix"), cmd.String("wrapProductionsSuffix"))
	}
	if cmd.String("wrapTagsPrefix") != "" || cmd.String("wrapTagsSuffix") != "" {
		p = WrapProductions(p, cmd.String("wrapTagsPrefix"), cmd.String("wrapTagsSuffix"))
	}
	if cmd.String("collectTagsChar") != "" {
		p = CollectTags(p, cmd.String("collectTagsChar"))
	}
	if cmd.Bool("removeTags") {
		p = RemoveTags(p)
	}
	if cmd.Bool("removeMultiSpaces") {
		p = RemoveMultipleSpaces(p)
	}
	if cmd.Bool("removeEndSpaces") {
		p = RemoveEndSpaces(p)
	}
	if cmd.Bool("renderNewlines") {
		p = RenderNewLines(p)
	}
	if cmd.Bool("renderTabs") {
		p = RenderTabs(p)
	}
	return p
}

func BuildGrammar(cmd *cli.Command) (Grammar, error) {
	var (
		g   Grammar = NewGrammar()
		s   *bufio.Scanner
		err error
	)

	s, err = FileScanner(cmd.Args().First())
	if err != nil {
		log.Fatal(err)
	}
	lex := NewJSGFLexer(cmd.String("quoteChar"))
	g, err = FomJSGF(g, s, lex)
	if err != nil {
		log.Fatal(err)
	}
	err = ValidateGrammarCompleteness(g)
	if err != nil {
		namespace, err := CreateNameSpace(cmd.String("inFile"), cmd.String("ext"))
		if err != nil {
			log.Fatal(err)
		}
		g = ImportNameSpace(g, namespace, lex)
	}

	if cmd.Bool("minimize") {
		for k, v := range g.Rules {
			v.Graph = Minimize(v.Graph)
			g.Rules[k] = v
		}
	}

	g, err = ResolveRules(g, lex)
	if err != nil {
		log.Fatal(err)
	}

	if cmd.Bool("minimize") {
		for k, v := range g.Rules {
			if v.IsPublic {
				v.Graph = Minimize(v.Graph)
				g.Rules[k] = v
			}
		}
	}
	return g, nil
}

func PrepareContext(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	if cmd.Args().Get(0) == "" {
		cli.ShowSubcommandHelpAndExit(cmd, 0)
	}

	cmd.Set("inFile", cmd.Args().Get(0))
	cmd.Set("ext", filepath.Ext(cmd.Args().Get(0)))
	if cmd.Bool("singleQuote") {
		cmd.Set("quoteChar", "'")
	}
	if cmd.Int("nProductions") == -1 {
		cmd.Set("nProductions", "1")
	}

	return ctx, nil
}
