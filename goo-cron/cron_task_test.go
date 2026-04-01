package goo_cron

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	key := "sc:cron:task"

	//r := redis.NewClient(&redis.Options{
	//	Addr:     "redis.in:20063",
	//	Password: "7fdacd2183ab",
	//	DB:       0,
	//})

	tasks := map[string]TaskFunc{
		"a1": func(task *TaskData) {
			fmt.Println(task.Code, time.Now().Format("15:04:05"), task.Data)
		},
		"a2": func(task *TaskData) {
			fmt.Println(task.Code, time.Now().Format("15:04:05"), task.Data)
		},
	}

	//c := New(key, tasks, WithRedis(r))
	c := New(key, tasks)

	c.Add(&TaskData{
		Code:   "a1",
		Spec:   "*/2 * * * * *",
		Status: 1,
	})
	c.Add(&TaskData{
		Code:   "a2",
		Spec:   "*/3 * * * * *",
		Status: 1,
	})

	c.Run()
}
