package handlers

import "github.com/tomdoesdev/knox/kit/ast"

type HandlerFunc func(args ...any) (ast.Node, error)
