package main

import (
    "bytes"
    "context"
    "log/slog"
    "testing"
)

func TestLogging(t *testing.T) {
    var buf bytes.Buffer
    opts := PrettyHandlerOptions{
        SlogOpts: slog.HandlerOptions{Level: slog.LevelDebug},
    }
    handler := NewPrettyHandler(&buf, opts)
    logger := slog.New(&ContextHandler{Handler: handler})

    ctx := AppendCtx(context.Background(), slog.String("request_id", "req-999"))
    logger.InfoContext(ctx, "Test message", slog.String("key", "value"))

    output := buf.String()
    if !bytes.Contains([]byte(output), []byte("INFO")) || !bytes.Contains([]byte(output), []byte("request_id")) {
        t.Errorf("Expected INFO log with request_id, got: %s", output)
    }
}

func TestUserLogValue(t *testing.T) {
    user := User{ID: "user-123", Email: "test@example.com", Password: "secret"}
    val := user.LogValue()
    attrs := val.Group()

    foundID := false
    for _, attr := range attrs {
        if attr.Key == "id" && attr.Value.String() == "user-123" {
            foundID = true
        }
        if attr.Key == "email" && attr.Value.String() == "[REDACTED]" {
            continue
        }
        if attr.Key == "password" {
            t.Errorf("Password field should not be logged")
        }
    }
    if !foundID {
        t.Errorf("Expected user ID in log value")
    }
}