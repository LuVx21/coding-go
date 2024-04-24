package common_kv

import (
    . "github.com/luvx21/coding-go/coding-common/common_x/pairs"
    "luvx/gin/dao/common_kv"
    "luvx/gin/model"
)

type CommonKVBizType int32
type KeyWithType = Pair[string, int]

const (
    UNKNOWN CommonKVBizType = iota
    LONG
    STRING
    MAP
    BEAN
    LIST
    ARRAY
    INDEX_SPIDER
)

var statusDescriptions = map[CommonKVBizType]string{}

func Get(bizType CommonKVBizType, keys ...string) map[string]*model.CommonKeyValue {
    kvs := common_kv.Get(int32(bizType), keys...)
    m := make(map[string]*model.CommonKeyValue)
    for _, kv := range kvs {
        m[kv.CommonKey] = kv
    }
    return m
}
