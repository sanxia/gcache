package gcache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

/* ================================================================================
 * cache client
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * LocalCache
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type LocalCache struct {
	sync.Mutex
	list       *list.List
	elements   map[string]*list.Element
	expiration time.Duration
	maxLength  int
	packer     IPack
}

type CacheData struct {
	key   string
	value []byte
	date  time.Time
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化LocalCache
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewLocalCache(maxLength int, expiration time.Duration) *LocalCache {
	return &LocalCache{
		list:       list.New(),
		elements:   make(map[string]*list.Element, maxLength),
		expiration: expiration,
		maxLength:  maxLength,
		packer:     NewMessagePack(),
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *LocalCache) Get(key string, dest interface{}) error {
	s.Lock()
	defer s.Unlock()

	element, isOk := s.elements[key]
	if !isOk {
		return errors.New("key not found")
	}

	cacheData := element.Value.(*CacheData)
	if time.Since(cacheData.date) > s.expiration {
		s.deleteElement(element)
		return errors.New("data is expire")
	}

	s.list.MoveToFront(element)

	s.packer.Unmarshal(cacheData.value, &dest)

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 添加数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *LocalCache) Set(key string, dest interface{}, second int) error {
	s.Lock()
	defer s.Unlock()

	data, err := s.packer.Marshal(dest)
	if err != nil {
		return err
	}

	if element, isOk := s.elements[key]; isOk {
		//更新
		oldCacheData := element.Value.(*CacheData)
		oldCacheData.value = data
		oldCacheData.date = time.Now().Add(time.Duration(second) * time.Second)

		s.list.MoveToFront(element)
	} else {
		//新增
		cacheData := &CacheData{
			key:   key,
			value: data,
			date:  time.Now().Add(time.Duration(second) * time.Second),
		}

		s.elements[key] = s.list.PushFront(cacheData)

		//淘汰数据
		for s.list.Len() > s.maxLength {
			s.deleteElement(s.list.Back())
		}
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *LocalCache) Remove(keys ...string) bool {
	s.Lock()
	defer s.Unlock()

	for _, key := range keys {
		if element, isOk := s.elements[key]; isOk {
			s.deleteElement(element)
		} else {
			return false
		}
	}

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取缓存长度
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *LocalCache) Len() int {
	return s.list.Len()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除元素
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *LocalCache) deleteElement(element *list.Element) {
	s.list.Remove(element)
	delete(s.elements, element.Value.(*CacheData).key)
}
