// -*- coding: utf-8 -*-

// Created on Mon Dec 30 05:00:34 PM EST 2024
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"slices"
	"testing"
)

func TestToEdgeList(t *testing.T) {
	lexer := NewJSGFLexer("\"")
	table := []struct {
		r       string
		want    EdgeList
		wantErr bool
	}{
		{
			r:       "",
			want:    EdgeList{},
			wantErr: true,
		},
		{
			r:       "=",
			want:    EdgeList{},
			wantErr: true,
		},
		{
			r:       "<>=;",
			want:    EdgeList{},
			wantErr: true,
		},
		{
			r:       "public <test> = ;",
			want:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			wantErr: false,
		},
		{
			r:       "public <test> = one two three;",
			want:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			wantErr: false,
		},
		{
			r: "public <test> = four|five|six;",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 3, Weight: 1.0},
				{From: 0, To: 5, Weight: 1.0},
				{From: 1, To: 6, Weight: 1.0},
				{From: 3, To: 6, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = six[ seven][ eight];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 4, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 5, To: 7, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = eight( nine)( ten);",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [12|13|14];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 (12|13|14);",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 3, To: 8, Weight: 1.0},
				{From: 5, To: 8, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 (((12)));",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 ((12)(13)(14));",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 ((12)|(13)|(14));",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 ([[12]]);",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 ([12][13][14]);",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 ([12]|[13]|[14]);",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [((12))];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [(12)(13)(14)];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [(12)|(13)|(14)];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [[[12]]];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 2, To: 8, Weight: 1.0},
				{From: 3, To: 7, Weight: 1.0},
				{From: 4, To: 6, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [[12][13][14]];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 2, To: 12, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 9, To: 11, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [[12]|[13]|[14]];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 14, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 14, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 2, To: 14, Weight: 1.0},
				{From: 3, To: 5, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 11, To: 13, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 ((12)|(13)[14]) 15;",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 13, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 10, To: 12, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = 11 [(12)|[13]14][15];",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 11, Weight: 1.0},
				{From: 2, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 7, To: 9, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 2, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
				{From: 12, To: 14, Weight: 1.0},
				{From: 15, To: 16, Weight: 1.0},
			},
			wantErr: false,
		},
		{
			r: "public <test> = [(11)12[13](14)] 15;",
			want: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 1, To: 2, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 3, To: 4, Weight: 1.0},
				{From: 4, To: 5, Weight: 1.0},
				{From: 5, To: 6, Weight: 1.0},
				{From: 6, To: 7, Weight: 1.0},
				{From: 7, To: 8, Weight: 1.0},
				{From: 8, To: 9, Weight: 1.0},
				{From: 9, To: 10, Weight: 1.0},
				{From: 10, To: 11, Weight: 1.0},
				{From: 11, To: 12, Weight: 1.0},
				{From: 12, To: 13, Weight: 1.0},
				{From: 13, To: 14, Weight: 1.0},
				{From: 1, To: 12, Weight: 1.0},
				{From: 6, To: 8, Weight: 1.0},
				{From: 14, To: 15, Weight: 1.0},
			},
			wantErr: false,
		},
	}
	for i, test := range table {
		_, v, err := ParseRule(test.r, lexer)
		got := ToEdgeList(ToTokens(v.exp, lexer))
		if !slices.Equal(Sort(got), Sort(test.want)) {
			t.Errorf("test %v: %v.toArray(lexer)\nGOT %v\nWANT %v", i, test.r, got, test.want)
		}
		if (err != nil) != test.wantErr {
			t.Errorf("test %v: %v.toArray(lexer).err\nGOT %v\nWANT %v", i, test.r, err, test.wantErr)
		}
	}
}

func TestEdgeListSort(t *testing.T) {
	table := []struct {
		e    EdgeList
		want EdgeList
	}{
		{
			e:    EdgeList{},
			want: EdgeList{},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 2, To: 3, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
	}
	for i, test := range table {
		got := Sort(test.e)
		if !slices.Equal(got, test.want) {
			t.Errorf("test %v: %v.Sort()\nGOT %v\nWANT %v", i, test.e, got, test.want)
		}
	}
}

func TestEdgeListUnique(t *testing.T) {
	table := []struct {
		e    EdgeList
		want EdgeList
	}{
		{
			e:    EdgeList{},
			want: EdgeList{},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 1, To: 2, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 2, To: 3, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}, {From: 0, To: 1, Weight: 1.0}},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
		{
			e: EdgeList{
				{From: 0, To: 1, Weight: 1.0},
				{From: 0, To: 1, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
				{From: 2, To: 3, Weight: 1.0},
			},
			want: EdgeList{{From: 0, To: 1, Weight: 1.0}, {From: 2, To: 3, Weight: 1.0}},
		},
	}
	for i, test := range table {
		got := Unique(test.e)
		if !slices.Equal(Sort(got), Sort(test.want)) {
			t.Errorf("test %v: %v.Unique()\nGOT %v\nWANT %v", i, test.e, got, test.want)
		}
	}
}

func TestEdgeListMax(t *testing.T) {
	table := []struct {
		e    EdgeList
		want int
	}{
		{
			e:    EdgeList{},
			want: 0,
		},
		{
			e:    EdgeList{{From: 0, To: 0, Weight: 1.0}},
			want: 0,
		},
		{
			e:    EdgeList{{From: 1, To: 1, Weight: 1.0}, {From: 1, To: 1, Weight: 1.0}},
			want: 1,
		},
		{
			e: EdgeList{
				{From: 10, To: 12, Weight: 1.0},
				{From: 55, To: 94, Weight: 1.0},
				{From: 0, To: 15, Weight: 1.0},
				{From: 1, To: 1, Weight: 1.0},
			},
			want: 94,
		},
		{
			e: EdgeList{
				{From: -1, To: 1651, Weight: 1.0},
				{From: 55, To: 65, Weight: 1.0},
				{From: 10, To: 1, Weight: 1.0},
				{From: 15, To: 99, Weight: 1.0},
				{From: 65, To: 54, Weight: 1.0},
				{From: 1000000000, To: 0, Weight: 1.0},
				{From: 0, To: 8, Weight: 1.0},
				{From: 15, To: 44, Weight: 1.0},
			},
			want: 1000000000,
		},
	}
	for i, test := range table {
		got := test.e.max()
		if got != test.want {
			t.Errorf("test %v: %v.Max()\nGOT %v\nWANT %v", i, test.e, got, test.want)
		}
	}
}

func TestEdgeListIncrement(t *testing.T) {
	table := []struct {
		e    EdgeList
		n    int
		want EdgeList
	}{
		{
			e:    EdgeList{},
			n:    0,
			want: EdgeList{},
		},
		{
			e:    EdgeList{{From: 0, To: 0, Weight: 1.0}},
			n:    0,
			want: EdgeList{{From: 0, To: 0, Weight: 1.0}},
		},
		{
			e:    EdgeList{{From: 1, To: 1, Weight: 1.0}, {From: 1, To: 1, Weight: 1.0}},
			n:    1,
			want: EdgeList{{From: 2, To: 2, Weight: 1.0}, {From: 2, To: 2, Weight: 1.0}},
		},
		{
			e: EdgeList{
				{From: 10, To: 12, Weight: 1.0},
				{From: 55, To: 94, Weight: 1.0},
				{From: 0, To: 15, Weight: 1.0},
				{From: 1, To: 1, Weight: 1.0},
			},
			n: -1,
			want: EdgeList{
				{From: 9, To: 11, Weight: 1.0},
				{From: 54, To: 93, Weight: 1.0},
				{From: -1, To: 14, Weight: 1.0},
				{From: 0, To: 0, Weight: 1.0},
			},
		},
		{
			e: EdgeList{
				{From: -1, To: 1651, Weight: 1.0},
				{From: 55, To: 65, Weight: 1.0},
				{From: 10, To: 1, Weight: 1.0},
				{From: 15, To: 99, Weight: 1.0},
				{From: 65, To: 54, Weight: 1.0},
				{From: 1000000000, To: 0, Weight: 1.0},
				{From: 0, To: 8, Weight: 1.0},
				{From: 15, To: 44, Weight: 1.0},
			},
			n: 10,
			want: EdgeList{
				{From: 9, To: 1661, Weight: 1.0},
				{From: 65, To: 75, Weight: 1.0},
				{From: 20, To: 11, Weight: 1.0},
				{From: 25, To: 109, Weight: 1.0},
				{From: 75, To: 64, Weight: 1.0},
				{From: 1000000010, To: 10, Weight: 1.0},
				{From: 10, To: 18, Weight: 1.0},
				{From: 25, To: 54, Weight: 1.0},
			},
		},
	}
	for i, test := range table {
		got := increment(test.e, test.n)
		if !slices.Equal(Sort(got), Sort(test.want)) {
			t.Errorf("test %v: %v,Increment(%v)\nGOT %v\nWANT %v", i, test.e, test.n, got, test.want)
		}
	}
}
