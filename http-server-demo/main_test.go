package main

import (
    "bufio"
    "net"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestServerRootRoute(t *testing.T) {
    go main() // Start server in background
    time.Sleep(100 * time.Millisecond) // Wait for server to start

    conn, err := net.Dial("tcp", "localhost:4221")
    assert.NoError(t, err)
    defer conn.Close()

    _, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"))
    assert.NoError(t, err)

    reader := bufio.NewReader(conn)
    line, err := reader.ReadString('\n')
    assert.NoError(t, err)
    assert.Equal(t, "HTTP/1.1 200 OK\r\n", line)

    // Read headers until empty line
    for {
        line, err = reader.ReadString('\n')
        assert.NoError(t, err)
        if line == "\r\n" {
            break
        }
    }
}

func TestServerEchoRoute(t *testing.T) {
    go main()
    time.Sleep(100 * time.Millisecond)

    conn, err := net.Dial("tcp", "localhost:4221")
    assert.NoError(t, err)
    defer conn.Close()

    _, err = conn.Write([]byte("GET /echo/hello HTTP/1.1\r\nHost: localhost\r\n\r\n"))
    assert.NoError(t, err)

    reader := bufio.NewReader(conn)
    line, err := reader.ReadString('\n')
    assert.NoError(t, err)
    assert.Equal(t, "HTTP/1.1 200 OK\r\n", line)

    // Read headers
    headers := make(map[string]string)
    for {
        line, err = reader.ReadString('\n')
        assert.NoError(t, err)
        if line == "\r\n" {
            break
        }
        parts := strings.SplitN(line, ":", 2)
        assert.Len(t, parts, 2)
        headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
    }

    assert.Equal(t, "text/plain", headers["Content-Type"])
    assert.Equal(t, "5", headers["Content-Length"])

    // Read body
    body, err := reader.ReadString('\n')
    assert.NoError(t, err)
    assert.Equal(t, "hello", body)
}

func TestServerNotFound(t *testing.T) {
    go main()
    time.Sleep(100 * time.Millisecond)

    conn, err := net.Dial("tcp", "localhost:4221")
    assert.NoError(t, err)
    defer conn.Close()

    _, err = conn.Write([]byte("GET /unknown HTTP/1.1\r\nHost: localhost\r\n\r\n"))
    assert.NoError(t, err)

    reader := bufio.NewReader(conn)
    line, err := reader.ReadString('\n')
    assert.NoError(t, err)
    assert.Equal(t, "HTTP/1.1 404 Not Found\r\n", line)
}

func TestInvalidRequest(t *testing.T) {
    go main()
    time.Sleep(100 * time.Millisecond)

    conn, err := net.Dial("tcp", "localhost:4221")
    assert.NoError(t, err)
    defer conn.Close()

    _, err = conn.Write([]byte("INVALID\r\n"))
    assert.NoError(t, err)

    reader := bufio.NewReader(conn)
    line, err := reader.ReadString('\n')
    assert.NoError(t, err)
    assert.Equal(t, "HTTP/1.1 400 Bad Request\r\n", line)
}