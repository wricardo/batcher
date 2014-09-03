package redis_flusher

import (
	"github.com/garyburd/redigo/redis"
	"github.com/wricardo/batcher"
)

type LPushFlusher struct {
	List string
}

func (this LPushFlusher) Flush(conn redis.Conn, f batcher.Flushable) {
	strings := f.Strings()
	commands := make([]interface{}, len(strings)+1)
	commands[0] = this.List
	for k, v := range strings {
		commands[k+1] = v
	}
	_, err := conn.Do("lpush", commands...)
	if err != nil {
		panic(err)
	}
}

func NewLPushFlusher(list_name string, flush_buffer_size int, pool RedisPool) batcher.Flusher {
	rlf := LPushFlusher{List: list_name}
	rf := NewRedisFlusher(pool, rlf)
	return batcher.NewDefaultFlusher(flush_buffer_size, rf)
}
