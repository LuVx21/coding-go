package syncx

import "testing"

func Test_sync_atomic_00(t *testing.T) {
    v := NewAtomicValue(123)
    println(v.Load())
}
