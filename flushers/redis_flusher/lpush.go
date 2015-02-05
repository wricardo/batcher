package redis_flusher

import (
	"github.com/garyburd/redigo/redis"
	"github.com/wricardo/batcher/buffer"
)

type LPushFlusher struct {
	List string
}

func (this LPushFlusher) Flush(conn redis.Conn, f buffer.Flushable) {
	Batch(conn, "LPUSH", this.List, f)
}

func NewLPushFlusher(list_name string, flush_buffer_size int, pool RedisPool) buffer.Flusher {
	rlf := LPushFlusher{List: list_name}
	rf := NewRedisFlusher(pool, rlf)
	return buffer.NewDefaultFlusher(flush_buffer_size, rf)
}
