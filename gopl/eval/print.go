package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// func (v Var) String() string {
// 	return string(v)
// }
//
// func (v literal) String() string {
// 	return strconv.FormatFloat(float64(v), 'G', 5, 64)
// }
//
// func (v unary) String() string {
// 	return fmt.Sprintf("%c%s", v.op, printExpr(v.x))
// }
//
// func (v binary) String() string {
// 	return fmt.Sprintf("(%s %c %s)", printExpr(v.x), v.op, printExpr(v.y))
// }
//
// func (v call) String() string {
// 	var argsBuf bytes.Buffer
// 	argsBuf.WriteByte('(')
// 	for i, arg := range v.args {
// 		if i > 0 {
// 			argsBuf.WriteByte(' ')
// 		}
// 		argsBuf.WriteString(printExpr(arg))
// 		if i != len(v.args)-1 {
// 			argsBuf.WriteByte(',')
// 		}
// 	}
// 	argsBuf.WriteByte(')')
// 	return fmt.Sprintf("%s%s", v.f, argsBuf.String())
// }

func printExpr(input expr) string {
	switch v := input.(type) {
	case Var:
		return string(v)
	case literal:
		return strconv.FormatFloat(float64(v), 'G', 5, 64)
	case unary:
		return fmt.Sprintf("%c%s", v.op, printExpr(v.x))
	case binary:
		return fmt.Sprintf("(%s %c %s)", printExpr(v.x), v.op, printExpr(v.y))
	case call:
		var argsBuf bytes.Buffer
		argsBuf.WriteByte('(')
		for i, arg := range v.args {
			if i > 0 {
				argsBuf.WriteByte(' ')
			}
			argsBuf.WriteString(printExpr(arg))
			if i != len(v.args)-1 {
				argsBuf.WriteByte(',')
			}
		}
		argsBuf.WriteByte(')')
		return fmt.Sprintf("%s%s", v.f, argsBuf.String())
	}
	return "unknown expression"
}

func (exp *Expression) String() string {
	str := printExpr(exp.e)
	if strings.HasPrefix(str, "(") {
		str = str[1 : len(str)-1]
	}
	return str
}

func (exp *Expression) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, exp.String())
	return err
}
