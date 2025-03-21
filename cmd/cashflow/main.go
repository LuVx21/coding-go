package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/luvx21/coding-go/coding-common/ios"
	"github.com/luvx21/coding-go/infra/infra_sql"
	"strings"
)

var (
	file = flag.String("f", "", "文件路径")

	host     = flag.String("h", "127.0.0.1", "数据库地址")
	port     = flag.Int("P", 3306, "端口")
	username = flag.String("u", "root", "用户名")
	password = flag.String("p", "1121", "密码")
	database = flag.String("d", "boot", "数据库")
)

func main() {
	flag.Parse()
	fmt.Println(*file, *host, *port, *username, *password, *database)
	if len(*file) == 0 {
		return
	}

	lines, _ := ios.ReadLines(*file)
	lines = readLines1(lines)
	if len(lines) == 0 {
		return
	}
	connect, _ := infra_sql.ConnectMySQL(*host, *port, *username, *password, *database)
	for _, line := range lines {
		row := strings.Split(line, ",")
		anies := cast(row)
		_, err := connect.Exec("insert into cashflow values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", anies...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cast(row []string) []any {
	anies := make([]any, len(row))
	for i := range row {
		anies[i] = row[i]
	}
	return anies
}

func readLines(lines []string) []string {
	var start = false
	rr := make([]string, 0)
	for _, line := range lines {
		if strings.Contains(line, "---------------------------------") {
			start = !start
		}

		if start {
			rr = append(rr, line)
		}
	}
	return rr
}
func readLines1(lines []string) []string {
	return lines[5 : len(lines)-7]
}
