package logs

// https://blog.schwarzeni.com/2019/11/02/logrus%E8%BE%93%E5%87%BA%E6%97%A5%E5%BF%97%E8%87%B3stdout%E5%92%8Cfile/

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_a(t *testing.T) {
	log.WithFields(log.Fields{
		"漫画id": "comicId",
		"响应码":  "resp.StatusCode",
	}).Info("获取所有章节接口:")

	log.WithField("json", "json1").Infoln("foo bar")
}
