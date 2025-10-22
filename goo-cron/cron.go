package goo_cron

import (
	"context"
	"fmt"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/robfig/cron/v3"
	"time"
)

type Crontab struct {
	c *cron.Cron
}

func New(opts ...cron.Option) *Crontab {
	return &Crontab{c: cron.New(opts...)}
}

func Default() *Crontab {
	return New(cron.WithSeconds())
}

func (c *Crontab) Run() {
	c.c.Start()

	<-goo_context.WithCancel().Done()
	goo_log.WithTag("goo-cron").Debug("系统退出，等待全部任务执行结束...")

	<-c.c.Stop().Done()
	goo_log.WithTag("goo-cron").Debug("系统退出成功，全部任务执行结束")

	time.Sleep(time.Second)
}

func (c *Crontab) Start() {
	c.c.Start()
}

func (c *Crontab) Stop() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-goo_context.WithCancel().Done()
		goo_log.WithTag("goo-cron").Debug("系统退出，等待全部任务执行结束...")

		<-c.c.Stop().Done()
		goo_log.WithTag("goo-cron").Debug("系统退出成功，全部任务执行结束")

		time.Sleep(time.Second)

		cancel()
	}()
	return ctx
}

func (c *Crontab) AddFunc(spec string, fn ...func()) *Crontab {
	for _, f := range fn {
		c.c.AddFunc(spec, f)
	}
	return c
}

func (c *Crontab) AddJob(spec string, job ...cron.Job) *Crontab {
	for _, j := range job {
		c.c.AddJob(spec, j)
	}
	return c
}

// 每天0点0分0秒执行
func (c *Crontab) Day(fn ...func()) *Crontab {
	return c.AddFunc("0 0 0 * * *", fn...)
}

// 每天x点0分0秒执行
func (c *Crontab) DayHour(hour int, fn ...func()) *Crontab {
	return c.AddFunc(fmt.Sprintf("0 0 %d * * *", hour), fn...)
}

// 每天x点x分0秒执行
func (c *Crontab) DayHourMinute(hour, minute int, fn ...func()) *Crontab {
	return c.AddFunc(fmt.Sprintf("0 %d %d * * *", minute, hour), fn...)
}

// 每小时执行
func (c *Crontab) Hour(fn ...func()) *Crontab {
	return c.AddFunc("0 0 */1 * * *", fn...)
}

// 每隔x小时执行
func (c *Crontab) HourX(x int, fn ...func()) *Crontab {
	return c.AddFunc(fmt.Sprintf("0 0 */%d * * *", x), fn...)
}

// 每分钟执行
func (c *Crontab) Minute(fn ...func()) *Crontab {
	return c.AddFunc("0 */1 * * * *", fn...)
}

// 每隔x分钟执行
func (c *Crontab) MinuteX(x int, fn ...func()) *Crontab {
	return c.AddFunc(fmt.Sprintf("0 */%d * * * *", x), fn...)
}

// 每秒钟执行
func (c *Crontab) Second(fn ...func()) *Crontab {
	return c.AddFunc("* * * * * *", fn...)
}

// 每隔x秒执行
func (c *Crontab) SecondX(x int, fn ...func()) *Crontab {
	return c.AddFunc(fmt.Sprintf("*/%d * * * * *", x), fn...)
}
