package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/elliotchance/sshtunnel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/ssh"
)

func Test_sshTunnel_00(t *testing.T) {
	tunnel, err := sshtunnel.NewSSHTunnel(
		"root@mini.local:50022",

		// sshtunnel.PrivateKeyFile("/Users/renxie/.ssh/id_rsa"),
		ssh.Password("xxxx"),

		"postgresql-master:5432",
		"44444",
	)
	go func() {
		if err := tunnel.Start(); err != nil {
			slog.Error("tunnel建立失败", "error", err)
		}
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("隧道信息", tunnel.Local)

	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Asia/Shanghai", "postgres", "xxxx", tunnel.Local.Host, tunnel.Local.Port, "boot"))
	if err != nil {
		slog.Error("连接数据库失败", "error", err)
	}
	var s string
	err = db.QueryRow("select version();").Scan(&s)
	fmt.Println("查询错误", err, s)
}
