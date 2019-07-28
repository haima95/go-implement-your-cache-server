package cacheClient

import (
	"github.com/go-redis/redis"
)

type redisClient struct {
	*redis.Client
}

func (r *redisClient) get(key string) (string,error) {
	res,err := r.Get(key).Result()
	if err == redis.Nil {
		return "",nil
	}
	return res,err
}

func (r *redisClient) set(key,value string) error {
	return r.Set(key,value,0).Err()
}

func (r *redisClient) del(key string) error {
	return r.Del(key).Err()
}

func (r *redisClient) Run(cmd *Cmd) {

	if cmd.Name == "get" {
		cmd.Value ,cmd.Error= r.get(cmd.Key)
		return
	}
	if cmd.Name == "set" {
		cmd.Error = r.set(cmd.Key,cmd.Value)
		return
	}
	if cmd.Name == "del" {
		cmd.Error = r.del(cmd.Key)
		return
	}
	panic("unknown cmd name "+cmd.Name)
}

func (r *redisClient) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	pipe := r.Pipeline()   // 批量处理
	cmders := make([]redis.Cmder,len(cmds))
	for i,c := range cmds {
		if c.Name == "get" {
			cmders[i] = pipe.Get(c.Key)
		}else if c.Name == "set" {
			cmders[i] = pipe.Set(c.Key,c.Value,0)
		}else if c.Name == "del" {
			cmders[i] = pipe.Del(c.Key)
		}else {
			panic("unknown cmd name "+c.Name)
		}
	}
	_ , err := pipe.Exec()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	for i,c := range cmds {
		if c.Name == "get" {
			value,err := cmders[i].(*redis.StringCmd).Result()
			if err == nil {
				value,err = "",nil
			}
			c.Value,c.Error = value,err
		}else {
			c.Error = cmders[i].Err()
		}
	}
}

func newRedisClient(server string) *redisClient {
	return &redisClient{redis.NewClient(&redis.Options{Addr:server+":6379",ReadTimeout:-1})}
}
