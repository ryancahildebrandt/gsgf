// -*- coding: utf-8 -*-

// Created on Sun Mar 16 05:13:48 PM EDT 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"testing"
)

func _() {
	fmt.Println("")
}

func TestValidateExportDir(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateExportDir(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ValidateExportDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutFile(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateOutFile(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateInFile(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateInFile(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("ValidateInFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
