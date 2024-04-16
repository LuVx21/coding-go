package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type R struct {
    WeatherInfo Weather `json:"weatherinfo"`
}

type Weather struct {
    city string
    temp string
}

func aaa() {
    url := "http://www.weather.com.cn/data/sk/101010100.html"
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    fmt.Println(string(body))

    s := `{"weatherinfo":{"city":"北京","cityid":"101010100","temp":"27.9","WD":"南风","WS":"小于3级","SD":"28%","AP":"1002hPa","njd":"暂无实况","WSE":"<3","time":"17:55","sm":"2.1","isRadar":"1","Radar":"JC_RADAR_AZ9010_JB"}}`
    var f R
    json.Unmarshal([]byte(s), &f)
    fmt.Printf("%+v\n", f)
}
