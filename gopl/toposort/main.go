package main

import (
	"fmt"
	"slices"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
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
  var visitAll func(input []string)

  visitAll = func (input []string) {
    for _, item := range input {
      if _, ok := seen[item]; !ok {
        seen[item] = true
        visitAll(m[item])
        order = append(order, item)
      }
    }
  }

  var keys []string
  for k := range m {
    keys = append(keys, k)
  }

  slices.Sort(keys)

  visitAll(keys)
  return order
}

func main() {
  order := topoSort(prereqs)
  for i, course := range order {
    fmt.Printf("%-2d: %s\n", i+1, course)
  }
}
