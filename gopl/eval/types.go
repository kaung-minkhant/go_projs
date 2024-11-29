package main

type expr interface {
	eval(env Env) float64
	check(vars map[Var]bool) error
}

type Var string

type literal float64

type unary struct {
	op rune
	x  expr
}

type binary struct {
	op   rune
	x, y expr
}

type call struct {
	f    string
	args []expr
}

type Env map[Var]float64

type Expression struct {
	e expr
  vars map[Var]bool
}
