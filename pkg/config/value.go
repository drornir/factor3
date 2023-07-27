package config

//go:generate stringer -type=ValueType -linecomment
//go:generate jsonenums -type=ValueType
type ValueType uint

const (
	ValueTypeString ValueType = iota + 1 // <string>
	ValueTypeNumber                      // <number>
)
