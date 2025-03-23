package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
)

func grep(pattern, path string, recursive bool) error {
    re, err := regexp.Compile(pattern)
    if err != nil {
        return fmt.Errorf("invalid regex: %v", err)
    }

    return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() && !recursive {
            return filepath.SkipDir
        }
        if info.IsDir() {
            return nil
        }

        file, err := os.Open(filePath)
        if err != nil {
            return err
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        lineNum := 0
        for scanner.Scan() {
            lineNum++
            line := scanner.Text()
            if re.MatchString(line) {
                fmt.Printf("%s:%d: %s\n", filePath, lineNum, line)
            }
        }
        return scanner.Err()
    })
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run grep.go <pattern> <path> [-r]")
        os.Exit(1)
    }
    pattern := os.Args[1]
    path := os.Args[2]
    recursive := len(os.Args) > 3 && os.Args[3] == "-r"
    if err := grep(pattern, path, recursive); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}