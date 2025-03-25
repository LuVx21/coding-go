package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/luvx21/coding-go/coding-common/slices_x"
)

const (
	AI_KEY = ""
)

type ServiceProvider struct {
	BaseUrl  string   `json:"baseUrl"`
	ApiKey   string   `json:"apiKey"`
	ModelIds []string `json:"modelIds"`
}

func NewServiceProvider(key, token string, modelIds ...string) *ServiceProvider {
	return &ServiceProvider{BaseUrl: key, ApiKey: token, ModelIds: modelIds}
}

func (sp *ServiceProvider) ToModels() []Model {
	a := slices_x.Transfer(func(id string) Model { return Model{id, sp} }, sp.ModelIds...)
	return a
}

func (sp *ServiceProvider) ToModelsMap() map[string]Model {
	m := make(map[string]Model)
	for _, id := range sp.ModelIds {
		m[id] = Model{id, sp}
	}
	return m
}

func (sp *ServiceProvider) ApiUrl() string {
	return ApiUrl(sp.BaseUrl)
}

func (sp *ServiceProvider) String() string {
	return sp.ApiUrl() + " 模型:\n" + strings.Join(sp.ModelIds, "\n")
}

type Model struct {
	Id string
	Sp *ServiceProvider
}

func (m *Model) Request(question string) (*http.Response, error) {
	r := &ChatRequest{
		Model:       m.Id,
		Messages:    []Message{{"user", question}},
		Stream:      true,
		MaxTokens:   2048,
		Temperature: 0.7,
	}
	payload, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", m.Sp.ApiUrl(), bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+m.Sp.ApiKey)
	return http.DefaultClient.Do(req)
}

func ApiUrl(api string) string {
	if api == "" {
		return ""
	}
	if strings.HasSuffix(api, "#") {
		// 以#结尾, 使用原始地址
		return api[0 : len(api)-1]
	}
	if strings.HasSuffix(api, "/") {
		// 以/结尾, 忽略v1版本
		return api[0:len(api)-1] + CHAT_API_PATH
	}
	return api + OPENAI_CHAT_POSTFIX
}
