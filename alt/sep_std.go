package alt

import "fmt"

type StrSeparator string

func (c StrSeparator) AppendToPrefix(prefix string, key interface{}) string {
	if prefix == "" {
		return fmt.Sprintf("%v", key)
	}
	return fmt.Sprintf("%s%v%v", prefix, c, key)
}
