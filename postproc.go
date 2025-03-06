// -*- coding: utf-8 -*-

// Created on Wed Mar  5 11:33:35 AM EST 2025
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"fmt"
	"regexp"
	"strings"
)

func RemoveEndSpaces(p []string) []string {
	for i := range p {
		p[i] = strings.Trim(p[i], "\t\r\n ")
	}
	return p
}

func RemoveMultipleSpaces(p []string) []string {
	// will eat \n\t etc and replace w " "
	for i := range p {
		p[i] = strings.Join(strings.Fields(p[i]), " ")
	}
	return p
}

func RenderTabs(p []string) []string {
	for i := range p {
		p[i] = strings.Replace(p[i], `\t`, "\t", -1)
	}
	return p
}

func RenderNewLines(p []string) []string {
	for i := range p {
		p[i] = strings.Replace(p[i], `\n`, "\n", -1)
	}
	return p
}

func RemoveTags(p []string) []string {
	for i := range p {
		p[i] = regexp.MustCompile(`\{.*?\}`).ReplaceAllString(p[i], "")
	}
	return p
}

func WrapTags(p []string, prefix string, suffix string) []string {
	var seen map[string]struct{} = make(map[string]struct{})

	for i := range p {
		tags := regexp.MustCompile(`(\{.*?\})`).FindAllString(p[i], -1)
		for _, tag := range tags {
			_, ok := seen[tag]
			if !ok {
				seen[tag] = struct{}{}
				p[i] = strings.ReplaceAll(p[i], tag, fmt.Sprint(prefix, tag, suffix))
			}
		}
	}
	return p
}

func CollectTags(p []string, c string) []string {
	var b strings.Builder

	for i := range p {
		b.WriteString(regexp.MustCompile(`\{.*?\}`).ReplaceAllString(p[i], ""))
		tags := regexp.MustCompile(`(\{.*?\})`).FindAllString(p[i], -1)
		if len(tags) > 0 {
			b.WriteString(c)
		}
		for _, tag := range tags {
			b.WriteString(tag)
			b.WriteString(",")
		}
		p[i] = b.String()
		b.Reset()
	}

	return p
}

func WrapProductions(p []string, prefix string, suffix string) []string {
	for i := range p {
		p[i] = fmt.Sprint(prefix, p[i], suffix)
	}
	return p
}
