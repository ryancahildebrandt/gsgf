// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:15 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"slices"
	"sort"
	"testing"
)

func TestCreateNameSpace(t *testing.T) {
	table := []struct {
		d       string
		e       string
		r       map[string]string
		wantErr bool
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
			wantErr: false,
		},
		{
			d: "data/tests/test0.jjsgf",
			e: ".jjsgf",
			r: map[string]string{
				"<order>":   "i'd like [to order|a|<quant>];",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<quant>":   "some|a (cup|glass) of;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
			wantErr: false,
		},
		{
			d: "data/tests/test1.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			wantErr: false,
		},
		{
			d:       "data/tests/test2.jsgf",
			e:       ".jsgf",
			r:       map[string]string{},
			wantErr: true,
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
			wantErr: false,
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
			wantErr: false,
		},
		{
			d: "data/tests/test4.jjsgf",
			e: ".jjsgf",
			r: map[string]string{
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
		},
		{
			d:       "data/tests/a.jsgf",
			e:       ".jsgf",
			r:       map[string]string{},
			wantErr: false,
		},
		{
			d: "data/tests/b.jsgf",
			e: ".jsgf",
			r: map[string]string{
				"<brew>":    "(make|brew|whip up) <quant>;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<quant>":   "some|a (cup|glass) of;",
			},
			wantErr: false,
		},
		{
			d:       "data/tests/dir0/c.jsgf",
			e:       ".jsgf",
			r:       map[string]string{},
			wantErr: false,
		},
		{
			d:       "data/tests/dir0/c.jjsgf",
			e:       ".jjsgf",
			r:       map[string]string{},
			wantErr: false,
		},
		{
			d:       "data/tests/dir0/dir1/d.jsgf",
			e:       ".jsgf",
			r:       map[string]string{},
			wantErr: true,
		},
		{
			d:       "data/tests/dir0/dir1/d.jjsgf",
			e:       ".jjsgf",
			r:       map[string]string{},
			wantErr: true,
		},
		{
			d:       "data/tests/dir0/dir1/dir2/e.jsgf",
			e:       ".jsgf",
			r:       map[string]string{},
			wantErr: false,
		},
	}
	for i, test := range table {
		rules, err := CreateNameSpace(test.d, test.e)
		if len(rules) != len(test.r) {
			t.Errorf("test %v: CreateNameSpace(%v, %v).rules\nGOT %v\nWANT %v", i, test.d, test.e, rules, test.r)
		}
		for k1, v1 := range rules {
			v2, ok := test.r[k1]
			if !ok {
				t.Errorf("test %v: CreateNameSpace(%v, %v).rules\nGOT %v\nWANT %v", i, test.d, test.e, v1, v2)
			}
			if v1 != v2 {
				t.Errorf("test %v: CreateNameSpace(%v, %v).rules\nGOT %v\nWANT %v", i, test.d, test.e, v1, v2)
			}
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: CreateNameSpace(%v, %v).err\nGOT %v\nWANT %v", i, test.d, test.e, err, test.wantErr)
		}
	}
}

func TestFindGrammar(t *testing.T) {
	table := []struct {
		p       string
		t       string
		e       string
		want    string
		wantErr bool
	}{
		{p: "./data/tests/.jsgf", t: "test0", e: ".jsgf", want: "data/tests/test0.jsgf", wantErr: false},
		{p: "./data/tests/.jsgf", t: "test0", e: ".jjsgf", want: "data/tests/test0.jjsgf", wantErr: false},
		{p: "./data/tests/test0.jsgf", t: "test0", e: ".jsgf", want: "data/tests/test0.jsgf", wantErr: false},
		{p: "./data/tests/test0.jsgf", t: "a", e: ".jsgf", want: "data/tests/a.jsgf", wantErr: false},
		{p: "./data/tests/test0.jsgf", t: "a", e: ".jjsgf", want: "data/tests/a.jjsgf", wantErr: false},
		{p: "./data/tests/test0.jsgf", t: "e", e: ".jsgf", want: "data/tests/dir0/dir1/dir2/e.jsgf", wantErr: false},
		{p: "./data/tests/a.jsgf", t: "a", e: ".jsgf", want: "data/tests/a.jsgf", wantErr: false},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "d", e: ".jsgf", want: "data/tests/dir0/dir1/d.jsgf", wantErr: false},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "e", e: ".jsgf", want: "data/tests/dir0/dir1/dir2/e.jsgf", wantErr: false},
		{p: "./data/tests/dir0/dir1/dir2/e.jsgf", t: "e", e: ".jsgf", want: "data/tests/dir0/dir1/dir2/e.jsgf", wantErr: false},
		{p: "./data/tests", t: "test0", e: ".jsgf", want: "data/tests/test0.jsgf", wantErr: false},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "b", e: ".jsgf", want: "", wantErr: true},
		{p: "./data/tests/dir0/dir1/c.jsgf", t: "b", e: ".jjsgf", want: "", wantErr: true},
		{p: "./data/tests/test0.jsgf", t: "f", e: ".jsgf", want: "", wantErr: true},
		{p: "./data/tests/dir0/dir1/dir2/e.jsgf", t: "d", e: ".jsgf", want: "", wantErr: true},
	}
	for i, test := range table {
		got, err := findGrammar(test.p, test.t, test.e)
		if got != test.want {
			t.Errorf("test %v: FindGrammar(%v, %v, %v)\nGOT %v\nWANT %v", i, test.p, test.t, test.e, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: FindGrammar(%v, %v, %v).err\nGOT %v\nWANT %v", i, test.p, test.t, test.e, err, test.wantErr)
		}
	}
}

func TestImportOrder(t *testing.T) {
	table := []struct {
		p       string
		e       string
		want    []string
		wantErr bool
	}{
		{p: "./data/tests", e: ".jsgf", want: []string{}, wantErr: true},
		{p: "./data/tests/.jsgf", e: ".jsgf", want: []string{}, wantErr: true},
		{p: "./data/tests/.jjsgf", e: ".jjsgf", want: []string{}, wantErr: true},
		{p: "./data/tests/test0.jsgf", e: ".jsgf", want: []string{"import <a.*>;"}, wantErr: false},
		{p: "./data/tests/test1.jsgf", e: ".jsgf", want: []string{"import <c.brew>;"}, wantErr: false},
		{p: "./data/tests/test1.jjsgf", e: ".jjsgf", want: []string{"import <c.brew>;"}, wantErr: false},
		{p: "./data/tests/test3.jsgf", e: ".jsgf", want: []string{"import <e.dne>;"}, wantErr: false},
		{p: "./data/tests/test4.jsgf", e: ".jsgf", want: []string{"import <a.order>;", "import <c.teatype>;", "import <d.*>;"}, wantErr: false},
		{p: "./data/tests/test4.jsgf", e: ".jsgf", want: []string{"import <a.order>;", "import <c.teatype>;", "import <d.*>;"}, wantErr: false},
		{p: "./data/tests/test5.jsgf", e: ".jsgf", want: []string{"import <b.request>;", "import <c.brew>;"}, wantErr: false},
		{p: "./data/tests/a.jsgf", e: ".jsgf", want: []string{}, wantErr: false},
		{p: "./data/tests/b.jsgf", e: ".jsgf", want: []string{"import <c.brew>;"}, wantErr: false},
		{p: "./data/tests/dir0/c.jsgf", e: ".jsgf", want: []string{}, wantErr: false},
		{p: "./data/tests/dir0/c.jjsgf", e: ".jjsgf", want: []string{}, wantErr: false},
		{p: "./data/tests/dir0/dir1/dir2/e.jsgf", e: ".jsgf", want: []string{}, wantErr: false},
		{p: "./data/tests/test2.jsgf", e: ".jsgf", want: []string{}, wantErr: true},
		{p: "./data/tests/test2.jjsgf", e: ".jjsgf", want: []string{}, wantErr: true},
		{p: "./data/tests/dir0/dir1/d.jsgf", e: ".jsgf", want: []string{}, wantErr: true},
	}
	for i, test := range table {
		got, err := getImportOrder(test.p, test.e)
		sort.Strings(test.want)
		sort.Strings(got)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: ImportOrder(%v, %v)\nGOT %v\nWANT %v", i, test.p, test.e, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ImportOrder(%v, %v).err\nGOT %v\nWANT %v", i, test.p, test.e, err, test.wantErr)
		}
	}
}

func TestPeekGrammar(t *testing.T) {
	table := []struct {
		p       string
		n       string
		imports []string
		rules   map[string]string
	}{
		{
			p:       "data/tests/test0.jsgf",
			n:       "test0",
			imports: []string{"import <a.*>;"},
			rules: map[string]string{
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
		},
		{
			p:       "data/tests/test0.jjsgf",
			n:       "test0",
			imports: []string{"import <a.*>;"},
			rules: map[string]string{
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
		},
		{
			p:       "data/tests/test1.jsgf",
			n:       "test1",
			imports: []string{"import <c.brew>;"},
			rules: map[string]string{
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
			},
		},
		{
			p:       "data/tests/test2.jsgf",
			n:       "test2",
			imports: []string{"import <a1.*>;"},
			rules: map[string]string{
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
		},
		{
			p:       "data/tests/test3.jsgf",
			n:       "test3",
			imports: []string{"import <e.dne>;"},
			rules: map[string]string{
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
		},
		{
			p:       "data/tests/test4.jsgf",
			n:       "test4",
			imports: []string{"import <d.*>;"},
			rules: map[string]string{"<main>": "(<request>|<order>) <quant> <teatype> tea;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<quant>":   "some|a (cup|glass) of;",
				"<brew>":    "(make|brew|whip up) <quant>;"},
		},
		{
			p:       "data/tests/test5.jsgf",
			n:       "test5",
			imports: []string{"import <b.request>;"},
			rules: map[string]string{"<main>": "(<request>|<order>) <quant> <teatype> tea;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;"},
		},
		{
			p:       "data/tests/a.jsgf",
			n:       "a",
			imports: []string{},
			rules: map[string]string{"<request>": "[(could|will|would) you] please <brew>;",
				"<order>": "i'd like [to order|a|<quant>];",
				"<brew>":  "(make|brew|whip up) <quant>;",
				"<quant>": "some|a (cup|glass) of;"},
		},
		{
			p:       "data/tests/a.jjsgf",
			n:       "a",
			imports: []string{},
			rules: map[string]string{"<request>": "[(could|will|would) you] please <brew>;",
				"<order>": "i'd like [to order|a|<quant>];",
				"<brew>":  "(make|brew|whip up) <quant>;",
				"<quant>": "some|a (cup|glass) of;"},
		},
		{
			p:       "data/tests/b.jsgf",
			n:       "b",
			imports: []string{"import <c.brew>;"},
			rules: map[string]string{"<request>": "[(could|will|would) you] please <brew>;",
				"<order>": "i'd like [to order|a|<quant>];",
				"<quant>": "some|a (cup|glass) of;"},
		},
		{
			p:       "data/tests/dir0/c.jsgf",
			n:       "c",
			imports: []string{},
			rules: map[string]string{"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":  "(make|brew|whip up) <quant>;",
				"<quant>": "some|a (cup|glass) of;"},
		},
		{
			p: "data/tests/dir0/dir1/d.jsgf",
			n: "d",
			imports: []string{"import <c.teatype>;",
				"import <a.order>;"},
			rules: map[string]string{},
		},
		{
			p:       "data/tests/dir0/dir1/dir2/e.jsgf",
			n:       "e",
			imports: []string{},
			rules: map[string]string{
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
		},
		{
			p:       "data/tests/dir0/dir1/dir2/e.jjsgf",
			n:       "e",
			imports: []string{},
			rules: map[string]string{
				"<main>":    "(<request>|<order>) <quant> <teatype> tea;",
				"<request>": "[(could|will|would) you] please <brew>;",
				"<order>":   "i'd like [to order|a|<quant>];",
				"<quant>":   "some|a (cup|glass) of;",
				"<teatype>": "red|sweet|green|jasmine|milk;",
				"<brew>":    "(make|brew|whip up) <quant>;",
			},
		},
	}
	for i, test := range table {
		name, imports, rules, err := peekGrammar(test.p)
		if err != nil {
			t.Errorf("test %v: PeekGrammar(%v)\nGOT error %v", i, test.p, err)
		}
		if name != test.n {
			t.Errorf("test %v: PeekGrammar(%v).imports\nGOT %v\nWANT %v", i, test.p, name, test.n)
		}
		sort.Strings(imports)
		sort.Strings(test.imports)
		if !slices.Equal(imports, test.imports) {
			t.Errorf("test %v: PeekGrammar(%v)imports\nGOT %v\nWANT %v", i, test.p, imports, test.imports)
		}
		if len(rules) != len(test.rules) {
			t.Errorf("test %v: PeekGrammar(%v).rules\nGOT %v\nWANT %v", i, test.p, rules, test.rules)
		}
		for k1, v1 := range rules {
			v2, ok := test.rules[k1]
			if !ok {
				t.Errorf("test %v: PeekGrammar(%v).rules\nGOT %v\nWANT %v", i, test.p, ok, k1)
			}
			if v1 != v2 {
				t.Errorf("test %v: PeekGrammar(%v).rules\nGOT %v\nWANT %v", i, test.p, v1, v2)
			}
		}
	}
}

func TestWrapRule(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: "<>"},
		{s: " ", want: "< >"},
		{s: "<>", want: "<<>>"},
		{s: "abc", want: "<abc>"},
		{s: "<abc>", want: "<<abc>>"},
		{s: "abc def", want: "<abc def>"},
	}
	for i, test := range tests {
		got := wrapRule(test.s)
		if got != test.want {
			t.Errorf("test %v: WrapRule(%v)\nGOT %v\nWANT %v", i, test.s, got, test.want)
		}
	}
}

func TestUnwrapRule(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: ""},
		{s: " ", want: ""},
		{s: "public<>", want: ""},
		{s: "public abc", want: "abc"},
		{s: "<abc>", want: "abc"},
		{s: "  <abc> ", want: "abc"},
		{s: "public <def>", want: "def"},
		{s: "public <<def>>", want: "<def>"},
	}
	for i, test := range tests {
		got := unwrapRule(test.s)
		if got != test.want {
			t.Errorf("test %v: UnwrapRule(%v)\nGOT %v\nWANT %v", i, test.s, got, test.want)
		}
	}
}

func TestCleanImportStatement(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: ""},
		{s: " ", want: ""},
		{s: "import<>;", want: ""},
		{s: "import abc;", want: "abc"},
		{s: "<abc.r>", want: "abc.r"},
		{s: "  <abc.>", want: "abc."},
		{s: "import <.def>;", want: ".def"},
		{s: "import <<.def.abc>>;", want: "<.def.abc>"},
	}
	for i, test := range tests {
		got := cleanImportStatement(test.s)
		if got != test.want {
			t.Errorf("test %v: CleanImportStatement(%v)\nGOT %v\nWANT %v", i, test.s, got, test.want)
		}
	}
}

func TestCleanGrammarStatement(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: ""},
		{s: " ", want: ""},
		{s: "grammar<>;", want: "<>"},
		{s: "grammar abc;", want: "abc"},
		{s: "<abc.r>", want: "<abc.r>"},
		{s: "  <abc.>", want: "<abc.>"},
		{s: "grammar <.def>;", want: "<.def>"},
		{s: "grammar <<.def.abc>>;", want: "<<.def.abc>>"},
	}
	for i, test := range tests {
		got := cleanGrammarStatement(test.s)
		if got != test.want {
			t.Errorf("test %v: CleanGrammarStatement(%v)\nGOT %v\nWANT %v", i, test.s, got, test.want)
		}
	}
}

func TestValidateJSGFName(t *testing.T) {
	tests := []struct {
		s       string
		wantErr bool
	}{
		{s: "", wantErr: true},
		{s: " ", wantErr: true},
		{s: "grammar<>;", wantErr: true},
		{s: "grammar abc;", wantErr: false},
		{s: "<abc.r>", wantErr: true},
		{s: "  <abc.>", wantErr: true},
		{s: "grammar <.def>;", wantErr: false},
		{s: "grammar <<.def.abc>>;", wantErr: false},
	}
	for i, test := range tests {
		err := ValidateJSGFName(test.s)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateJSGFName(%v)\nGOT %v\nWANT %v", i, test.s, err, test.wantErr)
		}
	}
}

func TestValidateJSGFImport(t *testing.T) {
	tests := []struct {
		s       string
		wantErr bool
	}{
		{s: "", wantErr: true},
		{s: " ", wantErr: true},
		{s: "import<>;", wantErr: true},
		{s: "import abc;", wantErr: true},
		{s: "<abc.r>", wantErr: true},
		{s: "  <abc.>", wantErr: true},
		{s: "import <.def>;", wantErr: false},
		{s: "import <<.def.abc>>;", wantErr: false},
	}
	for i, test := range tests {
		err := ValidateJSGFImport(test.s)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateJSGFImport(%v)\nGOT %v\nWANT %v", i, test.s, err, test.wantErr)
		}
	}
}
func TestValidateJSGFRule(t *testing.T) {
	table := []struct {
		l       string
		wantErr bool
	}{
		{l: "", wantErr: true},
		{l: ";", wantErr: true},
		{l: "=;", wantErr: true},
		{l: "<>=;", wantErr: true},
		{l: "public<>=;", wantErr: true},
		{l: "public <>=;", wantErr: true},
		{l: "< > = <>; ", wantErr: true},
		{l: "< > = <>;", wantErr: false},
		{l: "< >=;", wantErr: false},
		{l: "public < >=;", wantErr: false},
		{l: "public < > = ;", wantErr: false},
		{l: "<abc> = def <ghi>;", wantErr: false},
		{l: "<abc> = def = <ghi>;", wantErr: false},
		{l: "<abc> = \"def\" = <ghi>;", wantErr: false},
		{l: "<abc> = def <ghi>;;", wantErr: false},
	}
	for i, test := range table {
		err := ValidateJSGFRule(test.l)
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ValidateJSGFRule(%v)\nGOT %v\nWANT %v", i, test.l, err, test.wantErr)
		}
	}
}
