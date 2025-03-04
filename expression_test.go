// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:12:27 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"slices"
	"testing"
)

func TestToTokens(t *testing.T) {
	lexer := NewJSGFLexer()
	table := []struct {
		e    Expression
		want []Expression
	}{
		{
			e:    "",
			want: []Expression{}},
		{
			e:    ";",
			want: []Expression{"<SOS>", ";", "<EOS>"}},
		{
			e:    " ",
			want: []Expression{"<SOS>", " ", "<EOS>"}},
		{
			e:    "test expression 123",
			want: []Expression{"<SOS>", "test expression 123", "<EOS>"}},
		{
			e:    "test expression 123;",
			want: []Expression{"<SOS>", "test expression 123", ";", "<EOS>"}},
		{
			e:    "test expression 123 (abc);",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 [abc];",
			want: []Expression{"<SOS>", "test expression 123 ", "[", "abc", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 [(abc)];",
			want: []Expression{"<SOS>", "test expression 123 ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (abc) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (ab|c) | [de|f];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab", "|", "c", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 <rule> (abc) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 <rule> (abc) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 <rule> (abc) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "<rule>", " ", "(", "abc", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123// [(abc)];",
			want: []Expression{"<SOS>", "test expression 123// ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 // [(abc)];",
			want: []Expression{"<SOS>", "test expression 123 // ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 /0.0/ [(abc)];",
			want: []Expression{"<SOS>", "test expression 123 /0.0/", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123{} [(abc)];",
			want: []Expression{"<SOS>", "test expression 123{}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 {} [(abc)];",
			want: []Expression{"<SOS>", "test expression 123 {}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 {tag} [(abc)];",
			want: []Expression{"<SOS>", "test expression 123 {tag}", " ", "[", "(", "abc", ")", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (abc /1.0/) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc /1.0/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (abc /1000/) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc /1000/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (abc /-0.0001/) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc /-0.0001/", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (abc { }) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc { }", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (abc {_}) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc {_}", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (abc {_t_a_g_}) [def];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "abc {_t_a_g_}", ")", " ", "[", "def", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (ab/1.0/|c/1.0/) | [de|f];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab/1.0/", "|", "c/1.0/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (ab/1000/|c/1000/) | [de|f];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab/1000/", "|", "c/1000/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (ab/-0.0001/|c/-0.0001/) | [de|f];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab/-0.0001/", "|", "c/-0.0001/", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (ab{1}|c{1}) | [de|f];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1}", "|", "c{1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (ab{1.1}|c{1.1}) | [de|f];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1}", "|", "c{1.1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
		{
			e:    "test expression 123 (ab{1.1/1}|c{1.1/1}) | [de|f];",
			want: []Expression{"<SOS>", "test expression 123 ", "(", "ab{1.1/1}", "|", "c{1.1/1}", ")", " ", "|", " ", "[", "de", "|", "f", "]", ";", "<EOS>"},
		},
	}
	for i, test := range table {
		got := ToTokens(test.e, lexer)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.ToTokens(jsgflexer)\nGOT %v\nWANT %v", i, test.e, got, test.want)
		}
	}
}

func TestParseWeight(t *testing.T) {
	table := []struct {
		e       Expression
		ee      Expression
		ef      float64
		wantErr bool
	}{
		{e: "", ee: "", ef: 0.0, wantErr: true},
		{e: "/", ee: "/", ef: 0.0, wantErr: true},
		{e: "//", ee: "//", ef: 0.0, wantErr: true},
		{e: "/.//", ee: "/.//", ef: 0.0, wantErr: true},
		{e: "abc", ee: "abc", ef: 0.0, wantErr: true},
		{e: "abc//", ee: "abc//", ef: 0.0, wantErr: true},
		{e: "/aaa/", ee: "/aaa/", ef: 0.0, wantErr: true},
		{e: "/1.0/", ee: "", ef: 1.0, wantErr: false},
		{e: "abc/1.0/", ee: "abc", ef: 1.0, wantErr: false},
		{e: "abc/0.1/", ee: "abc", ef: 0.1, wantErr: false},
		{e: "abc/0.001/", ee: "abc", ef: 0.001, wantErr: false},
		{e: "abc/100000/", ee: "abc", ef: 100000.0, wantErr: false},
		{e: "abc/100000.0/", ee: "abc", ef: 100000.0, wantErr: false},
		{e: "abc/1/", ee: "abc", ef: 1, wantErr: false},
		{e: "abc/-1/", ee: "abc", ef: -1.0, wantErr: false},
		{e: "abc/-1.0/", ee: "abc", ef: -1.0, wantErr: false},
		{e: "abc/0/", ee: "abc", ef: 0.0, wantErr: false},
	}
	for i, test := range table {
		e, f, err := ParseWeight(test.e)
		if test.ee != e {
			t.Errorf("test %v: ParseWeight(%v)\nGOT %v\nWANT %v", i, test.e, e, test.ee)
		}
		if test.ef != f {
			t.Errorf("test %v: ParseWeight(%v)\nGOT %v\nWANT %v", i, test.e, f, test.ef)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: ParseWeight(%v)\nGOT %v\nWANT %v", i, test.e, err, test.wantErr)
		}
	}
}

func TestIsWeighted(t *testing.T) {
	tests := []struct {
		e    Expression
		want bool
	}{
		{e: "/0/", want: true},
		{e: " /0.0/", want: true},
		{e: "abc/999/", want: true},
		{e: "\abc/./", want: true},
		{e: "abc/.0/abc", want: true},
		{e: "", want: false},
		{e: " ", want: false},
		{e: "abc", want: false},
		{e: "\abc", want: false},
		{e: "//abc", want: false},
		{e: "/0.0a/abc", want: false},
	}
	for i, test := range tests {
		got := IsWeighted(test.e)
		if got != test.want {
			t.Errorf("test %v: IsWeighted(%v)\nGOT %v\nWANT %v", i, test.e, got, test.want)
		}
	}
}
