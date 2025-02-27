// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:24 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
)

type Stack []int
//type Stack = []int // will need to change methods to functions

func (s Stack) Push(v int) Stack {
	return append(s, v)
}

func (s Stack) Peek() (t int, err error) {
	if len(s) == 0 {
		return 0, errors.New("no top value in empty stack")
	}
	t = s[len(s)-1]
	return t, nil
}

func (s Stack) Pop() (t int, b Stack, err error) {
	t, err = s.Peek()
	if err != nil {
		return 0, Stack{}, errors.New("cannot pop from empty stack")
	}
	b = s[:len(s)-1]
	return t, b, nil
}

func (s Stack) Drop(v int) (out Stack) {
	for _, i := range s {
		if i != v {
			out = out.Push(i)
		}
	}
	return out
}
