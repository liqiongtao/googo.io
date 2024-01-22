package goo_cron

import (
	"fmt"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/robfig/cron/v3"
	"time"
)

type crontab struct {
	c *cron.Cron
}

func New(opts ...cron.Option) *crontab {
	return &crontab{c: cron.New(opts...)}
}

func Default() *crontab {
	return New(cron.WithSeconds())
}

func (c *crontab) Run() {
	c.c.Start()
	c.Stop()
}

func (c *crontab) Start() {
	c.c.Start()
}

func (c *crontab) Stop() {
	<-goo_context.Cancel().Done()
	goo_log.WithTag("goo-cron").Debug("系统退出，等待全部任务执行结束...")

	<-c.c.Stop().Done()
	goo_log.WithTag("goo-cron").Debug("系统退出成功，全部任务执行结束")

	time.Sleep(time.Second)
}

func (c *crontab) AddFunc(spec string, fn ...func()) *crontab {
	for _, f := range fn {
		c.c.AddFunc(spec, f)
	}
	return c
}

func (c *crontab) AddJob(spec string, job ...cron.Job) *crontab {
	for _, j := range job {
		c.c.AddJob(spec, j)
	}
	return c
}

// 每天0点0分0秒执行
func (c *crontab) Day(fn ...func()) *crontab {
	return c.AddFunc("0 0 0 * * *", fn...)
}

// 每天x点0分0秒执行
func (c *crontab) DayHour(hour int, fn ...func()) *crontab {
	return c.AddFunc(fmt.Sprintf("0 0 %d * * *", hour), fn...)
}

// 每天x点x分0秒执行
func (c *crontab) DayHourMinute(hour, minute int, fn ...func()) *crontab {
	return c.AddFunc(fmt.Sprintf("0 %d %d * * *", minute, hour), fn...)
}

// 每小时执行
func (c *crontab) Hour(fn ...func()) *crontab {
	return c.AddFunc("0 0 */1 * * *", fn...)
}

// 每隔x小时执行
func (c *crontab) HourX(x int, fn ...func()) *crontab {
	return c.AddFunc(fmt.Sprintf("0 0 */%d * * *", x), fn...)
}

// 每分钟执行
func (c *crontab) Minute(fn ...func()) *crontab {
	return c.AddFunc("0 */1 * * * *", fn...)
}

// 每隔x分钟执行
func (c *crontab) MinuteX(x int, fn ...func()) *crontab {
	return c.AddFunc(fmt.Sprintf("0 */%d * * * *", x), fn...)
}

// 每秒钟执行
func (c *crontab) Second(fn ...func()) *crontab {
	return c.AddFunc("* * * * * *", fn...)
}

// 每隔x秒执行
func (c *crontab) SecondX(x int, fn ...func()) *crontab {
	return c.AddFunc(fmt.Sprintf("*/%d * * * * *", x), fn...)
}
