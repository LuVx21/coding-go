package bili

import (
	"fmt"
	"testing"
)

func Test_00(t *testing.T) {
	fmt.Println(append(getFollows(43510), getFollows(-10)...))
	fmt.Println(getCollections())
	timeFlow()
}
