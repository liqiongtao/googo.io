package goo_cron

import (
	"fmt"
	"testing"
	"time"
)

func TestCron_Second(t *testing.T) {
	c := Default()
	c.SecondX(10, func() {
		fmt.Println("--11---", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(5 * time.Second)
		fmt.Println("--12---", time.Now().Format("2006-01-02 15:04:05"))
	})
	c.Run()
}
