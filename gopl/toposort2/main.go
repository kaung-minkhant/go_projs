package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

var prereqs = map[string][]string{
	"algorithms":     {"data structures"},
	"calculus":       {"linear algebra"},
	// "linear algebra": {"data structures"},
	"linear algebra": {"calculus"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(input []string, parents []string)

	visitAll = func(input []string, parents []string) {
		for _, item := range input {
			if _, ok := seen[item]; !ok {
				// fmt.Println("Visiting: ", item, " with parents ", parents)
				seen[item] = true
				for _, dep := range m[item] {
          index := slices.Index(parents, dep)
					if index != -1 {
            // fmt.Println("index: ", index)
            cycle := strings.Join(parents[index:], " -> ")
            cycle += fmt.Sprintf(" -> %s -> %s", item, parents[index])
            log.Fatal("Contains cycles: ", cycle)
					}
				}
				visitAll(m[item], append(parents, item))
				order = append(order, item)
			}
		}
	}

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	visitAll(keys, nil)
	return order
}

func main() {
	order := topoSort(prereqs)
	for i, course := range order {
		fmt.Printf("%-2d: %s\n", i+1, course)
	}
}
