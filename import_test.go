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
			d:   "data/tests/test6.jsgf",
			e:   ".jsgf",
			r:   map[string]string{},
			err: dummyError,
		},
		{
			d:   "data/tests/a.jsgf",
			e:   ".jsgf",
			r:   map[string]string{},
			err: nil,
		},
		{
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
		{p: "./data/tests/.jsgf", t: "test0", e: ".jsgf", exp: "data/tests/test0.jsgf", err: nil},
		{p: "./data/tests/test0.jsgf", t: "test0", e: ".jsgf", exp: "data/tests/test0.jsgf", err: nil},
		{p: "./data/tests/test0.jsgf", t: "a", e: ".jsgf", exp: "data/tests/a.jsgf", err: nil},
		{p: "./data/tests/test0.jsgf", t: "e", e: ".jsgf", exp: "data/tests/dir0/dir1/dir2/e.jsgf", err: nil},
		{p: "./data/tests/a.jsgf", t: "a", e: ".jsgf", exp: "data/tests/a.jsgf", err: nil},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "d", e: ".jsgf", exp: "data/tests/dir0/dir1/d.jsgf", err: nil},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "e", e: ".jsgf", exp: "data/tests/dir0/dir1/dir2/e.jsgf", err: nil},
		{p: "./data/tests/dir0/dir1/dir2/e.jsgf", t: "e", e: ".jsgf", exp: "data/tests/dir0/dir1/dir2/e.jsgf", err: nil},

		{p: "./data/tests", t: "test0", e: ".jsgf", exp: "", err: dummyError},
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
		{p: "./data/tests/test0.jsgf", e: ".jsgf", exp: []string{"import <a.*>;"}, err: nil},
		{p: "./data/tests/test1.jsgf", e: ".jsgf", exp: []string{"import <c.brew>;"}, err: nil},
		{p: "./data/tests/test3.jsgf", e: ".jsgf", exp: []string{"import <e.dne>;"}, err: nil},
		{
			p: "./data/tests/test4.jsgf", e: ".jsgf",
			exp: []string{"import <a.order>;", "import <c.teatype>;", "import <d.*>;"}, err: nil,
		},
		{p: "./data/tests/test5.jsgf", e: ".jsgf", exp: []string{"import <b.request>;", "import <c.brew>;"}, err: nil},
		{p: "./data/tests/a.jsgf", e: ".jsgf", exp: []string{}, err: nil},
		{p: "./data/tests/b.jsgf", e: ".jsgf", exp: []string{"import <c.brew>;"}, err: nil},
		{p: "./data/tests/dir0/c.jsgf", e: ".jsgf", exp: []string{}, err: nil},
		{p: "./data/tests/dir0/dir1/dir2/e.jsgf", e: ".jsgf", exp: []string{}, err: nil},
		{p: "./data/tests/test2.jsgf", e: ".jsgf", exp: []string{}, err: dummyError},
		{p: "./data/tests/dir0/dir1/d.jsgf", e: ".jsgf", exp: []string{}, err: dummyError},
	}
	for i, test := range table {
		res, err := GetImportOrder(test.p, test.e)
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

func TestPeekGrammar(t *testing.T) {
	table := []struct {
		p       string
		n       string
		imports []string
		rules   map[string][]string
	}{
		{
			p: "data/tests/test0.jsgf", n: "test0", imports: []string{"import <a.*>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "quant": {}, "teatype": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test1.jsgf", n: "test1", imports: []string{"import <c.brew>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {},
			},
		},
		{
			p: "data/tests/test2.jsgf", n: "test2", imports: []string{"import <a1.*>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test3.jsgf", n: "test3", imports: []string{"import <e.dne>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test4.jsgf", n: "test4", imports: []string{"import <d.*>;"}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "quant": {}, "brew": {"quant"},
			},
		},
		{
			p: "data/tests/test5.jsgf", n: "test5", imports: []string{"import <b.request>;"},
			rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "order": {"quant"}, "quant": {}, "teatype": {},
				"brew": {"quant"},
			},
		},
		{
			p: "data/tests/a.jsgf", n: "a", imports: []string{},
			rules: map[string][]string{"request": {"brew"}, "order": {"quant"}, "brew": {"quant"}, "quant": {}},
		},
		{
			p: "data/tests/b.jsgf", n: "b", imports: []string{"import <c.brew>;"},
			rules: map[string][]string{"request": {"brew"}, "order": {"quant"}, "quant": {}},
		},
		{
			p: "data/tests/dir0/c.jsgf", n: "c", imports: []string{},
			rules: map[string][]string{"teatype": {}, "brew": {"quant"}, "quant": {}},
		},
		{
			p: "data/tests/dir0/dir1/d.jsgf", n: "d", imports: []string{"import <c.teatype>;", "import <a.order>;"},
			rules: map[string][]string{},
		},
		{
			p: "data/tests/dir0/dir1/dir2/e.jsgf", n: "e", imports: []string{}, rules: map[string][]string{
				"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {},
				"teatype": {}, "brew": {"quant"},
			},
		},
	}
	for i, test := range table {
		name, imports, rules, err := PeekGrammar(test.p)
		if err != nil {
			t.Errorf("test %v: Grammar(%v).Peek()\nGOT error %v", i, test.p, err)
		}
		if name != test.n {
			t.Errorf("test %v: Grammar(%v).Peek().imports\nGOT %v\nEXP %v", i, test.p, name, test.n)
		}
		sort.Strings(imports)
		sort.Strings(test.imports)
		if !slices.Equal(imports, test.imports) {
			t.Errorf("test %v: Grammar(%v).Peek().imports\nGOT %v\nEXP %v", i, test.p, imports, test.imports)
		}
		if len(rules) != len(test.rules) {
			t.Errorf("test %v: Grammar(%v).Peek().rules\nGOT %v\nEXP %v", i, test.p, rules, test.rules)
		}
		// for k1, v1 := range rules {
		// 	v2, ok := test.rules[k1]
		// 	if !ok {
		// 		t.Errorf("test %v: Grammar(%v).Peek().rules\nGOT %v\nEXP %v", i, test.p, rules, test.rules)
		// 	}
		// 	fmt.Println(v1, v2)
		// 	// if !slices.Equal(v1, v2) {
		// 	// 	t.Errorf("test %v: Grammar(%v).Peek().rules\nGOT %v\nEXP %v", i, test.p, rules, test.rules)
		// 	// }
		// }
	}
}

func TestWrapRule(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapRule(tt.args.s); got != tt.want {
				t.Errorf("WrapRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnwrapRule(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnwrapRule(tt.args.s); got != tt.want {
				t.Errorf("UnwrapRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanImportStatement(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanImportStatement(tt.args.s); got != tt.want {
				t.Errorf("CleanImportStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanGrammarStatement(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanGrammarStatement(tt.args.s); got != tt.want {
				t.Errorf("CleanGrammarStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateJSGFName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateJSGFName(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSGFName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJSGFImport(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateJSGFImport(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSGFImport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestValidateJSGFRule(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		l   string
		err error
	}{
		{l: "", err: dummyError},
		{l: ";", err: dummyError},
		{l: "=;", err: dummyError},
		{l: "<>=;", err: dummyError},
		{l: "public<>=;", err: dummyError},
		{l: "public <>=;", err: dummyError},
		{l: "< > = <>; ", err: dummyError},
		{l: "< > = <>;", err: nil},
		{l: "< >=;", err: nil},
		{l: "public < >=;", err: nil},
		{l: "public < > = ;", err: nil},
		{l: "<abc> = def <ghi>;", err: nil},
		{l: "<abc> = def = <ghi>;", err: nil},
		{l: "<abc> = \"def\" = <ghi>;", err: nil},
		{l: "<abc> = def <ghi>;;", err: nil},
	}
	for i, test := range table {
		err := ValidateJSGFRule(test.l)
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: ValidateJSGF(%v)\nGOT %v\nEXP %v", i, test.l, err, test.err)
		}
	}
}
