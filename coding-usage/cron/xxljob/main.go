package main

import (
    "fmt"
    "github.com/xxl-job/xxl-job-executor-go"
    "log"
    "strconv"
)

const (
    AppName      string = "testJob"
    Token        string = "default_token"
    ClientPort   int    = 8899
    AdminAddress string = "http://127.0.0.1:8090/xxl-job-admin"
)

type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
    fmt.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
    log.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func StartXXLJobClient() {
    executor := xxl.NewExecutor(
        xxl.RegistryKey(AppName),     //执行器名称
        xxl.ServerAddr(AdminAddress), //地址
        xxl.AccessToken(Token),       //请求令牌(默认为空)
        //xxl.ExecutorIp("127.0.0.1"),                //可自动获取
        xxl.ExecutorPort(strconv.Itoa(ClientPort)), //默认9999（非必填）
        xxl.SetLogger(&logger{}),                   //自定义日志
    )
    executor.Init()
    //设置日志查看handler
    executor.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
        return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
            FromLineNum: req.FromLineNum,
            ToLineNum:   2,
            LogContent:  "这个是自定义日志handler",
            IsEnd:       true,
        }}
    })
    //注册任务handler
    executor.RegTask("TestJob1", TestJob1)
    executor.RegTask("TestJob2", TestJob2)
    log.Fatal(executor.Run())
}

func main() {
    StartXXLJobClient()
}
