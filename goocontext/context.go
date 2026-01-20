package goocontext

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Key string

const (
	TraceIdKey Key = "trace_id"
)

func Default(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx
}

func WithValue[T any](ctx context.Context, key string, value T) context.Context {
	return context.WithValue(Default(ctx), Key(key), value)
}

func Value[T any](ctx context.Context, key string) (T, bool) {
	if ctx == nil {
		var zero T
		return zero, false
	}

	v := ctx.Value(Key(key))
	if v == nil {
		var zero T
		return zero, false
	}

	val, ok := v.(T)
	return val, ok
}

func ValueWithDefault[T any](ctx context.Context, key string, defValue any) T {
	val, ok := Value[T](ctx, key)
	if ok {
		return val
	}

	v, ok := defValue.(T)
	if ok {
		return v
	}

	var zero T
	return zero
}

func StringValue(ctx context.Context, key string, value string) string {
	return ValueWithDefault[string](ctx, key, value)
}

func Int64Value(ctx context.Context, key string, value string) int64 {
	return ValueWithDefault[int64](ctx, key, value)
}

func Int32Value(ctx context.Context, key string, value string) int32 {
	return ValueWithDefault[int32](ctx, key, value)
}

func IntValue(ctx context.Context, key string, value string) int {
	return ValueWithDefault[int](ctx, key, value)
}

func Float64Value(ctx context.Context, key string, value string) float64 {
	return ValueWithDefault[float64](ctx, key, value)
}

func Float32Value(ctx context.Context, key string, value string) float32 {
	return ValueWithDefault[float32](ctx, key, value)
}

func BoolValue(ctx context.Context, key string, value string) bool {
	return ValueWithDefault[bool](ctx, key, value)
}

func WithTraceId(ctx context.Context, traceId string) context.Context {
	return WithValue(ctx, string(TraceIdKey), traceId)
}

func WithGenerateTraceId(ctx context.Context) context.Context {
	return WithTraceId(ctx, uuid.New().String())
}

func TraceId(ctx context.Context) string {
	return ValueWithDefault[string](ctx, string(TraceIdKey), "")
}

func WithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(Default(ctx))
}

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(Default(ctx), timeout)
}

func WithDeadline(ctx context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(Default(ctx), d)
}

func WithSignalNotify(ctx context.Context, signals ...os.Signal) context.Context {
	if len(signals) == 0 {
		signals = []os.Signal{
			syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP,
			syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
		}
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, signals...)

	ctx, cancel := WithCancel(ctx)

	go func() {
		select {
		case <-signalCh:
			cancel()

		case <-ctx.Done():
		}

		signal.Stop(signalCh)
		close(signalCh)
	}()

	return ctx
}
