package main

import (
	"errors"
	"go/ast"
	"go/constant"
	"go/parser"
	"strconv"
)

type ExprRunner struct {
	expr ast.Expr
}

type Vars map[string]string

type ExprContext struct {
	vars Vars
}

func parseExpr(e string) (ast.Expr, error) {
	expr, err := parser.ParseExpr(e)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func evalExpr(expr ast.Expr, ctx *ExprContext) (constant.Value, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return evalBasicLit(e, ctx)
	case *ast.BinaryExpr:
		return evalBinaryExpr(e, ctx)
	case *ast.Ident:
		return evalIdent(e, ctx)
	case *ast.ParenExpr:
		return evalExpr(e.X, ctx)
	}

	return constant.MakeUnknown(), errors.New("unknown expr")
}

func evalBasicLit(expr *ast.BasicLit, ctx *ExprContext) (constant.Value, error) {
	return constant.MakeFromLiteral(expr.Value, expr.Kind, 0), nil
}

func evalBinaryExpr(expr *ast.BinaryExpr, ctx *ExprContext) (constant.Value, error) {
	x, err := evalExpr(expr.X, ctx)
	if err != nil {
		return constant.MakeUnknown(), err
	}

	y, err := evalExpr(expr.Y, ctx)
	if err != nil {
		return constant.MakeUnknown(), err
	}

	return constant.BinaryOp(x, expr.Op, y), nil
}

func evalIdent(expr *ast.Ident, ctx *ExprContext) (constant.Value, error) {
	name, ok := ctx.vars[expr.Name]
	if !ok {
		return constant.MakeUnknown(), errors.New("unknown variable name")
	}

	n, err := strconv.ParseFloat(name, 64)
	if err != nil {
		return constant.MakeUnknown(), errors.New("variable type must be numeric")
	}

	return constant.MakeFloat64(n), nil
}
