package utils

type StringMapping[T comparable] struct {
	ToValue  map[string]T
	ToString map[T]string
}

func NewStringMapping[T comparable](mapping map[T]string) *StringMapping[T] {
	ret := &StringMapping[T]{make(map[string]T), mapping}
	for val, str := range mapping {
		ret.ToValue[str] = val
	}
	return ret
}
