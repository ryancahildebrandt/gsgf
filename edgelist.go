// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:25 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"log"
	"slices"
	"sort"
)

// Contains weight information between 2 graph nodes
type Edge struct {
	From   int
	To     int
	Weight float64
}

// Checks if both e.From and e.To are the default int value 0
func (e Edge) isEmpty() bool {
	return e.From == 0 && e.To == 0
}

// Convenience type definition for a slice of edges, used to convert token sequences to graph
type EdgeList []Edge

// Returns EdgeList e sorted by From and To values
func Sort(e EdgeList) EdgeList {
	sort.Slice(e, func(i, j int) bool {
		switch {
		case e[i].From < e[j].From:
			return true
		case e[i].From == e[j].From && e[i].To < e[j].To:
			return true
		default:
			return false
		}
	})

	return e
}

// Returns EdgeList with all From and To values increased by n. Used mainly in graph composition
func increment(e EdgeList, n int) EdgeList {
	var e1 EdgeList
	e1 = append(e1, e...)

	for i := range e1 {
		e1[i].From += n
		e1[i].To += n
	}

	return e1
}

// Checks if first edge is empty or EdgeList has length 0
func (e EdgeList) isEmpty() bool {
	return len(e) == 0 || e[0].isEmpty()
}

// Returns highest value from collected From and To values
func (e EdgeList) max() int {
	if e.isEmpty() {
		return 0
	}

	var arr []int
	for _, edg := range e {
		arr = append(arr, edg.From)
		arr = append(arr, edg.To)
	}

	return slices.Max(arr)
}

// Returns a slice of all unique edges in EdgeList e
func Unique(e EdgeList) EdgeList {
	var out EdgeList
	var seen map[string]int = make(map[string]int)

	for i := range e {
		seen[fmt.Sprint(e[i])] = i
	}
	for _, v := range seen {
		out = append(out, e[v])
	}

	return out
}

// Converts a slice of tokens/Expressions to an EdgeList
// Uses flow control tokens (), [], | to capture possible state transitions between tokens
func ToEdgeList(arr []Expression) EdgeList {
	var (
		edges      EdgeList
		err        error
		from       int
		group      int
		groupStack stack
		groupMap   = make(map[int][]int)
	)

	for i, token := range arr {
		switch token {
		case "<EOS>":
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
		case "<SOS>":
			from = i
			groupStack = groupStack.push(i)
			groupMap[i] = []int{}
		case ";":
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			for _, v := range groupMap {
				for _, v1 := range v {
					edges = append(edges, Edge{From: v1, To: i, Weight: 1.0})
				}
			}
			from = i
		case "(", "[":
			groupStack = groupStack.push(i)
			groupMap[i] = []int{}
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			from = i
		case ")":
			group, err = groupStack.top()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[group] {
				edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
			}
			groupStack = groupStack.drop(group)
			delete(groupMap, group)
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			from = i
		case "]":
			group, err := groupStack.top()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[group] {
				edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
			}
			groupStack = groupStack.drop(group)
			delete(groupMap, group)
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			edges = append(edges, Edge{From: group, To: i, Weight: 1.0})
			from = i
		case "|":
			group, groupStack, err = groupStack.pop()
			if err != nil {
				log.Fatal(err)
			}
			groupStack = groupStack.push(group)
			groupMap[group] = append(groupMap[group], from)
			from = group
		default:
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			from = i
		}
	}

	return Unique(edges)
}
