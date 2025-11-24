package main

import (
	"log"
	"sync"
	"testing"
)

func Test_WaitGroup_01(t *testing.T) {
	var wg sync.WaitGroup
	wg.Go(func() { Task() })
	log.Print("111")
	wg.Wait()
	log.Print("main done")
}
