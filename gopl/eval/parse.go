package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

type lexer struct {
	s     scanner.Scanner
	token rune
}

type lexPanic string

func (l *lexer) next() {
	l.token = l.s.Scan()
}

func (l *lexer) text() string {
	return l.s.TokenText()
}

func (l *lexer) describe() string {
	switch l.token {
	case scanner.EOF:
		return fmt.Sprint("end of input")
	case scanner.Ident:
		return fmt.Sprintf("identifier: %q", l.text())
	case scanner.Int, scanner.Float:
		return fmt.Sprintf("number: %q", l.text())
	}
	return fmt.Sprintf("%c", rune(l.token))
}

func newLexer(src io.Reader) *lexer {
	newScanner := new(scanner.Scanner)
	newScanner = newScanner.Init(src)
	newScanner.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanInts | scanner.SkipComments
	return &lexer{
		s: *newScanner,
	}
}

func precedence(operator rune) int {
	switch operator {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

func Parse(inputExpr string) (_ *Expression, err error) {
	defer func() {
		switch e := recover().(type) {
		case nil:
		case lexPanic:
			err = fmt.Errorf("parse error: %s", e)
		default:
			panic(e)
		}
	}()
	l := newLexer(strings.NewReader(inputExpr))
	l.next()
	e := l.exprParse()
	if l.token != scanner.EOF {
		return nil, fmt.Errorf("parsing failed: unexpected: %s", l.describe())
	}
	vars := make(map[Var]bool)
	if err := e.check(vars); err != nil {
		return nil, err
	}
	return &Expression{
		e:    e,
		vars: vars,
	}, nil
}

func (l *lexer) exprParse() expr {
	return l.binaryParse(0)
}

func (l *lexer) binaryParse(preced int) expr {
	lhs := l.unaryParse()
	// fmt.Printf("lhs: %#v with current token: %q\n", lhs, l.token)
	for prece := precedence(l.token); prece > preced; prece-- {
		for precedence(l.token) == prece {
			op := l.token
			l.next()
			rhs := l.binaryParse(prece)
			// fmt.Printf("rhs: %#v with current token: %q\n", rhs, l.token)
			lhs = binary{
				op: op,
				x:  lhs,
				y:  rhs,
			}
		}
	}
	return lhs
}

func (l *lexer) unaryParse() expr {
	if l.token == '+' || l.token == '-' {
		op := l.token
		l.next()
		return unary{op, l.unaryParse()}
	}
	return l.primaryParse()
}

func (l *lexer) primaryParse() expr {
	switch l.token {
	case scanner.Ident:
		txt := l.text()
		l.next()
		if l.token != '(' {
			return Var(txt)
		}
		var args []expr
		l.next()
		if l.token != ')' {
			for {
				e := l.exprParse()
				args = append(args, e)
				if l.token != ',' {
					break
				}
				l.next()
			}
		}
		if l.token != ')' {
			msg := fmt.Sprintf("want ')', got %s", l.describe())
			panic(lexPanic(msg))
		}
		l.next()

		return call{
			f:    txt,
			args: args,
		}

	case scanner.Int, scanner.Float:
		f, err := strconv.ParseFloat(l.text(), 64)
		if err != nil {
			panic(lexPanic(err.Error()))
		}
		l.next()
		return literal(f)
	case '(':
		l.next()
		e := l.exprParse()
		if l.token != ')' {
			msg := fmt.Sprintf("want ')', got %s", l.describe())
			panic(lexPanic(msg))
		}
		l.next()
		return e
	}
	msg := fmt.Sprintf("unexpected input: %s", l.describe())
	panic(lexPanic(msg))
}
