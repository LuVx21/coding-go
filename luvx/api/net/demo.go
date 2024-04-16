package main

import (
    "encoding/json"
    "fmt"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    _ "luvx/api/logs"
    "net/http"
    "time"
)

/**
https://www.cnblogs.com/zhaof/p/11346412.html
*/
func main() {
    //ids := [9]int32{2, 3, 4, 10, 11, 17, 18, 20, 22}
    ids := [1]int32{26}
    for _, id := range ids {
        chapterIdsJson := getAllChaptersByComicId(id)
        var f []interface{}
        json.Unmarshal(chapterIdsJson, &f)
        for _, v := range f {
            if kvs, ok := v.(map[string]interface{}); ok {
                getAllChapterPhotos(id, kvs)
            }
        }
    }
}

func getAllComicList() {
    url := "http://m.dolulu.com.cn/query/querycartoon/QueryCartoonListnolimit?apikey=xxxx"
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}

func getAllChaptersByComicId(comicId int32) []byte {
    url := "http://m.dolulu.com.cn/query/querycartoon/querychapterlistnolimit?apikey=xxxx&cartoon_id=%d"
    url = fmt.Sprintf(url, comicId)
    resp, _ := http.Get(url)

    log.WithFields(log.Fields{
        "漫画id": comicId,
        "响应码":  resp.StatusCode,
    }).Info("获取所有章节接口:")

    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    time.Sleep(time.Duration(2) * time.Second)

    return body
}

func getAllChapterPhotos(comicId int32, kvs map[string]interface{}) {
    chapterId := int64(kvs["chapter_id"].(float64))
    url := "http://m.dolulu.com.cn/query/querycartoon/querychaptercontentnolimit?apikey=xxxx&cartoon_id=%d&chapter_id=%d"
    url = fmt.Sprintf(url, comicId, chapterId)
    resp, _ := http.Get(url)

    log.WithFields(log.Fields{
        "响应码":  resp.StatusCode,
        "漫画id": comicId,
        "章节id": chapterId,
        "顺序":   kvs["order"],
    }).Info("章节信息->")

    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    var f []interface{}
    json.Unmarshal(body, &f)

    for i, v := range f {
        if kvs1, ok := v.(map[string]interface{}); ok {
            url := kvs1["url"].(string)
            //log.WithFields(log.Fields{
            //   "响应码":     resp.StatusCode,
            //   "漫画id":    comicId,
            //   "章节id":    chapterId,
            //   "顺序":      kvs["order"],
            //   "图片index": i+1,
            //   "图片url":   url,
            //}).Info("章节图片->")

            r, _ := http.Get(url)
            if r.StatusCode >= 300 {
                log.WithFields(log.Fields{
                    "响应码":     r.StatusCode,
                    "漫画id":    comicId,
                    "章节id":    chapterId,
                    "顺序":      kvs["order"],
                    "图片index": i + 1,
                    "图片url":   url,
                }).Info("访问不通的图片->")
            }
        }
    }

    time.Sleep(time.Duration(2) * time.Second)
}
