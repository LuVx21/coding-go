package life

import (
	"context"
	"encoding/json"
	"log/slog"
	"luvx/gin/common/consts"
	"luvx/gin/db"
	"luvx/gin/service"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	reportJsonCli = db.MongoMainCon().Database("life").Collection("report_json")
	reportCli     = db.MongoMainCon().Database("life").Collection("report")

	fields = sets.NewSet("LabRepItemName", "ResultType", "ResultText1", "ReferenceText", "LabRepItemUnit", "DangerFlag", "ChangeFlag", "LabFlow", "LabTime")
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		service.NewRunner("read_report", "job not need", time.Minute*10, readLifeJson),
	}
}

func readLifeJson() {
	opts := options.Find().SetSort(bson.M{"day": -1})
	rows, _ := mongodb.RowsMap(context.TODO(), reportJsonCli, bson.M{}, opts)

	reportCli.DeleteMany(context.TODO(), bson.M{})
	for i, row := range *rows {
		day, name, _json := row["day"].(string), row["name"].(string), row["json"].(string)
		if len(_json) == 0 {
			continue
		}
		slog.Info("report进度...", "No", i+1, "day", day, "name", name)
		datas := readData(day, name, _json)
		_, err := reportCli.InsertMany(context.TODO(), datas)
		if err != nil {
			slog.Error("mongo写入错误", "error", err)
		}
	}
}

func readData(day, name, _json string) []any {
	var root a.Row
	_ = json.Unmarshal([]byte(_json), &root)

	result := make([]any, 0, 16)
	for _, a := range root["result"].([]any) {
		m := a.(map[string]any)
		maps_x.RemoveIf(m, func(k string, _ any) bool { return !fields.Contains(k) })
		m["_id"], m["day"], m["name"] = consts.IdWorker.NextId(), day, name
		result = append(result, m)
	}

	return result
}
