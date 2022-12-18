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
		result := parse.Parse(obj)
		fmt.Println(result)
	}
}
