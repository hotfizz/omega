package alt

type Separator interface {
	AppendToPrefix(prefix string, key interface{}) string
}
