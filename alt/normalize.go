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
	Parse(data interface{}) (result []map[string]interface{}, err []error)
}

type Separator interface {
	AppendToPrefix(prefix string, key interface{}) string
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type StrSeparator string

func (c StrSeparator) AppendToPrefix(prefix string, key interface{}) string {
	if prefix == "" {
		return fmt.Sprintf("%v", key)
	}
	return fmt.Sprintf("%s%v%v", prefix, c, key)
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

var _, _ = NewDataEtlParser(
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

func (c *dataEtl) Parse(data interface{}) (result []map[string]interface{}, err []error) {
	if data == nil {
		return nil, append(err, fmt.Errorf("invalid data"))
	}
	return c.normalize(data, "", make(map[string]interface{}), 0, nil)
}

func (c *dataEtl) normalize(
	data interface{},
	prefix string,
	currentMap map[string]interface{},
	depth int,
	err []error,
) (result []map[string]interface{}, resErr []error) {
	if _, find := c.ignore[prefix]; data == nil || (depth > c.maxDepth && c.maxDepth != Infinity) || find {
		msg := fmt.Sprintf(
			"invalid data %v or depth <%d> is bigger than max_depth <%d>, more: if max_depth is <%d> is infinity"+
				"or mose key is ignore",
			data == nil, depth, c.maxDepth, Infinity)
		c.logger.Warn(msg)
		return []map[string]interface{}{cpm(currentMap)}, append(err, fmt.Errorf(msg))
	}

	switch reflect.TypeOf(data).Kind() {
	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
		tmp := cpm(currentMap)
		tmp[prefix] = data
		c.logger.Debug(fmt.Sprintf("current type is %v, value is %v", reflect.TypeOf(data).Kind(), data))
		return append(result, tmp), err
	case reflect.Slice:
		var newList []map[string]interface{}
		var s = reflect.ValueOf(data)
		c.logger.Debug(fmt.Sprintf("type is slice len is %d", s.Len()))
		for i := 0; i < s.Len(); i++ {
			var _res, tmpErr = c.normalize(s.Index(i).Interface(), prefix, cpm(currentMap), depth+1, err)
			newList = append(newList, _res...)
			err = append(err, tmpErr...)
		}
		// 如果列表为空, 返回 currentMap 对象本身的数组
		if len(newList) == 0 {
			newList = append(newList, cpm(currentMap))
		}
		return newList, err
	case reflect.Map:
		c.logger.Debug("type is map ", data)
		var newList = []map[string]interface{}{cpm(currentMap)}
		for mr := reflect.ValueOf(data).MapRange(); mr.Next(); {
			_key := mr.Key().Interface()
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
			case reflect.Map, reflect.Slice, reflect.Struct:
				var copyList = make([]map[string]interface{}, 0)
				c.logger.Debug(fmt.Sprintf("type is %v", reflect.TypeOf(_v)))
				for _, _v2 := range newList {
					var _res, tmpErr = c.normalize(_v, c.separator.AppendToPrefix(prefix, _key), _v2, depth+1, err)
					copyList = append(copyList, _res...)
					err = append(err, tmpErr...)
				}
				newList = copyList
			default:
				c.logger.Warn("unknown type ", reflect.TypeOf(_v).Kind())
			}
		}
		return newList, err
	case reflect.Struct:
		c.logger.Debug("type is struct  ", data)
		var newList = []map[string]interface{}{cpm(currentMap)}
		rv := reflect.ValueOf(data)
		rt := reflect.TypeOf(data)
		for i := 0; i < rv.NumField(); i++ {
			_key := rt.Field(i).Name
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
			case reflect.Map, reflect.Slice, reflect.Struct:
				var copyList = make([]map[string]interface{}, 0)
				c.logger.Debug(fmt.Sprintf("type is %v", reflect.TypeOf(_v)))
				for _, _v2 := range newList {
					var _res, tmpErr = c.normalize(_v, c.separator.AppendToPrefix(prefix, _key), _v2, depth+1, err)
					copyList = append(copyList, _res...)
					err = append(err, tmpErr...)
				}
				newList = copyList
				// 不支持的类型
			case reflect.Chan, reflect.Pointer, reflect.Func, reflect.Complex64, reflect.Complex128, reflect.UnsafePointer:
				c.logger.Warn("unknown type i ", reflect.TypeOf(_v).Kind())
			default:
				c.logger.Warn("unknown type ", reflect.TypeOf(_v).Kind())
			}
		}
		return newList, err

	case reflect.Chan, reflect.Pointer, reflect.Func, reflect.Complex64, reflect.Complex128, reflect.UnsafePointer:
		c.logger.Warn("unknown type i ", reflect.TypeOf(data).Kind())
	default:
		c.logger.Warn("unknown type ", reflect.TypeOf(data).Kind())
	}
	return []map[string]interface{}{cpm(currentMap)}, err
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
