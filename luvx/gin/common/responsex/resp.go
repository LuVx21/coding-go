package responsex

import (
	"context"
	"fmt"
	"net/http"

	"github.com/luvx21/coding-go/coding-common/web"
	"luvx/gin/common/errorx"

	"github.com/gin-gonic/gin"
)

var info = web.WebProvider()

type response struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    any            `json:"data"`
	Host    map[string]any `json:"host"`
	TraceId string         `json:"traceId"`
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
			Host:    info,
		}
	case error:
		return &response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Host:    info,
		}
	}

	return &response{
		Message: http.StatusText(http.StatusOK),
		Data:    data,
		Host:    info,
	}
}
