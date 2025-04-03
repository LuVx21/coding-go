package common_kv_dao

type CommonKVBizType = int32

const (
	UNKNOWN CommonKVBizType = iota
	LONG
	STRING
	MAP
	BEAN
	LIST
	ARRAY
	INDEX_SPIDER
	INDEX_SPIDER_NEW
	LOCK
)
