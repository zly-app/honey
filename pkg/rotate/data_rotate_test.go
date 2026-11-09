package rotate

import (
	"testing"
	"time"
)

func waitBatch(t *testing.T, ch chan []interface{}, wantLen int) []interface{} {
	t.Helper()
	select {
	case got := <-ch:
		if len(got) != wantLen {
			t.Fatalf("回调条数 = %d, 期望 %d", len(got), wantLen)
		}
		return got
	case <-time.After(2 * time.Second):
		t.Fatalf("2s 内未收到回调（期望 %d 条）", wantLen)
		return nil
	}
}

func TestRotateByBatchSize(t *testing.T) {
	ch := make(chan []interface{}, 1)
	r := NewRotate(func(values []interface{}) { ch <- values }, WithBatchSize(2))

	r.Add("a")
	r.Add("b")

	got := waitBatch(t, ch, 2)
	if got[0] != "a" || got[1] != "b" {
		t.Fatalf("批次内容不符: %v", got)
	}
}

func TestRotateImmediately(t *testing.T) {
	ch := make(chan []interface{}, 1)
	r := NewRotate(func(values []interface{}) { ch <- values }, WithBatchSize(10))

	r.Add("x")
	r.Rotate()

	got := waitBatch(t, ch, 1)
	if got[0] != "x" {
		t.Fatalf("批次内容不符: %v", got)
	}
}

func TestRotateEmptyNoCallback(t *testing.T) {
	called := make(chan struct{}, 1)
	r := NewRotate(func(values []interface{}) { called <- struct{}{} }, WithBatchSize(10))

	r.Rotate()

	select {
	case <-called:
		t.Fatal("空批次不应触发回调")
	case <-time.After(300 * time.Millisecond):
	}
}

func TestApplyEnforcesMinimums(t *testing.T) {
	r := &rotate{}
	r.Apply([]Option{WithBatchSize(0), WithAutoRotateTime(time.Millisecond)})
	if r.batchSize != minBatchSize {
		t.Fatalf("batchSize 下限未生效: %d", r.batchSize)
	}
	if r.autoRotateTime != minAutoRotateTime {
		t.Fatalf("autoRotateTime 下限未生效: %v", r.autoRotateTime)
	}
}

func TestApplyAutoRotateTime(t *testing.T) {
	ch := make(chan []interface{}, 1)
	r := NewRotate(func(values []interface{}) { ch <- values },
		WithBatchSize(100),
		WithAutoRotateTime(minAutoRotateTime),
	)

	r.Add("tick")
	got := waitBatch(t, ch, 1)
	if got[0] != "tick" {
		t.Fatalf("自动旋转批次内容不符: %v", got)
	}
}
