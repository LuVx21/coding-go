package mongodb

import (
	"context"
	"log"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/event"
)

var LogMonitor = event.CommandMonitor{
	Started: func(ctx context.Context, event *event.CommandStartedEvent) {
		if event.CommandName == "ping" {
			return
		}
		log.Println("---------------------------------------Started---------------------------------------")
		g := gjson.Parse(event.Command.String())
		table := g.Get(event.CommandName).String()
		log.Printf("命令: %s 表: %s.%s sql: %+v", event.CommandName, event.DatabaseName, table, "event.Command")

		switch event.CommandName {
		case "insert":
		case "delete", "update":
			s := g.Get(event.CommandName + "s")
			for _, q := range s.Get("#.q").Array() {
				for k, v := range q.Map() {
					log.Printf("%s -> %s", k, v)
				}
			}
			for _, u := range s.Get("#.u").Array() {
				for k, v := range u.Map() {
					log.Printf("%s -> %s", k, v)
				}
			}
		case "find":
			log.Printf("\nfilter: %s\nlimit: %s\nsort: %s", g.Get("filter"), g.Get("limit"), g.Get("sort"))
		}
		log.Println("-------------------------------------------------------------------------------------")
	},
	Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
		if event.CommandName != "ping" {
			log.Println("---------------------------------------Succeed---------------------------------------")
			log.Printf("查询语句: %s 耗时: %dms", event.CommandName, event.Duration/1000/1000)
			log.Println("-------------------------------------------------------------------------------------")
		}
	},
	Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
		log.Fatalf("查询语句: %s 耗时: %dms", event.CommandName, event.Duration/1000/1000)
	},
}
