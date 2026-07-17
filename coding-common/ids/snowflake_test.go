package ids

import (
	"fmt"
	"testing"
)

func Test_01(t *testing.T) {
	node, _ := NewSnowflakeIdWorker(0, 0)
	for range 2 {
		fmt.Println(node.NextId())
	}
}
