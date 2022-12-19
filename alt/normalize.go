package alt

import (
	"fmt"
	"os"
	"reflect"
)

const (
	Infinity = -1
)

type Parser interface {
	Parse(data interface{}) (result []map[string]interface{})
}

func NewDataEtlParser(opt ...OptionFunc) Parser {
	v := &dataEtl{
		maxDepth:  Infinity,
		ignore:    make(map[string]struct{}),
		logger:    NewStdLogger(LevelDebug, os.Stdout),
		separator: StrSeparator("."),
	}
	for _, f := range opt {
		f(v)
	}
	return v
}

type OptionFunc func(c *dataEtl)

func SetLogger(logger Logger) OptionFunc {
	return func(c *dataEtl) {
		c.logger = logger
	}
}

func SetMaxDepth(maxDepth int) OptionFunc {
	return func(c *dataEtl) {
		c.maxDepth = maxDepth
	}
}

func SetIgnore(ignore map[string]struct{}) OptionFunc {
	return func(c *dataEtl) {
		c.ignore = ignore
	}
}

func SetSeparator(separator Separator) OptionFunc {
	return func(c *dataEtl) {
		c.separator = separator
	}
}

var _ = NewDataEtlParser(
	SetMaxDepth(Infinity),
	SetLogger(NewStdLogger(LevelDebug, os.Stdout)),
	SetSeparator(StrSeparator("_")),
	SetIgnore(nil),
).Parse(nil)

type dataEtl struct {
	separator Separator
	logger    Logger
	maxDepth  int
	ignore    map[string]struct{}
}

func (c *dataEtl) Parse(data interface{}) (result []map[string]interface{}) {
	if data == nil {
		c.logger.Debug("parser data is nil, return empty list")
		return make([]map[string]interface{}, 0)
	}
	return c.normalize(data, "", make(map[string]interface{}), 0)
}

func (c *dataEtl) normalize(
	data interface{},
	prefix string,
	currentMap map[string]interface{},
	depth int,
) (result []map[string]interface{}) {
	if data == nil {
		ve := fmt.Sprintf("current key is %s, but value is <nil>", prefix)
		c.logger.Debug(ve)
		return []map[string]interface{}{cpm(currentMap)}
	}

	if _, find := c.ignore[prefix]; find {
		ig := fmt.Sprintf("key is ignored %s", prefix)
		c.logger.Debug(ig)
		return []map[string]interface{}{cpm(currentMap)}
	}

	if c.maxDepth != Infinity && depth > c.maxDepth {
		msg := fmt.Sprintf("exclude depth %s, currentDepth:%d, maxDepth:%d", prefix, depth, c.maxDepth)
		c.logger.Warn(msg)
		return []map[string]interface{}{cpm(currentMap)}
	}

	c.logger.Debug(fmt.Sprintf("parser  type %v, data %v", reflect.TypeOf(data).Kind(), data))
	switch reflect.TypeOf(data).Kind() {
	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
		return c.parsePrimitive(data, prefix, currentMap, depth)

	case reflect.Slice, reflect.Array:
		return c.parseSlice(data, prefix, currentMap, depth)

	case reflect.Map:

		return c.parseMap(data, prefix, currentMap, depth)

	case reflect.Struct:

		return c.parseStruct(data, prefix, currentMap, depth)

	default:

		c.logger.Warn("normalize: unknown type ", reflect.TypeOf(data).Kind())
	}

	return []map[string]interface{}{cpm(currentMap)}
}

// primitive  判断是否为基础类型
func (c *dataEtl) primitive(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
		return true
	default:
		return false
	}
}

// 处理基础类型
func (c *dataEtl) parsePrimitive(
	data interface{},
	prefix string,
	currentMap map[string]interface{},
	depth int,
) (result []map[string]interface{}) {
	tmp := cpm(currentMap)
	if _, ok := c.ignore[prefix]; !ok &&
		((c.maxDepth == Infinity) || (c.maxDepth != Infinity && c.maxDepth > depth)) &&
		c.primitive(reflect.TypeOf(data).Kind()) {
		c.logger.Debug(fmt.Sprintf("current type is %v, value is %v", reflect.TypeOf(data).Kind(), data))
		tmp[prefix] = data
	}
	return append(result, tmp)
}

// 处理切片类型
func (c *dataEtl) parseSlice(
	data interface{},
	prefix string,
	currentMap map[string]interface{},
	depth int,
) (result []map[string]interface{}) {
	var newList []map[string]interface{}
	var s = reflect.ValueOf(data)
	c.logger.Debug(fmt.Sprintf("type is slice len is %d", s.Len()))
	for i := 0; i < s.Len(); i++ {
		var _res = c.normalize(s.Index(i).Interface(), prefix, cpm(currentMap), depth+1)
		newList = append(newList, _res...)
	}
	// 如果列表为空, 返回 currentMap 对象本身的数组
	if len(newList) == 0 {
		newList = append(newList, cpm(currentMap))
	}
	return newList
}

// 解析 map 类型
func (c *dataEtl) parseMap(
	data interface{},
	prefix string,
	currentMap map[string]interface{},
	depth int,
) (result []map[string]interface{}) {
	c.logger.Debug("type is map ", data)
	var newList = []map[string]interface{}{cpm(currentMap)}
	if c.maxDepth != Infinity && depth > c.maxDepth {
		c.logger.Debug(fmt.Sprintf("map: current depth:%d is more than maxDepth:%d", depth, c.maxDepth))
		return newList
	}
	for mr := reflect.ValueOf(data).MapRange(); mr.Next(); {
		_key := mr.Key().Interface()
		// 处理忽略键对象
		if _, ok := c.ignore[c.separator.AppendToPrefix(prefix, _key)]; ok {
			c.logger.Debug("map: current key is ignore %s", c.separator.AppendToPrefix(prefix, _key))
			continue
		}
		_v := mr.Value().Interface()
		c.logger.Debug(fmt.Sprintf("rec map key:%v, value:%v", _key, _v))
		// NOTE 必须保证判断是有效的
		// this case { "data": null }
		if _v == nil {
			c.logger.Warn(_key, " nil type ")
			continue
		}

		switch reflect.TypeOf(_v).Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Bool, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
			// 如果当前的值 数值，整型，布尔时，填充到所有已经遍历的对象
			for i := range newList {
				newList[i][c.separator.AppendToPrefix(prefix, _key)] = _v
			}
		case reflect.Map, reflect.Slice, reflect.Struct, reflect.Array:
			var copyList = make([]map[string]interface{}, 0)
			c.logger.Debug(fmt.Sprintf("type is %v", reflect.TypeOf(_v)))
			for _, _v2 := range newList {
				var _res = c.normalize(_v, c.separator.AppendToPrefix(prefix, _key), _v2, depth+1)
				copyList = append(copyList, _res...)
			}
			newList = copyList
		default:
			c.logger.Warn("unknown type ", reflect.TypeOf(_v).Kind())
		}
	}
	return newList
}

// 解析 struct 类型
func (c *dataEtl) parseStruct(
	data interface{},
	prefix string,
	currentMap map[string]interface{},
	depth int,
) (result []map[string]interface{}) {
	c.logger.Debug("type is struct  ", data)
	var newList = []map[string]interface{}{cpm(currentMap)}
	if c.maxDepth != Infinity && depth > c.maxDepth {
		c.logger.Debug(fmt.Sprintf("struct: current depth:%d is more than maxDepth:%d", depth, c.maxDepth))
		return newList
	}
	rv := reflect.ValueOf(data)
	rt := reflect.TypeOf(data)
	for i := 0; i < rv.NumField(); i++ {
		_key := rt.Field(i).Name
		if _, ok := c.ignore[c.separator.AppendToPrefix(prefix, _key)]; ok {
			c.logger.Debug("struct: current key is ignore %s", c.separator.AppendToPrefix(prefix, _key))
			continue
		}
		_v := rv.Field(i).Interface()
		c.logger.Debug(fmt.Sprintf("rec map key:%v, value:%v", _key, _v))
		// NOTE 必须保证判断是有效的
		// this case { "data": null }
		if _v == nil {
			c.logger.Warn(_key, " nil type ")
			continue
		}

		switch reflect.TypeOf(_v).Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Bool, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
			// 如果当前的值 数值，整型，布尔时，填充到所有已经遍历的对象
			for i := range newList {
				newList[i][c.separator.AppendToPrefix(prefix, _key)] = _v
			}
		case reflect.Map, reflect.Slice, reflect.Struct, reflect.Array:
			var copyList = make([]map[string]interface{}, 0)
			c.logger.Debug(fmt.Sprintf("type is %v", reflect.TypeOf(_v)))
			for _, _v2 := range newList {
				var _res = c.normalize(_v, c.separator.AppendToPrefix(prefix, _key), _v2, depth+1)
				copyList = append(copyList, _res...)
			}
			newList = copyList
			// 不支持的类型
		case reflect.Chan, reflect.Func, reflect.Complex64, reflect.Complex128:
			c.logger.Warn("unknown type i ", reflect.TypeOf(_v).Kind())
		default:
			c.logger.Warn("unknown type ", reflect.TypeOf(_v).Kind())
		}
	}
	return newList
}

// copy map object
// todo this function do not a deep copy
func cpm(src map[string]interface{}) (dst map[string]interface{}) {
	dst = make(map[string]interface{}, len(src))
	for _k, _v := range src {
		dst[_k] = _v
	}
	return dst
}
