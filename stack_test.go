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
		s       Stack
		want    int
		wantErr bool
	}{
		{s: Stack{}, want: 0, wantErr: true},
		{s: Stack{0}, want: 0, wantErr: false},
		{s: Stack{0, 1, 2}, want: 2, wantErr: false},
		{s: Stack{100, 99, 98, 97, 96, 95}, want: 95, wantErr: false},
		{s: Stack{2, 1, 0, -1, -2}, want: -2, wantErr: false},
	}
	for i, test := range table {
		got, err := test.s.Top()
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
		s       Stack
		t       int
		b       Stack
		wantErr bool
	}{
		{s: Stack{}, t: 0, b: Stack{}, wantErr: true},
		{s: Stack{0}, t: 0, b: Stack{}, wantErr: false},
		{s: Stack{0, 1, 2}, t: 2, b: Stack{0, 1}, wantErr: false},
		{s: Stack{100, 99, 98, 97, 96, 95}, t: 95, b: Stack{100, 99, 98, 97, 96}, wantErr: false},
		{s: Stack{2, 1, 0, -1, -2}, t: -2, b: Stack{2, 1, 0, -1}, wantErr: false},
	}
	for i, test := range table {
		top, bot, err := test.s.Pop()
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
		s    Stack
		n    int
		want Stack
	}{
		{s: Stack{}, n: 0, want: Stack{}},
		{s: Stack{0}, n: 0, want: Stack{}},
		{s: Stack{0, 1, 2}, n: 1, want: Stack{0, 2}},
		{s: Stack{0, 1, 2, 2}, n: 2, want: Stack{0, 1}},
		{s: Stack{100, 99, 98, 97, 96, 95}, n: 98, want: Stack{100, 99, 97, 96, 95}},
		{s: Stack{2, 1, 0, -1, -2}, n: -2, want: Stack{2, 1, 0, -1}},
	}
	for i, test := range table {
		got := test.s.Drop(test.n)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.Drop()\nGOT %v\nWANT %v", i, test.s, got, test.want)
		}
	}
}
