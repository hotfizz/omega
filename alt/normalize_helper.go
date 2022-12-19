package alt

import (
	"fmt"
	"reflect"
	"sort"
)

type pair struct {
	Key   string
	Value interface{}
}

type kvPairs []pair

func (c kvPairs) Len() int {
	return len(c)
}

func (c kvPairs) Less(i, j int) bool {
	if c[i].Key < c[j].Key {
		return true
	}
	if c[i].Key > c[j].Key {
		return false
	}

	// 相等
	if c[i].Value == nil {
		return true
	}
	if c[j].Value == nil {
		return false
	}

	if reflect.TypeOf(c[i].Value).Kind() != reflect.TypeOf(c[j].Value).Kind() {
		panic(fmt.Sprintf("compare in different type %s %s", reflect.TypeOf(c[i].Value).Kind(), reflect.TypeOf(c[j].Value).Kind()))
	}
	return objectCompare(c[i].Value, c[j].Value)
}

func (c kvPairs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type matrixKvPairs [][]pair

func (c matrixKvPairs) Len() int {
	return len(c)
}

func (c matrixKvPairs) Less(i, j int) bool {
	ii, jj := 0, 0
	for ; ii < len(c[i]) && jj < len(c[j]) && c[i][ii].Key == c[j][jj].Key && reflect.DeepEqual(c[i][ii].Value, c[j][jj].Value); ii, jj = ii+1, jj+1 {

	}

	if c[i][ii].Key < c[j][jj].Key {
		return true
	} else if c[i][ii].Key > c[j][jj].Key {
		return false
	}

	return objectCompare(c[i][ii].Value, c[j][jj].Value)
}

func (c matrixKvPairs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func covertHelper(objList []map[string]interface{}) (res matrixKvPairs) {
	for _, obj := range objList {
		var pairs = make(kvPairs, 0, len(obj))
		for k, v := range obj {
			pairs = append(pairs, pair{Key: k, Value: v})
		}
		sort.Sort(pairs)
		res = append(res, pairs)
	}
	sort.Sort(res)
	return res
}

func objectCompare(o1, o2 interface{}) bool {
	if o1 == nil {
		return true
	}
	if o2 == nil {
		return false
	}
	// 两者都不为空
	if reflect.TypeOf(o1).Kind() != reflect.TypeOf(o2).Kind() {
		panic(fmt.Sprintf("can not compare with two type %s %s", reflect.TypeOf(o1).Kind(), reflect.TypeOf(o2).Kind()))
	}

	fmt.Println("o1 == ", o1, ", o2 == ", o2)
	switch v1 := o1.(type) {
	case bool:
		switch v2 := o2.(type) {
		case bool:
			return v1 == false || v1 == v2
		}

	case int:
		switch v2 := o2.(type) {
		case int:
			return v1 <= v2
		}

	case int8:
		switch v2 := o2.(type) {
		case int8:
			return v1 <= v2
		}

	case int16:
		switch v2 := o2.(type) {
		case int16:
			return v1 <= v2
		}

	case int32:
		switch v2 := o2.(type) {
		case int32:
			return v1 <= v2
		}

	case int64:
		switch v2 := o2.(type) {
		case int64:
			return v1 <= v2
		}
	case uint:
		switch v2 := o2.(type) {
		case uint:
			return v1 <= v2
		}
	case uint8:
		switch v2 := o2.(type) {
		case uint8:
			return v1 <= v2
		}
	case uint16:
		switch v2 := o2.(type) {
		case uint16:
			return v1 <= v2
		}
	case uint32:
		switch v2 := o2.(type) {
		case uint32:
			return v1 <= v2
		}

	case float32:
		switch v2 := o2.(type) {
		case float32:
			return v1 <= v2
		}

	case float64:
		switch v2 := o2.(type) {
		case float64:
			return v1 <= v2
		}

	case uint64:
		switch v2 := o2.(type) {
		case uint64:
			return v1 <= v2
		}

	case string:
		switch v2 := o2.(type) {
		case string:
			return v1 <= v2
		}
	}
	return false
}
