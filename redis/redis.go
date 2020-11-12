package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisService struct {
	IP            string
	RedisPool     *redis.Pool
	ExpireTime    int
	BlackListTime int
}

func ConnectRedis(ip string, expireTime, blackListTime int) *RedisService {
	redisPool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ip)
		},
	}
	rs := &RedisService{
		IP:         ip,
		RedisPool:  redisPool,
		ExpireTime: expireTime,
	}
	return rs
}

func (r *RedisService) AddKey(key string, value string) error {
	conn := r.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	ret, err := redis.Int(conn.Do("EXPIRE", key, strconv.Itoa(r.ExpireTime)))
	if err != nil {
		return err
	}
	if ret != 1 {
		return errors.New("expire session error")
	}
	return nil
}

func (r *RedisService) GetKey(key string) (string, error) {
	conn := r.RedisPool.Get()
	defer conn.Close()
	jsonStr, err := redis.String(conn.Do("GET", key))
	return jsonStr, err
}

func (r *RedisService) DeleteKey(key string) error {
	conn := r.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PUBLISH", "logout", key)
	if err != nil {
		return err
	}
	_, err = conn.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) ExpireKey(key string) error {
	conn := r.RedisPool.Get()
	defer conn.Close()
	ret, err := redis.Int(conn.Do("EXPIRE", key, strconv.Itoa(r.ExpireTime)))
	if err != nil {
		return err
	}
	if ret != 1 {
		return errors.New("expire session error")
	}
	return nil
}
