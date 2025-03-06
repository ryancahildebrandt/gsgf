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
		// TODO: Add test cases.
	}
	for i, test := range tests {
		got := WrapTags(test.p, test.prefix, test.suffix)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RemoveEndSpaces(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
	}
}

func TestCollectestags(t *testing.T) {
	tests := []struct {
		p    []string
		c    string
		want []string
	}{
		// TODO: Add test cases.
	}
	for i, test := range tests {
		got := CollectTags(test.p, test.c)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RemoveEndSpaces(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
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
		// TODO: Add test cases.
	}
	for i, test := range tests {
		got := WrapProductions(test.p, test.prefix, test.suffix)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: RemoveEndSpaces(%v)\nGOT %v\nWANT %v", i, test.p, got, test.want)
		}
	}
}
