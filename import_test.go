// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:15 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestCreateNameSpace(t *testing.T) {
	dummy_error := errors.New("")
	table := []struct {
		d     string
		e     string
		rules map[string]string
		err   error
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
	for _, test := range table {
		rules, err := CreateNameSpace(test.d, test.e)

		if !reflect.DeepEqual(rules, test.rules) {
			t.Errorf("CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", test.d, test.e, rules, test.rules)
		}

		// for k, v_res := range rules {
		// 	v_exp, ok := test.rules[k]
		// 	if !ok {
		// 		t.Errorf("CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", test.d, test.e, v_res, v_exp)
		// 	}
		// 	for kk, vv_res := range v_exp {
		// 		vv_exp, ok := v_exp[kk]
		// 		if !ok {
		// 			t.Errorf("CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", test.d, test.e, vv_res, vv_exp)
		// 		}
		// 		sort.Strings(vv_exp)
		// 		sort.Strings(vv_res)
		// 		if fmt.Sprint(vv_exp) != fmt.Sprint(vv_res) {
		// 			t.Errorf("CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", test.d, test.e, vv_res, vv_exp)
		// 		}
		// 	}
		// }

		if test.err != nil && err == nil {
			t.Errorf("CreateNameSpace(%v, %v).err\nGOT %v\nEXP %v", test.d, test.e, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("CreateNameSpace(%v, %v).err\nGOT %v\nEXP %v", test.d, test.e, err, test.err)
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
	for _, test := range table {
		res, err := FindGrammar(test.p, test.t, test.e)
		if res != test.exp {
			t.Errorf("FindGrammar(%v, %v, %v)\nGOT %v\nEXP %v", test.p, test.t, test.e, res, test.exp)
		}
		if test.err != nil && err == nil {
			t.Errorf("FindGrammar(%v, %v, %v).err\nGOT %v\nEXP %v", test.p, test.t, test.e, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("FindGrammar(%v, %v, %v).err\nGOT %v\nEXP %v", test.p, test.t, test.e, err, test.err)
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
	for _, test := range table {
		res, err := ImportOrder(test.p, test.e)
		sort.Strings(test.exp)
		sort.Strings(res)

		if fmt.Sprint(res) != fmt.Sprint(test.exp) {
			t.Errorf("ImportOrder(%v, %v)\nGOT %v\nEXP %v", test.p, test.e, res, test.exp)
		}
		if test.err != nil && err == nil {
			t.Errorf("ImportOrder(%v, %v).err\nGOT %v\nEXP %v", test.p, test.e, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("ImportOrder(%v, %v).err\nGOT %v\nEXP %v", test.p, test.e, err, test.err)
		}

	}
}
