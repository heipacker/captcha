// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type redisStore struct {
	sync.RWMutex
	client *redis.Client
	// Expiration time of captchas.
	expiration time.Duration
}

func NewRedisStore(expiration time.Duration) Store {
	s := new(redisStore)
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB

	})
	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println(pong)
	}
	s.expiration = expiration
	s.client = client
	return s
}

func (s *redisStore) Set(id string, digits []byte) {
	s.Lock()
	err := s.client.Set("captcha_"+id, string(digits), s.expiration).Err()
	if err != nil {
		panic(err)
	}
	s.Unlock()
}

func (s *redisStore) Get(id string, clear bool) (digits []byte) {
	val, err := s.client.Get("captcha_" + id).Result()
	if err != nil {
		panic(err)
	}
	digits = []byte(val)
	return
}
