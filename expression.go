// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:03 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bzick/tokenizer"
)

// Type alias for rule definitions/tokens
type Expression = string

// Splits an expression into a slice of expression/tokens using tokenizer
func ToTokens(e Expression, lex *tokenizer.Tokenizer) []Expression {
	if e == "" {
		return []Expression{}
	}

	var (
		res     string
		builder strings.Builder
		out     []Expression      = []Expression{"<SOS>"}
		stream  *tokenizer.Stream = lex.ParseString(e)
	)

	for stream.IsValid() {
		switch {
		case stream.CurrentToken().Is(SquareOpen, SquareClose, ParenthesisOpen, ParenthesisClose, Alternate, Semicolon):
			builder, out = flushBuilder(builder, out)
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, res)
			stream.GoNext()
		case stream.CurrentToken().Is(BackSlash):
			stream.GoNext()
			builder.WriteString(stream.CurrentToken().ValueUnescapedString())
			stream.GoNext()
		case stream.CurrentToken().Is(ForwardSlash):
			stream.GoNext()
			builder.WriteString("/")
			res, _ = captureString(stream, "/", true)
			builder.WriteString(res)
			builder, out = flushBuilder(builder, out)
			stream.GoNext()
		case stream.CurrentToken().Is(AngleOpen):
			builder, out = flushBuilder(builder, out)
			res, _ = captureString(stream, ">", true)
			out = append(out, res)
			stream.GoNext()
		case stream.CurrentToken().Is(CurlyOpen):
			res, _ = captureString(stream, "}", true)
			builder.WriteString(res)
			builder, out = flushBuilder(builder, out)
			stream.GoNext()
		default:
			builder.WriteString(stream.CurrentToken().ValueUnescapedString())
			stream.GoNext()
		}
	}
	builder, out = flushBuilder(builder, out)
	out = append(out, "<EOS>")

	return out
}

// Helper function to get current contents of strings.Builder and reset
func flushBuilder(b strings.Builder, o []Expression) (strings.Builder, []Expression) {
	var str string = b.String()

	b.Reset()
	if str != "" {
		o = append(o, str)
	}

	return b, o
}

// Check if an espression has a weight defined by /[0-9\.]+/
func isWeighted(e Expression) bool {
	return regexp.MustCompile(`/[0-9\.]+/`).MatchString(e)
}

// Splits a weighted expression into the base expression and its weight value
func ParseWeight(e Expression) (Expression, float64, error) {
	split := strings.Split(e, "/")
	if len(split) != 3 {
		return e, 0.0, fmt.Errorf("error when calling ParseWeight(%v), split into %v:\n%+w", e, split, errors.New("expression e not separable into expected 3 parts exp/weight/end. e may have the incorrect number of /"))
	}
	weight, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return e, 0.0, fmt.Errorf("error when calling ParseWeight(%v), strconv.ParseFloat(%v, 64):\n%+w", e, split[1], err)
	}

	return split[0], weight, nil
}
