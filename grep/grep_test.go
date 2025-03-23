package main

import (
    "testing"
)

func TestMatchPattern(t *testing.T) {
    tests := []struct {
        pattern string
        text    string
        want    bool
    }{
        {"hello", "hello", true},
        {"h.llo", "hello", true},
        {"hel*o", "helllo", true},
        {"hel+o", "hello", true},
        {"hel?o", "helo", true},
        {"[a-z]o", "bo", true},
        {"h.*o", "heo", true},
        {"xyz", "abc", false},
    }

    for _, tt := range tests {
        t.Run(tt.pattern+"_"+tt.text, func(t *testing.T) {
            got := matchPattern(tt.pattern, tt.text, false)
            if got != tt.want {
                t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.pattern, tt.text, got, tt.want)
            }
        })
    }
}