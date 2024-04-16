package io

import (
    "fmt"
    "testing"
    "time"
)

func Test_read(t *testing.T) {
    //r1()
    //r2()
    //r3()
    //r4()
    r5()
}

func Test_write(t *testing.T) {
    start := time.Now()
    //w1()
    w2()
    //w3()
    //w4()
    t0 := time.Now()
    fmt.Printf("Cost time %v\n", t0.Sub(start))
}

//----------------
