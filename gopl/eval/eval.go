package main

import (
	"fmt"
	"log"
	"math"
)

func (v Var) eval(env Env) float64 {
	return env[v]
}

func (l literal) eval(_ Env) float64 {
	return float64(l)
}

func (u unary) eval(env Env) float64 {
	op := u.op
	switch op {
	case '+':
		return +u.x.eval(env)
	case '-':
		return -u.x.eval(env)
	}
	log.Fatalf("Unsupported unary operator: %q\n", op)
	return 0
}

func (b binary) eval(env Env) float64 {
	op := b.op
	switch op {
	case '+':
		return b.x.eval(env) + b.y.eval(env)
	case '-':
		return b.x.eval(env) - b.y.eval(env)
	case '*':
		return b.x.eval(env) * b.y.eval(env)
	case '/':
		return b.x.eval(env) / b.y.eval(env)
	}
	log.Fatalf("unsupported binary operator: %q\n", op)
	return 0
}

func (c call) eval(env Env) float64 {
	f := c.f
	switch f {
	case "pow":
		return math.Pow(c.args[0].eval(env), c.args[1].eval(env))
	case "sin":
		return math.Sin(c.args[0].eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].eval(env))
  case "min":
    return math.Min(c.args[0].eval(env), c.args[1].eval(env))
  case "max":
    return math.Max(c.args[0].eval(env), c.args[1].eval(env))
	}
	log.Fatalf("unsupported function call: %q\n", f)
	return 0
}

func (exp *Expression) CheckEnvAndEval(env Env) (float64, error) {
	if exp == nil {
		return 0, fmt.Errorf("nil expression")
	}
	if err := exp.CheckEnv(env); err != nil {
		return 0, err
	}
	return exp.Eval(env)
}

func (exp *Expression) CheckExpWithEnvAndEval(env Env) (float64, error) {
	if exp == nil {
		return 0, fmt.Errorf("nil expression")
	}
	if err := exp.CheckExpFromEnv(env); err != nil {
		return 0, err
	}
	return exp.Eval(env)
}

func (exp *Expression) Eval(env Env) (float64, error) {
	if exp == nil {
		return 0, fmt.Errorf("nil expression")
	}
	return exp.e.eval(env), nil
}
