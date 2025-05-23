package termlog

import "context"

// contextKey is the type used for storing termlog in context
type contextKey struct{}

// NewContext creates a new context with an included Logger
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// FromContext retrieves a Logger from a context. If no logger is present, we
// return a new silenced logger that will produce no output.
func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(contextKey{}).(Logger)
	if !ok {
		l := NewLog()
		l.Quiet()
		return l
	}
	return logger
}
