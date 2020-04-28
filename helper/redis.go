// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/3/11

package helper

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	EXPIRE  = "EXPIRE"
	HMSET   = "HMSET"
	HGETALL = "HGETALL"
	EXISTS  = "EXISTS"
)

type RedisOption struct {
	Address  string
	Password string
}

func NewRedisPool(opt RedisOption) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", opt.Address)
			if err != nil {
				return nil, err
			}
			if opt.Password != "" {
				if _, err := c.Do("AUTH", opt.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", "0"); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
