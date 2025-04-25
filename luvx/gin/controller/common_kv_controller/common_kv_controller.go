package common_kv_controller

import (
	"luvx/gin/common/responsex"
	"luvx/gin/dao/common_kv_dao"
	"luvx/gin/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

//	func GetCommonKeyValue(c *gin.Context) {
//		kvs := common_kv_dao.Get(5)
//		responsex.R(c, kvs)
//	}
//

func GetCommonKeyValue(c *gin.Context) {
	cursorIDStr := c.Query("cursorId")
	cursorID := 0
	if cursorIDStr != "" {
		if id, err := strconv.Atoi(cursorIDStr); err == nil {
			cursorID = id
		}
	}

	limitStr := c.Query("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	kvs, nextCursorID, err := common_kv_dao.GetByCursor(cursorID, limit, 5)
	if err != nil {
		responsex.R(c, gin.H{"error": err.Error()})
		return
	}

	responsex.R(c, gin.H{
		"data":           kvs,
		"next_cursor_id": nextCursorID,
	})
}

func CreateCommonKeyValue(c *gin.Context) {
	var kv model.CommonKeyValue
	if err := c.ShouldBindJSON(&kv); err != nil {
		responsex.R(c, gin.H{"error": err.Error()})
		return
	}
	if err := common_kv_dao.Create(&kv); err != nil {
		responsex.R(c, gin.H{"error": err.Error()})
		return
	}
	responsex.R(c, gin.H{"message": "success"})
}

func DeleteCommonKeyValue(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		responsex.R(c, gin.H{"error": err.Error()})
		return
	}
	if err := common_kv_dao.Delete(ids); err != nil {
		responsex.R(c, gin.H{"error": err.Error()})
		return
	}
	responsex.R(c, gin.H{"message": "success"})
}

func UpdateCommonKeyValue(c *gin.Context) {
	var kv model.CommonKeyValue
	if err := c.ShouldBindJSON(&kv); err != nil {
		responsex.R(c, gin.H{"error": err.Error()})
		return
	}
	if err := common_kv_dao.Update(&kv); err != nil {
		responsex.R(c, gin.H{"error": err.Error()})
		return
	}
	responsex.R(c, gin.H{"message": "success"})
}
