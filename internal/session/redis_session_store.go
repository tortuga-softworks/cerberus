package session

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisSessionStore struct {
	client          *redis.Client
	sessionDuration uint64
}

func NewRedisSessionStore(addr, password string, sessionDuration uint64) (*RedisSessionStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, CacheError{err.Error()}
	}

	return &RedisSessionStore{client, sessionDuration}, nil
}

func (rss *RedisSessionStore) CreateSession(username string) (*Session, error) {
	sessionID, sessionIDErr := generateSessionID()

	if sessionIDErr != nil {
		return nil, CacheError{sessionIDErr.Error()}
	}

	creationTime := time.Now()
	session := Session{ID: sessionID, Username: username, CreationTime: creationTime}
	sessionKey := "session:" + sessionID
	expiration := rss.sessionDuration * uint64(time.Second)

	err := rss.client.Set(sessionKey, session, time.Duration(expiration)).Err()
	if err != nil {
		return nil, CacheError{err.Error()}
	}

	return &session, nil
}

func (rss *RedisSessionStore) RefreshSession(sessionID string) error {
	sessionKey := "session:" + sessionID
	expiration := time.Duration(rss.sessionDuration) * time.Second

	refreshed, err := rss.client.Expire(sessionKey, expiration).Result()

	if err != nil {
		return CacheError{err.Error()}
	}

	if !refreshed {
		return &SessionNotFoundError{sessionID}
	}

	return nil
}
