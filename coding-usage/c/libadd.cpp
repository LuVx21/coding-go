#include <iostream>

#ifdef __cplusplus
extern "C" {
#endif

// 使用 extern "C" 禁用名称修饰，导出标准的 C 符号
int add(int a, int b) {
    // 内部可以使用 C++ 逻辑，但接口必须是 C 风格
    std::cout << "[C++] Called add function" << std::endl;
    return a + b;
}

#ifdef __cplusplus
}
#endif