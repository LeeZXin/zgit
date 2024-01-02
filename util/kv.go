package util

type KeyVal interface {
	Key() string
	Val() string
	FromStore(string) error
}
