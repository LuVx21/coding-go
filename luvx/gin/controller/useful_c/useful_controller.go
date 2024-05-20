package useful_c

import (
    "github.com/gin-gonic/gin"
    . "github.com/luvx21/coding-go/coding-common/common_x/alias_x"
    "github.com/luvx21/coding-go/coding-common/slices_x"
    "luvx/gin/common/responsex"
)

func Compare(c *gin.Context) {
    _json := make(map[string]any)
    _ = c.BindJSON(&_json)
    left, right := _json["left"].([]any), _json["right"].([]any)

    intersect := slices_x.Intersect(left, right)
    a := slices_x.Diff(left, right)
    b := slices_x.Diff(right, left)
    responsex.R(c, MapStr2Any{
        "left_right": a,
        "join":       intersect,
        "right_left": b,
    })
}
