package model

type CommonKeyValue struct {
    Id          int64 `gorm:"primary_key"`
    BizType     int32
    CommonKey   string
    CommonValue string
    Invalid     int32
    CreateTime  int64
    UpdateTime  int64
}
