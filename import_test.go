// -*- coding: utf-8 -*-

// Created on Sat Feb 15 05:19:15 PM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

// func TestImportPeek(t *testing.T) {
// 	table := []struct {
// 		p           string
// 		exp_imports []string
// 		exp_rules   map[string][]string
// 	}{
// 		{"data/tests/dir2/dir1/dir0/test0.jsgf", []string{"import <a.*>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/test1.jsgf", []string{"import <dir0/c.brew>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/test2.jsgf", []string{"import <a1.*>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/test3.jsgf", []string{"import <e.dne>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/test4.jsgf", []string{"import <dir0/dir1/d.*>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "quant": {}, "brew": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/test5.jsgf", []string{"import <b.request>"}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/a.jsgf", []string{}, map[string][]string{"request": {"brew"}, "order": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/b.jsgf", []string{"import <dir0/c.brew>"}, map[string][]string{"request": {"brew"}, "order": {"quant"}, "quant": {}}},
// 		{"data/tests/dir2/dir1/dir0/dir0/c.jsgf", []string{}, map[string][]string{"teatype": {}, "brew": {"quant"}}},
// 		{"data/tests/dir2/dir1/dir0/dir0/dir1/d.jsgf", []string{"import <../c.teatype>", "import <../../a.order>"}, map[string][]string{}},
// 		{"data/tests/dir2/dir1/dir0/dir0/dir1/dir2/e.jsgf", []string{}, map[string][]string{"main": {"request", "order", "quant", "teatype"}, "request": {"brew"}, "order": {"quant"}, "quant": {}, "teatype": {}, "brew": {"quant"}}},
// 	}
// 	for _, test := range table {
// 		res_imports, res_rules, err := NewImport(test.p, ".jsgf").Peek()
// 		if err != nil {
// 			t.Errorf("Import.Peek(%v)\nGOT error %v", test.p, err)
// 		}
// 		sort.Strings(res_imports)
// 		sort.Strings(test.exp_imports)

// 		if fmt.Sprint(res_imports) != fmt.Sprint(test.exp_imports) {
// 			t.Errorf("Import.Peek(%v).imports\nGOT %v\nEXP %v", test.p, res_imports, test.exp_imports)
// 		}

// 		for k, v_res := range res_rules {
// 			v_exp, ok := test.exp_rules[k]
// 			if !ok {
// 				t.Errorf("Import.Peek(%v).rules\nGOT %v\nEXP %v", test.p, res_rules, test.exp_rules)
// 			}
// 			sort.Strings(v_exp)
// 			sort.Strings(v_res)
// 			if fmt.Sprint((v_exp)) != fmt.Sprint((v_res)) {
// 				t.Errorf("Import.Peek(%v).rules\nGOT %v\nEXP %v", test.p, res_rules, test.exp_rules)
// 			}
// 		}
// 	}
// }

// func TestNewImport(t *testing.T) {
// 	table := []struct {
// 		p   string
// 		exp Import
// 	}{
// 		{"data/tests/dir2/dir1/dir0/test0.rule", Import{"data/tests/dir2/dir1/dir0/test0.rule", ".jsgf", "test0.jsgf", "test0.rule", "test0", "rule", "data/tests/dir2/dir1/dir0"}},
// 		{"data/tests/dir2/dir1/dir0/a.rule", Import{"data/tests/dir2/dir1/dir0/a.rule", ".jsgf", "a.jsgf", "a.rule", "a", "rule", "data/tests/dir2/dir1/dir0"}},
// 		{"data/tests/dir2/dir1/dir0/b.*", Import{"data/tests/dir2/dir1/dir0/b.*", ".jsgf", "b.jsgf", "b.*", "b", "*", "data/tests/dir2/dir1/dir0"}},
// 		{"data/tests/dir2/dir1/dir0/dir0/c.*", Import{"data/tests/dir2/dir1/dir0/dir0/c.*", ".jsgf", "c.jsgf", "c.*", "c", "*", "data/tests/dir2/dir1/dir0/dir0"}},
// 		{"data/tests/dir2/dir1/dir0/dir0/dir1/d.jsgf", Import{"data/tests/dir2/dir1/dir0/dir0/dir1/d.jsgf", ".jsgf", "d.jsgf", "d.jsgf", "d", "jsgf", "data/tests/dir2/dir1/dir0/dir0/dir1"}},
// 		{"data/tests/dir2/dir1/dir0/dir0/dir1/dir2/e.jsgf", Import{"data/tests/dir2/dir1/dir0/dir0/dir1/dir2/e.jsgf", ".jsgf", "e.jsgf", "e.jsgf", "e", "jsgf", "data/tests/dir2/dir1/dir0/dir0/dir1/dir2"}},

// 		{"import <test0.rule>", Import{"test0.rule", ".jsgf", "test0.jsgf", "test0.rule", "test0", "rule", "."}},
// 		{"import <../a.rule>", Import{"../a.rule", ".jsgf", "a.jsgf", "a.rule", "a", "rule", ".."}},
// 		{"import <../../b.*>", Import{"../../b.*", ".jsgf", "b.jsgf", "b.*", "b", "*", "../.."}},
// 		{"import <../../dir0/c.*>", Import{"../../dir0/c.*", ".jsgf", "c.jsgf", "c.*", "c", "*", "../../dir0"}},
// 		{"import <dir0/dir1/d.jsgf>", Import{"dir0/dir1/d.jsgf", ".jsgf", "d.jsgf", "d.jsgf", "d", "jsgf", "dir0/dir1"}},
// 		{"import <dir0/dir1/dir2/e.jsgf>", Import{"dir0/dir1/dir2/e.jsgf", ".jsgf", "e.jsgf", "e.jsgf", "e", "jsgf", "dir0/dir1/dir2"}},

// 		{"import <test0.>", Import{"test0.", ".jsgf", "test0.jsgf", "test0.", "test0", "", "."}},
// 		{"import <../a.>", Import{"../a.", ".jsgf", "a.jsgf", "a.", "a", "", ".."}},
// 		{"import <dir0/d.>", Import{"dir0/d.", ".jsgf", "d.jsgf", "d.", "d", "", "dir0"}},
// 		{"", Import{"", "", "", "", "", "", ""}},
// 		{"import <>", Import{"", "", "", "", "", "", ""}},
// 	}
// 	for _, test := range table {
// 		res := NewImport(test.p, ".jsgf")

// 		if res.path != test.exp.path {
// 			t.Errorf("NewImport(%v).path\nGOT %v\nEXP %v", test.p, res.path, test.exp.path)
// 		}
// 		if res.ext != test.exp.ext {
// 			t.Errorf("NewImport(%v).ext\nGOT %v\nEXP %v", test.p, res.ext, test.exp.ext)
// 		}
// 		if res.file != test.exp.file {
// 			t.Errorf("NewImport(%v).file\nGOT %v\nEXP %v", test.p, res.file, test.exp.file)
// 		}
// 		if res.target != test.exp.target {
// 			t.Errorf("NewImport(%v).target\nGOT %v\nEXP %v", test.p, res.target, test.exp.target)
// 		}
// 		if res.gram != test.exp.gram {
// 			t.Errorf("NewImport(%v).gram\nGOT %v\nEXP %v", test.p, res.gram, test.exp.gram)
// 		}
// 		if res.rule != test.exp.rule {
// 			t.Errorf("NewImport(%v).rule\nGOT %v\nEXP %v", test.p, res.rule, test.exp.rule)
// 		}
// 		if res.dir != test.exp.dir {
// 			t.Errorf("NewImport(%v).dir\nGOT %v\nEXP %v", test.p, res.dir, test.exp.dir)
// 		}
// 	}
// }
