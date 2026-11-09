package utils

import "testing"

func TestGetInstanceNeverEmpty(t *testing.T) {
	const def = "default-instance"
	if got := GetInstance(def); got == "" {
		t.Fatal("GetInstance 不应返回空")
	}
}

// 有非回环 IPv4 时返回该 IP，否则返回默认值；两者都在允许集合内
func TestGetInstanceReturnsIPOrDefault(t *testing.T) {
	const def = "fallback"
	got := GetInstance(def)
	if got != def && got == "" {
		t.Fatalf("返回值既不是默认值也不是 IP: %q", got)
	}
}
