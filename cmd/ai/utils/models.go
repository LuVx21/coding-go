package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/ios"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/infra/ai"
)

var (
	allModels     map[int32]*ai.Model
	allModelsInfo string
)

func SelectModel(curModel *ai.Model) *ai.Model {
	if curModel != nil {
		return curModel
	}
	if len(allModels) == 0 {
		allModels, allModelsInfo = loadModels("")
	}
	fmt.Println(allModelsInfo)

	for {
		fmt.Print("\n\033[32m请选择模型编号(默认1): \033[0m")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		input = strings.TrimSpace(input)

		no := common_x.IfThen(input == "", 1, cast_x.ToInt32(input))
		curModel = allModels[no]
		if curModel != nil {
			fmt.Println("模型已切换为:", curModel.Id, "服务商:", curModel.Sp.BaseUrl)
			break
		} else {
			fmt.Println("模型不存在, 请重新选择, No:", no)
			continue
		}
	}
	return curModel
}

func loadModels(path string) (map[int32]*ai.Model, string) {
	if path == "" {
		home, _ := os.UserHomeDir()
		pathList := []string{"", home + "/data/ai/config/", home + "/config/ai/config/"}
		for _, _path := range pathList {
			fileInfo, err := os.Stat(_path + "models.json")
			if err == nil && !fileInfo.IsDir() {
				path = _path + "models.json"
				break
			}
		}
	}
	fmt.Println("模型配置文件使用:", path)
	_json, _ := ios.ReadAll(path)
	var sps []ai.ServiceProvider
	json.Unmarshal([]byte(_json), &sps)

	rows, m := make([]table.Row, 0), make(map[int32]*ai.Model)
	var no int32 = 1
	for _, sp := range sps {
		sp.ApiKey = os_x.Getenv(sp.ApiKey)
		for _, id := range sp.ModelIds {
			m[no] = &ai.Model{Id: id, Sp: &sp}
			rows = append(rows, table.Row{no, id, sp.BaseUrl})
			no++
		}
	}

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	t.SetTitle("所有模型")
	t.AppendHeader(table.Row{"No", "模型", "服务商"})
	t.AppendRows(rows, rowConfigAutoMerge)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, AutoMerge: true, Align: text.AlignRight},
	})
	t.Style().Options.SeparateRows = true
	return m, t.Render()
}
