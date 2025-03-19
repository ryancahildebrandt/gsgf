// -*- coding: utf-8 -*-

// Created on Wed Mar  5 11:33:38 AM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"slices"
	"testing"
)

func TestRemoveEndSpaces(t *testing.T) {
	tests := []struct {
		p    []string
		want []string
	}{
		{p: []string{}, want: []string{}},
		{p: []string{""}, want: []string{""}},
		{p: []string{"{}"}, want: []string{"{}"}},
		{p: []string{"\t", "\n", "\r"}, want: []string{"", "", ""}},
		{p: []string{"\\t", "\\n", "\\r"}, want: []string{"\\t", "\\n", "\\r"}},
		{p: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{p: []string{" abc", " def ", "ghi "}, want: []string{"abc", "def", "ghi"}},
		{p: []string{"\tabc", "d\nef", "ghi\r"}, want: []string{"abc", "d\nef", "ghi"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, want: []string{"{}abc", "de{e}f", "ghi{ghi}"}},
	}
	for i, test := range tests {
		got := RemoveEndSpaces(test.p)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RemoveEndSpaces(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
	}
}

func TestRemoveMultipleSpaces(t *testing.T) {
	tests := []struct {
		p    []string
		want []string
	}{
		{p: []string{}, want: []string{}},
		{p: []string{""}, want: []string{""}},
		{p: []string{"{}"}, want: []string{"{}"}},
		{p: []string{"\t", "\n", "\r"}, want: []string{"", "", ""}},
		{p: []string{"\\t", "\\n", "\\r"}, want: []string{"\\t", "\\n", "\\r"}},
		{p: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{p: []string{" abc", " def ", "ghi "}, want: []string{"abc", "def", "ghi"}},
		{p: []string{"  abc", "  d  e  f  ", "ghi  "}, want: []string{"abc", "d e f", "ghi"}},

		{p: []string{"\tabc", "d\nef", "ghi\r"}, want: []string{"abc", "d ef", "ghi"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, want: []string{"{}abc", "de{e}f", "ghi{ghi}"}},
	}
	for i, test := range tests {
		got := RemoveMultipleSpaces(test.p)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RemoveMultipleSpaces(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
	}
}

func TestRenderTabs(t *testing.T) {
	tests := []struct {
		p    []string
		want []string
	}{
		{p: []string{}, want: []string{}},
		{p: []string{""}, want: []string{""}},
		{p: []string{"{}"}, want: []string{"{}"}},
		{p: []string{"\t", "\n", "\r"}, want: []string{"\t", "\n", "\r"}},
		{p: []string{"\\t", "\\n", "\\r"}, want: []string{"\t", "\\n", "\\r"}},
		{p: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{p: []string{" abc", " def ", "ghi "}, want: []string{" abc", " def ", "ghi "}},
		{p: []string{"\tabc", "d\nef", "ghi\r"}, want: []string{"\tabc", "d\nef", "ghi\r"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, want: []string{"{}abc", "de{e}f", "ghi{ghi}"}},
	}
	for i, test := range tests {
		got := RenderTabs(test.p)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RenderTabs(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
	}
}

func TestRenderNewLines(t *testing.T) {
	tests := []struct {
		p    []string
		want []string
	}{
		{p: []string{}, want: []string{}},
		{p: []string{""}, want: []string{""}},
		{p: []string{"{}"}, want: []string{"{}"}},
		{p: []string{"\t", "\n", "\r"}, want: []string{"\t", "\n", "\r"}},
		{p: []string{"\\t", "\\n", "\\r"}, want: []string{"\\t", "\n", "\\r"}},
		{p: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{p: []string{" abc", " def ", "ghi "}, want: []string{" abc", " def ", "ghi "}},
		{p: []string{"\tabc", "d\nef", "ghi\r"}, want: []string{"\tabc", "d\nef", "ghi\r"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, want: []string{"{}abc", "de{e}f", "ghi{ghi}"}},
	}
	for i, test := range tests {
		got := RenderNewLines(test.p)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RenderNewLines(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
	}
}

func TestRemoveTags(t *testing.T) {
	tests := []struct {
		p    []string
		want []string
	}{
		{p: []string{}, want: []string{}},
		{p: []string{""}, want: []string{""}},
		{p: []string{"{}"}, want: []string{""}},
		{p: []string{"\t", "\n", "\r"}, want: []string{"\t", "\n", "\r"}},
		{p: []string{"\\t", "\\n", "\\r"}, want: []string{"\\t", "\\n", "\\r"}},
		{p: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{p: []string{" abc", " def ", "ghi "}, want: []string{" abc", " def ", "ghi "}},
		{p: []string{"\tabc", "d\nef", "ghi\r"}, want: []string{"\tabc", "d\nef", "ghi\r"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, want: []string{"abc", "def", "ghi"}},
	}
	for i, test := range tests {
		got := RemoveTags(test.p)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RemoveTags(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
	}
}

func TestWrapTags(t *testing.T) {
	tests := []struct {
		p      []string
		prefix string
		suffix string
		want   []string
	}{
		{p: []string{}, prefix: "", suffix: "", want: []string{}},
		{p: []string{""}, prefix: " ", suffix: "", want: []string{""}},
		{p: []string{"{}"}, prefix: "", suffix: " ", want: []string{"{} "}},
		{p: []string{"\t", "\n", "\r"}, prefix: "-", suffix: "_", want: []string{"\t", "\n", "\r"}},
		{p: []string{"\\t", "\\n", "\\r"}, prefix: "abc", suffix: "def", want: []string{"\\t", "\\n", "\\r"}},
		{p: []string{"a", "b", "c"}, prefix: "<", suffix: ">", want: []string{"a", "b", "c"}},
		{p: []string{" abc", " def ", "ghi "}, prefix: "{}", suffix: "{}", want: []string{" abc", " def ", "ghi "}},
		{p: []string{"\tabc", "d\nef", "ghi\r"}, prefix: "\n", suffix: "\n", want: []string{"\tabc", "d\nef", "ghi\r"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, prefix: "\t", suffix: "\t", want: []string{"\t{}\tabc", "de\t{e}\tf", "ghi\t{ghi}\t"}},
	}
	for i, test := range tests {
		got := WrapTags(test.p, test.prefix, test.suffix)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: WrapTags(%v, %v, %v)\nGOT %v\nWANT %v", i, test.p, test.prefix, test.suffix, got, test.want)
		}
	}
}

func TestWrapProductions(t *testing.T) {
	tests := []struct {
		p      []string
		prefix string
		suffix string
		want   []string
	}{
		{p: []string{}, prefix: "", suffix: "", want: []string{}},
		{p: []string{""}, prefix: " ", suffix: "", want: []string{" "}},
		{p: []string{"{}"}, prefix: "", suffix: " ", want: []string{"{} "}},
		{p: []string{"\t", "\n", "\r"}, prefix: "-", suffix: "_", want: []string{"-\t_", "-\n_", "-\r_"}},
		{p: []string{"\\t", "\\n", "\\r"}, prefix: "abc", suffix: "def", want: []string{"abc\\tdef", "abc\\ndef", "abc\\rdef"}},
		{p: []string{"a", "b", "c"}, prefix: "<", suffix: ">", want: []string{"<a>", "<b>", "<c>"}},
		{p: []string{" abc", " def ", "ghi "}, prefix: "{}", suffix: "{}", want: []string{"{} abc{}", "{} def {}", "{}ghi {}"}},
		{p: []string{"\tabc", "d\nef", "ghi\r"}, prefix: "\n", suffix: "\n", want: []string{"\n\tabc\n", "\nd\nef\n", "\nghi\r\n"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, prefix: "\t", suffix: "\t", want: []string{"\t{}abc\t", "\tde{e}f\t", "\tghi{ghi}\t"}},
	}
	for i, test := range tests {
		got := WrapProductions(test.p, test.prefix, test.suffix)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: WrapProductions(%v, %v, %v)\nGOT %v\nWANT %v", i, test.p, test.prefix, test.suffix, got, test.want)
		}
	}
}

func TestCollectTags(t *testing.T) {
	tests := []struct {
		p    []string
		c    string
		want []string
	}{
		{p: []string{}, c: "", want: []string{}},
		{p: []string{""}, c: "#", want: []string{""}},
		{p: []string{"{}"}, c: " ", want: []string{" {},"}},
		{p: []string{"\t", "\n", "\r"}, c: "_", want: []string{"\t", "\n", "\r"}},
		{p: []string{"\\t", "\\n", "\\r"}, c: "def", want: []string{"\\t", "\\n", "\\r"}},
		{p: []string{"a", "b", "c"}, c: "// ", want: []string{"a", "b", "c"}},
		{p: []string{" abc", " def ", "ghi "}, c: "{}", want: []string{" abc", " def ", "ghi "}},
		{p: []string{"\tabc", "d\nef", "ghi\r"}, c: "\n", want: []string{"\tabc", "d\nef", "ghi\r"}},
		{p: []string{"{}abc", "de{e}f", "ghi{ghi}"}, c: "\t", want: []string{"abc\t{},", "def\t{e},", "ghi\t{ghi},"}},
	}
	for i, test := range tests {
		got := CollectTags(test.p, test.c)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: CollectTags(%v, %v)\nGOT %v\nWANT %v", i, test.p, test.c, got, test.want)
		}
	}
}
