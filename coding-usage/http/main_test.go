package main

import (
    "fmt"
    "github.com/go-resty/resty/v2"
    "github.com/levigross/grequests"
    "github.com/parnurzeal/gorequest"
    "github.com/valyala/fasthttp"
    "log"
    "os"
    "testing"
    "time"
)

func Test00(t *testing.T) {
    request := gorequest.New()
    _, body, _ := request.Get("https://httpbin.org/get").
        //RedirectPolicy(redirectPolicyFunc).
        Set("Accept", "application/json").
        End()
    fmt.Println(body)
}

func Test01(t *testing.T) {
    resp, err := grequests.Get("https://httpbin.org/get", nil)
    if err != nil {
        log.Fatalln("Unable to make request:", err)
    }

    fmt.Println(resp.String())
}

func Test02(t *testing.T) {
    client := resty.New()
    resp, err := client.R().
        EnableTrace().
        Get("https://httpbin.org/get")

    fmt.Println("Response Info:")
    fmt.Println("  Error      :", err)
    fmt.Println("  Status Code:", resp.StatusCode())
    fmt.Println("  Status     :", resp.Status())
    fmt.Println("  Proto      :", resp.Proto())
    fmt.Println("  Time       :", resp.Time())
    fmt.Println("  Received At:", resp.ReceivedAt())
    fmt.Println("  Body       :\n", resp)
    fmt.Println()

    fmt.Println("Request Trace Info:")
    ti := resp.Request.TraceInfo()
    fmt.Println("  DNSLookup     :", ti.DNSLookup)
    fmt.Println("  ConnTime      :", ti.ConnTime)
    fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
    fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
    fmt.Println("  ServerTime    :", ti.ServerTime)
    fmt.Println("  ResponseTime  :", ti.ResponseTime)
    fmt.Println("  TotalTime     :", ti.TotalTime)
    fmt.Println("  IsConnReused  :", ti.IsConnReused)
    fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
    fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
    fmt.Println("  RequestAttempt:", ti.RequestAttempt)
    fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}

func Test_03(t *testing.T) {
    readTimeout, _ := time.ParseDuration("5000ms")
    writeTimeout, _ := time.ParseDuration("5000ms")
    maxIdleConnDuration, _ := time.ParseDuration("1h")
    client := &fasthttp.Client{
        ReadTimeout:                   readTimeout,
        WriteTimeout:                  writeTimeout,
        MaxIdleConnDuration:           maxIdleConnDuration,
        NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
        DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
        DisablePathNormalizing:        true,
        // increase DNS cache time to an hour instead of default minute
        Dial: (&fasthttp.TCPDialer{
            Concurrency:      4096,
            DNSCacheDuration: time.Hour,
        }).Dial,
    }

    req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
    defer func() {
        fasthttp.ReleaseRequest(req)
        fasthttp.ReleaseResponse(resp)
    }()
    req.SetRequestURI("https://httpbin.org/get")
    req.Header.SetMethod(fasthttp.MethodGet)
    err := client.Do(req, resp)

    if err == nil {
        fmt.Printf("DEBUG Response: %s\n", resp.Body())
    } else {
        _, _ = fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
    }
}
