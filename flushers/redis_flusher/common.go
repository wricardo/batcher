package redis_flusher

import(
	"github.com/garyburd/redigo/redis"
	"github.com/wricardo/batcher/buffer"
)

func Batch(conn redis.Conn, redis_command string, redis_key string, f buffer.Flushable) {
	strings := f.Strings()
	commands := make([]interface{}, len(strings)+1)
	commands[0] = redis_key
	for k, v := range strings {
		commands[k+1] = v
	}
	_, err := conn.Do(redis_command, commands...)
	if err != nil {
		panic(err)
	}
}
