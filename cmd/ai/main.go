package main

// https://api-docs.deepseek.com/zh-cn/
import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/luvx21/coding-go/cmds/ai/utils"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/strings_x"
	"github.com/luvx21/coding-go/infra/ai"
)

var (
	curModel  *ai.Model
	file      *os.File
	stream    = cast_x.ToBool(strings_x.FirstNonEmpty(os.Getenv("STREAM"), "true"))
	render, _ = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
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
	newFile("chat")
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
			newFile(strings.TrimSpace(strings.TrimPrefix(question, "session")))
			continue
		}

		file.WriteString(fmt.Sprintf("---\n\n### %s(%s) %s\n", question, curModel.Id, time.Now().Format("2006-01-02 15:04:05")))

		res, err := curModel.Request(question, stream)
		if err != nil {
			fmt.Printf("请求接口失败,模型:%s, 服务商:%s\n", curModel.Id, curModel.Sp.BaseUrl)
			continue
		}
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			line := scanner.Text()
			content := ai.ParseLineContent(line)
			if content != "" {
				file.WriteString(content)
				if !stream {
					content, _ = render.Render(content)
				}
				fmt.Print(content)
			}
		}
		if err := scanner.Err(); err != nil {
			slog.Warn("读取失败:" + err.Error())
		}
	}
}
