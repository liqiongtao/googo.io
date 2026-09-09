package goocontext

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"google.golang.org/grpc/metadata"
)

type Key string

const (
	TraceIdKey     Key = "trace-id"
	ServiceNameKey Key = "service-name"
)

func Default(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx
}

func WithValue(ctx context.Context, key string, value any) context.Context {
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

func MetaDataValue(ctx context.Context, key string) []string {
	if ctx == nil {
		return []string{}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return []string{}
	}

	v, ok := md[key]
	if !ok {
		return []string{}
	}

	return v
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

func ValueString(ctx context.Context, key string) string {
	return ValueWithDefault[string](ctx, key, "")
}

func ValueInt64(ctx context.Context, key string) int64 {
	return ValueWithDefault[int64](ctx, key, 0)
}

func ValueInt32(ctx context.Context, key string) int32 {
	return ValueWithDefault[int32](ctx, key, 0)
}

func ValueInt(ctx context.Context, key string) int {
	return ValueWithDefault[int](ctx, key, 0)
}

func ValueFloat64(ctx context.Context, key string) float64 {
	return ValueWithDefault[float64](ctx, key, 0)
}

func ValueFloat32(ctx context.Context, key string) float32 {
	return ValueWithDefault[float32](ctx, key, 0)
}

func ValueBool(ctx context.Context, key string) bool {
	return ValueWithDefault[bool](ctx, key, false)
}

func WithTraceId(ctx context.Context, traceId string) context.Context {
	return WithValue(ctx, string(TraceIdKey), traceId)
}

func WithGenerateTraceId(ctx context.Context) context.Context {
	return WithTraceId(ctx, uuid.New().String())
}

func TraceId(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v, ok := md[string(TraceIdKey)]; ok {
			return strings.Join(v, " ")
		}
	}
	return ValueString(ctx, string(TraceIdKey))
}

func WithServiceName(ctx context.Context, serviceName string) context.Context {
	return WithValue(ctx, string(ServiceNameKey), serviceName)
}

func ServiceName(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v, ok := md[string(ServiceNameKey)]; ok {
			return strings.Join(v, " ")
		}
	}
	return ValueString(ctx, string(ServiceNameKey))
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

func Log(ctx context.Context) *goo_log.Entry {
	log := goo_log.WithField("trace-id", TraceId(ctx))

	if v := ServiceName(ctx); v != "" {
		log = log.WithField("service-name", v)
	}
	if v := ValueString(ctx, "request-uri"); v != "" {
		log = log.WithField("request-uri", v)
	}

	return log
}
