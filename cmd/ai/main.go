package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/luvx21/coding-go/cmds/ai/utils"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/strings_x"
	"github.com/luvx21/coding-go/infra/ai"
	"github.com/luvx21/coding-go/infra/logs/slogs"
	"github.com/tursodatabase/go-libsql"
)

var (
	dbName       = os_x.Getenv("TURSO_DB") + "-" + os_x.Getenv("TURSO_ORG")
	tursoDbToken = os_x.Getenv("LIBSQL_TOKEN")
)

const (
	sqlite_ddl = `
	CREATE TABLE IF NOT EXISTS ai_record
	(
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		session  TEXT    NOT NULL DEFAULT '',
		question TEXT    NOT NULL DEFAULT '',
		answer   TEXT    NOT NULL DEFAULT '',
		model    TEXT    NOT NULL DEFAULT '',
		time     TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
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
		fmt.Print("\n\n\033[32m请输入您的问题(q:退出,models:切换模型,session:切换保存文件)" + curModel.Id + ": \033[0m")
		question, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "q" {
			break
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
		file.WriteString(fmt.Sprintf("---\n\n### %s(%s) %s\n", question, curModel.Id, now))

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
				if !stream {
					content, _ = render.Render(content)
				}
				fmt.Print(content)
			}
		}
		if err := scanner.Err(); err != nil {
			slog.Warn("读取失败:" + err.Error())
			continue
		}
		runs.Go(func() {
			saveToLibsql(curSession, question, answer.String(), curModel.Id, now)
		})
	}
}

func saveToLibsql(session, question, answer, model, now string) {
	if len(dbName) == 0 || dbName == "-" {
		return
	}
	if db == nil {
		home, _ := common_x.Dir()
		dbPath := filepath.Join(home+"/data/sqlite/libsql", dbName)

		connector, _ := libsql.NewEmbeddedReplicaConnector(dbPath, fmt.Sprintf("libsql://%s.turso.io", dbName),
			libsql.WithAuthToken(tursoDbToken), libsql.WithSyncInterval(time.Second*10),
		)

		db = sql.OpenDB(connector)
	}
	_, err := db.Exec("insert into ai_record(session, question, answer, model, time) values(?,?,?,?,?)", session, question, answer, model, now)
	if err != nil {
		logger.Warn("写数据异常", "err", err.Error())
		if strings.Contains(err.Error(), "no such table") {
			db.Exec(sqlite_ddl)
		}
	}
}
