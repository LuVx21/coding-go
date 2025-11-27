package os_x

import (
	"fmt"
	"strings"
	"testing"
)

func Test_env_00(t *testing.T) {
	r, b := Command("go", "env", "GOMOD")
	fmt.Println(b, r)

	r, b = Command("sh", "-c", "go list -m -f {{.Dir}}")
	fmt.Println(b, strings.TrimSpace(r))
}
