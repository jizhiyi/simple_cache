package timingtasker

import (
	"container/list"
)

// 单个任务的结构
type task struct {
	key      string
	callback func()

	// 相对于插入槽位代表的时间的倒计时, 单位是最小刻度
	afterScale int64

	// 控制删除
	l    *list.List
	elem *list.Element
}
