package session

import (
	"context"
	"encoding/json"
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

func (store *RedisSessionStore) Create(ctx context.Context, userID string) (*Session, error) {
	sessionID, sessionIDErr := generateSessionID()

	if sessionIDErr != nil {
		return nil, CacheError{sessionIDErr.Error()}
	}

	creationTime := time.Now()
	session := Session{ID: sessionID, UserID: userID, CreationTime: creationTime}
	sessionKey := "session:" + sessionID
	expiration := time.Duration(store.sessionDuration * uint64(time.Second))

	err := store.client.Set(sessionKey, session, expiration).Err()
	if err != nil {
		return nil, CacheError{err.Error()}
	}

	return &session, nil
}

func (store *RedisSessionStore) FindByID(ctx context.Context, sessionID string) (*Session, error) {
	sessionKey := "session:" + sessionID
	data, err := store.client.Get(sessionKey).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, &SessionNotFoundError{SessionID: sessionID}
		} else if err != nil {
			return nil, &CacheError{err.Error()}
		}
	}

	var session Session
	err = json.Unmarshal([]byte(data), &session)

	if err != nil {
		return nil, &CacheError{err.Error()}
	}

	return &session, nil
}

func (store *RedisSessionStore) Refresh(ctx context.Context, sessionID string) error {
	sessionKey := "session:" + sessionID
	expiration := time.Duration(store.sessionDuration * uint64(time.Second))

	refreshed, err := store.client.Expire(sessionKey, expiration).Result()

	if err != nil {
		return CacheError{err.Error()}
	}

	if !refreshed {
		return &SessionNotFoundError{sessionID}
	}

	return nil
}

func (store *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	sessionKey := "session:" + sessionID

	err := store.client.Del(sessionKey).Err()

	if err != nil {
		return CacheError{err.Error()}
	}

	return nil
}
