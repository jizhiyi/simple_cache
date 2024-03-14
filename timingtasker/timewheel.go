package timingtasker

import (
	"container/list"
	"time"
)

// 时间轮结构
type timeWheel struct {
	// 最小的表盘
	minWheel *wheel

	// 精度
	interval time.Duration

	// 定时器
	ticker *time.Ticker
	// 记录任务 key -> timeWheelTask
	record map[string]*task
	// 修改任务
	addChan chan *task
	delChan chan string
	// 暂停
	stopChannel chan struct{}
}

func newTimeWheel(interval time.Duration, wheelSizes []int) *timeWheel {
	t := &timeWheel{
		minWheel: newWheel(wheelSizes),
		interval: interval,
		ticker:   time.NewTicker(interval),
		record:   make(map[string]*task),
		addChan:  make(chan *task),
		delChan:  make(chan string),
	}
	return t
}

func (t *timeWheel) run() {
	for {
		select {
		case <-t.ticker.C:
			t.execTask()
		case task := <-t.addChan:
			t.addTask(task)
		case key := <-t.delChan:
			t.delTask(key)
		case <-t.stopChannel:
			return
		}
	}
}

func (t *timeWheel) execTask() {
	for tmpWheel := t.minWheel; tmpWheel != nil; tmpWheel = tmpWheel.nextWheel {
		if tmpWheel.level == 0 {
			// 第一轮的任务， 直接丢去执行
			tmpWheel.runCurTask(func(key string) {
				delete(t.record, key)
			})
			tmpWheel.advance()
		} else {
			tmpWheel.diffuseTask()
			tmpWheel.advance()
		}
		// tmpWheel.curPos = 0 表示需要去下一层, 扩散任务
		if tmpWheel.curPos != 0 {
			break
		}
	}
}

// addTask
func (t *timeWheel) addTask(task *task) {
	if _, ok := t.record[task.key]; ok {
		t.delTask(task.key)
	}
	var sumScale int64
	for tmpWheel := t.minWheel; tmpWheel != nil; tmpWheel = tmpWheel.nextWheel {
		pushPos := -1
		for i := tmpWheel.curPos; i < tmpWheel.wheelSize; i++ {
			sumScale += tmpWheel.scale
			if sumScale >= task.afterScale {
				pushPos = i
				break
			}
		}
		if pushPos != -1 {
			if tmpWheel.slots[pushPos] == nil {
				tmpWheel.slots[pushPos] = list.New()
			}
			// 修改afterScale 为相对于要插入的槽位的倒计时
			task.afterScale = task.afterScale - sumScale + tmpWheel.scale
			task.l = tmpWheel.slots[pushPos]
			task.elem = tmpWheel.slots[pushPos].PushBack(task)
			break
		}
	}
	t.record[task.key] = task
}

// delTask
func (t *timeWheel) delTask(key string) {
	if task, ok := t.record[key]; ok {
		if task.l != nil && task.elem != nil {
			task.l.Remove(task.elem)
		}
		delete(t.record, key)
	}
}

// AddTask 添加任务, 只是添加到 chan 上
func (t *timeWheel) AddTask(key string, after time.Duration, callback func()) {
	if callback == nil {
		return
	}
	if after < t.interval {
		go callback()
		return
	}
	t.addChan <- &task{key: key, afterScale: int64(after / t.interval), callback: callback}
}

// DelTask 添加任务, 添加到 chan 上, 没有上时间轮
func (t *timeWheel) DelTask(key string) {
	t.delChan <- key
}

func (t *timeWheel) Stop() {
	t.stopChannel <- struct{}{}
}
