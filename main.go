package batcher
import(
	"github.com/wricardo/batcher/buffer"
	"github.com/wricardo/batcher/flushers/redis_flusher"
	"github.com/wricardo/batcher/flushers/function_flusher"
	"time"
)

type Collector interface{
	buffer.Collector
}



func NewRPush(pool redis_flusher.RedisPool, list_name string, buffer_size int, timeout time.Duration) *buffer.BufferedCollector {
	config, _ := buffer.NewBufferedCollectorConfiguration(buffer_size, timeout)
	return buffer.NewBufferedCollector(config, redis_flusher.NewRPushFlusher(list_name, buffer_size, pool))
}

func NewLPush(pool redis_flusher.RedisPool, list_name string, buffer_size int, timeout time.Duration) *buffer.BufferedCollector {
	config, _ := buffer.NewBufferedCollectorConfiguration(buffer_size, timeout)
	return buffer.NewBufferedCollector(config, redis_flusher.NewLPushFlusher(list_name, buffer_size, pool))
}

func NewSAdd(pool redis_flusher.RedisPool, set_name string, buffer_size int, timeout time.Duration) *buffer.BufferedCollector {
	config, _ := buffer.NewBufferedCollectorConfiguration(buffer_size, timeout)
	return buffer.NewBufferedCollector(config, redis_flusher.NewSAddFlusher(set_name, buffer_size, pool))
}

func NewFunction(f function_flusher.FlushFunction, buffer_size int, timeout time.Duration) *buffer.BufferedCollector {
	config, _ := buffer.NewBufferedCollectorConfiguration(buffer_size, timeout)
	return buffer.NewBufferedCollector(config, function_flusher.NewFunctionFlusher(buffer_size, f))
}
