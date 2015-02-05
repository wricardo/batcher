package redis_flusher

import (
	"github.com/wricardo/batcher/buffer"
	"github.com/garyburd/redigo/redis"
)

type SAddFlusher struct {
	Set string
}

func (this SAddFlusher) Flush(conn redis.Conn, f buffer.Flushable) {
	Batch(conn, "SADD", this.Set, f)
}

func NewSAddFlusher(list_name string, flush_buffer_size int, pool RedisPool) buffer.Flusher {
	rlf := SAddFlusher{Set: list_name}
	rf := NewRedisFlusher(pool, rlf)
	return buffer.NewDefaultFlusher(flush_buffer_size, rf)
}

