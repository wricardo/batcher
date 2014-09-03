package redis_flusher

import (
	"github.com/garyburd/redigo/redis"
	"github.com/wricardo/batcher"
)

type RPushFlusher struct {
	List string
}

func (this RPushFlusher) Flush(conn redis.Conn, f batcher.Flushable) {
	strings := f.Strings()
	commands := make([]interface{}, len(strings)+1)
	commands[0] = this.List
	for k, v := range strings {
		commands[k+1] = v
	}
	_, err := conn.Do("rpush", commands...)
	if err != nil {
		panic(err)
	}
}

func NewRPushFlusher(list_name string, flush_buffer_size int, pool RedisPool) batcher.Flusher {
	rlf := RPushFlusher{List: list_name}
	rf := NewRedisFlusher(pool, rlf)
	return batcher.NewDefaultFlusher(flush_buffer_size, rf)
}
