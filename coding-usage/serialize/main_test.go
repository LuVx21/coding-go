package serialize

import (
	"fmt"
	"testing"

	furygo "github.com/apache/fury/go/fury"
)

func Test_fury_00(t *testing.T) {
	type SomeClass struct {
		F2 map[string]string
		F3 map[string]string
	}
	fury := furygo.NewFury(true)
	if err := fury.RegisterTagType("example.SomeClass", SomeClass{}); err != nil {
		panic(err)
	}

	value := SomeClass{F2: map[string]string{"k1": "v1", "k2": "v2"}}
	value.F3 = value.F2
	bytes, _ := fury.Marshal(value)

	var newValue SomeClass
	// bytes can be data serialized by other languages.
	if err := fury.Unmarshal(bytes, &newValue); err != nil {
		panic(err)
	}
	fmt.Println(newValue)
}
