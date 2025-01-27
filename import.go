// -*- coding: utf-8 -*-

// Created on Thu Jan 23 08:51:30 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"errors"
	"strings"
)

func ParseImport(s string) ([]string, error) {
	var rule string
	var gram string
	var root string

	split := strings.Split(s, ".")
	if len(split) < 2 {
		return []string{}, errors.New("rule address is too short to contain the required grammar and rule specifications")
	}
	rule = split[len(split)-1]
	gram = split[len(split)-2]
	root = strings.Join(split[0:len(split)-2], "/")
	root = ExpandRoot(root)
	out := []string{root, gram, rule}
	return out, nil
}

func ExpandRoot(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	split := strings.Split(s, "/")
	for i := 0; i < len(split); i++ {
		b.WriteString("../")
	}
	b.WriteString(s)
	return b.String()
}
