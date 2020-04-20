package gcache

/* ================================================================================
 * cache client
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ICache接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	ICache interface {
		Get(key string, dest interface{}) error
		Set(key string, dest interface{}, second int) error
		Remove(keys ...string) error
	}

	IPack interface {
		Marshal(dest interface{}) ([]byte, error)
		Unmarshal(data []byte, dest interface{}) error
	}
)
