package mock
import(
	"github.com/garyburd/redigo/redis"
	"time"
	"encoding/json"
)
type MockRedisPool struct {
	Dos           chan string
	CommandsDelay time.Duration
}

func NewMockRedisPool() *MockRedisPool {
	m := new(MockRedisPool)
	m.Dos = make(chan string, 10)
	return m
}

func (this *MockRedisPool) Get() redis.Conn {
	rc := NewMockRedisConnection(this.Dos)
	if int64(this.CommandsDelay) != 0 {
		rc.CommandsDelay = this.CommandsDelay
	}
	return rc
}

type MockRedisConnection struct {
	Dos           chan string
	CommandsDelay time.Duration
}

func NewMockRedisConnection(Dos chan string) *MockRedisConnection {
	m := new(MockRedisConnection)
	m.Dos = Dos
	return m
}

func (this *MockRedisConnection) Close() error {
	return nil
}

func (this *MockRedisConnection) Err() error {
	return nil
}

func (this *MockRedisConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if int64(this.CommandsDelay) != 0 {
		time.Sleep(this.CommandsDelay)
	}
	s := make([]interface{}, len(args)+1)
	s[0] = commandName
	for i, v := range args {
		s[i+1] = v
	}
	encoded, _ := json.Marshal(s)
	this.Dos <- string(encoded)
	return nil, nil
}

func (this *MockRedisConnection) Send(commandName string, args ...interface{}) error {
	return nil
}

func (this *MockRedisConnection) Flush() error {
	return nil
}

func (this *MockRedisConnection) Receive() (reply interface{}, err error) {
	return nil, nil
}

type MockRedisPoolDiscard struct {
	Dos chan string
}

func NewMockRedisPoolDiscard() *MockRedisPoolDiscard {
	m := new(MockRedisPoolDiscard)
	m.Dos = make(chan string, 10)
	return m
}

func (this *MockRedisPoolDiscard) Get() redis.Conn {
	rc := NewMockRedisConnectionDiscard(this.Dos)
	return rc
}

type MockRedisConnectionDiscard struct {
	Dos chan string
}

func NewMockRedisConnectionDiscard(Dos chan string) *MockRedisConnectionDiscard {
	m := new(MockRedisConnectionDiscard)
	m.Dos = Dos
	return m
}

func (this *MockRedisConnectionDiscard) Close() error {
	return nil
}

func (this *MockRedisConnectionDiscard) Err() error {
	return nil
}

func (this *MockRedisConnectionDiscard) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return nil, nil
}

func (this *MockRedisConnectionDiscard) Send(commandName string, args ...interface{}) error {
	return nil
}

func (this *MockRedisConnectionDiscard) Flush() error {
	return nil
}

func (this *MockRedisConnectionDiscard) Receive() (reply interface{}, err error) {
	return nil, nil
}

type TestStructs []TestStruct

func (this TestStructs) Strings() []string {
	tmp := make([]string, 0)
	for _, v := range this {
		tmp2 := v.String()
		tmp = append(tmp, tmp2)
	}
	return tmp
}

type TestStruct struct {
	Field1 string
	Field2 int
}

func (this TestStruct) String() string {
	encoded, _ := json.Marshal(this)
	return string(encoded)
}
