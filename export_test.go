// -*- coding: utf-8 -*-// Created on Sat Feb 22 10:19:34 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt
package main

import (
	"reflect"
	"testing"
)

func TestRuleToJSON(t *testing.T) {
	type args struct {
		r Rule
	}
	tests := []struct {
		name string
		args args
		exp  RuleJSON
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := RuleToJSON(test.args.r); !reflect.DeepEqual(got, test.exp) {
				t.Errorf("RuleToJSON() = %v, exp %v", got, test.exp)
			}
		})
	}
}

func TestEdgeToJSON(t *testing.T) {
	type args struct {
		e Edge
	}
	tests := []struct {
		name string
		args args
		exp  EdgeJSON
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := EdgeToJSON(test.args.e); !reflect.DeepEqual(got, test.exp) {
				t.Errorf("EdgeToJSON() = %v, exp %v", got, test.exp)
			}
		})
	}
}

func TestGraphToJSON(t *testing.T) {
	type args struct {
		g Graph
	}
	tests := []struct {
		name string
		args args
		exp  GraphJSON
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := GraphToJSON(test.args.g); !reflect.DeepEqual(got, test.exp) {
				t.Errorf("GraphToJSON() = %v, exp %v", got, test.exp)
			}
		})
	}
}

func TestGrammarToJSON(t *testing.T) {
	type args struct {
		g Grammar
	}
	tests := []struct {
		name string
		args args
		exp  GrammarJSON
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := GrammarToJSON(test.args.g); !reflect.DeepEqual(got, test.exp) {
				t.Errorf("GrammarToJSON() = %v, exp %v", got, test.exp)
			}
		})
	}
}

func TestGraphToTXT(t *testing.T) {
	type args struct {
		g Graph
	}
	tests := []struct {
		name string
		args args
		exp  string
		exp1 string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, got1 := GraphToTXT(test.args.g)
			if got != test.exp {
				t.Errorf("GraphToTXT() got = %v, exp %v", got, test.exp)
			}
			if got1 != test.exp1 {
				t.Errorf("GraphToTXT() got1 = %v, exp %v", got1, test.exp1)
			}
		})
	}
}

func TestGraphToDOT(t *testing.T) {
	type args struct {
		g Graph
	}
	tests := []struct {
		name string
		args args
		exp  string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := GraphToDOT(test.args.g); got != test.exp {
				t.Errorf("GraphToDOT() = %v, exp %v", got, test.exp)
			}
		})
	}
}

func TestReferencesToDOT(t *testing.T) {
	type args struct {
		g Grammar
	}
	tests := []struct {
		name string
		args args
		exp  string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ReferencesToDOT(test.args.g); got != test.exp {
				t.Errorf("ReferencesToDOT() = %v, exp %v", got, test.exp)
			}
		})
	}
}
