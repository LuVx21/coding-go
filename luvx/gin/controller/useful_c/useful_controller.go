package useful_c

import (
	"luvx/gin/common/responsex"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"github.com/luvx21/coding-go/coding-common/slices_x"
)

func Compare(c *gin.Context) {
	_json := make(map[string]any)
	_ = c.BindJSON(&_json)
	left, right := _json["left"].([]any), _json["right"].([]any)

	responsex.R(c, a.SAM{
		"left_right": slices_x.Diff(left, right),
		"join":       slices_x.Intersect(left, right),
		"right_left": slices_x.Diff(right, left),
	})
}
