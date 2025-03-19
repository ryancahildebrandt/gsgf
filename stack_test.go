// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:11:49 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"slices"
	"testing"
)

func TestStackTop(t *testing.T) {
	table := []struct {
		s       stack
		want    int
		wantErr bool
	}{
		{s: stack{}, want: 0, wantErr: true},
		{s: stack{0}, want: 0, wantErr: false},
		{s: stack{0, 1, 2}, want: 2, wantErr: false},
		{s: stack{100, 99, 98, 97, 96, 95}, want: 95, wantErr: false},
		{s: stack{2, 1, 0, -1, -2}, want: -2, wantErr: false},
	}
	for i, test := range table {
		got, err := test.s.top()
		if got != test.want {
			t.Errorf("test %v: %v.Top()\nGOT %v\nWANT %v", i, test.s, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.Top()\nGOT %v\nWANT %v", i, test.s, err, test.wantErr)
		}
	}
}

func TestStackPop(t *testing.T) {
	table := []struct {
		s       stack
		t       int
		b       stack
		wantErr bool
	}{
		{s: stack{}, t: 0, b: stack{}, wantErr: true},
		{s: stack{0}, t: 0, b: stack{}, wantErr: false},
		{s: stack{0, 1, 2}, t: 2, b: stack{0, 1}, wantErr: false},
		{s: stack{100, 99, 98, 97, 96, 95}, t: 95, b: stack{100, 99, 98, 97, 96}, wantErr: false},
		{s: stack{2, 1, 0, -1, -2}, t: -2, b: stack{2, 1, 0, -1}, wantErr: false},
	}
	for i, test := range table {
		top, bot, err := test.s.pop()
		if top != test.t {
			t.Errorf("test %v: %v.Pop()\nGOT %v\nWANT %v", i, test.s, t, test.t)
		}
		if !slices.Equal(bot, test.b) {
			t.Errorf("test %v: %v.Pop()\nGOT %v\nWANT %v", i, test.s, bot, test.b)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.Pop()\nGOT %v\nWANT %v", i, test.s, err, test.wantErr)
		}
	}
}

func TestStackDrop(t *testing.T) {
	table := []struct {
		s    stack
		n    int
		want stack
	}{
		{s: stack{}, n: 0, want: stack{}},
		{s: stack{0}, n: 0, want: stack{}},
		{s: stack{0, 1, 2}, n: 1, want: stack{0, 2}},
		{s: stack{0, 1, 2, 2}, n: 2, want: stack{0, 1}},
		{s: stack{100, 99, 98, 97, 96, 95}, n: 98, want: stack{100, 99, 97, 96, 95}},
		{s: stack{2, 1, 0, -1, -2}, n: -2, want: stack{2, 1, 0, -1}},
	}
	for i, test := range table {
		got := test.s.drop(test.n)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.Drop()\nGOT %v\nWANT %v", i, test.s, got, test.want)
		}
	}
}
