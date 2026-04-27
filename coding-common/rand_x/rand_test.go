package rand_x

import (
	"fmt"
	"testing"
)

func Test_rand_00(t *testing.T) {
	for range 10 {
		fmt.Println(RandomInt(2, 8))
	}
}
