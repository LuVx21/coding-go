package iterators

import (
    "github.com/luvx21/coding-go/coding-common/reflects"
)

type CursorIterator[ITEM, ID any] struct {
    InitCursor       ID
    CheckFirstCursor bool
    //RateLimiter
    DataAccessor    func(ID) []ITEM
    CursorExtractor func([]ITEM) ID
    EndChecker      func(ID) bool
    //-----
    currentCursor    ID
    currentData      []ITEM
    currentDataIndex int
}

func NewCursorIterator[ITEM, ID any](
    initCursor ID,
    checkFirstCursor bool,
    dataAccessor func(ID) []ITEM,
    cursorExtractor func([]ITEM) ID,
    endChecker func(ID) bool,
) *CursorIterator[ITEM, ID] {
    result := &CursorIterator[ITEM, ID]{
        InitCursor:       initCursor,
        CheckFirstCursor: checkFirstCursor,
        DataAccessor:     dataAccessor,
        CursorExtractor:  cursorExtractor,
        EndChecker:       endChecker,
    }
    // 检查游标
    if checkFirstCursor && endChecker(initCursor) {
        result.currentCursor = initCursor
        result.currentData = make([]ITEM, 0)
        result.currentDataIndex = -1
        return result
    }

    // 读取数据
    currentCursor := initCursor
    currentData := dataAccessor(currentCursor)
    var currentDataIndex = -1
    if !reflects.IsNil(currentData) {
        if len(currentData) > 0 {
            currentDataIndex = 0
        }
        // 数据结果提取下次执行的游标
        currentCursor = cursorExtractor(currentData)
    }
    result.currentCursor = currentCursor
    result.currentData = currentData
    result.currentDataIndex = currentDataIndex
    return result
}

func (it *CursorIterator[ITEM, ID]) HasNext() bool {
    _index := it.currentDataIndex
    if it.currentDataIndex >= 0 && _index < len(it.currentData) {
        return true
    }
    // 上次拉取的数据已经迭代结束, 需再次拉取数据
    it.roll()
    index := it.currentDataIndex
    return index >= 0 && index < len(it.currentData)
}

func (it *CursorIterator[ITEM, ID]) Next() ITEM {
    item := it.currentData[it.currentDataIndex]
    it.currentDataIndex++
    return item
}

func (it *CursorIterator[ITEM, ID]) roll() {
    // 再次拉取前检查游标
    if it.EndChecker(it.currentCursor) {
        it.currentData = make([]ITEM, 0)
        it.currentDataIndex = -1
        return
    }
    it.currentData = it.DataAccessor(it.currentCursor)
    if reflects.IsNil(it.currentData) {
        it.currentDataIndex = -1
    } else {
        if len(it.currentData) > 0 {
            it.currentDataIndex = 0
        }
        it.currentCursor = it.CursorExtractor(it.currentData)
    }
}

func (it *CursorIterator[ITEM, ID]) ForEachRemaining(f func(ITEM)) {
    for it.HasNext() {
        f(it.Next())
    }
}
