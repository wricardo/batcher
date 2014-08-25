package function_flusher
import(
	"github.com/wricardo/batcher"
)

type FlushFunction func(batcher.Flushable) error

type FunctionFlusher struct {
	f FlushFunction
}

func NewFunctionFlusher(flush_buffer_size int, f FlushFunction) batcher.Flusher{
	ff := FunctionFlusher{
		f: f,
	}
	return batcher.NewDefaultFlusher(flush_buffer_size, ff)
}

func (this FunctionFlusher) Flush(to_flush batcher.Flushable) (err error) {
	this.f(to_flush)
	return err
}
