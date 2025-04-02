package tools

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Test_validator_00(t *testing.T) {
	validate = validator.New()

	validateMap()
}
func validateMap() {
	user := map[string]any{"name": "Arshiya Kiani", "email": "zytel3301@gmail.com"}
	rules := map[string]any{"name": "required,min=8,max=32", "email": "omitempty,required,email"}

	errs := validate.ValidateMap(user, rules)

	if len(errs) > 0 {
		fmt.Println(errs)
	}
}
