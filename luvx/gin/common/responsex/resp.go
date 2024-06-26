package responsex

import (
    "context"
    "fmt"
    "luvx/gin/common/errorx"
    "net/http"

    "github.com/gin-gonic/gin"
)

type response struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data"`
    TraceId string `json:"traceId"`
}

func NoMethod(ctx *gin.Context) {
    R(ctx, errorx.NewCodeError(http.StatusMethodNotAllowed))
}

func NoRoute(ctx *gin.Context) {
    R(ctx, errorx.NewCodeError(http.StatusNotFound))
}

func ServiceUnavailable(ctx *gin.Context) {
    R(ctx, errorx.NewCodeError(http.StatusServiceUnavailable))
}

func R(ctx *gin.Context, data any) {
    result := newResponse(ctx, data)
    result.TraceId = fmt.Sprintf("%s", ctx.Value("traceId"))
    httpStatus := http.StatusOK
    if result.Code >= 100 && result.Code <= 511 {
        httpStatus = result.Code
    }
    ctx.JSON(httpStatus, result)
}

func newResponse(ctx context.Context, data any) *response {
    switch value := data.(type) {
    case *errorx.CodeError:
        return &response{
            Code:    value.Code,
            Message: value.Msg,
        }
    case error:
        return &response{
            Code:    http.StatusInternalServerError,
            Message: http.StatusText(http.StatusInternalServerError),
        }
    }

    return &response{
        Message: http.StatusText(http.StatusOK),
        Data:    data,
    }
}
