package goo_task_queue

import (
	"errors"
	"fmt"
	"time"
)

// retryAfterError 业务通过 RetryAfter 返回，框架据此延迟再次入队。
type retryAfterError struct {
	after time.Duration
	err   error
}

// RetryAfter 在业务失败时指定多久后再可被调度。
// 未包装时，框架默认以当前时间（立即）重新入队。
func RetryAfter(d time.Duration, err error) error {
	if err == nil {
		err = fmt.Errorf("retry after %s", d)
	}
	return &retryAfterError{after: d, err: err}
}

func (e *retryAfterError) Error() string {
	return e.err.Error()
}

func (e *retryAfterError) Unwrap() error {
	return e.err
}

func retryAfterDuration(err error) (time.Duration, bool) {
	var e *retryAfterError
	if errors.As(err, &e) {
		return e.after, true
	}
	return 0, false
}
