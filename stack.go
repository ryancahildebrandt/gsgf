// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:24 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"fmt"
)

type stack []int

func (s stack) push(v int) stack {
	return append(s, v)
}

func (s stack) top() (int, error) {
	var top int

	if len(s) == 0 {
		return 0, fmt.Errorf("error when calling Stack(%v).Top():\n%+w", s, errors.New("no top value in empty stack"))
	}
	top = s[len(s)-1]

	return top, nil
}

func (s stack) pop() (int, stack, error) {
	var top int
	var bot stack

	top, err := s.top()
	if err != nil {
		return 0, stack{}, fmt.Errorf("in Stack(%v).Pop():\n%+w", s, err)
	}
	bot = s[:len(s)-1]

	return top, bot, nil
}

// Returns stack with all instances of value v removed
func (s stack) drop(v int) stack {
	var s1 stack

	for _, i := range s {
		if i != v {
			s1 = s1.push(i)
		}
	}

	return s1
}
