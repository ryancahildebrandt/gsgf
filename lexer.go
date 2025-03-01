// -*- coding: utf-8 -*-

// Created on Sat Aug 24 01:56:01 PM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"strings"

	"github.com/bzick/tokenizer"
)

const (
	AngleOpen = iota + 1
	AngleClose
	SquareOpen
	SquareClose
	ParenthesisOpen
	ParenthesisClose
	CurlyOpen
	CurlyClose
	Alternate
	Comment
	Assignment
	Semicolon
	Modifier
	SequenceStart
	SequenceEnd
	DoubleQuote
	BackSlash
	ForwardSlash
)

func NewJSGFLexer() *tokenizer.Tokenizer {
	var lexer *tokenizer.Tokenizer = tokenizer.New()

	lexer.SetWhiteSpaces([]byte{})
	lexer.DefineStringToken(DoubleQuote, `"`, `"`).SetEscapeSymbol(BackSlash).AddSpecialStrings(tokenizer.DefaultSpecialString)
	// lexer.DefineStringToken(SingleQuote, `'`, `'`).SetSpecialSymbols(tokenizer.DefaultStringEscapes).SetEscapeSymbol(BackSlash)

	lexer.DefineTokens(AngleOpen, []string{"<"})
	lexer.DefineTokens(AngleClose, []string{">"})
	lexer.DefineTokens(CurlyOpen, []string{"{"})
	lexer.DefineTokens(CurlyClose, []string{"}"})
	lexer.DefineTokens(SquareOpen, []string{"["})
	lexer.DefineTokens(SquareClose, []string{"]"})
	lexer.DefineTokens(ParenthesisOpen, []string{"("})
	lexer.DefineTokens(ParenthesisClose, []string{")"})
	lexer.DefineTokens(Alternate, []string{"|"})
	lexer.DefineTokens(Comment, []string{"//", "/*", "*/"})
	lexer.DefineTokens(Semicolon, []string{";", " ;", "; ", " ; "})
	lexer.DefineTokens(Assignment, []string{"="})
	lexer.DefineTokens(Modifier, []string{"*", "+"})
	lexer.DefineTokens(SequenceStart, []string{"<SOS>"})
	lexer.DefineTokens(SequenceEnd, []string{"<EOS>"})
	lexer.DefineTokens(BackSlash, []string{"\\"})
	lexer.DefineTokens(ForwardSlash, []string{"/"})

	return lexer
}

func CaptureString(s *tokenizer.Stream, end string, includeEnd bool) (string, error) {
	var builder strings.Builder
	var remainder string = s.GetSnippetAsString(0, 1000000, 0)
	// really high value here because s doesn't have a "show me whats left in the string" method

	if !strings.Contains(remainder, end) {
		return "", errors.New("close token not found in remaining string")
	}
	if remainder == "" {
		return "", errors.New("cannot capture from empty string")
	}

	for s.IsValid() {
		if s.CurrentToken().ValueUnescapedString() == end {
			if includeEnd {
				builder.WriteString(s.CurrentToken().ValueUnescapedString())
			}

			break
		}
		builder.WriteString(s.CurrentToken().ValueUnescapedString())
		s.GoNext()
	}

	return builder.String(), nil
}
