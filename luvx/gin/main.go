package gin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/logs"
	"luvx/gin/config"
	"luvx/gin/middleware"
	"luvx/gin/router"
	"luvx/gin/runner"
)

func WebStart() {
	logs.Log.Infoln("ʕ◔ϖ◔ʔ 启动... ʕ◔ϖ◔ʔ")
	runner.Start()

	r := gin.Default()

	router.Register(r)
	middleware.RegisterGlobalMiddlewares(r)

	port := config.AppConfig.Server.Port
	if err := r.Run(port); err != nil {
		fmt.Printf("Start server error,err=%v", err)
	}

	srv := &http.Server{
		Addr:           port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.Log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefulExitWeb(srv)
}

func gracefulExitWeb(srv *http.Server) {
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) // 此处不会阻塞
	sig := <-quit                                                         // 阻塞在此，当接收到上述两种信号时才会往下执行
	logs.Log.Infof("退出服务%v\n", sig)

	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		logs.Log.Warnln("退出时异常:", err)
	}

	logs.Log.Infof("ʕ◔ϖ◔ʔ 已退出... ʕ◔ϖ◔ʔ 耗时:%v\n", time.Since(now))
}
