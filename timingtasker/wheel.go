package timingtasker

import "container/list"

type wheel struct {
	level     int
	wheelSize int
	// 相当于多少个最小刻度
	scale int64
	// 记录上下级
	prevWheel *wheel
	nextWheel *wheel
	// 每个刻度里保存的任务
	slots []*list.List
	// 当前指针在那个刻度
	curPos int
}

func newWheel(sizes []int) *wheel {
	var newWheel, prevWheel *wheel
	for i, size := range sizes {
		if i == 0 {
			newWheel = newOneWheel(i, size, nil)
			newWheel.scale = 1
			prevWheel = newWheel
		} else {
			tmpWheel := newOneWheel(i, size, prevWheel)
			prevWheel.nextWheel = tmpWheel
			tmpWheel.scale = prevWheel.scale * int64(size)
			prevWheel = tmpWheel
		}
	}
	return newWheel
}

func newOneWheel(level int, size int, prevWheel *wheel) *wheel {
	newWheel := &wheel{
		level:     level,
		wheelSize: size,
		prevWheel: prevWheel,
		slots:     make([]*list.List, size),
		curPos:    0,
	}
	return newWheel
}

func (w *wheel) runCurTask(deleteCallback func(string)) {
	curTaskList := w.slots[w.curPos]
	if curTaskList == nil {
		return
	}
	for elem := curTaskList.Front(); elem != nil; elem = elem.Next() {
		task, ok := elem.Value.(*task)
		if !ok {
			continue
		}
		if task.callback != nil {
			go task.callback()
		}
		deleteCallback(task.key)
	}
	w.slots[w.curPos] = nil
}

// advance 指针步进
func (w *wheel) advance() {
	w.curPos++
	w.curPos %= w.wheelSize
}

func (w *wheel) diffuseTask() {
	curTaskList := w.slots[w.curPos]
	if curTaskList == nil {
		return
	}
	for elem := curTaskList.Front(); elem != nil; elem = elem.Next() {
		task, ok := elem.Value.(*task)
		if !ok {
			continue
		}
		// afterScale 1-60 是在一个槽位
		prevWheel := w.prevWheel
		pushPos := (task.afterScale - 1) / prevWheel.scale
		if prevWheel.slots[pushPos] == nil {
			prevWheel.slots[pushPos] = list.New()
		}
		// afterScale 60 是最后一个槽位
		task.afterScale = (task.afterScale-1)%prevWheel.scale + 1
		task.l = prevWheel.slots[pushPos]
		task.elem = prevWheel.slots[pushPos].PushBack(task)
	}
	w.slots[w.curPos] = nil
}
