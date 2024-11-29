package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
)

type Person struct {
	Name string
	Age  int
}

type ByColumn struct {
	p          []*Person
	columns    []compareFunc
	maxColumns int
}

type compareFunc func(a, b *Person) comparasm

type comparasm int

const (
	eq comparasm = iota
	lt
	gt
)

func (s *ByColumn) ByName(a, b *Person) comparasm {
	if a.Name == b.Name {
		return eq
	}
	if a.Name < b.Name {
		return lt
	}
	return gt
}

func (s *ByColumn) ByAge(a, b *Person) comparasm {
	if a.Age == b.Age {
		return eq
	}
	if a.Age < b.Age {
		return lt
	}
	return gt
}

func (c *ByColumn) Compare(f compareFunc) {
	c.columns = append([]compareFunc{f}, c.columns...)

	if len(c.columns) > c.maxColumns {
		c.columns = c.columns[:c.maxColumns]
	}
}

func (c *ByColumn) Len() int {
	return len(c.p)
}

func (c *ByColumn) Swap(i, j int) {
	c.p[i], c.p[j] = c.p[j], c.p[i]
}

func (c *ByColumn) Less(i, j int) bool {
	for _, f := range c.columns {
		result := f(c.p[i], c.p[j])
		switch result {
		case eq:
			continue
		case lt:
			return true
		case gt:
			return false
		}
	}
	return false
}

func NewSortByColumn(p []*Person, maxColumn int) *ByColumn {
	return &ByColumn{
		p:          p,
		maxColumns: maxColumn,
	}
}

func printPersonTable(p []*Person) {
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	const format = "%v\t%v\t\n"
	fmt.Fprintf(tw, format, "Name", "Age")
	fmt.Fprintf(tw, format, "-----", "-----")
	for _, person := range p {
		fmt.Fprintf(tw, format, person.Name, person.Age)
	}
	tw.Flush()
	fmt.Println()
}

var p = []*Person{
	{"kaung", 25},
	{"shunn", 22},
	{"kaung", 24},
}

const htmlTemplate = `<html>
  <head>
    <title>Sort in go</title>
  </head>
  <body>
    <table>
      <thead>
        <tr>
          <th><a href="?sort=name">Name</a></th>
          <th><a href="?sort=age">Age</a></th>
        </tr>
      </thead>
      <tbody>
        {{range .}}
          <tr>
            <td>{{.Name}}</td>
            <td>{{.Age}}</td>
          </tr>
        {{end}}
      </tbody>
    </table>
  </body>
</html>`

var temp = template.Must(template.New("temp").Parse(htmlTemplate))

func (s *ByColumn) handleSort(w http.ResponseWriter, r *http.Request) {
	sortKey := r.URL.Query().Get("sort")
	switch sortKey {
	case "name":
		s.Compare(s.ByName)
		sort.Sort(s)
		temp.Execute(w, s.p)
	case "age":
		s.Compare(s.ByAge)
		sort.Sort(s)
		temp.Execute(w, s.p)
  default:
    temp.Execute(w, s.p)
	}
}

func main() {
	sortByColumn := NewSortByColumn(p, 2)

	http.HandleFunc("/", sortByColumn.handleSort)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
