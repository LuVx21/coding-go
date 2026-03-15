package main

// #include <stdio.h>
// #include <stdlib.h>
//
// // 声明 C 函数
// int add(int a, int b) {
//     return a + b;
// }
//
// void print_hello() {
//     printf("Hello from C!\n");
// }
import "C" // 必须紧跟在 C 代码注释后，不能有空行
import (
	"fmt"
)

func main() {
	// 调用 C 函数
	result := C.add(3, 4)
	fmt.Printf("3 + 4 = %d\n", result)

	C.print_hello()
}
