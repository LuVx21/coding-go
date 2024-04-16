package main

import "github.com/elastic/go-elasticsearch/v8"

var EsClient *elasticsearch.TypedClient

func NewEsClient() *elasticsearch.TypedClient {
    if EsClient != nil {
        return EsClient
    }

    cfg := elasticsearch.Config{
        Addresses: []string{
            "http://localhost:9200",
        },
    }

    EsClient, _ := elasticsearch.NewTypedClient(cfg)
    return EsClient
}
