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

type Edge struct {
	From   int
	To     int
	Weight float64
}

func (e Edge) IsEmpty() bool {
	return e.From == 0 && e.To == 0
}

type EdgeList []Edge

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

func Increment(e EdgeList, n int) EdgeList {
	var e1 EdgeList
	e1 = append(e1, e...)

	for i := range e1 {
		e1[i].From += n
		e1[i].To += n
	}

	return e1
}

func (e EdgeList) IsEmpty() bool {
	return len(e) == 0 || e[0].IsEmpty()
}

func (e EdgeList) Max() int {
	if e.IsEmpty() {
		return 0
	}

	var arr []int
	for _, edg := range e {
		arr = append(arr, edg.From)
		arr = append(arr, edg.To)
	}

	return slices.Max(arr)
}

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

func ToEdgeList(arr []Expression) EdgeList {
	var (
		edges      EdgeList
		err        error
		from       int
		group      int
		groupStack Stack
		groupMap   = make(map[int][]int)
	)

	for i, token := range arr {
		switch token {
		case "<EOS>":
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
		case "<SOS>":
			from = i
			groupStack = groupStack.Push(i)
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
			groupStack = groupStack.Push(i)
			groupMap[i] = []int{}
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			from = i
		case ")":
			group, err = groupStack.Top()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[group] {
				edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
			}
			groupStack = groupStack.Drop(group)
			delete(groupMap, group)
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			from = i
		case "]":
			group, err := groupStack.Top()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[group] {
				edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
			}
			groupStack = groupStack.Drop(group)
			delete(groupMap, group)
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			edges = append(edges, Edge{From: group, To: i, Weight: 1.0})
			from = i
		case "|":
			group, groupStack, err = groupStack.Pop()
			if err != nil {
				log.Fatal(err)
			}
			groupStack = groupStack.Push(group)
			groupMap[group] = append(groupMap[group], from)
			from = group
		default:
			edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
			from = i
		}
	}

	return Unique(edges)
}
