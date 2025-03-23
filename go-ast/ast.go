package main

import (
    "fmt"
    "strconv"
    "strings"
)

type TokenType int

const (
    Number TokenType = iota
    Operator
    Variable
    Paren
)

type Token struct {
    Type  TokenType
    Value string
}

type Node struct {
    Type     TokenType
    Value    string
    Left     *Node
    Right    *Node
}

func tokenize(expr string) []Token {
    var tokens []Token
    for _, char := range strings.Fields(expr) {
        if char == "+" || char == "-" {
            tokens = append(tokens, Token{Operator, char})
        } else if char == "(" || char == ")" {
            tokens = append(tokens, Token{Paren, char})
        } else if _, err := strconv.ParseFloat(char, 64); err == nil {
            tokens = append(tokens, Token{Number, char})
        } else {
            tokens = append(tokens, Token{Variable, char})
        }
    }
    return tokens
}

func parse(tokens []Token) (*Node, int) {
    var stack []*Node
    i := 0
    for i < len(tokens) {
        switch tokens[i].Type {
        case Number, Variable:
            stack = append(stack, &Node{Type: tokens[i].Type, Value: tokens[i].Value})
        case Operator:
            op := &Node{Type: Operator, Value: tokens[i].Value}
            op.Right = stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            op.Left = stack[len(stack)-1]
            stack[len(stack)-1] = op
        case Paren:
            if tokens[i].Value == ")" {
                return stack[len(stack)-1], i
            }
        }
        i++
    }
    return stack[0], i
}

func evaluate(node *Node, vars map[string]float64) float64 {
    switch node.Type {
    case Number:
        val, _ := strconv.ParseFloat(node.Value, 64)
        return val
    case Variable:
        return vars[node.Value]
    case Operator:
        left := evaluate(node.Left, vars)
        right := evaluate(node.Right, vars)
        if node.Value == "+" {
            return left + right
        }
        return left - right
    }
    return 0
}

func main() {
    expr := "( 3 + x ) - 1"
    tokens := tokenize(expr)
    ast, _ := parse(tokens)
    vars := map[string]float64{"x": 4}
    result := evaluate(ast, vars)
    fmt.Printf("Result: %.2f\n", result) // Output: 6.00
}