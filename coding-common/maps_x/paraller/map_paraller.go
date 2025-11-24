package paraller

import "sync"

// ForEachKeyValue 并放迭代
func ForEachKeyValue[M ~map[K]V, K comparable, V any](m M, iteratee func(K, V)) {
	ForEachKeyValueMax(m, 0, iteratee)
}

// ForEachKeyValueMax 并放迭代(有最大限制)
func ForEachKeyValueMax[M ~map[K]V, K comparable, V any](m M, max int, iteratee func(K, V)) {
	var ch chan struct{}
	if max != 0 {
		ch = make(chan struct{}, max)
	}

	var wg sync.WaitGroup
	wg.Add(len(m))

	for k, v := range m {
		if max != 0 {
			ch <- struct{}{}
		}

		go func(_k K, _v V) {
			defer func() {
				if max != 0 {
					<-ch
				}
				wg.Done()
			}()

			iteratee(_k, _v)
		}(k, v)
	}

	wg.Wait()
}
