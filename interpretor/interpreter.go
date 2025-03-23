package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

// TokenType represents the type of token
type TokenType int

const (
    TokenNumber TokenType = iota
    TokenPlus
    TokenMinus
    TokenMultiply
    TokenDivide
    TokenEOF
)

// Token represents a lexical token
type Token struct {
    Type    TokenType
    Literal string
}

// Lexer tokenizes the input string
type Lexer struct {
    input  string
    pos    int
    tokens []Token
}

func NewLexer(input string) *Lexer {
    return &Lexer{input: strings.TrimSpace(input), pos: 0}
}

func (l *Lexer) Lex() []Token {
    for l.pos < len(l.input) {
        switch l.input[l.pos] {
        case '+':
            l.tokens = append(l.tokens, Token{TokenPlus, "+"})
            l.pos++
        case '-':
            l.tokens = append(l.tokens, Token{TokenMinus, "-"})
            l.pos++
        case '*':
            l.tokens = append(l.tokens, Token{TokenMultiply, "*"})
            l.pos++
        case '/':
            l.tokens = append(l.tokens, Token{TokenDivide, "/"})
            l.pos++
        case ' ':
            l.pos++
        default:
            if isDigit(l.input[l.pos]) {
                start := l.pos
                for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
                    l.pos++
                }
                l.tokens = append(l.tokens, Token{TokenNumber, l.input[start:l.pos]})
            } else {
                fmt.Fprintf(os.Stderr, "Invalid character: %c\n", l.input[l.pos])
                os.Exit(1)
            }
        }
    }
    l.tokens = append(l.tokens, Token{TokenEOF, ""})
    return l.tokens
}

func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}

// NodeType represents the type of AST node
type NodeType int

const (
    NodeNumber NodeType = iota
    NodeBinaryOp
)

// Node represents an AST node
type Node struct {
    Type     NodeType
    Value    int       // For numbers
    Operator TokenType // For binary operations
    Left     *Node
    Right    *Node
}

// Parser builds an AST from tokens
type Parser struct {
    tokens []Token
    pos    int
}

func NewParser(tokens []Token) *Parser {
    return &Parser{tokens: tokens, pos: 0}
}

// ParseExpression handles precedence (multiplication/division before addition/subtraction)
func (p *Parser) ParseExpression() *Node {
    node := p.parseTerm()
    for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == TokenPlus || p.tokens[p.pos].Type == TokenMinus) {
        op := p.tokens[p.pos]
        p.pos++
        right := p.parseTerm()
        node = &Node{Type: NodeBinaryOp, Operator: op.Type, Left: node, Right: right}
    }
    return node
}

func (p *Parser) parseTerm() *Node {
    node := p.parseFactor()
    for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == TokenMultiply || p.tokens[p.pos].Type == TokenDivide) {
        op := p.tokens[p.pos]
        p.pos++
        right := p.parseFactor()
        node = &Node{Type: NodeBinaryOp, Operator: op.Type, Left: node, Right: right}
    }
    return node
}

func (p *Parser) parseFactor() *Node {
    token := p.tokens[p.pos]
    p.pos++
    if token.Type == TokenNumber {
        value, _ := strconv.Atoi(token.Literal)
        return &Node{Type: NodeNumber, Value: value}
    }
    fmt.Fprintf(os.Stderr, "Unexpected token: %v\n", token)
    os.Exit(1)
    return nil
}

// Interpreter evaluates the AST
type Interpreter struct{}

func (i *Interpreter) Evaluate(node *Node) int {
    if node.Type == NodeNumber {
        return node.Value
    }
    left := i.Evaluate(node.Left)
    right := i.Evaluate(node.Right)
    switch node.Operator {
    case TokenPlus:
        return left + right
    case TokenMinus:
        return left - right
    case TokenMultiply:
        return left * right
    case TokenDivide:
        if right == 0 {
            fmt.Fprintln(os.Stderr, "Error: Division by zero")
            os.Exit(1)
        }
        return left / right
    default:
        fmt.Fprintf(os.Stderr, "Unknown operator: %v\n", node.Operator)
        os.Exit(1)
        return 0
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run interpreter.go <expression>")
        fmt.Println("Example: go run interpreter.go \"2 + 3 * 4\"")
        os.Exit(1)
    }

    input := os.Args[1]
    lexer := NewLexer(input)
    tokens := lexer.Lex()
    parser := NewParser(tokens)
    ast := parser.ParseExpression()
    interpreter := Interpreter{}
    result := interpreter.Evaluate(ast)
    fmt.Println(result)
}