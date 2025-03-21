package reflects

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_00(t *testing.T) {

}

func add(a, b int) int {
	return a + b
}
func sum(numbers ...int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func Test_000(t *testing.T) {
	funcAdd := reflect.ValueOf(add)

	fnType := funcAdd.Type()
	numIn, numOut := fnType.NumIn(), fnType.NumOut()
	fmt.Println("入参数量:", numIn, "出参数量:", numOut)
	for i := range numIn {
		paramType := fnType.In(i)
		fmt.Printf("函数的第 %d 个参数类型: %v\n", i+1, paramType)
	}
	for i := range numOut {
		paramType := fnType.Out(i)
		fmt.Printf("函数的第 %d 个参数类型: %v\n", i+1, paramType)
	}

	args := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(5)}
	result := funcAdd.Call(args)
	fmt.Println("result:", result[0].Int())

	// ---------------------------------------------------

	funcSum := reflect.ValueOf(sum)
	fnType = funcSum.Type()
	numIn, numOut = fnType.NumIn(), fnType.NumOut()
	fmt.Println("入参数量:", numIn, "出参数量:", numOut, fnType.In(0))
	args = []reflect.Value{reflect.ValueOf(100), reflect.ValueOf(101)}
	result = funcSum.Call(args)
	fmt.Println("Sum result:", result[0].Int())
}

func Test_Call(t *testing.T) {
	CallFunc(add, 1, 2, 3, 4, 5, 6)
}

type Calculator struct{}

func (c Calculator) Add(a, b int) int {
	return a + b
}

func (c Calculator) Sum(numbers ...int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func Test_struct(t *testing.T) {
	calc := Calculator{}

	val := reflect.ValueOf(calc)
	methodAdd := val.MethodByName("Add")
	args := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(5)}
	resultAdd := methodAdd.Call(args)
	fmt.Println("Add result:", resultAdd[0].Int())

	// 不可导出的方法, 这里会提示invalid
	methodSum := val.MethodByName("Sum")
	args = []reflect.Value{reflect.ValueOf(100), reflect.ValueOf(101)}
	resultSum := methodSum.Call(args)
	fmt.Println("Sum result:", resultSum[0].Int())
}
