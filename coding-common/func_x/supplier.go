package func_x

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// CloseableSupplier 可关闭的 Supplier，支持懒加载和资源释放
type CloseableSupplier[T any] struct {
	delegate        func() T     // 值提供函数
	resetAfterClose bool         // 关闭后是否重置
	initialized     int32        // 原子标记，替代 volatile boolean, 0:未 1:已
	value           *T           // 缓存的值
	mu              sync.RWMutex // 读写锁，保证线程安全
}

func Lazy[T any](delegate func() T) *CloseableSupplier[T] {
	return LazyWith(delegate, true)
}

// LazyWith 创建新的 CloseableSupplier
func LazyWith[T any](delegate func() T, resetAfterClose bool) *CloseableSupplier[T] {
	return &CloseableSupplier[T]{
		delegate:        delegate,
		resetAfterClose: resetAfterClose,
	}
}

// Get 获取值，如果未初始化则懒加载
func (cs *CloseableSupplier[T]) Get() T {
	if atomic.LoadInt32(&cs.initialized) == 1 {
		cs.mu.RLock()
		defer cs.mu.RUnlock()
		return *cs.value
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// 双重检查锁
	if atomic.LoadInt32(&cs.initialized) == 0 {
		val := cs.delegate()
		cs.value = &val
		atomic.StoreInt32(&cs.initialized, 1)
	}

	return *cs.value
}

// IsInitialized 检查是否已初始化
func (cs *CloseableSupplier[T]) IsInitialized() bool {
	return atomic.LoadInt32(&cs.initialized) == 1
}

// IfPresent 如果值存在则执行消费者函数
func (cs *CloseableSupplier[T]) IfPresent(consumer func(T) error) error {
	if e, b := cs.Map(func(t T) any { return consumer(t) }); b {
		return nil
	} else {
		return e.(error)
	}
}

// Map 将当前值映射为另一种类型
func (cs *CloseableSupplier[T]) Map(mapper func(T) any) (any, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if atomic.LoadInt32(&cs.initialized) == 1 && cs.value != nil {
		return mapper(*cs.value), true
	}
	return nil, false
}

// TryClose 尝试关闭，释放资源
func (cs *CloseableSupplier[T]) TryClose() error {
	return cs.TryCloseWith(func(T) error { return nil })
}

// TryCloseWith 使用自定义关闭函数尝试关闭
func (cs *CloseableSupplier[T]) TryCloseWith(closer func(T) error) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if atomic.LoadInt32(&cs.initialized) == 1 && cs.value != nil {
		if err := closer(*cs.value); err != nil {
			return err
		}

		if cs.resetAfterClose {
			cs.value = nil
			atomic.StoreInt32(&cs.initialized, 0)
		}
	}
	return nil
}

// String 字符串表示
func (cs *CloseableSupplier[T]) String() string {
	if atomic.LoadInt32(&cs.initialized) == 1 {
		cs.mu.RLock()
		defer cs.mu.RUnlock()
		if cs.value != nil {
			return fmt.Sprintf("CloseableSupplier(%v)", *cs.value)
		}
	}
	return "CloseableSupplier(<uninitialized>)"
}
