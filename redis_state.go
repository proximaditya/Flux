package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Use your Redis URL
	})
}

// Mark user as online
func SetUserOnline(userID string) {
	rdb.Set(ctx, "presence:"+userID, "online", 0)
}

// Mark user as offline
func SetUserOffline(userID string) {
	rdb.Del(ctx, "presence:"+userID)
}

// Add user to a room and return the new total count
func JoinRoomState(roomID, userID string) int64 {
	rdb.SAdd(ctx, "room:"+roomID+":members", userID)
	count, _ := rdb.SCard(ctx, "room:"+roomID+":members").Result()
	return count
}

// Remove user from room
func LeaveRoomState(roomID, userID string) {
	rdb.SRem(ctx, "room:"+roomID+":members", userID)
}
