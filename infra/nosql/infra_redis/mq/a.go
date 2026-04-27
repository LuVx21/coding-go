package mq

import "github.com/redis/go-redis/v9"

// IsGroupExistsError 创建消费组时, 是否为已存在的错误
// go-redis 中，如果组已存在，错误消息包含 "BUSYGROUP"
func IsGroupExistsError(err error) bool {
	return err != nil && err.Error() != "" &&
		(err.Error() == "BUSYGROUP Consumer Group name already exists" || redis.HasErrorPrefix(err, "BUSYGROUP"))
}
