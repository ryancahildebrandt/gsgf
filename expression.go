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
	var res string
	var b strings.Builder
	out := []Expression{"<SOS>"}
	if e == "" {
		return []Expression{}
	}
	stream := lex.ParseString(e)
	for stream.IsValid() {
		switch {
		case stream.CurrentToken().Is(AngleOpen):
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			res, _ = captureString(stream, ">", true)
			out = append(out, Expression(res))
			stream.GoNext()
		case stream.CurrentToken().Is(SquareOpen):
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, Expression(res))
			stream.GoNext()
		case stream.CurrentToken().Is(SquareClose):
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, Expression(res))
			stream.GoNext()
		case stream.CurrentToken().Is(ParenthesisOpen):
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, Expression(res))
			stream.GoNext()
		case stream.CurrentToken().Is(ParenthesisClose):
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, Expression(res))
			stream.GoNext()
		case stream.CurrentToken().Is(Alternate):
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, Expression(res))
			stream.GoNext()
		case stream.CurrentToken().Is(Semicolon):
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			res = stream.CurrentToken().ValueUnescapedString()
			out = append(out, Expression(res))
			stream.GoNext()
		case stream.CurrentToken().Is(BackSlash):
			stream.GoNext()
			b.WriteString(stream.CurrentToken().ValueUnescapedString())
			stream.GoNext()
		case stream.CurrentToken().Is(ForwardSlash):
			stream.GoNext()
			b.WriteString("/")
			res, _ = captureString(stream, "/", true)
			b.WriteString(res)
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			stream.GoNext()
		case stream.CurrentToken().Is(CurlyOpen):
			res, _ = captureString(stream, "}", true)
			b.WriteString(res)
			res = b.String()
			b.Reset()
			if res != "" {
				out = append(out, Expression(res))
			}
			stream.GoNext()
		default:
			b.WriteString(stream.CurrentToken().ValueUnescapedString())
			stream.GoNext()
		}
	}
	res = b.String()
	if res != "" {
		out = append(out, Expression(res))
	}
	out = append(out, Expression("<EOS>"))

	return out
}

func IsWeighted(e Expression) bool {
	return regexp.MustCompile(`/[0-9\.]+/`).MatchString(e)
}

func ParseWeight(e Expression) (Expression, float64, error) {
	if strings.Contains(e, "//") {
		return e, 0.0, errors.New("empty weight // in Expression e")
	}
	split := strings.Split(e, "/")
	if len(split) != 3 {
		return e, 0.0, errors.New("expression e not separable into expected 3 parts exp/weight/end. e may have the incorrect number of /")
	}
	f, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return e, 0.0, errors.New("unable to parse weight in expression e to float64")
	}

	return Expression(split[0]), f, nil
}
