package goo

import (
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

/**
 * 每天执行一次
 */
func CronDay(fns ...func()) {
	goo_utils.AsyncFunc(func() {
		timer := time.NewTimer(0)
		for {
			select {
			case <-goo_context.Cancel().Done():
				return
			case <-timer.C:
				for _, fn := range fns {
					goo_utils.AsyncFunc(fn)
				}
				ti, _ := time.ParseInLocation("2006-01-02", time.Now().Add(24*time.Hour).Format("2006-01-02"), time.Local)
				timer.Reset(time.Duration(ti.Unix()-time.Now().Unix()) * time.Second)
			}
		}
	})
}

/**
 * 每小时执行一次
 */
func CronHour(fns ...func()) {
	goo_utils.AsyncFunc(func() {
		timer := time.NewTimer(0)
		for {
			select {
			case <-goo_context.Cancel().Done():
				return
			case <-timer.C:
				for _, fn := range fns {
					goo_utils.AsyncFunc(fn)
				}
				ti, _ := time.ParseInLocation("2006-01-02 15", time.Now().Add(1*time.Hour).Format("2006-01-02 15"), time.Local)
				timer.Reset(time.Duration(ti.Unix()-time.Now().Unix()) * time.Second)
			}
		}
	})
}

/**
 * 每分钟执行一次
 */
func CronMinute(fns ...func()) {
	goo_utils.AsyncFunc(func() {
		timer := time.NewTimer(0)
		for {
			select {
			case <-goo_context.Cancel().Done():
				return
			case <-timer.C:
				for _, fn := range fns {
					goo_utils.AsyncFunc(fn)
				}
				ti, _ := time.ParseInLocation("2006-01-02 15:04", time.Now().Add(60*time.Second).Format("2006-01-02 15:04"), time.Local)
				timer.Reset(time.Duration(ti.Unix()-time.Now().Unix()) * time.Second)
			}
		}
	})
}

/**
 * 定期执行任务
 */
func Crond(d time.Duration, fns ...func()) {
	goo_utils.AsyncFunc(func() {
		timer := time.NewTimer(0)
		for {
			select {
			case <-goo_context.Cancel().Done():
				return
			case <-timer.C:
				for _, fn := range fns {
					goo_utils.AsyncFunc(fn)
				}
				timer.Reset(d)
			}
		}
	})
}
