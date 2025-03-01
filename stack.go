// -*- coding: utf-8 -*-

// Created on Wed Jan 22 08:07:24 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
)

type Stack []int

func (s Stack) Push(v int) Stack {
	return append(s, v)
}

func (s Stack) Top() (int, error) {
	var top int

	if len(s) == 0 {
		return 0, errors.New("no top value in empty stack")
	}
	top = s[len(s)-1]

	return top, nil
}

func (s Stack) Pop() (int, Stack, error) {
	var top int
	var bot Stack

	top, err := s.Top()
	if err != nil {
		return 0, Stack{}, err
	}
	bot = s[:len(s)-1]

	return top, bot, nil
}

func (s Stack) Drop(v int) Stack {
	var s1 Stack

	for _, i := range s {
		if i != v {
			s1 = s1.Push(i)
		}
	}

	return s1
}
