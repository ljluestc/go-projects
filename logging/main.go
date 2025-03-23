package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "path/filepath"

    "github.com/mdobak/go-xerrors"
)

type User struct {
    ID        string
    Email     string
    Password  string // Sensitive field
}

func (u User) LogValue() slog.Value {
    return slog.GroupValue(
        slog.String("id", u.ID),
        slog.String("email", "[REDACTED]"), // Hide email
    )
}

func main() {
    // Determine environment
    appEnv := os.Getenv("APP_ENV")
    if appEnv == "" {
        appEnv = "development"
    }

    // Handler options
    opts := &slog.HandlerOptions{
        AddSource: true,
        Level:     slog.LevelDebug,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.LevelKey {
                if level, ok := a.Value.Any().(slog.Level); ok {
                    if label, exists := LevelNames[level]; exists {
                        a.Value = slog.StringValue(label)
                    }
                }
            }
            return a
        },
    }

    // Switch handler based on environment
    var handler slog.Handler
    if appEnv == "production" {
        handler = slog.NewJSONHandler(os.Stdout, opts)
    } else {
        handler = NewPrettyHandler(os.Stdout, PrettyHandlerOptions{SlogOpts: *opts})
    }

    // Create context-aware handler
    ctxHandler := &ContextHandler{Handler: handler}
    logger := slog.New(ctxHandler)
    slog.SetDefault(logger)

    // HTTP server setup
    http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
        ctx := AppendCtx(r.Context(), slog.String("request_id", "req-"+strconv.Itoa(rand.Intn(1000))))

        logger.InfoContext(ctx, "Incoming request",
            slog.String("method", r.Method),
            slog.String("path", r.URL.Path),
        )

        user := User{ID: "user-123", Email: "user@example.com", Password: "secret"}
        logger.DebugContext(ctx, "Processing user", slog.Any("user", user))

        // Simulate an error
        err := xerrors.New("Upload failed")
        logger.ErrorContext(ctx, "Error occurred",
            slog.Any("error", fmtErr(err)),
        )

        w.Write([]byte("Upload processed\n"))
    })

    logger.Info("Server starting", slog.String("port", "8080"))
    if err := http.ListenAndServe(":8080", nil); err != nil {
        logger.Error("Server failed", slog.Any("error", err))
        os.Exit(1)
    }
}

// ContextHandler for adding context attributes
type ContextHandler struct {
    slog.Handler
}

type ctxKey string

const slogFields ctxKey = "slog_fields"

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
    if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
        for _, v := range attrs {
            r.AddAttrs(v)
        }
    }
    return h.Handler.Handle(ctx, r)
}

func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
    if parent == nil {
        parent = context.Background()
    }
    if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
        v = append(v, attr)
        return context.WithValue(parent, slogFields, v)
    }
    v := []slog.Attr{}
    v = append(v, attr)
    return context.WithValue(parent, slogFields, v)
}

// Error formatting with stack traces
type stackFrame struct {
    Func   string `json:"func"`
    Source string `json:"source"`
    Line   int    `json:"line"`
}

func fmtErr(err error) slog.Value {
    var groupValues []slog.Attr
    groupValues = append(groupValues, slog.String("msg", err.Error()))

    if frames := marshalStack(err); frames != nil {
        groupValues = append(groupValues, slog.Any("trace", frames))
    }
    return slog.GroupValue(groupValues...)
}

func marshalStack(err error) []stackFrame {
    trace := xerrors.StackTrace(err)
    if len(trace) == 0 {
        return nil
    }
    frames := trace.Frames()
    s := make([]stackFrame, len(frames))
    for i, v := range frames {
        s[i] = stackFrame{
            Source: filepath.Join(filepath.Base(filepath.Dir(v.File)), filepath.Base(v.File)),
            Func:   filepath.Base(v.Function),
            Line:   v.Line,
        }
    }
    return s
}