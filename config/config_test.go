package config

import "testing"

func TestNewConfigDefaults(t *testing.T) {
	conf := NewConfig()
	if conf.Inputs != DefaultInputs || conf.Outputs != DefaultOutputs {
		t.Fatalf("输入输出默认值不符: %q/%q", conf.Inputs, conf.Outputs)
	}
	if !conf.ThisLog.StopLogOutput {
		t.Fatal("StopLogOutput 默认应为 true")
	}
}

func TestCheckFillsDefaults(t *testing.T) {
	conf := NewConfig()
	if err := conf.Check("my-app"); err != nil {
		t.Fatalf("Check 失败: %v", err)
	}
	if conf.ThisLog.Env != DefaultThisLogEnv {
		t.Fatalf("Env 未填默认值: %q", conf.ThisLog.Env)
	}
	if conf.ThisLog.App != "my-app" {
		t.Fatalf("App 应使用注入的 appName: %q", conf.ThisLog.App)
	}
	if conf.ThisLog.Instance == "" {
		t.Fatal("Instance 未填默认值")
	}
	if conf.LogBatchSize != DefaultLogBatchSize || conf.AutoRotateTime != DefaultAutoRotateTime ||
		conf.MaxRotateThreadNum != DefaultMaxRotateThreadNum {
		t.Fatalf("数值默认值未补齐: %+v", conf)
	}
}

func TestCheckKeepsUserValues(t *testing.T) {
	conf := NewConfig()
	conf.ThisLog.Env = "prod"
	conf.ThisLog.App = "user-app"
	conf.ThisLog.Instance = "user-instance"
	conf.LogBatchSize = 1
	conf.AutoRotateTime = 2
	conf.MaxRotateThreadNum = 3
	if err := conf.Check("injected"); err != nil {
		t.Fatalf("Check 失败: %v", err)
	}
	if conf.ThisLog.Env != "prod" || conf.ThisLog.App != "user-app" || conf.ThisLog.Instance != "user-instance" {
		t.Fatalf("用户配置被覆盖: %+v", conf.ThisLog)
	}
	if conf.LogBatchSize != 1 || conf.AutoRotateTime != 2 || conf.MaxRotateThreadNum != 3 {
		t.Fatalf("用户数值被覆盖: %+v", conf)
	}
}

// 同一进程内多份配置互不影响（曾经是包级 Conf 全局）
func TestConfigInstancesAreIsolated(t *testing.T) {
	a := NewConfig()
	b := NewConfig()
	a.ThisLog.App = "a"
	b.ThisLog.App = "b"
	if err := a.Check("x"); err != nil {
		t.Fatal(err)
	}
	if err := b.Check("x"); err != nil {
		t.Fatal(err)
	}
	if a.ThisLog.App != "a" || b.ThisLog.App != "b" {
		t.Fatalf("配置实例互相污染: %q/%q", a.ThisLog.App, b.ThisLog.App)
	}
}
