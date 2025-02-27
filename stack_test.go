// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:11:49 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"slices"
	"testing"
)

func TestStackPeek(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		s   Stack
		exp int
		err error
	}{
		{s: Stack{}, exp: 0, err: dummyError},
		{s: Stack{0}, exp: 0, err: nil},
		{s: Stack{0, 1, 2}, exp: 2, err: nil},
		{s: Stack{100, 99, 98, 97, 96, 95}, exp: 95, err: nil},
		{s: Stack{2, 1, 0, -1, -2}, exp: -2, err: nil},
	}
	for i, test := range table {
		res, err := test.s.Peek()
		if res != test.exp {
			t.Errorf("test %v: %v.Peek()\nGOT %v\nEXP %v", i, test.s, res, test.exp)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.Peek()\nGOT %v\nEXP %v", i, test.s, err, test.err)
		}
	}
}

func TestStackPop(t *testing.T) {
	dummyError := errors.New("")
	table := []struct {
		s   Stack
		t   int
		b   Stack
		err error
	}{
		{s: Stack{}, t: 0, b: Stack{}, err: dummyError},
		{s: Stack{0}, t: 0, b: Stack{}, err: nil},
		{s: Stack{0, 1, 2}, t: 2, b: Stack{0, 1}, err: nil},
		{s: Stack{100, 99, 98, 97, 96, 95}, t: 95, b: Stack{100, 99, 98, 97, 96}, err: nil},
		{s: Stack{2, 1, 0, -1, -2}, t: -2, b: Stack{2, 1, 0, -1}, err: nil},
	}
	for i, test := range table {
		top, bot, err := test.s.Pop()
		if top != test.t {
			t.Errorf("test %v: %v.Pop()\nGOT %v\nEXP %v", i, test.s, t, test.t)
		}
		if !slices.Equal(bot, test.b) {
			t.Errorf("test %v: %v.Pop()\nGOT %v\nEXP %v", i, test.s, bot, test.b)
		}
		if (test.err != nil && err == nil) || (test.err == nil && err != nil) {
			t.Errorf("test %v: %v.Pop()\nGOT %v\nEXP %v", i, test.s, err, test.err)
		}
	}
}

func TestStackDrop(t *testing.T) {
	table := []struct {
		s   Stack
		n   int
		exp Stack
	}{
		{s: Stack{}, n: 0, exp: Stack{}},
		{s: Stack{0}, n: 0, exp: Stack{}},
		{s: Stack{0, 1, 2}, n: 1, exp: Stack{0, 2}},
		{s: Stack{0, 1, 2, 2}, n: 2, exp: Stack{0, 1}},
		{s: Stack{100, 99, 98, 97, 96, 95}, n: 98, exp: Stack{100, 99, 97, 96, 95}},
		{s: Stack{2, 1, 0, -1, -2}, n: -2, exp: Stack{2, 1, 0, -1}},
	}
	for i, test := range table {
		res := test.s.Drop(test.n)
		if !slices.Equal(res, test.exp) {
			t.Errorf("test %v: %v.Drop()\nGOT %v\nEXP %v", i, test.s, res, test.exp)
		}
	}
}
