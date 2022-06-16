package goo

import (
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"time"
)

/**
 * 每天00时执行一次
 */
func CronDay(fns ...func()) {
	goo_utils.AsyncFunc(func() {
		nw := time.Now()

		d := time.Date(nw.Year(), nw.Month(), nw.Day(),
			0, 0, 0, 0, time.Local).
			Add(24 * time.Hour).Sub(nw)

		time.Sleep(d)

		goo_utils.AsyncFuncGroup(fns...)

		Cron(24*time.Hour, fns...)
	})
}

/**
 * 每小时00分执行一次
 */
func CronHour(fns ...func()) {
	goo_utils.AsyncFunc(func() {
		nw := time.Now()

		d := time.Date(nw.Year(), nw.Month(), nw.Day(),
			nw.Hour(), 0, 0, 0, time.Local).
			Add(time.Hour).Sub(nw)

		time.Sleep(d)

		goo_utils.AsyncFuncGroup(fns...)

		Cron(time.Hour, fns...)
	})
}

/**
 * 每分钟00秒执行一次
 */
func CronMinute(fns ...func()) {
	goo_utils.AsyncFunc(func() {
		nw := time.Now()

		d := time.Date(nw.Year(), nw.Month(), nw.Day(),
			nw.Hour(), nw.Minute(), 0, 0, time.Local).
			Add(time.Minute).Sub(nw)

		time.Sleep(d)

		goo_utils.AsyncFuncGroup(fns...)

		Cron(time.Minute, fns...)
	})
}

/**
 * 没多长时间执行任务
 */
func Cron(d time.Duration, fns ...func()) {
	goo_utils.AsyncFunc(func() {
		ticker := time.NewTicker(d)

		for {
			select {
			case <-goo_context.Cancel().Done():
				ticker.Stop()
				return

			case <-ticker.C:
				time.Sleep(time.Second)
				goo_utils.AsyncFuncGroup(fns...)
			}
		}
	})
}
