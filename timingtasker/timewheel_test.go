package timingtasker

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_newTimeWheel_run(t *testing.T) {

	tWheel := newTimeWheel(time.Second, []int{8, 8, 8})
	count := make(map[int]struct{})
	lock := &sync.Mutex{}
	var waits []*sync.WaitGroup
	runTimes := 8 * 8 * 8
	for i := 1; i <= runTimes; i++ {
		wait := &sync.WaitGroup{}
		wait.Add(1)
		waits = append(waits, wait)
		tWheel.addTask(&task{
			key: strconv.Itoa(i),
			callback: func() {
				defer lock.Unlock()
				defer wait.Done()
				lock.Lock()
				count[i] = struct{}{}

			},
			afterScale: int64(i),
		})
	}
	// 模拟每一秒
	for i := 1; i <= runTimes; i++ {
		tWheel.execTask()
		waits[i-1].Wait()
	}

	for i := 1; i <= runTimes; i++ {
		if _, ok := count[i]; !ok {
			t.Errorf("i = %d", i)
		}
	}
}

func getCountMap(tWheel *timeWheel) map[string]int {
	mp := make(map[string]int)
	for oneWheel := tWheel.minWheel; oneWheel != nil; oneWheel = oneWheel.nextWheel {
		for i := 0; i < oneWheel.wheelSize; i++ {
			taskList := oneWheel.slots[i]
			if taskList != nil {
				for elem := taskList.Front(); elem != nil; elem = elem.Next() {
					task, ok := elem.Value.(*task)
					if !ok {
						continue
					}
					mp[task.key]++
				}
			}
		}
	}
	return mp
}

func Test_newTimeWheel_del(t *testing.T) {
	tWheel := newTimeWheel(time.Second, []int{3, 3, 3})
	runTimes := 3 * 3 * 3
	for i := 1; i <= runTimes; i++ {
		tWheel.addTask(&task{
			key:        strconv.Itoa(i),
			callback:   func() {},
			afterScale: int64(i),
		})
	}

	// 模拟重复添加
	for i := 1; i <= runTimes; i++ {
		tWheel.addTask(&task{
			key:        strconv.Itoa(i),
			callback:   func() {},
			afterScale: int64(i),
		})
	}

	// 检查是否重复
	mp := getCountMap(tWheel)
	for i := 1; i <= runTimes; i++ {
		if count := mp[strconv.Itoa(i)]; count != 1 {
			t.Errorf("i = %d, count = %d", i, count)
		}
		if _, ok := tWheel.record[strconv.Itoa(i)]; !ok {
			t.Errorf("i = %d", i)
		}
	}

	// 模拟删除
	for i := 1; i <= runTimes; i = i + 2 {
		tWheel.delTask(strconv.Itoa(i))
	}
	// 检查是否删除
	mp1 := getCountMap(tWheel)
	for i := 1; i <= runTimes; i = i + 2 {
		if count := mp1[strconv.Itoa(i)]; count == 1 {
			t.Errorf("i = %d, count = %d", i, count)
		}
		if _, ok := tWheel.record[strconv.Itoa(i)]; ok {
			t.Errorf("i = %d", i)
		}
	}
}
