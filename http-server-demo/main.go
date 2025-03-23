package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
)

const (
    port = "4221"
)

type Request struct {
    Method  string
    Path    string
    Version string
    Headers map[string]string
}

type Response struct {
    StatusCode int
    StatusText string
    Headers    map[string]string
    Body       string
}

// parseRequest parses an HTTP request from a reader.
func parseRequest(reader *bufio.Reader) (*Request, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return nil, err
    }
    parts := strings.Fields(line)
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid request line: %s", line)
    }

    req := &Request{
        Method:  parts[0],
        Path:    parts[1],
        Version: parts[2],
        Headers: make(map[string]string),
    }

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return nil, err
        }
        line = strings.TrimSpace(line)
        if line == "" {
            break
        }
        headerParts := strings.SplitN(line, ":", 2)
        if len(headerParts) != 2 {
            continue
        }
        req.Headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
    }

    return req, nil
}

// handleRequest processes an HTTP request and returns a response.
func handleRequest(req *Request) *Response {
    resp := &Response{
        Headers: make(map[string]string),
    }

    switch req.Path {
    case "/":
        resp.StatusCode = 200
        resp.StatusText = "OK"
    case "/echo":
        resp.StatusCode = 404
        resp.StatusText = "Not Found"
    default:
        if strings.HasPrefix(req.Path, "/echo/") {
            resp.StatusCode = 200
            resp.StatusText = "OK"
            resp.Body = strings.TrimPrefix(req.Path, "/echo/")
            resp.Headers["Content-Type"] = "text/plain"
            resp.Headers["Content-Length"] = fmt.Sprintf("%d", len(resp.Body))
        } else {
            resp.StatusCode = 404
            resp.StatusText = "Not Found"
        }
    }

    return resp
}

// formatResponse formats an HTTP response as a string.
func formatResponse(resp *Response) string {
    var builder strings.Builder
    builder.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", resp.StatusCode, resp.StatusText))
    for key, value := range resp.Headers {
        builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
    }
    builder.WriteString("\r\n")
    builder.WriteString(resp.Body)
    return builder.String()
}

// handleConnection handles a single TCP connection.
func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    req, err := parseRequest(reader)
    if err != nil {
        log.Printf("Error parsing request: %v", err)
        conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
        return
    }

    resp := handleRequest(req)
    conn.Write([]byte(formatResponse(resp)))
}

// StartServer starts the HTTP server on the specified port.
func StartServer(port string) error {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return fmt.Errorf("failed to listen on port %s: %v", port, err)
    }
    defer listener.Close()

    log.Printf("Server listening on port %s", port)
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        go handleConnection(conn)
    }
}

func main() {
    if err := StartServer(port); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}