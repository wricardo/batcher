package main

import (
"sync"
	"time"
	"github.com/wricardo/batcher"
	"github.com/garyburd/redigo/redis"
)

type MyString string
func (this MyString) String() string{
	return string(this)
}

func main() {
	pool := newPool(":6379")
	rpush := batcher.NewRPush(pool, "list", 1000, time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for x := 0; x< 10000; x++ {
			rpush.Collect(MyString("df"))
		}
		wg.Done()
	}()
	wg.Wait()
	rpush.Shutdown()
}

func newPool(server string) *redis.Pool {
    return &redis.Pool{
        MaxIdle: 3,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
}
