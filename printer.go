package gocheck

import (
    "bytes"
    "go/ast"
    "go/parser"
    "go/token"
    "go/printer"
    "os"
)

func indent(s, with string) (r string) {
    eol := true
    for i := 0; i != len(s); i++ {
        c := s[i]
        switch {
        case eol && c == '\n' || c == '\r':
        case c == '\n' || c == '\r':
            eol = true
        case eol:
            eol = false
            s = s[:i] + with + s[i:]
            i += len(with)
        }
    }
    return s
}

func printLine(filename string, line int) (string, os.Error) {
    fset := token.NewFileSet()
    file, err := os.Open(filename, os.O_RDONLY, 0)
    if err != nil {
        return "", err
    }
    fnode, err := parser.ParseFile(fset, filename, file, 0)
    if err != nil {
        return "", err
    }
    config := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
    lp := &linePrinter{fset: fset, line: line, config: config}
    ast.Walk(lp, fnode)
    return lp.output.String(), nil
}

type linePrinter struct {
    fset *token.FileSet
    line int
    output bytes.Buffer
    stmt ast.Stmt
    config *printer.Config
}

func (lp *linePrinter) Visit(n ast.Node) (w ast.Visitor) {
    if n == nil {
        return nil
    }
    if stmt, ok := n.(ast.Stmt); ok {
        lp.stmt = stmt
    }
    p := n.Pos()
    line := lp.fset.Position(p).Line
    if line == lp.line {
        lp.trim(lp.stmt)
        lp.config.Fprint(&lp.output, lp.fset, lp.stmt)
        return nil
    }
    return lp
}

func (lp *linePrinter) trim(n ast.Node) bool {
    stmt, ok := n.(ast.Stmt)
    if !ok {
        return true
    }
    p := n.Pos()
    line := lp.fset.Position(p).Line
    if line != lp.line {
        return false
    }
    switch stmt := stmt.(type) {
    case *ast.IfStmt:
        stmt.Body = lp.trimBlock(stmt.Body)
    case *ast.SwitchStmt:
        stmt.Body = lp.trimBlock(stmt.Body)
    case *ast.TypeSwitchStmt:
        stmt.Body = lp.trimBlock(stmt.Body)
    case *ast.CaseClause:
        stmt.Body = lp.trimList(stmt.Body)
    case *ast.TypeCaseClause:
        stmt.Body = lp.trimList(stmt.Body)
    case *ast.CommClause:
        stmt.Body = lp.trimList(stmt.Body)
    case *ast.BlockStmt:
        stmt.List = lp.trimList(stmt.List)
    }
    return true
}

func (lp *linePrinter) trimBlock(stmt *ast.BlockStmt) *ast.BlockStmt {
    if !lp.trim(stmt) {
        return lp.emptyBlock(stmt)
    }
    stmt.Rbrace = stmt.Lbrace
    return stmt
}

func (lp *linePrinter) trimList(stmts []ast.Stmt) []ast.Stmt {
    for i := 0; i != len(stmts); i++ {
        if !lp.trim(stmts[i]) {
            stmts[i] = lp.emptyStmt(stmts[i])
            break
        }
    }
    return stmts
}

func (lp *linePrinter) emptyStmt(n ast.Node) *ast.ExprStmt {
    return &ast.ExprStmt{&ast.Ellipsis{n.Pos(), nil}}
}

func (lp *linePrinter) emptyBlock(n ast.Node) *ast.BlockStmt {
    p := n.Pos()
    return &ast.BlockStmt{p, []ast.Stmt{lp.emptyStmt(n)}, p}
}
