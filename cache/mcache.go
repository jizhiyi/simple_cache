package cache

import (
	"cache/timingtasker"
	"cache/util"
	"container/list"
	"log"
	"sync"
	"time"

	"github.com/DmitriyVTitov/size"
)

type mCacheItem struct {
	key        string
	value      any
	size       int64 // 数据大小
	expireTime int64 // 过期时间
}

type mCache struct {
	mutex *sync.RWMutex

	maxBytes     int64 // 内存限制
	currentBytes int64 // 当前内存

	dataMap  map[string]*list.Element
	dataList *list.List

	timingTasker timingtasker.TimingTasker
}

func NewMCache() Cache {
	return &mCache{
		mutex:        &sync.RWMutex{},
		dataMap:      make(map[string]*list.Element),
		dataList:     list.New(),
		timingTasker: timingtasker.NewTimingTasker(),
	}
}

func (m *mCache) SetMaxMemory(size string) bool {
	bytes, err := util.ParseMemorySize(size)
	if err != nil {
		log.Printf("invalid memory size string: %v", err)
		return false
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.maxBytes = bytes
	return true
}

func (m *mCache) Set(key string, val interface{}, expire time.Duration) {
	// 数据复杂的话会慢点，放前面
	dataSize := size.Of(val)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 如果存在就删除, 从list里删除, 删除定时任务, 更新当前内存
	if elem, ok := m.dataMap[key]; ok {
		m.dataList.Remove(elem)
		oldItem := elem.Value.(*mCacheItem)
		if oldItem.expireTime != 0 {
			m.timingTasker.DelTask(key)
		}
		m.currentBytes -= oldItem.size
	}

	newItem := &mCacheItem{
		key:   key,
		value: val,
		size:  int64(dataSize),
	}
	if expire > 0 {
		newItem.expireTime = time.Now().Add(expire).Unix()
		m.timingTasker.AddTask(key, expire, func() {
			m.Del(key)
		})
	}
	m.currentBytes += newItem.size
	m.dataMap[key] = m.dataList.PushBack(newItem)

	// lru
	for m.maxBytes != 0 && m.currentBytes > m.maxBytes {
		elem := m.dataList.Front()
		if elem == nil {
			break
		}
		item := elem.Value.(*mCacheItem)
		m.currentBytes -= item.size
		delete(m.dataMap, item.key)
		m.dataList.Remove(elem)
	}
}

func (m *mCache) Get(key string) (interface{}, bool) {
	// 要修改 数据在dataList里的位置，还是的用写锁
	m.mutex.Lock()
	defer m.mutex.Unlock()

	elem, ok := m.dataMap[key]
	if !ok {
		return nil, false
	}
	// 判断下是否过期
	if item := elem.Value.(*mCacheItem); item.expireTime != 0 && time.Now().Unix() > item.expireTime {
		return nil, false
	}
	// 移动到后面
	m.dataList.MoveToBack(elem)
	item := elem.Value.(*mCacheItem)
	return item.value, true
}

func (m *mCache) Del(key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	elem, ok := m.dataMap[key]
	if !ok {
		return false
	}
	m.dataList.Remove(elem)
	delete(m.dataMap, key)
	return true
}

func (m *mCache) Exists(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	elem, ok := m.dataMap[key]
	if !ok {
		return false
	}
	// 判断下是否过期
	if item := elem.Value.(*mCacheItem); item.expireTime != 0 && time.Now().Unix() > item.expireTime {
		return false
	}
	return true
}

func (m *mCache) Flush() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.currentBytes = 0
	m.dataMap = make(map[string]*list.Element)
	m.dataList.Init()
	// 定时器停掉, 直接弄个新的就好了
	m.timingTasker.Stop()
	m.timingTasker = timingtasker.NewTimingTasker()
	return true
}

func (m *mCache) Keys() int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return int64(len(m.dataMap))
}
