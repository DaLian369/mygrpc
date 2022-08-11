package logic

import (
	"fmt"
	"mygrpc/proto"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	DEFAULT_POOLSIZE = 5
)

type RedisManager struct {
	host      string
	auth      string
	redisPool chan redis.Conn
	pool      *redis.Pool
	timeout   time.Duration
}

func InitCache(cfg *proto.RedisSt, proc *ProcessSt) (err error) {
	hostport := cfg.Hostport
	auth := cfg.Auth
	poolsize := cfg.Poolsize
	timeout := cfg.Timeout
	if hostport == "" {
		return fmt.Errorf("redis hostport is nil")
	}
	proc.cache, err = NewRedisManager(hostport, auth, poolsize, time.Duration(timeout)*time.Second)
	if err != nil {
		return
	}

	return
}

func NewRedisManager(host, auth string, poolsize int, timeout time.Duration) (mgr *RedisManager, err error) {
	if poolsize == 0 {
		poolsize = DEFAULT_POOLSIZE
	}
	mgr = &RedisManager{
		host:      host,
		auth:      auth,
		timeout:   timeout,
		redisPool: make(chan redis.Conn, poolsize),
	}
	mgr.pool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", host,
			redis.DialConnectTimeout(timeout),
			redis.DialReadTimeout(timeout),
			redis.DialWriteTimeout(timeout))
	}, poolsize)
	for i := 0; i < poolsize; i++ {
		conn, err := redis.Dial("tcp", host,
			redis.DialConnectTimeout(timeout),
			redis.DialReadTimeout(timeout),
			redis.DialWriteTimeout(timeout))
		if err != nil {
			return nil, err
		}
		mgr.redisPool <- conn
	}
	return
}

func (w *RedisManager) getConn() redis.Conn {
	return <-w.redisPool
}

func (w *RedisManager) putConn(conn redis.Conn) {
	w.redisPool <- conn
}

// 必须使用redis.String或其他函数包装一下conn.Do的结果才能返回redis.ErrNil错误
func (w *RedisManager) do(action string, args ...interface{}) (string, error) {
	conn := w.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do(action, args...))
}

func (w *RedisManager) GetKeyInterface(key string) (res interface{}, err error) {
	// conn := w.getConn()
	// defer w.putConn(conn)

	res, err = w.do("GET", key)
	return
}

func (w *RedisManager) GetData(key string) (bs []byte, err error) {
	res, err := w.GetKeyInterface(key)
	if err == redis.ErrNil {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	bs, err = redis.Bytes(res, err)
	if err == redis.ErrNil {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return
}

func (w *RedisManager) Set(args ...interface{}) (err error) {
	_, err = w.do("set", args...)
	return err
}

func (w *RedisManager) Del(args ...interface{}) (err error) {
	_, err = w.do("del", args...)
	return err
}

func (w *RedisManager) Eval(script string, keynum int, args ...interface{}) (result interface{}, err error) {
	a := append([]interface{}{script, keynum}, args...)
	_, err = w.do("eval", a...)
	if err == redis.ErrNil {
		err = nil
	}
	return
}
