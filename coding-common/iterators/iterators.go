package iterators

import (
    "github.com/luvx21/coding-go/coding-common/func_x"
    "github.com/luvx21/coding-go/coding-common/reflects"
)

type CursorIterator[ITEM, ID, ITEMS any] struct {
    InitCursor       ID
    CheckFirstCursor bool
    //RateLimiter
    DataAccessor    func(curId ID) ITEMS           // 根据游标查询数据
    CursorExtractor func(curId ID, items ITEMS) ID // 根据当前的游标和查询结果,计算下次迭代的游标
    DataExtractor   func(items ITEMS) []ITEM       // 根据查询结果,提取其中需要的数据
    EndChecker      func(curId ID) bool            // 检查游标是否终止迭代
    //-----
    currentCursor    ID
    currentData      []ITEM
    currentDataIndex int
}

func NewPageIteratorSimple[ITEM any](
    initCursor int, checkFirstCursor bool,
    dataAccessor func(curId int) []ITEM,
    endChecker func(curId int) bool,
) *CursorIterator[ITEM, int, []ITEM] {
    return NewPageIterator[ITEM, []ITEM](initCursor, checkFirstCursor, dataAccessor, func_x.Identity[[]ITEM], endChecker)
}

func NewPageIterator[ITEM, ITEMS any](
    initCursor int, checkFirstCursor bool,
    dataAccessor func(curId int) ITEMS,
    dataExtractor func(items ITEMS) []ITEM,
    endChecker func(curId int) bool,
) *CursorIterator[ITEM, int, ITEMS] {
    cursorExtractor := func(id int, items ITEMS) int {
        return id + 1
    }
    return NewCursorIterator[ITEM, int, ITEMS](initCursor, checkFirstCursor, dataAccessor, cursorExtractor, dataExtractor, endChecker)
}

func NewCursorIteratorSimple[ITEM, ID any](
    initCursor ID, checkFirstCursor bool,
    dataAccessor func(curId ID) []ITEM,
    cursorExtractor func(curId ID, items []ITEM) ID,
    endChecker func(curId ID) bool,
) *CursorIterator[ITEM, ID, []ITEM] {
    return NewCursorIterator[ITEM, ID, []ITEM](initCursor, checkFirstCursor, dataAccessor, cursorExtractor,
        func_x.Identity[[]ITEM], endChecker,
    )
}

// NewCursorIterator ITEM: 切片中具体存储的元素, ID:查询元素所用的游标,
// ITEMS:包含了ITEM的切片,还可能包含下次迭代的游标,如果没有则下次迭代游标是计算得到, 此时这个泛型是不需要的
func NewCursorIterator[ITEM, ID, ITEMS any](
    initCursor ID, checkFirstCursor bool,
    dataAccessor func(curId ID) ITEMS,
    cursorExtractor func(curId ID, items ITEMS) ID,
    dataExtractor func(items ITEMS) []ITEM,
    endChecker func(curId ID) bool,
) *CursorIterator[ITEM, ID, ITEMS] {
    result := &CursorIterator[ITEM, ID, ITEMS]{
        InitCursor:       initCursor,
        CheckFirstCursor: checkFirstCursor,
        DataAccessor:     dataAccessor,
        CursorExtractor:  cursorExtractor,
        DataExtractor:    dataExtractor,
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
    items := dataAccessor(currentCursor)
    currentData := dataExtractor(items)
    var currentDataIndex = -1
    if !reflects.IsNil(currentData) {
        if len(currentData) > 0 {
            currentDataIndex = 0
        }
        // 数据结果提取下次执行的游标
        currentCursor = cursorExtractor(currentCursor, items)
    }
    result.currentCursor = currentCursor
    result.currentData = currentData
    result.currentDataIndex = currentDataIndex
    return result
}

func (it *CursorIterator[ITEM, ID, ITEMS]) HasNext() bool {
    _index := it.currentDataIndex
    if it.currentDataIndex >= 0 && _index < len(it.currentData) {
        return true
    }
    // 当前拉取的数据已经迭代结束, 需再次拉取数据
    it.roll()
    index := it.currentDataIndex
    return index >= 0 && index < len(it.currentData)
}

func (it *CursorIterator[ITEM, ID, ITEMS]) Next() ITEM {
    item := it.currentData[it.currentDataIndex]
    it.currentDataIndex++
    return item
}

func (it *CursorIterator[ITEM, ID, ITEMS]) roll() {
    // 再次拉取前检查游标
    if it.EndChecker(it.currentCursor) {
        it.currentData = make([]ITEM, 0)
        it.currentDataIndex = -1
        return
    }
    items := it.DataAccessor(it.currentCursor)
    it.currentData = it.DataExtractor(items)
    if reflects.IsNil(it.currentData) {
        it.currentDataIndex = -1
    } else {
        if len(it.currentData) > 0 {
            it.currentDataIndex = 0
        }
        it.currentCursor = it.CursorExtractor(it.currentCursor, items)
    }
}

func (it *CursorIterator[ITEM, ID, ITEMS]) ForEachRemaining(f func(ITEM)) {
    for it.HasNext() {
        f(it.Next())
    }
}
