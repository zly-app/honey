package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zly-app/honey/pkg/rotate"
)

type rotateWaitGroup struct {
	r  rotate.Rotator
	wg sync.WaitGroup
}

type rotateEnvGroup struct {
	creator func(env string) rotate.Rotator // 旋转器建造者
	wgs     map[string]*rotateWaitGroup
	mx      sync.RWMutex
}

// 获取rotate
func (r *rotateEnvGroup) GetRotate(env string) rotate.Rotator {
	r.mx.RLock()
	wg, ok := r.wgs[env]
	r.mx.RUnlock()

	if ok {
		wg.wg.Wait()
		return wg.r
	}

	r.mx.Lock()

	// 再获取一次, 它可能在获取锁的过程中完成了
	if wg, ok = r.wgs[env]; ok {
		r.mx.Unlock()
		wg.wg.Wait()
		return wg.r
	}

	// 占位置
	wg = new(rotateWaitGroup)
	wg.wg.Add(1)
	r.wgs[env] = wg
	r.mx.Unlock()

	// 创建
	wg.r = r.creator(env)
	wg.wg.Done()
	return wg.r
}

// 获取所有rotate
func (r *rotateEnvGroup) GetAllRotate() []rotate.Rotator {
	r.mx.Lock()
	defer r.mx.Unlock()

	rotates := make([]rotate.Rotator, len(r.wgs))
	index := 0
	for _, wg := range r.wgs {
		wg.wg.Wait()
		rotates[index] = wg.r
		index++
	}
	return rotates
}

func newRotateGroup(creator func(env string) rotate.Rotator) *rotateEnvGroup {
	return &rotateEnvGroup{
		creator: creator,
		wgs:     make(map[string]*rotateWaitGroup),
	}
}

// 旋转器建造者
func (h *Honey) rotateCreator(env string) rotate.Rotator {
	opts := []rotate.Option{
		rotate.WithBatchSize(h.conf.BatchSize),
		rotate.WithAutoRotateTime(time.Duration(h.conf.AutoRotateTime) * time.Second),
	}
	callback := func(values []interface{}) {
		h.RotateCallback(env, values)
	}
	return rotate.NewRotate(callback, opts...)
}

// 旋转器回调
func (h *Honey) RotateCallback(env string, a []interface{}) {
	// todo 待实现
	for _, v := range a {
		fmt.Println(env, v)
	}
}
