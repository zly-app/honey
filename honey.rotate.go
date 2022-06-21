package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zly-app/zapp/component/gpool"

	"github.com/zly-app/honey/log_data"
	"github.com/zly-app/honey/pkg/rotate"
)

type rotateWaitGroup struct {
	r  rotate.IRotator
	wg sync.WaitGroup
}

type rotateEnvGroup struct {
	creator func(env, app, instance string) rotate.IRotator // 旋转器建造者
	wgs     map[string]*rotateWaitGroup
	mx      sync.RWMutex
}

// 获取rotate
func (r *rotateEnvGroup) GetRotate(env, app, instance string) rotate.IRotator {
	name := fmt.Sprintf("%s|%s|%s", env, app, instance)

	r.mx.RLock()
	wg, ok := r.wgs[name]
	r.mx.RUnlock()

	if ok {
		wg.wg.Wait()
		return wg.r
	}

	r.mx.Lock()

	// 再获取一次, 它可能在获取锁的过程中完成了
	if wg, ok = r.wgs[name]; ok {
		r.mx.Unlock()
		wg.wg.Wait()
		return wg.r
	}

	// 占位置
	wg = new(rotateWaitGroup)
	wg.wg.Add(1)
	r.wgs[name] = wg
	r.mx.Unlock()

	// 创建
	wg.r = r.creator(env, app, instance)
	wg.wg.Done()
	return wg.r
}

// 获取所有rotate
func (r *rotateEnvGroup) GetAllRotate() []rotate.IRotator {
	r.mx.Lock()
	defer r.mx.Unlock()

	rotates := make([]rotate.IRotator, len(r.wgs))
	index := 0
	for _, wg := range r.wgs {
		wg.wg.Wait()
		rotates[index] = wg.r
		index++
	}
	return rotates
}

// 生成旋转组
func (h *Honey) MakeRotateGroup() {
	h.rotateGroup = &rotateEnvGroup{
		creator: h.rotateCreator,
		wgs:     make(map[string]*rotateWaitGroup),
	}
	h.rotateGPool = gpool.NewGPool(&gpool.GPoolConfig{
		ThreadCount: h.conf.MaxRotateThreadNum,
	})
}

// 旋转器建造者
func (h *Honey) rotateCreator(env, app, instance string) rotate.IRotator {
	opts := []rotate.Option{
		rotate.WithBatchSize(h.conf.LogBatchSize),
		rotate.WithAutoRotateTime(time.Duration(h.conf.AutoRotateTime) * time.Second),
	}
	callback := func(values []interface{}) {
		h.rotateGPool.Go(func() error {
			h.RotateCallback(env, app, instance, values)
			return nil
		}, nil)
	}
	return rotate.NewRotate(callback, opts...)
}

// 旋转器回调
func (h *Honey) RotateCallback(env, app, instance string, a []interface{}) {
	data := make([]*log_data.LogData, len(a))
	for i, v := range a {
		data[i] = v.(*log_data.LogData)
	}

	// 输出
	for _, out := range h.outputs {
		out.Out(env, app, instance, data)
	}
}
