package goo_log

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNormalizeFilePath(t *testing.T) {
	if got := normalizeFilePath(""); got != "logs/" {
		t.Fatalf("empty path: got %q", got)
	}
	if got := normalizeFilePath("  "); got != "logs/" {
		t.Fatalf("blank path: got %q", got)
	}
	if got := normalizeFilePath("tmp/logs"); got != "tmp/logs/" {
		t.Fatalf("no slash: got %q", got)
	}
	if got := normalizeFilePath("tmp/logs/"); got != "tmp/logs/" {
		t.Fatalf("with slash: got %q", got)
	}
}

func TestFileAdapterSyncClose(t *testing.T) {
	dir := t.TempDir()
	fa := NewFileAdapter(FilePathOption(dir), FileMaxSizeOption(1<<20))
	l := New(fa)

	l.Info("sync-close-test")
	if err := l.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}

	// Close 后再写应安全丢弃
	l.Info("after-close")

	matches, err := filepath.Glob(filepath.Join(dir, "*.log"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("expected log file")
	}
	b, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("expected log content")
	}
}

func TestRotateCountUsesMaxSuffix(t *testing.T) {
	dir := t.TempDir() + "/"
	ymd := "20260110"
	// 故意制造空洞：只有 _1 和 _90
	for _, n := range []int{1, 90} {
		name := fmt.Sprintf("%s%s_%d.log", dir, ymd, n)
		if err := os.WriteFile(name, []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	fa := &FileAdapter{filepath: dir}
	if got := fa.rotateCount(ymd); got != 90 {
		t.Fatalf("rotateCount=%d, want 90", got)
	}
}

func TestCutFileSequentialUnderConcurrency(t *testing.T) {
	dir := t.TempDir()
	l := NewFileLog(FilePathOption(dir), FileMaxSizeOption(300))

	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				l.Info(strings.Repeat("x", 40))
			}
		}()
	}
	wg.Wait()
	if err := l.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*_*.log"))
	if err != nil {
		t.Fatal(err)
	}
	max := 0
	for _, f := range matches {
		base := strings.TrimSuffix(filepath.Base(f), ".log")
		i := strings.LastIndex(base, "_")
		n, convErr := strconv.Atoi(base[i+1:])
		if convErr != nil {
			t.Fatal(convErr)
		}
		if n > max {
			max = n
		}
	}
	// 单 worker 下序号应连续：文件数 == 最大序号
	if max > 0 && len(matches) != max {
		t.Fatalf("rotated files not sequential: count=%d maxSuffix=%d files=%v", len(matches), max, matches)
	}
}

func TestFileDropWhenFull(t *testing.T) {
	dir := t.TempDir()
	fa := NewFileAdapter(
		FilePathOption(dir),
		FileQueueSizeOption(1),
		FileDropWhenFullOption(true),
	)
	l := New(fa)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			l.Info("drop-when-full", i)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Write blocked under DropWhenFull")
	}

	_ = l.Sync()
	_ = l.Close()
}

func TestEntryTraceNotLeakedOnReuse(t *testing.T) {
	var got [][]string
	l := New(&captureAdapter{onWrite: func(msg *Message) {
		var tr []string
		if msg.Entry != nil {
			tr = append(tr, msg.Entry.Trace...)
		}
		got = append(got, tr)
	}})

	e := NewEntry(l)
	e.Trace = []string{"fake 1L"} // 模拟上次 WARN+ 残留
	e.Info("with-stale-trace")
	e.Info("after-clear")

	if len(got) != 2 {
		t.Fatalf("writes=%d, want 2", len(got))
	}
	if len(got[0]) != 1 || got[0][0] != "fake 1L" {
		t.Fatalf("first write trace=%v", got[0])
	}
	if len(got[1]) != 0 {
		t.Fatalf("second write leaked trace: %v", got[1])
	}
}

type captureAdapter struct {
	onWrite func(msg *Message)
}

func (c *captureAdapter) Write(msg *Message) {
	if c.onWrite != nil {
		c.onWrite(msg)
	}
}
