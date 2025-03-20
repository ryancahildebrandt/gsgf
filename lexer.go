// -*- coding: utf-8 -*-

// Created on Sat Aug 24 01:56:01 PM EDT 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/bzick/tokenizer"
)

// Token types significant for the jsgf specification and jsgfLexer implemented here
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
	SequenceStart
	SequenceEnd
	DoubleQuote
	SingleQuote
	BackSlash
	ForwardSlash
)

// Tokens that can be ignored during graph traversal and production collection
var jsgfFilter []string = []string{"(", ")", "[", "]", "<SOS>", ";", "|", "<EOS>", ""}

// Returns a tokenizer for jsgf files with the specified quote token
func NewJSGFLexer(q string) *tokenizer.Tokenizer {
	var lexer *tokenizer.Tokenizer = tokenizer.New()

	if q != "" {
		lexer.DefineStringToken(DoubleQuote, `"`, `"`).SetEscapeSymbol(BackSlash).AddSpecialStrings(tokenizer.DefaultSpecialString)
	}

	lexer.SetWhiteSpaces([]byte{})

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
	lexer.DefineTokens(Semicolon, []string{";"})
	lexer.DefineTokens(Assignment, []string{"="})
	lexer.DefineTokens(SequenceStart, []string{"<SOS>"})
	lexer.DefineTokens(SequenceEnd, []string{"<EOS>"})
	lexer.DefineTokens(BackSlash, []string{"\\"})
	lexer.DefineTokens(ForwardSlash, []string{"/"})

	return lexer
}

// Returns a string beginning from s.CurrentToken and ending at the first occurrence of the ending string, optionally including the end token in the returned string
func captureString(s *tokenizer.Stream, end string, includeEnd bool) (string, error) {
	var builder strings.Builder
	var remainder string = s.GetSnippetAsString(0, 1000000, 0) // really high value here because s doesn't have a "show me whats left in the string" method
	if !strings.Contains(remainder, end) {
		return "", fmt.Errorf("error when calling CaptureString(%v, %v, %v), remainder %v:\n%+w", s, end, includeEnd, remainder, errors.New("close token not found in remaining string"))
	}

	if remainder == "" {
		return "", fmt.Errorf("error when calling CaptureString(%v, %v, %v), remainder %v:\n%+w", s, end, includeEnd, remainder, errors.New("cannot capture from empty string"))
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

// Checks that the provided string can be consumed by a tokenizer (is not empty and does not contain byte \x00)
func ValidateLexerString(s string) error {
	if s == "" {
		return fmt.Errorf("error when calling ValidateLexerString(%v):\n%+w", s, errors.New("cannot tokenize empty string"))
	}
	if bytes.Contains([]byte(s), []byte("\x00")) {
		return fmt.Errorf("error when calling ValidateLexerString(%v):\n%+w", s, errors.New("cannot tokenize string containing null char \x00"))
	}

	return nil
}
