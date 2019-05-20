package evalstack

import (
	"errors"
	"fmt"

	"github.com/naphaso/nlisp/pkg/sexp"
)

type EvalFrame struct {
	args []sexp.Sexp
	exp  sexp.Sexp
}

type EvalStack struct {
	stack []*EvalFrame
	last  int
}

func NewEvalStack() *EvalStack {
	stack := make([]*EvalFrame, 1000)
	for i := 0; i < 1000; i++ {
		stack[i] = &EvalFrame{args: make([]sexp.Sexp, 0, 10)}
	}
	return &EvalStack{
		stack: stack,
	}
}

func (s *EvalStack) Eval(e sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	s.stack[s.last].exp = e
	return s.Run(env)
}

func (s *EvalStack) Run(env *sexp.Env) (sexp.Sexp, error) {
	for {
		//fmt.Printf("run, stack size %d\n", s.last+1)
		lastFrame := s.stack[s.last]
		if lastFrame.exp == nil {
			r, err := lastFrame.Run()
			if err != nil {
				return nil, err
			}

			if s.last == 0 {
				return r, nil
			}

			s.last--
			s.stack[s.last].PushArg(r)
			continue
		}

		switch e := lastFrame.exp.(type) {
		case *sexp.Pair:
			lastFrame.exp = e.Tail
			// check head for simple types
			// TODO: check to special forms?
			switch e.Head.(type) {
			case sexp.Int64:
				lastFrame.PushArg(e.Head)
			default:
				s.last++
				s.stack[s.last].exp = e.Head
			}
		case sexp.Int64:
			if s.last == 0 {
				return lastFrame.exp, nil
			}
			s.last--
			s.stack[s.last].PushArg(lastFrame.exp)
		case sexp.Symbol:
			// TODO: symbol resolve
			if e.Name() == "plus" {
				if s.last == 0 {
					return lastFrame.exp, nil
				}
				s.last--
				s.stack[s.last].PushArg(lastFrame.exp)
			} else {
				return nil, errors.New("unsupported symbol yet: " + e.Name())
			}
		default:
			return nil, errors.New("unsupported type yet: " + fmt.Sprintf("%T", lastFrame.exp))
		}
	}
}

func (s *EvalFrame) PushArg(e sexp.Sexp) {
	s.args = append(s.args, e)
}

func (s *EvalFrame) Run() (sexp.Sexp, error) {
	if op, ok := s.args[0].(sexp.Symbol); ok {
		if op.Name() == "plus" {
			var r int64
			for i := 1; i < len(s.args); i++ {
				if v, ok := s.args[i].(sexp.Int64); ok {
					r += int64(v)
				} else {
					return nil, errors.New("invalid argument for plus op")
				}
			}
			s.args = s.args[:0]
			return sexp.NewInt64(r), nil
		}
	}

	return nil, errors.New("invalid op: " + sexp.ToString(s.args[0]))
}
