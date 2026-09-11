package goo_log

import (
	"sync"
	"testing"
)

type recordingAdapter struct {
	mu   sync.Mutex
	msgs []*Message
}

func (a *recordingAdapter) Write(msg *Message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.msgs = append(a.msgs, msg)
}

func fieldMap(msg *Message) map[string]any {
	m := map[string]any{}
	if msg == nil || msg.Entry == nil {
		return m
	}
	for _, f := range msg.Entry.Data {
		m[f.Field] = f.Value
	}
	return m
}

func TestEntryWithFieldCopyOnWrite(t *testing.T) {
	adapter := &recordingAdapter{}
	l := New(adapter)

	base := NewEntry(l).WithTag("svc").WithField("trace-id", "t1")
	child := base.WithField("msg", "hello")

	base.Info("base")
	child.Info("child")

	if len(adapter.msgs) != 2 {
		t.Fatalf("want 2 messages, got %d", len(adapter.msgs))
	}

	baseFields := fieldMap(adapter.msgs[0])
	if _, ok := baseFields["msg"]; ok {
		t.Fatal("base log should not contain temporary field msg")
	}
	if baseFields["trace-id"] != "t1" {
		t.Fatalf("base missing trace-id: %#v", baseFields)
	}
	if len(adapter.msgs[0].Entry.Tags) != 1 || adapter.msgs[0].Entry.Tags[0] != "svc" {
		t.Fatalf("base tags = %#v", adapter.msgs[0].Entry.Tags)
	}

	childFields := fieldMap(adapter.msgs[1])
	if childFields["msg"] != "hello" {
		t.Fatalf("child missing msg: %#v", childFields)
	}
	if childFields["trace-id"] != "t1" {
		t.Fatalf("child missing trace-id: %#v", childFields)
	}
}

func TestEntryReuseDoesNotAccumulateTemporaryFields(t *testing.T) {
	adapter := &recordingAdapter{}
	l := New(adapter)

	log := NewEntry(l).WithTag("kafka").WithField("trace-id", "t2")
	log.WithField("执行时间", "1.2").Error("fail")
	log.Debug("ok")

	if len(adapter.msgs) != 2 {
		t.Fatalf("want 2 messages, got %d", len(adapter.msgs))
	}

	errFields := fieldMap(adapter.msgs[0])
	if errFields["执行时间"] != "1.2" {
		t.Fatalf("error log missing 执行时间: %#v", errFields)
	}

	okFields := fieldMap(adapter.msgs[1])
	if _, ok := okFields["执行时间"]; ok {
		t.Fatal("reused base log should not keep temporary 执行时间")
	}
	if okFields["trace-id"] != "t2" {
		t.Fatalf("reused base missing trace-id: %#v", okFields)
	}
}

func TestEntryCloneDoesNotShareBackingArray(t *testing.T) {
	adapter := &recordingAdapter{}
	l := New(adapter)

	base := NewEntry(l).WithTag("a").WithField("k", 1)
	child := base.WithTag("b").WithField("k2", 2)

	// 再改 child 不应影响已产出的 base 快照能力
	_ = child.WithField("k3", 3)
	base.Info("base")
	child.Info("child")

	if len(adapter.msgs[0].Entry.Tags) != 1 || adapter.msgs[0].Entry.Tags[0] != "a" {
		t.Fatalf("base tags polluted: %#v", adapter.msgs[0].Entry.Tags)
	}
	if len(adapter.msgs[1].Entry.Tags) != 2 {
		t.Fatalf("child tags = %#v", adapter.msgs[1].Entry.Tags)
	}
	bf, cf := fieldMap(adapter.msgs[0]), fieldMap(adapter.msgs[1])
	if _, ok := bf["k2"]; ok {
		t.Fatalf("base has child field: %#v", bf)
	}
	if cf["k2"] != 2 || cf["k"] != 1 {
		t.Fatalf("child fields = %#v", cf)
	}
}
