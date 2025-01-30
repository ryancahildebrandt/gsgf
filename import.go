// -*- coding: utf-8 -*-

// Created on Thu Jan 23 08:51:30 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

func ParseImportPath(s string) (string, string, error) {
	ind := strings.LastIndex(s, ".")
	if ind == -1 {
		return "", "", errors.New("rule specification is too short to contain the required grammar and rule")
	}
	return s[:ind], s[ind+1:], nil
}

func ReadRule(s *bufio.Scanner, r string) (string, error) {
	target := fmt.Sprint("public <", r, ">")
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, target) {
			return line, nil
		}
	}
	return "", errors.New("target rule does not exist in grammar or is not public")
}

func ReadAllRules(s *bufio.Scanner) []string {
	out := []string{}
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "public <") {
			out = append(out, line)
		}
	}
	return out
}
