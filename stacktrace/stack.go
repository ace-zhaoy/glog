package stacktrace

import (
	"bytes"
	"runtime"
	"strconv"
)

const (
	_pcSize     = 64
	_bufferSize = 1024
)

var _stackPool = NewPool(func() *Stack {
	return &Stack{
		pcs: make([]uintptr, 0, _pcSize),
	}
})

var _formatterPool = NewPool(func() *Formatter {
	return newFormatter(bytes.NewBuffer(make([]byte, 0, _bufferSize)))
})

type Stack struct {
	pcs    []uintptr
	frames *runtime.Frames
}

func (st *Stack) Free() {
	st.frames = nil
	st.pcs = st.pcs[:0]
	_stackPool.Put(st)
}

func (st *Stack) Count() int {
	return len(st.pcs)
}

func (st *Stack) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

func (st *Stack) String() string {
	formatter := GetFormatter()
	defer formatter.Free()

	formatter.FormatStack(st)
	return formatter.String()
}

type Depth int

func (d Depth) Limit() Depth {
	if d < First || d > MaxDepth {
		return MaxDepth
	}
	return d
}

const (
	Full Depth = iota
	First
)

var MaxDepth Depth = 1024

func Capture(skip int, depth Depth) *Stack {
	stack, maxDepth := _stackPool.Get(), int(depth.Limit())

	if maxDepth <= cap(stack.pcs) {
		stack.pcs = stack.pcs[:maxDepth]
	} else {
		stack.pcs = stack.pcs[:cap(stack.pcs)]
	}

	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	if numFrames < len(stack.pcs) || numFrames == maxDepth {
		stack.pcs = stack.pcs[:numFrames]
	} else {
		pcs, pcsLen := stack.pcs, len(stack.pcs)
		for numFrames == len(pcs) && pcsLen < maxDepth {
			pcsLen = len(pcs) * 2
			if pcsLen > maxDepth {
				pcsLen = maxDepth
			}
			pcs = make([]uintptr, pcsLen)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		stack.pcs = pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

func Take(skip int) string {
	stack := Capture(skip+1, Full)
	defer stack.Free()

	return stack.String()
}

type Formatter struct {
	b        *bytes.Buffer
	nonEmpty bool
}

func GetFormatter() *Formatter {
	return _formatterPool.Get()
}

func newFormatter(b *bytes.Buffer) *Formatter {
	return &Formatter{b: b}
}

func (sf *Formatter) FormatStack(stack *Stack) {
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
	}
}

func (sf *Formatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty {
		sf.b.WriteByte('\n')
	}
	sf.nonEmpty = true
	sf.b.WriteString(frame.Function)
	sf.b.WriteByte('\n')
	sf.b.WriteByte('\t')
	sf.b.WriteString(frame.File)
	sf.b.WriteByte(':')
	sf.b.WriteString(strconv.FormatInt(int64(frame.Line), 10))
}

func (sf *Formatter) Free() {
	sf.b.Reset()
	sf.nonEmpty = false
	_formatterPool.Put(sf)
}

func (sf *Formatter) String() string {
	return sf.b.String()
}
