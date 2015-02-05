package redis_flusher

import (
	"github.com/garyburd/redigo/redis"
	"github.com/wricardo/batcher/buffer"
	//	"errors"
)

type RedisFlusherImplementation interface {
	Flush(redis.Conn, buffer.Flushable)
}

type RedisFlusher struct {
	redis_pool     RedisPool
	implementation RedisFlusherImplementation
}

func NewRedisFlusher(redis_pool RedisPool, implementation RedisFlusherImplementation) *RedisFlusher {
	rf := new(RedisFlusher)
	rf.redis_pool = redis_pool
	rf.implementation = implementation
	return rf
}

func (this *RedisFlusher) Flush(to_flush buffer.Flushable) (err error) {
	conn := this.redis_pool.Get()
	this.implementation.Flush(conn, to_flush)
	return err
}

type RedisPool interface {
	Get() redis.Conn
}
