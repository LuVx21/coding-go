package main

import (
	"bufio"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/glamour"
	"github.com/go-sql-driver/mysql"
	"github.com/luvx21/coding-go/cmds/ai/utils"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/dbs"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/strings_x"
	"github.com/luvx21/coding-go/infra/ai"
	"github.com/luvx21/coding-go/infra/logs/slogs"
)

var ()

const (
	sql_ddl = `
	CREATE TABLE IF NOT EXISTS ai_record
	(
		id       BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
		session  VARCHAR(255)        NOT NULL COMMENT '',
		question VARCHAR(2046)       NOT NULL COMMENT '',
		answer   TEXT                NOT NULL COMMENT '',
		model    VARCHAR(255)        NOT NULL COMMENT '',
		time     VARCHAR(255)        NOT NULL COMMENT '',
		PRIMARY KEY (id)
	) ENGINE = InnoDB
	DEFAULT CHARSET = utf8mb4 comment 'ai chat'
	`
)

var (
	curSession string = "chat"
	curModel   *ai.Model
	file       *os.File
	db         *sql.DB
	stream     = cast_x.ToBool(strings_x.FirstNonEmpty(os.Getenv("STREAM"), "true"))
	render, _  = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	logger = slogs.GetLogger()
)

func newFile(session string) *os.File {
	home, _ := os.UserHomeDir()
	var err error
	file, err = os.OpenFile(home+"/data/ai/"+session+".md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	return file
}

func main() {
	newFile(curSession)
	defer file.Close()

	for {
		curModel = utils.SelectModel(curModel)
		// fmt.Print("\n\n\033[32m请输入您的问题(q:退出,models:切换模型,session:切换保存文件)" + curModel.Id + ": \033[0m")
		// question, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		question := prompt.Input("请输入您的问题(q:退出,models:切换模型,session:切换保存文件)"+curModel.Id+": ", func(d prompt.Document) []prompt.Suggest { return nil })
		question = strings.TrimSpace(question)

		if question == "q" {
			break
		} else if question == ".help" {
			fmt.Println("=============================")
			fmt.Print("q: 退出\nmodels: 切换模型\nsession xxx: 切换保存文件\n")
			continue
		} else if question == "" {
			continue
		} else if question == "models" {
			curModel = nil
			continue
		} else if strings.HasPrefix(question, "session") {
			file.Close()
			curSession = strings.TrimSpace(strings.TrimPrefix(question, "session"))
			newFile(curSession)
			continue
		}

		now := time.Now().Format("2006-01-02 15:04:05")
		file.WriteString(fmt.Sprintf("\n---\n\n### %s(%s) %s\n", question, curModel.Id, now))

		res, err := curModel.Request(stream, question)
		if err != nil || res.StatusCode >= 300 {
			fmt.Printf("请求接口失败,模型:%s, 服务商:%s, 响应:%v, err: %v\n", curModel.Id, curModel.Sp.BaseUrl, res, err)
			continue
		}
		defer res.Body.Close()

		var answer strings.Builder
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			line := scanner.Text()
			content := ai.ParseLineContent(line)
			if content != "" {
				answer.WriteString(content)
				file.WriteString(content)
				if rand.Intn(3)%3 <= 1 {
					fmt.Print("\033[H\033[2J") // 清屏
					content, _ = render.Render(answer.String())
					fmt.Print(content)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			slog.Warn("读取失败:" + err.Error())
			continue
		}
		runs.Go(func() {
			saveToDb(curSession, question, answer.String(), curModel.Id, now)
		})
	}
}

func saveToDb(session, question, answer, model, now string) {
	if db == nil || db.Ping() != nil {
		tidbHost, tidbPort := os_x.Getenv("TIDB_HOST"), os_x.Getenv("TIDB_PORT")
		tidbUsername, tidbPassword := os_x.Getenv("TIDB_USERNAME"), os_x.Getenv("TIDB_PASSWORD")

		mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: tidbHost,
		})

		dsn := dbs.MySQLConnect(tidbHost, cast_x.ToInt(tidbPort), tidbUsername, tidbPassword, "boot", map[string]string{"tls": "tidb"})
		var err error
		db, err = sql.Open(dbs.DriverMysql, dsn)
		if err != nil {
			slog.Error("TIDB连接错误", "err", err)
			return
		}
	}
	_, err := db.Exec("insert into ai_record(session, question, answer, model, time) values(?,?,?,?,?)", session, question, answer, model, now)
	if err != nil {
		logger.Warn("写数据异常", "err", err.Error())
		if strings.Contains(err.Error(), "no such table") || strings.Contains(err.Error(), "doesn't exist") {
			db.Exec(sql_ddl)
		}
	}
}
