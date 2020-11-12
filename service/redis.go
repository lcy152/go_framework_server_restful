package service

import (
	"context"
	"log"
	"tumor_server/model"

	redisService "tumor_server/redis"

	"github.com/garyburd/redigo/redis"
)

func UpdateUser(ctx context.Context, user *model.User) error {
	sc := GetContainerInstance()
	err := sc.DB.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	err = sc.RedisService.DeleteKey(user.ID.String())
	if err != nil {
		return err
	}
	return nil
}

func InitRedis(r *redisService.RedisService) {
	go func() {
		conn := r.RedisPool.Get()
		defer conn.Close()
		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe("logout")
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Printf("%s: logout session_id: %s\n", v.Channel, v.Data)
				// TODO: send message to client
			case redis.Subscription:
				break
			case error:
				return
			}
		}
	}()
	go func() {
		conn := r.RedisPool.Get()
		defer conn.Close()
		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe("__keyevent@0__:expired")
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				flag := TokenTAG
				sessionID := string(v.Data)
				length := len(sessionID)
				if length > len(flag) {

				}
			case redis.Subscription:
				break
			case error:
				return
			}
		}
	}()
}
