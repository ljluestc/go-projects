package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
)

// matchPattern checks if a string matches an advanced regex pattern
func matchPattern(pattern, text string, caseInsensitive bool) bool {
    if caseInsensitive {
        pattern = strings.ToLower(pattern)
        text = strings.ToLower(text)
    }
    p := []rune(pattern)
    t := []rune(text)
    return matchHelper(p, t, 0, 0)
}

// matchHelper handles recursive pattern matching
func matchHelper(pattern, text []rune, pIdx, tIdx int) bool {
    // Base case: both pattern and text exhausted
    if pIdx >= len(pattern) {
        return tIdx >= len(text)
    }

    // Handle ^ at the start
    if pIdx == 0 && pattern[pIdx] == '^' {
        return matchHelper(pattern, text, pIdx+1, 0)
    }

    // If text is exhausted but pattern isn't, check for optional endings
    if tIdx >= len(text) {
        if pIdx+1 < len(pattern) && (pattern[pIdx+1] == '*' || pattern[pIdx+1] == '?') {
            return matchHelper(pattern, text, pIdx+2, tIdx)
        }
        if pIdx == len(pattern)-1 && pattern[pIdx] == '$' {
            return true
        }
        return false
    }

    // Handle $ at the end
    if pIdx == len(pattern)-1 && pattern[pIdx] == '$' {
        return tIdx == len(text)
    }

    // Handle character class [a-z]
    if pIdx < len(pattern) && pattern[pIdx] == '[' {
        endIdx := pIdx + 1
        for endIdx < len(pattern) && pattern[endIdx] != ']' {
            endIdx++
        }
        if endIdx >= len(pattern) {
            return false // Malformed pattern
        }
        rangeStart, rangeEnd := pattern[pIdx+1], pattern[pIdx+3]
        if tIdx < len(text) && text[tIdx] >= rangeStart && text[tIdx] <= rangeEnd {
            return matchHelper(pattern, text, endIdx+1, tIdx+1)
        }
        return false
    }

    // Look ahead for quantifiers
    nextQuantifier := false
    if pIdx+1 < len(pattern) && (pattern[pIdx+1] == '*' || pattern[pIdx+1] == '+' || pattern[pIdx+1] == '?') {
        nextQuantifier = true
    }

    if nextQuantifier {
        switch pattern[pIdx+1] {
        case '*': // Zero or more
            return matchHelper(pattern, text, pIdx+2, tIdx) ||
                ((pattern[pIdx] == '.' || pattern[pIdx] == text[tIdx]) && matchHelper(pattern, text, pIdx, tIdx+1))
        case '+': // One or more
            return (pattern[pIdx] == '.' || pattern[pIdx] == text[tIdx]) &&
                (matchHelper(pattern, text, pIdx+2, tIdx+1) || matchHelper(pattern, text, pIdx, tIdx+1))
        case '?': // Zero or one
            return matchHelper(pattern, text, pIdx+2, tIdx) ||
                ((pattern[pIdx] == '.' || pattern[pIdx] == text[tIdx]) && matchHelper(pattern, text, pIdx+2, tIdx+1))
        }
    }

    // Normal character or '.'
    if pattern[pIdx] == '.' || pattern[pIdx] == text[tIdx] {
        return matchHelper(pattern, text, pIdx+1, tIdx+1)
    }

    return false
}

// grepFile searches for matching lines
func grepFile(pattern, filename string, caseInsensitive bool) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        if matchPattern(pattern, line, caseInsensitive) {
            fmt.Printf("Line %d: %s\n", lineNum, line)
        }
    }
    return scanner.Err()
}

func main() {
    caseInsensitive := flag.Bool("i", false, "Case-insensitive matching")
    flag.Parse()

    if len(flag.Args()) != 2 {
        fmt.Println("Usage: go run grep.go [-i] <pattern> <filename>")
        fmt.Println("Example: go run grep.go -i '^[a-z]+o$' input.txt")
        os.Exit(1)
    }

    pattern := flag.Arg(0)
    filename := flag.Arg(1)

    if err := grepFile(pattern, filename, *caseInsensitive); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}