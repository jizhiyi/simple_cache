package timingtasker

import "time"

type TimingTasker interface {
	// AddTask 添加定时任务 expire 后执行
	AddTask(key string, after time.Duration, callback func())
	// DelTask 删除定时任务
	DelTask(key string)
	// Stop 停止
	Stop()
}

// NewTimingTasker 创建定时任务器 时间刻度是1s, 有效时间是 60*60*24*365*100 (大约100年)
func NewTimingTasker() TimingTasker {
	timeWheel := newTimeWheel(time.Second, []int{60, 60, 24, 365, 100})
	go timeWheel.run()
	return timeWheel
}

// NewTimingTaskerWithInterval 自定义
func NewTimingTaskerWithInterval(interval time.Duration, wheelSizes []int) TimingTasker {
	timeWheel := newTimeWheel(interval, wheelSizes)
	go timeWheel.run()
	return timeWheel
}
