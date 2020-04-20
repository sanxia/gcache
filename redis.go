package gcache

import (
	"github.com/golang/groupcache/singleflight"
	"github.com/sanxia/gredis"
)

/* ================================================================================
 * cache client
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * RedisCache
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type RedisCache struct {
	sf             singleflight.Group
	Redis          gredis.IRedis
	packer         IPack
	isSingleFlight bool
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化RedisCache
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRedisCache(ip string, port int, password string, db, timeout int, prefix string, args ...bool) ICache {
	isSingleFlight := true
	if len(args) > 0 {
		isSingleFlight = args[0]
	}

	return &RedisCache{
		Redis:          gredis.NewRedis(ip, port, password, db, timeout, prefix),
		packer:         NewMessagePack(),
		isSingleFlight: isSingleFlight,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RedisCache) Get(key string, dest interface{}) error {
	var data interface{}
	var err error

	if s.isSingleFlight {
		data, err = s.sf.Do(key, func() (interface{}, error) {
			return s.Redis.Get(key)
		})
	} else {
		data, err = s.Redis.Get(key)
	}

	if err != nil {
		return err
	}

	s.packer.Unmarshal(data.([]byte), &dest)

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 添加数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RedisCache) Set(key string, dest interface{}, second int) error {
	data, err := s.packer.Marshal(dest)
	if err != nil {
		return err
	}

	return s.Redis.Set(key, data, second)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RedisCache) Remove(keys ...string) error {
	return s.Redis.Del(keys...)
}
