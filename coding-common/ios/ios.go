package ios

import (
	"bufio"
	"io"
	"os"
)

// os.ReadFile: 小文件, 代码简洁

// ReadLines Scanner处理读取结果
// 大文件或逐行处理, 文件较大时，避免一次性加载全部内容到内存; 内存效率高
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ReadLines1 Reader处理读取结果
func ReadLines1(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	br := bufio.NewReader(file)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

// ReadAll 全部内容
func ReadAll(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), err
}
