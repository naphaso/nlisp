package sexp

import (
	"fmt"
)

type Env struct {
	parent *Env
	data   map[uint64]Sexp
}

func NewEnv() *Env {
	return &Env{data: map[uint64]Sexp{}}
}

func (s *Env) Wrap() *Env {
	return &Env{parent: s, data: map[uint64]Sexp{}}
}

func (s *Env) WrapMap(m map[uint64]Sexp) *Env {
	return &Env{parent: s, data: m}
}

func (s *Env) Unwrap() *Env {
	if s.parent != nil {
		return s.parent
	}

	panic("unwrap global env")
}

func (s *Env) GetRawData() map[uint64]Sexp {
	env := s
	data := map[uint64]Sexp{}
	for env != nil {
		for k, v := range env.data {
			data[k] = v
		}
		env = env.parent
	}
	return data
}

func (s *Env) String() string {
	return fmt.Sprintf("envlen:%d", len(s.data))
}

func (s *Env) SetGlobal(key *Symbol, value Sexp) {
	ss := s
	for ss.parent != nil {
		ss = ss.parent
	}
	ss.data[key.Code] = value
}

func (s *Env) SetLocal(key *Symbol, value Sexp) {
	s.data[key.Code] = value
}

func (s *Env) Get(key *Symbol) Sexp {
	if v, ok := s.data[key.Code]; ok {
		return v
	}

	if s.parent != nil {
		return s.parent.Get(key)
	}

	return nil
}
