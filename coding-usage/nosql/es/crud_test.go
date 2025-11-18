package main

import (
	"fmt"
	"github.com/cch123/elasticsql"
	"testing"
)

var client = NewEsClient()

func Test_createIndex(t *testing.T) {
	createIndex(client)
}

func Test_dsl(t *testing.T) {
	var _sql = `
select * from t_user
where a=1 and x = '三个男人'
and create_time between '2015-01-01T00:00:00+0800' and '2016-01-01T00:00:00+0800'
and process_id > 1 order by id desc limit 100,10
`
	dsl, esType, _ := elasticsql.Convert(_sql)
	fmt.Println(dsl)
	fmt.Println(esType)
}
