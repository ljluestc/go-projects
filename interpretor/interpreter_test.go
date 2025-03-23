package main

import (
    "testing"
)

func TestInterpreter(t *testing.T) {
    tests := []struct {
        input string
        want  int
    }{
        {"2 + 3", 5},
        {"3 * 4", 12},
        {"2 + 3 * 4", 14},    // Precedence: 3 * 4 = 12, then 2 + 12 = 14
        {"10 - 2 / 2", 9},    // Precedence: 2 / 2 = 1, then 10 - 1 = 9
        {"5 * 2 + 3", 13},    // Precedence: 5 * 2 = 10, then 10 + 3 = 13
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            lexer := NewLexer(tt.input)
            tokens := lexer.Lex()
            parser := NewParser(tokens)
            ast := parser.ParseExpression()
            interpreter := Interpreter{}
            got := interpreter.Evaluate(ast)
            if got != tt.want {
                t.Errorf("Evaluate(%q) = %d, want %d", tt.input, got, tt.want)
            }
        })
    }
}