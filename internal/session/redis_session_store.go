package session

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisSessionStore struct {
	client          *redis.Client
	sessionDuration int
}

func NewRedisSessionStore(addr, password string, sessionDuration int) (*RedisSessionStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return &RedisSessionStore{client, sessionDuration}, nil
}

func (rss *RedisSessionStore) CreateSession(username string) (*Session, error) {
	sessionID := generateSessionID()
	creationTime := time.Now()
	session := Session{ID: sessionID, Username: username, CreationTime: creationTime}
	sessionKey := "session:" + sessionID
	expiration := rss.sessionDuration * int(time.Second)

	err := rss.client.Set(sessionKey, session, time.Duration(expiration)).Err()
	if err != nil {
		return nil, err
	}

	return &session, nil
}
