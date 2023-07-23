package factor3

//go:generate stringer -type=ValueType
//go:generate jsonenums -type=ValueType
type ValueType uint

const (
	ValueTypeString ValueType = iota + 1
	ValueTypeNumber
)
