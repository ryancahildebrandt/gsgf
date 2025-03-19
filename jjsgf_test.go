// -*- coding: utf-8 -*-

// Created on Fri Mar 14 09:33:39 PM EDT 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"strings"
	"testing"
)

func TestJJSGFToJSGF(t *testing.T) {
	table := []struct {
		j    JJSGFGrammarJSON
		want string
	}{
		{
			j: JJSGFGrammarJSON{
				Name:    "",
				Public:  map[string]string{},
				Rules:   map[string]string{},
				Imports: []string{},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar ;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name:    "abc",
				Public:  map[string]string{},
				Rules:   map[string]string{"<a>": ""},
				Imports: []string{"def"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar abc;\nimport <def>;\n<<a>> = ;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name:    "",
				Public:  map[string]string{"<a>": "<b><c>"},
				Rules:   map[string]string{},
				Imports: []string{"def"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar ;\nimport <def>;\npublic <<a>> = <b><c>;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name:    "b",
				Public:  map[string]string{"<a>": ""},
				Rules:   map[string]string{},
				Imports: []string{"a.a"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar b;\nimport <a.a>;\npublic <<a>> = ;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name:    "grammarb",
				Public:  map[string]string{"<a>": "<b><c>"},
				Rules:   map[string]string{"<a>": "<b><c>"},
				Imports: []string{"test0.*"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar grammarb;\nimport <test0.*>;\npublic <<a>> = <b><c>;\n<<a>> = <b><c>;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name:   "name",
				Public: map[string]string{},
				Rules: map[string]string{
					"<a>": "a",
					"<b>": "b",
				},
				Imports: []string{"import"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar name;\nimport <import>;\n<<a>> = a;\n<<b>> = b;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name: "name name name",
				Public: map[string]string{
					"<a>": "",
					"<b>": "",
				},
				Rules:   map[string]string{},
				Imports: []string{"import", "import", "import"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar name name name;\nimport <import>;\nimport <import>;\nimport <import>;\npublic <<a>> = ;\npublic <<b>> = ;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name:    "",
				Public:  map[string]string{"<a>": "<b>"},
				Rules:   map[string]string{"<c>": "d"},
				Imports: []string{"empty"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar ;\nimport <empty>;\npublic <<a>> = <b>;\n<<c>> = d;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name:   ".<>.",
				Public: map[string]string{},
				Rules: map[string]string{
					"<a>": "<c>",
					"<b>": "<c>",
					"<c>": "<d>",
					"<d>": "",
				},
				Imports: []string{"does not exist"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar .<>.;\nimport <does not exist>;\n<<a>> = <c>;\n<<b>> = <c>;\n<<c>> = <d>;\n<<d>> = ;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name: "name",
				Public: map[string]string{"<a>": "<c>",
					"<b>": "<c>",
					"<c>": "<d>",
					"<d>": ""},
				Rules:   map[string]string{},
				Imports: []string{"name"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar name;\nimport <name>;\npublic <<a>> = <c>;\npublic <<b>> = <c>;\npublic <<c>> = <d>;\npublic <<d>> = ;\n",
		},
		{
			j: JJSGFGrammarJSON{
				Name: "grammar",
				Public: map[string]string{
					"<a>": "<c>",
					"<b>": "<c>",
					"<c>": "<d>",
					"<d>": ""},
				Rules: map[string]string{
					"<a>": "<c>",
					"<b>": "<c>",
					"<c>": "<d>",
					"<d>": "",
				},
				Imports: []string{"grammar.*"},
			},
			want: "#JSGF V1.0 ISO8859-1 en;\ngrammar grammar;\nimport <grammar.*>;\npublic <<a>> = <c>;\npublic <<b>> = <c>;\npublic <<c>> = <d>;\npublic <<d>> = ;\n<<a>> = <c>;\n<<b>> = <c>;\n<<c>> = <d>;\n<<d>> = ;\n",
		},
	}
	for i, test := range table {
		got := JJSGFToJSGF(test.j)
		if got != test.want {
			t.Errorf("test %v: JJSGFtoJSGF(%v)\nGOT %v\nWANT %v", i, test.j, strings.ReplaceAll(got, "\n", "\\n"), strings.ReplaceAll(test.want, "\n", "\\n"))
		}
	}
}
