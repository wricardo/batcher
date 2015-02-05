package redis_flusher

import (
	"github.com/garyburd/redigo/redis"
	"github.com/wricardo/batcher/buffer"
)

type RPushFlusher struct {
	List string
}

func (this RPushFlusher) Flush(conn redis.Conn, f buffer.Flushable) {
	Batch(conn, "RPUSH", this.List, f)
}

func NewRPushFlusher(list_name string, flush_buffer_size int, pool RedisPool) buffer.Flusher {
	rlf := RPushFlusher{List: list_name}
	rf := NewRedisFlusher(pool, rlf)
	return buffer.NewDefaultFlusher(flush_buffer_size, rf)
}

