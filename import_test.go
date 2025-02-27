// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:15 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"
	"sort"
	"testing"
)

func TestCreateNameSpace(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		d   string
		e   string
		r   map[string]string
		err error
	}{
		{
			d: "data/tests/test0.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<order>":   "i'd like [to order|a|<quant>];",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<quant>":   "some|a (cup|glass) of;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
			err: nil,
		},
		{
			d: "data/tests/test1.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			err: nil,
		},
		{
			d:   "data/tests/test2.jsgf",
			e:   ".jsgf",
			r:   map[string]string{},
			err: dummyError,
		},
		{
			d: "data/tests/test3.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<request>": "[(could|will|would) you] please <brew>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
			},
			err: nil,
		},
		{
			d: "data/tests/test4.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
			err: nil,
		},
		{
			d: "data/tests/test5.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<request>": "[(could|will|would) you] please <brew>;",
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			err: nil,
		},
		{
			d: "data/tests/test6.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<request>": "[(could|will|would) you] please <brew>;",
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			err: nil,
		},

		{
			d:   "data/tests/a.jsgf",
			e:   ".jsgf",
			r:   map[string]string{},
			err: nil,
		}, {
			d: "data/tests/b.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			err: nil,
		},
		{
			d:   "data/tests/dir0/c.jsgf",
			e:   ".jsgf",
			r:   map[string]string{},
			err: nil,
		},
		{
			d:   "data/tests/dir0/dir1/d.jsgf",
			e:   ".jsgf",
			r:   map[string]string{},
			err: dummyError,
		},
		{
			d:   "data/tests/dir0/dir1/dir2/e.jsgf",
			e:   ".jsgf",
			r:   map[string]string{},
			err: nil,
		},
	}
	for i, test := range table {
		rules, err := CreateNameSpace(test.d, test.e)
		if len(rules) != len(test.r) {
			t.Errorf("test %v: CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", i, test.d, test.e, rules, test.r)
		}
		for k1, v1 := range rules {
			v2, ok := test.r[k1]
			if !ok {
				t.Errorf("test %v: CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", i, test.d, test.e, v1, v2)
			}
			if v1 != v2 {
				t.Errorf("test %v: CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", i, test.d, test.e, v1, v2)
			}
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: CreateNameSpace(%v, %v).err\nGOT %v\nEXP %v", i, test.d, test.e, err, test.err)
		}
	}
}

func TestFindGrammar(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		p   string
		t   string
		e   string
		exp string
		err error
	}{
		{p: "./data/tests", t: "test0", e: ".jsgf", exp: "data/tests/test0.jsgf", err: nil},
		{p: "./data/tests/.jsgf", t: "test0", e: ".jsgf", exp: "data/tests/test0.jsgf", err: nil},
		{p: "./data/tests/test0.jsgf", t: "test0", e: ".jsgf", exp: "data/tests/test0.jsgf", err: nil},
		{p: "./data/tests/test0.jsgf", t: "a", e: ".jsgf", exp: "data/tests/a.jsgf", err: nil},
		{p: "./data/tests/test0.jsgf", t: "e", e: ".jsgf", exp: "data/tests/dir0/dir1/dir2/e.jsgf", err: nil},
		{p: "./data/tests/a.jsgf", t: "a", e: ".jsgf", exp: "data/tests/a.jsgf", err: nil},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "d", e: ".jsgf", exp: "data/tests/dir0/dir1/d.jsgf", err: nil},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "e", e: ".jsgf", exp: "data/tests/dir0/dir1/dir2/e.jsgf", err: nil},
		{
			p: "./data/tests/dir0/dir1/dir2/e.jsgf", t: "e", e: ".jsgf", exp: "data/tests/dir0/dir1/dir2/e.jsgf",
			err: nil,
		},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "b", e: ".jsgf", exp: "", err: dummyError},
		{p: "./data/tests/test0.jsgf", t: "f", e: ".jsgf", exp: "", err: dummyError},
		{p: "./data/tests/dir0/dir1/dir2/e.jsgf", t: "d", e: ".jsgf", exp: "", err: dummyError},
	}
	for i, test := range table {
		res, err := FindGrammar(test.p, test.t, test.e)
		if res != test.exp {
			t.Errorf("test %v: FindGrammar(%v, %v, %v)\nGOT %v\nEXP %v", i, test.p, test.t, test.e, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: FindGrammar(%v, %v, %v).err\nGOT %v\nEXP %v", i, test.p, test.t, test.e, err, test.err)
		}
	}
}

func TestImportOrder(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		p   string
		e   string
		exp []string
		err error
	}{
		{p: "./data/tests", e: ".jsgf", exp: []string{}, err: dummyError},
		{p: "./data/tests/.jsgf", e: ".jsgf", exp: []string{}, err: dummyError},
		{p: "./data/tests/test0.jsgf", e: ".jsgf", exp: []string{"import <a.*>"}, err: nil},
		{p: "./data/tests/test1.jsgf", e: ".jsgf", exp: []string{"import <c.brew>"}, err: nil},
		{p: "./data/tests/test3.jsgf", e: ".jsgf", exp: []string{"import <e.dne>"}, err: nil},
		{
			p: "./data/tests/test4.jsgf", e: ".jsgf",
			exp: []string{"import <a.order>", "import <c.teatype>", "import <d.*>"}, err: nil,
		},
		{p: "./data/tests/test5.jsgf", e: ".jsgf", exp: []string{"import <b.request>", "import <c.brew>"}, err: nil},
		{p: "./data/tests/a.jsgf", e: ".jsgf", exp: []string{}, err: nil},
		{p: "./data/tests/b.jsgf", e: ".jsgf", exp: []string{"import <c.brew>"}, err: nil},
		{p: "./data/tests/dir0/c.jsgf", e: ".jsgf", exp: []string{}, err: nil},
		{p: "./data/tests/dir0/dir1/dir2/e.jsgf", e: ".jsgf", exp: []string{}, err: nil},

		{p: "./data/tests/test2.jsgf", e: ".jsgf", exp: []string{}, err: dummyError},
		{p: "./data/tests/dir0/dir1/d.jsgf", e: ".jsgf", exp: []string{}, err: dummyError},
	}
	for i, test := range table {
		res, err := ImportOrder(test.p, test.e)
		sort.Strings(test.exp)
		sort.Strings(res)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: ImportOrder(%v, %v)\nGOT %v\nEXP %v", i, test.p, test.e, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ImportOrder(%v, %v).err\nGOT %v\nEXP %v", i, test.p, test.e, err, test.err)
		}
	}
}
