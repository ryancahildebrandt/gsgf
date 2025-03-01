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
		res     string
		builder strings.Builder
		out     = []Expression{"<SOS>"}
	)

	if e == "" {
		return []Expression{}
	}
	stream := lex.ParseString(e)
	for stream.IsValid() {
		switch {
		case stream.CurrentToken().Is(SquareOpen, SquareClose, ParenthesisOpen, ParenthesisClose, Alternate, Semicolon):
			builder, out = FlushBuilder(builder, out)
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
			res, _ = CaptureString(stream, "/", true)
			builder.WriteString(res)
			builder, out = FlushBuilder(builder, out)
			stream.GoNext()
		case stream.CurrentToken().Is(AngleOpen):
			builder, out = FlushBuilder(builder, out)
			res, _ = CaptureString(stream, ">", true)
			out = append(out, res)
			stream.GoNext()
		case stream.CurrentToken().Is(CurlyOpen):
			res, _ = CaptureString(stream, "}", true)
			builder.WriteString(res)
			builder, out = FlushBuilder(builder, out)
			stream.GoNext()
		default:
			builder.WriteString(stream.CurrentToken().ValueUnescapedString())
			stream.GoNext()
		}
	}
	builder, out = FlushBuilder(builder, out)
	out = append(out, "<EOS>")

	return out
}

func FlushBuilder(b strings.Builder, o []Expression) (strings.Builder, []Expression) {
	var str string = b.String()

	b.Reset()
	if str != "" {
		o = append(o, str)
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
	weight, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return e, 0.0, errors.New("unable to parse weight in expression e to float64")
	}

	return split[0], weight, nil
}
