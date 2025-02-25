// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:25 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"log"
	"sort"
)

type Edge struct {
	From   int
	To     int
	Weight float64
}

func (e Edge) Copy() Edge {
	return Edge{e.From, e.To, e.Weight}
}

func (e Edge) IsEmpty() bool {
	return e.From == 0 && e.To == 0
}

type EdgeList []Edge

func (e EdgeList) Sort() EdgeList {
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

func (e EdgeList) Copy() EdgeList {
	var e1 EdgeList
	for _, edge := range e {
		e1 = append(e1, edge.Copy())
	}
	return e1
}

func (e EdgeList) Increment(n int) EdgeList {
	e1 := e.Copy()
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
	if e.IsEmpty() {
		return 0
	}
	var arr []int
	for _, i := range e {
		arr = append(arr, i.From)
		arr = append(arr, i.To)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr[len(arr)-1]
}

func (e EdgeList) Unique() (out EdgeList) {
	e1 := e.Copy()
	for i := range e1.Sort() {
		if i+1 == len(e1) {
			out = append(out, e1[i])
			break
		}
		if fmt.Sprint(e1[i]) != fmt.Sprint(e1[i+1]) {
			out = append(out, e1[i])
		}
	}
	return out
}

func BuildEdgeList(arr []Expression) (edges EdgeList) {
	var err error
	var f int
	var g int
	var groupStack Stack
	var groupMap = make(map[int][]int)

	for i, token := range arr {
		switch token {
		case "<EOS>":
			edges = append(edges, Edge{f, i, 1.0})
		case "<SOS>":
			f = i
			groupStack = groupStack.Push(i)
			groupMap[i] = []int{}
		case ";":
			edges = append(edges, Edge{f, i, 1.0})
			for _, values := range groupMap {
				for _, v := range values {
					edges = append(edges, Edge{v, i, 1.0})
				}
			}
			f = i
		case "(", "[":
			groupStack = groupStack.Push(i)
			groupMap[i] = []int{}
			edges = append(edges, Edge{f, i, 1.0})
			f = i
		case ")":
			g, err = groupStack.Peek()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[g] {
				edges = append(edges, Edge{v, i, 1.0})
			}
			groupStack = groupStack.Drop(g)
			delete(groupMap, g)
			edges = append(edges, Edge{f, i, 1.0})
			f = i
		case "]":
			g, err := groupStack.Peek()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range groupMap[g] {
				edges = append(edges, Edge{v, i, 1.0})
			}
			groupStack = groupStack.Drop(g)
			delete(groupMap, g)
			edges = append(edges, Edge{f, i, 1.0})
			edges = append(edges, Edge{g, i, 1.0})
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
			edges = append(edges, Edge{f, i, 1.0})
			f = i
		}
	}
	return edges.Unique()
}
