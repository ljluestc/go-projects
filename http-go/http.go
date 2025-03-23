package main

import (
    "fmt"
    "log"
    "net"
    "strings"
)

type Handler func(string) string

var routes = map[string]Handler{
    "/": func(_ string) string {
        return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
    },
    "/about": func(_ string) string {
        return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nAbout Page"
    },
}

func handleConn(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    n, _ := conn.Read(buf)
    request := string(buf[:n])
    path := strings.Split(strings.Split(request, " ")[1], "?")[0]
    handler, ok := routes[path]
    if !ok {
        conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
        return
    }
    response := handler(request)
    conn.Write([]byte(response))
}

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
    fmt.Println("HTTP Server running on :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleConn(conn)
    }
}