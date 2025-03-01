// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:34 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"
	"testing"
)

func TestBuildEdgeList(t *testing.T) {
	dummyError := errors.New("")
	lexer := NewJSGFLexer()
	table := []struct {
		r   string
		exp EdgeList
		err error
	}{
		{
			r:   "",
			exp: EdgeList{},
			err: dummyError,
		},
		{
			r:   "=",
			exp: EdgeList{},
			err: dummyError,
		},
		{
			r:   "<>=;",
			exp: EdgeList{},
			err: dummyError,
		},
		{
			r:   "public <test> = ;",
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			err: nil,
		},
		{
			r:   "public <test> = one two three;",
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			err: nil,
		},
		{
			r: "public <test> = four|five|six;",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = six[ seven][ eight];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = eight( nine)( ten);",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [12|13|14];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 (12|13|14);",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 (((12)));",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 ((12)(13)(14));",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 ((12)|(13)|(14));",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 ([[12]]);",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 ([12][13][14]);",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 ([12]|[13]|[14]);",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [((12))];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [(12)(13)(14)];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [(12)|(13)|(14)];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [[[12]]];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [[12][13][14]];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [[12]|[13]|[14]];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 ((12)|(13)[14]) 15;",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 13, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 10, To: 12, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = 11 [(12)|[13]14][15];",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 11, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 12, To: 14, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			err: nil,
		},
		{
			r: "public <test> = [(11)12[13](14)] 15;",
			exp: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 1, To: 12, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			},
			err: nil,
		},
	}
	for i, test := range table {
		_, v, err := ParseRule(test.r, lexer)
		res := BuildEdgeList(ToTokens(v.Exp, lexer))
		if !slices.Equal(Sort(res), Sort(test.exp)) {
			t.Errorf("test %v: %v.toArray(lexer)\nGOT %v\nEXP %v", i, test.r, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.toArray(lexer).err\nGOT %v\nEXP %v", i, test.r, err, test.err)
		}
	}
}

func TestEdgeListSort(t *testing.T) {
	table := []struct {
		e   EdgeList
		exp EdgeList
	}{
		{e: EdgeList{}, exp: EdgeList{}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, exp: EdgeList{{From: 0, To: 1, Weight: 1.0}}},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e:   EdgeList{{From: 2, To: 3, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
	}
	for i, test := range table {
		res := Sort(test.e)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.Sort()\nGOT %v\nEXP %v", i, test.e, res, test.exp)
		}
	}
}

func TestEdgeListUnique(t *testing.T) {
	table := []struct {
		e   EdgeList
		exp EdgeList
	}{
		{e: EdgeList{}, exp: EdgeList{}},
		{e: EdgeList{{From: 0, To: 1, Weight: 1.0}}, exp: EdgeList{{From: 0, To: 1, Weight: 1.0}}},
		{
			e:   EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e:   EdgeList{{From: 2, To: 3, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}},
			exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 1, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
			}, exp: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
	}
	for i, test := range table {
		res := Unique(test.e)
		if !slices.Equal(Sort(res), Sort(test.exp)) {
			t.Errorf("test %v: %v.Unique()\nGOT %v\nEXP %v", i, test.e, res, test.exp)
		}
	}
}

func TestEdgeListMax(t *testing.T) {
	table := []struct {
		e   EdgeList
		exp int
	}{
		{e: EdgeList{}, exp: 0},
		{e: EdgeList{{From: 0, To: 0, Weight: 1.0}}, exp: 0},
		{e: EdgeList{{From: 1, To: 1, Weight: 1.0}, {From: 1, To: 1, Weight: 1.0}}, exp: 1},
		{
			e: EdgeList{
				{From: 10, To: 12, Weight: 1.0},
				{From: 55, To: 94, Weight: 1.0},
				{From: 0, To: 15, Weight: 1.0},
				{From: 1, To: 1, Weight: 1.0},
			}, exp: 94,
		},
		{
			e: EdgeList{
				{From: -1, To: 1651, Weight: 1.0},
				{From: 55, To: 65, Weight: 1.0},
				{From: 10, To: 1, Weight: 1.0},
				{From: 15, To: 99, Weight: 1.0},
				{From: 65, To: 54, Weight: 1.0},
				{From: 1000000000, To: 0, Weight: 1.0},
				{From: 0, To: 8, Weight: 1.0},
				{From: 15, To: 44, Weight: 1.0},
			}, exp: 1000000000,
		},
	}
	for i, test := range table {
		res := test.e.Max()
		if res != test.exp {
			t.Errorf("test %v: %v.Max()\nGOT %v\nEXP %v", i, test.e, res, test.exp)
		}
	}
}

func TestEdgeListIncrement(t *testing.T) {
	table := []struct {
		e   EdgeList
		n   int
		exp EdgeList
	}{
		{e: EdgeList{}, n: 0, exp: EdgeList{}},
		{e: EdgeList{{From: 0, To: 0, Weight: 1.0}}, n: 0, exp: EdgeList{{From: 0, To: 0, Weight: 1.0}}},
		{
			e: EdgeList{{From: 1, To: 1, Weight: 1.0}, {From: 1, To: 1, Weight: 1.0}}, n: 1,
			exp: EdgeList{{From: 2, To: 2, Weight: 1.0}, {From: 2, To: 2, Weight: 1.0}},
		},
		{
			e: EdgeList{
				{From: 10, To: 12, Weight: 1.0},
				{From: 55, To: 94, Weight: 1.0},
				{From: 0, To: 15, Weight: 1.0},
				{From: 1, To: 1, Weight: 1.0},
			}, n: -1, exp: EdgeList{
				{From: 9, To: 11, Weight: 1.0},
				{From: 54, To: 93, Weight: 1.0},
				{From: -1, To: 14, Weight: 1.0},
				{From: 0, To: 0, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: -1, To: 1651, Weight: 1.0},
				{From: 55, To: 65, Weight: 1.0},
				{From: 10, To: 1, Weight: 1.0},
				{From: 15, To: 99, Weight: 1.0},
				{From: 65, To: 54, Weight: 1.0},
				{From: 1000000000, To: 0, Weight: 1.0},
				{From: 0, To: 8, Weight: 1.0},
				{From: 15, To: 44, Weight: 1.0},
			}, n: 10, exp: EdgeList{
				{From: 9, To: 1661, Weight: 1.0},
				{From: 65, To: 75, Weight: 1.0},
				{From: 20, To: 11, Weight: 1.0},
				{From: 25, To: 109, Weight: 1.0},
				{From: 75, To: 64, Weight: 1.0},
				{From: 1000000010, To: 10, Weight: 1.0},
				{From: 10, To: 18, Weight: 1.0},
				{From: 25, To: 54, Weight: 1.0},
			},
		},
	}
	for i, test := range table {
		res := Increment(test.e, test.n)
		if !slices.Equal(Sort(res), Sort(test.exp)) {
			t.Errorf("test %v: %v,Increment(%v)\nGOT %v\nEXP %v", i, test.e, test.n, res, test.exp)
		}
	}
}
