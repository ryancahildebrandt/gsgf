// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:11:49 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestStackPeek(t *testing.T) {
	dummy_error := errors.New("")
	table := []struct {
		s   Stack
		exp int
		err error
	}{
		{Stack{}, 0, dummy_error},
		{Stack{0}, 0, nil},
		{Stack{0, 1, 2}, 2, nil},
		{Stack{100, 99, 98, 97, 96, 95}, 95, nil},
		{Stack{2, 1, 0, -1, -2}, -2, nil},
	}
	for _, test := range table {
		res, err := test.s.Peek()
		if fmt.Sprint(test.exp) != fmt.Sprint(res) {
			t.Errorf("%v.Peek()\nGOT %v\nEXP %v", test.s, res, test.exp)
		}
		if test.err != nil && err == nil {
			t.Errorf("%v.Peek()\nGOT %v\nEXP %v", test.s, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("%v.Peek()\nGOT %v\nEXP %v", test.s, err, test.err)
		}
	}
}

func TestStackPop(t *testing.T) {
	dummy_error := errors.New("")
	table := []struct {
		s     Stack
		exp_t int
		exp_b Stack
		err   error
	}{
		{Stack{}, 0, Stack{}, dummy_error},
		{Stack{0}, 0, Stack{}, nil},
		{Stack{0, 1, 2}, 2, Stack{0, 1}, nil},
		{Stack{100, 99, 98, 97, 96, 95}, 95, Stack{100, 99, 98, 97, 96}, nil},
		{Stack{2, 1, 0, -1, -2}, -2, Stack{2, 1, 0, -1}, nil},
	}
	for _, test := range table {
		res_t, res_b, err := test.s.Pop()
		if fmt.Sprint(test.exp_t) != fmt.Sprint(res_t) {
			t.Errorf("%v.Pop()\nGOT %v\nEXP %v", test.s, res_t, test.exp_t)
		}
		if fmt.Sprint(test.exp_b) != fmt.Sprint(res_b) {
			t.Errorf("%v.Pop()\nGOT %v\nEXP %v", test.s, res_b, test.exp_b)
		}
		if test.err != nil && err == nil {
			t.Errorf("%v.Pop()\nGOT %v\nEXP %v", test.s, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("%v.Pop()\nGOT %v\nEXP %v", test.s, err, test.err)
		}
	}
}

func TestStackDrop(t *testing.T) {
	table := []struct {
		s   Stack
		n   int
		exp Stack
	}{
		{Stack{}, 0, Stack{}},
		{Stack{0}, 0, Stack{}},
		{Stack{0, 1, 2}, 1, Stack{0, 2}},
		{Stack{0, 1, 2, 2}, 2, Stack{0, 1}},
		{Stack{100, 99, 98, 97, 96, 95}, 98, Stack{100, 99, 97, 96, 95}},
		{Stack{2, 1, 0, -1, -2}, -2, Stack{2, 1, 0, -1}},
	}
	for _, test := range table {
		res := test.s.Drop(test.n)
		if fmt.Sprint(test.exp) != fmt.Sprint(res) {
			t.Errorf("%v.Drop()\nGOT %v\nEXP %v", test.s, res, test.exp)
		}
	}
}
