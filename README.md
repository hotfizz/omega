# omega

## install

```shell
go get github.com/hotfizz/omega@latest
```

```shell
go install github.com/hotfizz/omega@latest
```

[examples](examples/simple/main.go)

```go
package main

import (
 "fmt"
 "os"

 "github.com/hotfizz/omega/alt"
)

func objectHelper() (objects []interface{}) {
 objects = append(objects, 1)
 objects = append(objects, 1, "foo")
 objects = append(objects, 1, "foo", "baz")
 objects = append(objects, 1, "foo", "baz", map[string]interface{}{"key": "base"})

 objects = append(objects, 1, "foo", "baz",
  map[string]interface{}{
   "key":  "base",
   "list": []interface{}{1, "99"},
  })

 return objects
}

func main() {
 objects := objectHelper()
 parse := alt.NewDataEtlParser(
  alt.SetIgnore(map[string]struct{}{}), // set ignore (key) list
  // set logger, level is debug, output device is console
  alt.SetLogger(alt.NewStdLogger(alt.LevelDebug, os.Stdout)),
  alt.SetMaxDepth(alt.Infinity),           // infinity depth, warn stack
  alt.SetSeparator(alt.StrSeparator("_")), // use underline as separator
 )

 for _, obj := range objects {
  result, err := parse.Parse(obj)
  fmt.Println(err)
  fmt.Println(result)
 }
}

```

用于将对象扁平化输出的库

常见的场景: 比如 elastic, mongo，通用 🕷 API 返回的 JSON 数据转为有结构化的数据

**特别是将数仓 ETL 过程中，将嵌套的数据结构** 保存到结构化存储引擎中

比如， 下面这个例子:

```json
{
  "name": "map",
  "data": [
    {
      "user_name": "小明",
      "age": 18,
      "province": "广东"
    },
    {
      "user_name": "小海",
      "age": 17,
      "province": "海南"
    }
  ],
  "data2": [
    {
      "persons": [
        {
          "address": "广东"
        },
        {
          "address": "海南"
        }
      ]
    }
  ]
}
```

这样的数据是没有办法入库的(clickhouse, doris, mysql) 这种需要预先定义结构的数据库，所以，需要将这种半结构化的数据转为结构化的

上面的数据会解析成下面这样(一共四条)

```json
[
  {
    "data.age": 18,
    "data.province": "广东",
    "data.user_name": "小明",
    "data2.persons.address": "广东",
    "name": "map"
  },
  {
    "data.age": 18,
    "data.province": "广东",
    "data.user_name": "小明",
    "data2.persons.address": "海南",
    "name": "map"
  },
  {
    "data.age": 17,
    "data.province": "海南",
    "data.user_name": "小海",
    "data2.persons.address": "广东",
    "name": "map"
  },
  {
    "data.age": 17,
    "data.province": "海南",
    "data.user_name": "小海",
    "data2.persons.address": "海南",
    "name": "map"
  }
]
```

只有转成这种结构之后，后面的处理程序才能更好地将这四条数据入库

[具体的实现](alt/doc/normalize_readme.md)
