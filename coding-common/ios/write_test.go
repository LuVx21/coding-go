package ios

import (
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x"
)

var home, _ = common_x.Dir()

func Test_WriteFile(t *testing.T) {
	data := make([]byte, 1<<30)
	_ = WriteFile(home+"/1.bin", data, &WriteOptions{
		BufferSize: 1 << 20,
	})
}

func Test_WriteStream(t *testing.T) {
}

func Test_WriteChunks(t *testing.T) {
	chunks := make(chan []byte, 10)

	go func() {
		for i := range 100 {
			chunks <- fmt.Appendf(nil, "chunk-%d\n", i)
		}
		close(chunks)
	}()

	_ = WriteChunks(home+"/3.log", chunks, &WriteOptions{
		Concurrency: 8,
	})

	time.Sleep(time.Second * 20)
}
