package main

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	"strings"
)

func (v Var) check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (l literal) check(_ map[Var]bool) error {
	return nil
}

func (u unary) check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected operator: %q", u.op)
	}
	return u.x.check(vars)
}

func (b binary) check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected operator: %q", b.op)
	}

	if e := b.x.check(vars); e != nil {
		return e
	}

	return b.y.check(vars)
}

func (c call) check(vars map[Var]bool) error {
	if !slices.Contains(supportedFunctions, c.f) {
		return fmt.Errorf("unsupported function: %q", c.f)
	}
	argNum := funcArgsNum[c.f]
	if argNum < 0 {
		// minimum args
    minimum := int(math.Abs(float64(argNum)))
		if len(c.args) < minimum {
			return fmt.Errorf("unmatched number of arguments: want minum %d, got %d", minimum, len(c.args))
		}
	}
	if argNum != len(c.args) {
		return fmt.Errorf("unmatched number of arguments: want %d, got %d", argNum, len(c.args))
	}
	for _, arg := range c.args {
		if e := arg.check(vars); e != nil {
			return e
		}

	}
	return nil
}

// check if the variable defined in env is sufficient for the expression
func (exp *Expression) CheckEnv(env Env) error {
	missing := []string{}
	for k := range exp.vars {
		if _, ok := env[k]; !ok {
			missing = append(missing, string(k))
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("incomplete environment: missing %s", printMissing(missing))
	}
	return nil
}

// check if the expression meets the required variables set in env
func (exp *Expression) CheckExpFromEnv(env Env) error {
	missing := []string{}
	for k := range env {
		if _, ok := exp.vars[k]; !ok {
			missing = append(missing, string(k))
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("incomplete expression: missing %s", printMissing(missing))
	}
	return nil
}

func printMissing(missing []string) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, item := range missing {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(item)
	}
	buf.WriteByte(']')
	return buf.String()
}
