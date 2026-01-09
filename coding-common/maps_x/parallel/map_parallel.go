package parallel

import "sync"

// ForEachKeyValue 并发迭代
func ForEachKeyValue[M ~map[K]V, K comparable, V any](m M, iteratee func(K, V)) {
	ForEachKeyValueMax(m, 0, iteratee)
}

// ForEachKeyValueMax 并发迭代(有最大限制)
func ForEachKeyValueMax[M ~map[K]V, K comparable, V any](m M, max uint, iteratee func(K, V)) {
	isLimit := max != 0
	var ch chan struct{}
	if isLimit {
		ch = make(chan struct{}, max)
	}

	var wg sync.WaitGroup
	wg.Add(len(m))

	for k, v := range m {
		if isLimit {
			ch <- struct{}{}
		}

		go func(_k K, _v V) {
			defer func() {
				if isLimit {
					<-ch
				}
				wg.Done()
			}()

			iteratee(_k, _v)
		}(k, v)
	}

	wg.Wait()
}

// Transfer 对key,value分别处理, 生成一个新map
func Transfer[M ~map[K1]V1, K1, K2 comparable, V1, V2 any](m M,
	keyMapper func(K1) K2,
	valueMapper func(V1) V2) map[K2]V2 {
	r := make(map[K2]V2, len(m))

	var mu sync.RWMutex
	var wg sync.WaitGroup
	wg.Add(len(m))

	for k, v := range m {
		go func(_k K1, _v V1) {
			defer wg.Done()
			__k, __v := keyMapper(_k), valueMapper(_v)
			mu.Lock()
			defer mu.Unlock()
			r[__k] = __v
		}(k, v)
	}

	wg.Wait()
	return r
}
