package main

import "testing"

var client = NewEsClient()

func Test_createIndex(t *testing.T) {
    createIndex(client)
}

