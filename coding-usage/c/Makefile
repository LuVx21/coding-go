build:
	mkdir -p lib
# linux
# 	gcc -shared -fPIC -o ./lib/libadd.so libadd.c
	gcc -shared -fPIC -o ./lib/libadd.dylib libadd.c

build-cpp:
	mkdir -p lib
	# Linux/macOS: 使用 g++ 或 clang++
	g++ -shared -fPIC -o ./lib/libadd-cpp.dylib libadd.cpp

test:
	CGO_ENABLED=0 go test -v -test.fullpath=true -run ^Test_main_00$$ main_test.go
