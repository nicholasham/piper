package core

type Iterator interface {
	HasNext() bool
	Next() interface{}
}
