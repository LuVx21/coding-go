package ios

import (
    "bufio"
    "io"
    "os"
)

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
