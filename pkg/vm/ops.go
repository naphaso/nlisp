package vm

type Op byte

const (
	OpLoadNil Op = iota
	OpLoadTrue
	OpLoadFalse
	OpLoadInt64
	OpLoadString
)
