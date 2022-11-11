package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hotfizz/omega/alt"
)

var (
	pf = flag.String("pf", "examples/alt/assets", "文件夹")
)

func main() {
	flag.Parse()
	var paths = make([]string, 0)
	_ = filepath.Walk(*pf, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		paths = append(paths, path)
		return err
	})
	etl := alt.NewDataEtlParser(alt.SetLogger(alt.NewStdLogger(alt.LevelWarn, os.Stdout)))
	for _, path := range paths {
		result, err := etl.Parse(readFile(path))
		fmt.Println("file ", path)
		for _, r := range result {

			fmt.Println(r)
		}
		fmt.Println(err)
		//fmt.Println(err)
	}

}

func readFile(path string) (v interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	if json.Valid(data) {
		//var l = make([]interface{}, 0)
		//if err := json.Unmarshal(data, &v); err == nil {
		//	fmt.Println(l)
		//	return l
		//}
		var m = make(map[string]interface{})
		if err := json.Unmarshal(data, &m); err == nil {
			//fmt.Println(m)
			return m
		}
	}
	return make(map[string]interface{})
}
