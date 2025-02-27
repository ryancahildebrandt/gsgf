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
	dummy_error := errors.New("")
	table := []struct {
		d   string
		e   string
		r   map[string]string
		err error
	}{
		{"data/tests/test0.jsgf",
			".jsgf",
			map[string]string{
				"<order>":   "i'd like [to order|a|<quant>];",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<quant>":   "some|a (cup|glass) of;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
			nil,
		},
		{"data/tests/test1.jsgf",
			".jsgf",
			map[string]string{
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			nil,
		},
		{"data/tests/test2.jsgf",
			".jsgf",
			map[string]string{},
			dummy_error,
		},
		{"data/tests/test3.jsgf",
			".jsgf",
			map[string]string{
				"<request>": "[(could|will|would) you] please <brew>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
			},
			nil,
		},
		{"data/tests/test4.jsgf",
			".jsgf",
			map[string]string{
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
			nil,
		},
		{"data/tests/test5.jsgf",
			".jsgf",
			map[string]string{
				"<request>": "[(could|will|would) you] please <brew>;",
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			nil,
		},
		{"data/tests/test6.jsgf",
			".jsgf",
			map[string]string{
				"<request>": "[(could|will|would) you] please <brew>;",
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			nil,
		},

		{"data/tests/a.jsgf",
			".jsgf",
			map[string]string{},
			nil,
		}, {"data/tests/b.jsgf",
			".jsgf",
			map[string]string{
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			nil,
		},
		{"data/tests/dir0/c.jsgf",
			".jsgf",
			map[string]string{},
			nil,
		},
		{"data/tests/dir0/dir1/d.jsgf",
			".jsgf",
			map[string]string{},
			dummy_error,
		},
		{"data/tests/dir0/dir1/dir2/e.jsgf",
			".jsgf",
			map[string]string{},
			nil,
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
	dummy_error := errors.New("")
	table := []struct {
		p   string
		t   string
		e   string
		exp string
		err error
	}{
		{"./data/tests", "test0", ".jsgf", "data/tests/test0.jsgf", nil},
		{"./data/tests/.jsgf", "test0", ".jsgf", "data/tests/test0.jsgf", nil},
		{"./data/tests/test0.jsgf", "test0", ".jsgf", "data/tests/test0.jsgf", nil},
		{"./data/tests/test0.jsgf", "a", ".jsgf", "data/tests/a.jsgf", nil},
		{"./data/tests/test0.jsgf", "e", ".jsgf", "data/tests/dir0/dir1/dir2/e.jsgf", nil},
		{"./data/tests/a.jsgf", "a", ".jsgf", "data/tests/a.jsgf", nil},
		{"./data/tests/dir0/dir1/c.jsgf", "d", ".jsgf", "data/tests/dir0/dir1/d.jsgf", nil},
		{"./data/tests/dir0/dir1/c.jsgf", "e", ".jsgf", "data/tests/dir0/dir1/dir2/e.jsgf", nil},
		{"./data/tests/dir0/dir1/dir2/e.jsgf", "e", ".jsgf", "data/tests/dir0/dir1/dir2/e.jsgf", nil},
		{"./data/tests/dir0/dir1/c.jsgf", "b", ".jsgf", "", dummy_error},
		{"./data/tests/test0.jsgf", "f", ".jsgf", "", dummy_error},
		{"./data/tests/dir0/dir1/dir2/e.jsgf", "d", ".jsgf", "", dummy_error},
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
	dummy_error := errors.New("")
	table := []struct {
		p   string
		e   string
		exp []string
		err error
	}{
		{"./data/tests", ".jsgf", []string{}, dummy_error},
		{"./data/tests/.jsgf", ".jsgf", []string{}, dummy_error},
		{"./data/tests/test0.jsgf", ".jsgf", []string{"import <a.*>"}, nil},
		{"./data/tests/test1.jsgf", ".jsgf", []string{"import <c.brew>"}, nil},
		{"./data/tests/test3.jsgf", ".jsgf", []string{"import <e.dne>"}, nil},
		{"./data/tests/test4.jsgf", ".jsgf", []string{"import <a.order>", "import <c.teatype>", "import <d.*>"}, nil},
		{"./data/tests/test5.jsgf", ".jsgf", []string{"import <b.request>", "import <c.brew>"}, nil},
		{"./data/tests/a.jsgf", ".jsgf", []string{}, nil},
		{"./data/tests/b.jsgf", ".jsgf", []string{"import <c.brew>"}, nil},
		{"./data/tests/dir0/c.jsgf", ".jsgf", []string{}, nil},
		{"./data/tests/dir0/dir1/dir2/e.jsgf", ".jsgf", []string{}, nil},

		{"./data/tests/test2.jsgf", ".jsgf", []string{}, dummy_error},
		{"./data/tests/dir0/dir1/d.jsgf", ".jsgf", []string{}, dummy_error},
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
