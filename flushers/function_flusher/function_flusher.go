package function_flusher

import (
	"github.com/wricardo/batcher/buffer"
)

type FlushFunction func(buffer.Flushable) error

type FunctionFlusher struct {
	f FlushFunction
}

func NewFunctionFlusher(flush_buffer_size int, f FlushFunction) buffer.Flusher {
	ff := FunctionFlusher{
		f: f,
	}
	return buffer.NewDefaultFlusher(flush_buffer_size, ff)
}

func (this FunctionFlusher) Flush(to_flush buffer.Flushable) (err error) {
	this.f(to_flush)
	return err
}
