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
		e1[i].From = e1[i].From + n
		e1[i].To = e1[i].To + n
	}

	return e1
}

func (e EdgeList) IsEmpty() bool {
	return len(e) == 0 || e[0].IsEmpty()
}

func (e EdgeList) Max() int {
	var arr []int
	if e.IsEmpty() {
		return 0
	}
	for _, i := range e {
		arr = append(arr, i.From)
		arr = append(arr, i.To)
	}
	slices.Sort(arr)

	return arr[len(arr)-1]
}

func Unique(e EdgeList) EdgeList {
	var out EdgeList
	seen := make(map[string]int)
	for i, edg := range e {
		seen[fmt.Sprint(edg)] = i
	}
	for _, v := range seen {
		out = append(out, e[v])
	}

	return out
}

func BuildEdgeList(arr []Expression) (edges EdgeList) {
	var (
		err        error
		f          int
		g          int
		groupStack Stack
	)

	groupMap := make(map[int][]int)
	for i, token := range arr {
		switch token {
		case "<EOS>":
			edges = append(edges, Edge{From: f, To: i, Weight: 1.0})
		case "<SOS>":
			f = i
			groupStack = groupStack.Push(i)
			groupMap[i] = []int{}
		case ";":
			edges = append(edges, Edge{From: f, To: i, Weight: 1.0})
			for _, values := range groupMap {
				for _, v := range values {
					edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
				}
			}
			f = i
		case "(", "[":
			groupStack = groupStack.Push(i)
			groupMap[i] = []int{}
			edges = append(edges, Edge{From: f, To: i, Weight: 1.0})
			f = i
		case ")":
			g, err = groupStack.Top()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[g] {
				edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
			}
			groupStack = groupStack.Drop(g)
			delete(groupMap, g)
			edges = append(edges, Edge{From: f, To: i, Weight: 1.0})
			f = i
		case "]":
			g, err := groupStack.Top()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[g] {
				edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
			}
			groupStack = groupStack.Drop(g)
			delete(groupMap, g)
			edges = append(edges, Edge{From: f, To: i, Weight: 1.0})
			edges = append(edges, Edge{From: g, To: i, Weight: 1.0})
			f = i
		case "|":
			g, groupStack, err = groupStack.Pop()
			if err != nil {
				log.Fatal(err)
			}
			groupStack = groupStack.Push(g)
			groupMap[g] = append(groupMap[g], f)
			f = g
		default:
			edges = append(edges, Edge{From: f, To: i, Weight: 1.0})
			f = i
		}
	}

	return Unique(edges)
}
