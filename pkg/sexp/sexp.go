package sexp

import (
	"errors"
	"unsafe"
)

type Sexp interface {
	SexpString() string
	Equal(o Sexp) bool
}

func ToString(v Sexp) string {
	if v == nil {
		return "nil"
	}

	return v.SexpString()
}

// deprecated
func IsNil(v Sexp) bool {
	// slow implementation with reflect
	//return i == nil || reflect.ValueOf(i).IsNil()

	// slow implementation with type switch
	//switch vv := v.(type) {
	//case *Pair:
	//	return vv == nil
	//case *String:
	//	return vv == nil
	//case *Bool:
	//	return vv == nil
	//case *Symbol:
	//	return vv == nil
	//case *Int64:
	//	return vv == nil
	//case *Float64:
	//	return vv == nil
	//case *Rune:
	//	return vv == nil
	//default:
	//	return i == nil || reflect.ValueOf(i).IsNil()
	//}

	// fast unsafe implementation
	return (*[2]uintptr)(unsafe.Pointer(&v))[1] == 0
}

func IsList() bool {
	panic("not implemented")
}

func ToArray(e Sexp) ([]Sexp, error) {
	var el []Sexp
	var curr *Pair
	var ok bool
	for e != nil {
		curr, ok = e.(*Pair)
		if !ok {
			return nil, errors.New("not a list")
		}

		el = append(el, curr.Head)
		e = curr.Tail
	}

	return el, nil
}

func NewList(elements []Sexp) Sexp {
	var root Sexp
	for i := len(elements) - 1; i >= 0; i-- {
		root = NewPair(elements[i], root)
	}
	return root
}

func List(args ...Sexp) Sexp {
	return NewList(args)
}

func NewListBuilder() *ListBuilder {
	return &ListBuilder{}
}

type ListBuilder struct {
	head *Pair
	tail *Pair
}

func (s *ListBuilder) Append(v Sexp) {
	node := NewPair(v, nil)
	if s.head == nil {
		s.head = node
		s.tail = node
	} else {
		s.tail.Tail = node
		s.tail = node
	}
}

func (s *ListBuilder) Build() Sexp {
	return s.head
}
