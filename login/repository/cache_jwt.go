// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/2/1

package repository

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"login-service/helper"
	"log"
	"time"
)

type redisRepository struct {
	sessions *redis.Pool
}

func NewSessionRepository(opt helper.RedisOption) *redisRepository {
	return &redisRepository{
		sessions: helper.NewRedisPool(opt),
	}
}

func (r *redisRepository) SetSession(sid string, fields map[string]interface{}, ttl int64) error {
	var err error
	conn := r.sessions.Get()

	_, err = conn.Do(helper.HMSET, redis.Args{}.Add(sid).AddFlat(fields)...)
	if err != nil {
		return err
	}
	if ttl > 0 {
		_, err = conn.Do(helper.EXPIRE, sid, ttl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *redisRepository) GetSession(sid string) (map[string]string, error) {
	var err error
	conn := r.sessions.Get()

	do, err := conn.Do(helper.HGETALL, sid)
	if err != nil {
		log.Fatalf("redis hmget %v fail err:%v", sid, err)
		panic(err)
	}
	logrus.Infof("cache session %v by %v", do, sid)
	result, err := redis.StringMap(do, err)
	if err != nil {
		log.Fatalf("redis hmget %v fail err:%v", sid, err)
		panic(err)
	}

	return result, nil
}

func (r *redisRepository) DoubleCheckSession(phone string) error {
	var err error
	conn := r.sessions.Get()

	exist, err := redis.Int(conn.Do(helper.EXISTS, phone))
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New("不可以重新请求")
	}
	return nil
}

// redis delete  key (key should be string string string ......)
func (r *redisRepository) Delete(key ...interface{}) error {
	var err error
	conn := r.sessions.Get()
	_, err = conn.Do("DEL", key...)
	// If delete fails, try again after 100 milliseconds
	if err != nil {
		time.Sleep(time.Millisecond * 100)
		_, err = conn.Do("DEL", key...)
		if err != nil {
			logrus.Warnf("delete redis key : %v failed : %v \n", key, err)
		}
	}
	return err
}

// redis delete  key (key should be string string string ......)
func (r *redisRepository) Get(key string) (string, error) {
	var err error
	conn := r.sessions.Get()
	res, err := conn.Do("GET", key)
	// If delete fails, try again after 100 milliseconds
	value, err := redis.String(res, err)
	if err != nil {
		logrus.Warnf("get key from redis failed , key : %v , error : 【%v】 \n", key, err.Error())
	}
	return value, err
}
