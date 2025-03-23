package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "sync"
)

func main() {
    fmt.Println("GoShell - Type 'exit' to quit")
    scanner := bufio.NewScanner(os.Stdin)
    var wg sync.WaitGroup

    for {
        fmt.Print("$ ")
        if !scanner.Scan() {
            break
        }
        input := strings.TrimSpace(scanner.Text())
        if input == "exit" {
            break
        }
        if input == "" {
            continue
        }

        // Split commands by pipe
        commands := strings.Split(input, "|")
        var prevCmd *exec.Cmd

        for i, cmdStr := range commands {
            args := strings.Fields(strings.TrimSpace(cmdStr))
            if len(args) == 0 {
                continue
            }

            // Handle background execution
            isBackground := false
            if args[len(args)-1] == "&" {
                isBackground = true
                args = args[:len(args)-1]
            }

            cmd := exec.Command(args[0], args[1:]...)
            if i > 0 {
                cmd.Stdin, _ = prevCmd.StdoutPipe()
            }
            if i < len(commands)-1 {
                cmd.Stdout = os.Stdout // Intermediate commands need piping
            } else {
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr
            }

            if isBackground {
                wg.Add(1)
                go func(c *exec.Cmd) {
                    defer wg.Done()
                    if err := c.Run(); err != nil {
                        fmt.Printf("Error: %v\n", err)
                    }
                }(cmd)
            } else {
                if err := cmd.Start(); err != nil {
                    fmt.Printf("Error: %v\n", err)
                    break
                }
                prevCmd = cmd
                if i == len(commands)-1 {
                    cmd.Wait()
                }
            }
        }
    }
    wg.Wait()
}