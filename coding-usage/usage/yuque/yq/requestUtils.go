package yq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 100)
var client = &http.Client{}

const namespace = "xieren/doc-sync"

func init() {
}

func createFile(createBody CreateBody) string {
	url := fmt.Sprintf("https://www.yuque.com/api/v2/repos/%s/docs", namespace)
	jsonBlob, _ := json.MarshalIndent(createBody, "", "    ")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBlob))
	return request(req)
}

func deleteFile(id uint64) {
	url := fmt.Sprintf("https://www.yuque.com/api/v2/repos/%s/docs/%d", namespace, id)
	req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(make([]byte, 0)))
	request(req)
}

func updateFile(id uint64, createBody CreateBody) {
	url := fmt.Sprintf("https://www.yuque.com/api/v2/repos/%s/docs/%d", namespace, id)
	jsonBlob, _ := json.MarshalIndent(createBody, "", "    ")
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBlob))
	request(req)
}

func readFile(slug string) string {
	url := fmt.Sprintf("https://www.yuque.com/api/v2/repos/%s/docs/%s", namespace, slug)
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(make([]byte, 0)))
	body := request(req)

	f := make(map[string]any)
	json.Unmarshal([]byte(body), &f)

	i2 := f["data"].(map[string]any)["body"]
	return i2.(string)
}

func request(req *http.Request) string {
	if !limiter.Allow() {
		//return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "doc-sync")
	req.Header.Set("X-Auth-Token", "UEQyqsil79bVwlylWWicJuvhGCi61lyHCPMTAlKV")

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

type CreateBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Slug  string `json:"slug"`

	Format   string `json:"format"`
	Public   uint32 `json:"public"`
	ForceAsl uint32 `json:"_force_asl"`
}

func NewCreateBody(title string, body string, slug string) CreateBody {
	return CreateBody{Title: title, Body: body, Slug: slug,
		Format: "markdown", Public: 0, ForceAsl: 1}
}
