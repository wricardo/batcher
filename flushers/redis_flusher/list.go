package redis_flusher

import (
	"github.com/garyburd/redigo/redis"
	"github.com/wricardo/batcher"
	//	"time"
	//	"errors"
)

type RedisListFlusher struct {
	List string
}

func (this RedisListFlusher) Flush(conn redis.Conn, f batcher.Flushable) {
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

func NewListFlusher(list_name string, flush_buffer_size int, pool RedisPool) batcher.Flusher {
	rlf := RedisListFlusher{List: list_name}
	rf := NewRedisFlusher(pool, rlf)
	return batcher.NewDefaultFlusher(flush_buffer_size, rf)
}
