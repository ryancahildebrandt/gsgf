// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:15 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestCreateNameSpace(t *testing.T) {
	table := []struct {
		d       string
		e       string
		imports map[string][]string
		rules   map[string]map[string][]string
		err     error
	}{
		{"data/tests/test0.jsgf",
			".jsgf",
			map[string][]string{
				"a":     {},
				"b":     {"import <c.brew>"},
				"c":     {},
				"d":     {"import <a.order>", "import <c.teatype>"},
				"e":     {},
				"test0": {"import <a.*>"},
				"test1": {"import <c.brew>"},
				"test2": {"import <a1.*>"},
				"test3": {"import <e.dne>"},
				"test4": {"import <d.*>"},
				"test5": {"import <b.request>"},
			},
			map[string]map[string][]string{
				"a":     {"order": {"quant"}, "request": {"brew"}},
				"b":     {"order": {"quant"}, "request": {"brew"}, "quant": {""}},
				"c":     {"teatype": {""}, "brew": {"quant"}},
				"d":     {},
				"e":     {"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
				"test0": {"main": {"request", "order", "quant", "teatype"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
				"test1": {"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {""}, "teatype": {""}},
				"test2": {"main": {"request", "order", "quant", "teatype"}, "order": {"quant"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
				"test3": {"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
				"test4": {"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "quant": {""}, "brew": {"quant"}},
				"test5": {"main": {"request", "order", "quant", "teatype"}, "order": {"quant"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
			},
			nil,
		},
		{"data/tests/dir0/c.jsgf",
			".jsgf",
			map[string][]string{
				"c": {},
				"d": {"import <a.order>", "import <c.teatype>"},
				"e": {},
			},
			map[string]map[string][]string{
				"c": {"teatype": {""}, "brew": {"quant"}},
				"d": {},
				"e": {"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
			},
			nil,
		},
		{"data/tests/dir0/dir1/d.jsgf",
			".jsgf",
			map[string][]string{
				"d": {"import <a.order>", "import <c.teatype>"},
				"e": {},
			},
			map[string]map[string][]string{
				"d": {},
				"e": {"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
			},
			nil,
		},
		{"data/tests/dir0/dir1/dir2/e.jsgf",
			".jsgf",
			map[string][]string{
				"e": {},
			},
			map[string]map[string][]string{
				"e": {"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {""}, "teatype": {""}, "brew": {"quant"}},
			},
			nil,
		},
	}
	for _, test := range table {
		imports, rules, err := CreateNameSpace(test.d, test.e)

		for k, v_res := range imports {
			v_exp, ok := test.imports[k]
			if !ok {
				t.Errorf("CreateNameSpace(%v, %v).imports\nGOT %v\nEXP %v", test.d, test.e, imports, test.imports)
			}
			sort.Strings(v_exp)
			sort.Strings(v_res)
			if fmt.Sprint((v_exp)) != fmt.Sprint((v_res)) {
				t.Errorf("CreateNameSpace(%v, %v).imports\nGOT %v\nEXP %v", test.d, test.e, imports, test.imports)
			}
		}

		for k, v_res := range rules {
			v_exp, ok := test.rules[k]
			if !ok {
				t.Errorf("CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", test.d, test.e, v_res, v_exp)
			}
			for kk, vv_res := range v_exp {
				vv_exp, ok := v_exp[kk]
				if !ok {
					t.Errorf("CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", test.d, test.e, vv_res, vv_exp)
				}
				sort.Strings(vv_exp)
				sort.Strings(vv_res)
				if fmt.Sprint(vv_exp) != fmt.Sprint(vv_res) {
					t.Errorf("CreateNameSpace(%v, %v).rules\nGOT %v\nEXP %v", test.d, test.e, vv_res, vv_exp)
				}
			}
		}

		if test.err != nil && err == nil {
			t.Errorf("CreateNameSpace(%v, %v).err\nGOT %v\nEXP %v", test.d, test.e, err, test.err)
		}
		if test.err == nil && err != nil {
			t.Errorf("CreateNameSpace(%v, %v).err\nGOT %v\nEXP %v", test.d, test.e, err, test.err)
		}

	}
}
