// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:03 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/bzick/tokenizer"
)

type Expression = string

func ToTokens(e Expression, lex *tokenizer.Tokenizer) []Expression {
	var (
		res string
		b   strings.Builder
		out = []Expression{"<SOS>"}
	)

	if e == "" {
		return []Expression{}
	}
	stream := lex.ParseString(e)
	for stream.IsValid() {
		switch {
		case stream.CurrentToken().Is(SquareOpen, SquareClose, ParenthesisOpen, ParenthesisClose, Alternate, Semicolon):
			b, out = FlushBuilder(b, out)
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, res)
			stream.GoNext()
		case stream.CurrentToken().Is(BackSlash):
			stream.GoNext()
			b.WriteString(stream.CurrentToken().ValueUnescapedString())
			stream.GoNext()
		case stream.CurrentToken().Is(ForwardSlash):
			stream.GoNext()
			b.WriteString("/")
			res, _ = CaptureString(stream, "/", true)
			b.WriteString(res)
			b, out = FlushBuilder(b, out)
			stream.GoNext()
		case stream.CurrentToken().Is(AngleOpen):
			b, out = FlushBuilder(b, out)
			res, _ = CaptureString(stream, ">", true)
			out = append(out, res)
			stream.GoNext()
		case stream.CurrentToken().Is(CurlyOpen):
			res, _ = CaptureString(stream, "}", true)
			b.WriteString(res)
			b, out = FlushBuilder(b, out)
			stream.GoNext()
		default:
			b.WriteString(stream.CurrentToken().ValueUnescapedString())
			stream.GoNext()
		}
	}
	b, out = FlushBuilder(b, out)
	out = append(out, "<EOS>")

	return out
}

func FlushBuilder(b strings.Builder, o []Expression) (strings.Builder, []Expression) {
	var s string = b.String()

	b.Reset()
	if s != "" {
		o = append(o, s)
	}

	return b, o
}

func IsWeighted(e Expression) bool {
	return regexp.MustCompile(`/[0-9\.]+/`).MatchString(e)
}

func ParseWeight(e Expression) (Expression, float64, error) {
	split := strings.Split(e, "/")
	if len(split) != 3 {
		return e, 0.0, errors.New("expression e not separable into expected 3 parts exp/weight/end. e may have the incorrect number of /")
	}
	f, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return e, 0.0, errors.New("unable to parse weight in expression e to float64")
	}

	return split[0], f, nil
}
